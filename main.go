package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/coredns/coredns/core/dnsserver"
	_ "github.com/coredns/coredns/core/plugin"
	"github.com/coredns/coredns/coremain"
	_ "github.com/tmeckel/coredns-finalizer"
	"gopkg.in/yaml.v3"
)

func init() {
	dnsserver.Directives = append(
		dnsserver.Directives,
		"finalize",
	)
}

func md5sum(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return hash.Sum(nil), nil
}

func watchFile(path string, callback func()) error {
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat file: %w", err)
	}

	lastModTime := stat.ModTime()
	lastMd5, err := md5sum(path)
	if err != nil {
		return fmt.Errorf("md5sum file: %w", err)
	}

	ticker := time.NewTicker(30 * time.Second)

	go func() {
		for {
			<-ticker.C
			stat, err := os.Stat(path)
			if err != nil {
				continue
			}

			if stat.ModTime().Compare(lastModTime) <= 0 {
				continue
			}

			newMd5, err := md5sum(path)
			if err != nil {
				continue
			}

			if bytes.Equal(lastMd5, newMd5) {
				continue
			}

			lastMd5 = newMd5
			lastModTime = stat.ModTime()
			callback()
		}
	}()

	return nil
}

func genConfigTo(tmplPath, dataPath, targetPath string) error {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("parse template from env[SLY64_TEMPLATE_PATH]=%s: %w", tmplPath, err)
	}

	var data any
	if dataPath != "" {
		rawData, err := os.ReadFile(dataPath)
		if err != nil {
			return fmt.Errorf("read data from env[SLY64_DATA_PATH]=%s: %w", dataPath, err)
		}

		if err := yaml.Unmarshal(rawData, &data); err != nil {
			return fmt.Errorf("invalid data from env[SLY64_DATA_PATH]=%s: %w", dataPath, err)
		}
	}

	f, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("create target file: %w", err)
	}
	defer func() { _ = f.Close() }()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("close target file: %w", err)
	}

	return nil
}

func prepareConfig() error {
	tmplPath := os.Getenv("SLY64_TEMPLATE_PATH")
	dataPath := os.Getenv("SLY64_DATA_PATH")
	if tmplPath == "" {
		return nil
	}

	workDir, err := os.MkdirTemp("", "sly64-*")
	if err != nil {
		return fmt.Errorf("create tempdir: %w", err)
	}

	if err := os.Chdir(workDir); err != nil {
		return fmt.Errorf("change working directory: %w", err)
	}

	targetPath := filepath.Join(workDir, "Corefile")
	if err := genConfigTo(tmplPath, dataPath, targetPath); err != nil {
		return fmt.Errorf("generate corefile: %w", err)
	}

	if dataPath != "" {
		err := watchFile(dataPath, func() {
			log.Println("data changes: generate new config")
			err := genConfigTo(tmplPath, dataPath, targetPath)
			if err != nil {
				log.Println("unable to generate corefile:", err)
				return
			}

			log.Println("data changes: send SIGUSR1")
			err = syscall.Kill(os.Getpid(), syscall.SIGUSR1)
			if err != nil {
				log.Println("failed to send SIGUSR1:", err)
			}
		})
		if err != nil {
			return fmt.Errorf("watch data file changes: %w", err)
		}
	}

	return nil
}

func main() {
	if err := prepareConfig(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to generate new config: %v\n", err)
		os.Exit(1)
	}

	coremain.Run()
}
