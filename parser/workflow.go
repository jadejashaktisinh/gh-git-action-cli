package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Workflow struct {
	Name string            `yaml:"name"`
	On   interface{}       `yaml:"on"`
	Env  map[string]string `yaml:"env"`
	Jobs map[string]Job    `yaml:"jobs"`
	Path string
}

type Job struct {
	Name   string            `yaml:"name"`
	RunsOn string            `yaml:"runs-on"`
	Env    map[string]string `yaml:"env"`
	Steps  []Step            `yaml:"steps"`
}

type Step struct {
	Name string            `yaml:"name"`
	Uses string            `yaml:"uses"`
	Run  string            `yaml:"run"`
	Env  map[string]string `yaml:"env"`
}

func ParseWorkflow(path string) (*Workflow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read workflow file: %w", err)
	}

	var workflow Workflow
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow: %w", err)
	}
	workflow.Path = path
	return &workflow, nil
}

func FindWorkflows() ([]string, error) {
	return filepath.Glob("./.github/workflows/*.yml")
}
