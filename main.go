package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/coredns/coredns/core/dnsserver"
	_ "github.com/coredns/coredns/core/plugin"
	"github.com/coredns/coredns/coremain"
	"github.com/fsnotify/fsnotify"
	_ "github.com/tmeckel/coredns-finalizer"
	"gopkg.in/yaml.v3"
)

func init() {
	dnsserver.Directives = append(
		dnsserver.Directives,
		"finalize",
	)
}

func watchFile(path string, callback func()) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("create fs watcher: %w", err)
	}

	err = watcher.Add(filepath.Dir(path))
	if err != nil {
		_ = watcher.Close()
		return fmt.Errorf("add path to watcher: %w", err)
	}

	go func() {
		defer watcher.Close()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Name != path {
					continue
				}

				if event.Op&(fsnotify.Create) == 0 {
					continue
				}

				callback()
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("watch fail:", err)
			}
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
