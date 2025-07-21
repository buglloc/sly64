package router

import (
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type WatchNotifier func()

type Source interface {
	Domains(ctx context.Context) ([]string, error)
	Watch(notify WatchNotifier)
	Close() error
}

var _ Source = (*StaticSource)(nil)

type StaticSource struct {
	domains []string
}

func NewStaticSource(domains []string) (*StaticSource, error) {
	return &StaticSource{
		domains: domains,
	}, nil
}

func (s *StaticSource) Domains(_ context.Context) ([]string, error) {
	return s.domains, nil
}

func (s *StaticSource) Watch(_ WatchNotifier) {}

func (s *StaticSource) Close() error {
	return nil
}

type FileSource struct {
	filepath       string
	domains        []string
	mu             sync.Mutex
	watchers       []WatchNotifier
	lastModTime    time.Time
	lastSum        [md5.Size]byte
	reloadInterval time.Duration
	closed         chan struct{}
	ctx            context.Context
	shutdownFn     context.CancelFunc
}

func NewFileSource(filepath string, opts ...FileSourceOption) (*FileSource, error) {
	ctx, shutdownFn := context.WithCancel(context.Background())
	s := &FileSource{
		filepath:   filepath,
		closed:     make(chan struct{}),
		ctx:        ctx,
		shutdownFn: shutdownFn,
	}

	for _, opt := range opts {
		opt(s)
	}

	if _, err := s.sync(); err != nil {
		return nil, fmt.Errorf("sync fail: %w", err)
	}

	go s.loop()
	return s, nil
}

func (s *FileSource) Domains(_ context.Context) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.domains, nil
}

func (s *FileSource) Watch(notify WatchNotifier) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.watchers = append(s.watchers, notify)
}

func (s *FileSource) Close() error {
	s.shutdownFn()
	<-s.closed
	return nil
}

func (s *FileSource) loop() {
	defer close(s.closed)

	if s.reloadInterval == 0 {
		return
	}

	ticker := time.NewTicker(s.reloadInterval)
	for {
		<-ticker.C

		notify, err := s.sync()
		if err != nil {
			log.Error().Err(err).Msgf("unable to sync file: %s", s.filepath)
			continue
		}

		if !notify {
			continue
		}

		log.Info().Str("filepath", s.filepath).Msg("source file changed: notify")
		// TODO(buglloc): racy
		for _, w := range s.watchers {
			w()
		}
	}
}

func (s *FileSource) sync() (bool, error) {
	stat, err := os.Stat(s.filepath)
	if err != nil {
		return false, fmt.Errorf("stat source file: %w", err)
	}

	if stat.ModTime().Compare(s.lastModTime) <= 0 {
		return false, nil
	}

	data, err := os.ReadFile(s.filepath)
	if err != nil {
		return false, fmt.Errorf("read source file: %w", err)
	}

	curSum := md5.Sum(data)
	if s.lastSum == curSum {
		return false, nil
	}

	var domains []string
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			continue
		}

		domains = append(domains, line)
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("scan error: %w", err)
	}

	s.mu.Lock()
	s.lastModTime = stat.ModTime()
	s.lastSum = curSum
	s.domains = domains
	s.mu.Unlock()

	return true, nil
}
