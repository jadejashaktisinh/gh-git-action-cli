package runner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jadejashaktisinh/gh-git-action-cli/parser"
)

type RunOptions struct {
	Env map[string]string
}

func RunJob(job parser.Job, opts RunOptions) error {
	fmt.Printf("🚀 Running job: %s\n", job.Name)

	for _, step := range job.Steps {
		if step.Run != "" {
			err := RunStep(step, opts)
			if err != nil {
				return fmt.Errorf("step failed: %w", err)
			}
		} else if step.Uses != "" {
			fmt.Printf("⚠️ Skipping third-party action: %s\n", step.Uses)
		}
	}

	return nil
}

func RunStep(step parser.Step, opts RunOptions) error {
	name := step.Name
	if name == "" {
		name = strings.Split(step.Run, "\n")[0]
	}
	fmt.Printf("🔄 Step: %s\n", name)

	cmd := exec.Command("sh", "-c", step.Run)

	// Combine environments: Host < Global Env < Job Env < Step Env
	cmd.Env = os.Environ()
	for k, v := range opts.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	for k, v := range step.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func RunWorkflow(wf parser.Workflow, opts RunOptions) error {
	fmt.Printf("🚀 Running Workflow: %s\n", wf.Path)

	for job := range wf.Jobs {
		err := RunJob(wf.Jobs[job], opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Job %s failed: %v\n", wf.Jobs[job].Name, err)
			return err
		}
	}

	return nil
}
