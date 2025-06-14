package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/coredns/coredns/core/dnsserver"
	_ "github.com/coredns/coredns/core/plugin"
	"github.com/coredns/coredns/coremain"
	_ "github.com/tmeckel/coredns-finalizer"
	"gopkg.in/yaml.v3"
)

func genConfig() (string, error) {
	tmplPath := os.Getenv("SLY64_TEMPLATE_PATH")
	if tmplPath == "" {
		return "", nil
	}

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", fmt.Errorf("parse template from env[SLY64_TEMPLATE_PATH]=%s: %w", tmplPath, err)
	}

	dataPath := os.Getenv("SLY64_DATA_PATH")
	var data any
	if dataPath != "" {
		rawData, err := os.ReadFile(dataPath)
		if err != nil {
			return "", fmt.Errorf("read data from env[SLY64_DATA_PATH]=%s: %w", dataPath, err)
		}

		if err := yaml.Unmarshal(rawData, &data); err != nil {
			return "", fmt.Errorf("invalid data from env[SLY64_DATA_PATH]=%s: %w", dataPath, err)
		}
	}

	dir, err := os.MkdirTemp("", "sly64-*")
	if err != nil {
		return "", fmt.Errorf("create tempdir: %w", err)
	}

	f, err := os.Create(filepath.Join(dir, "Corefile"))
	if err != nil {
		return "", fmt.Errorf("create Corefile: %w", err)
	}
	defer func() { _ = f.Close() }()

	if err := tmpl.Execute(f, data); err != nil {
		return "", fmt.Errorf("render template: %w", err)
	}

	if err := f.Close(); err != nil {
		return "", fmt.Errorf("close Corefile: %w", err)
	}

	return dir, nil
}

func init() {
	dnsserver.Directives = append(
		dnsserver.Directives,
		"finalize",
	)
}

func main() {
	newCwd, err := genConfig()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to generate new config: %v\n", err)
		os.Exit(1)
	}

	if newCwd != "" {
		if err := os.Chdir(newCwd); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to change working directory: %v\n", err)
			os.Exit(1)
		}

		defer func() { _ = os.RemoveAll(newCwd) }()
		fmt.Println("change workdir to", newCwd)
	}

	coremain.Run()
}
