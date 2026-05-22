package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/jadejashaktisinh/gh-git-action-cli/parser"
)

type workflowMsg []list.Item
type jobMsg []list.Item
type errorMsg error

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func fetchLocalWorkflows() tea.Cmd {
	return func() tea.Msg {
		files, err := parser.FindWorkflows()
		if err != nil {
			return errorMsg(err)
		}

		items := make([]list.Item, len(files))
		for i, f := range files {
			wf, err := parser.ParseWorkflow(f)
			title := f
			if err == nil && wf.Name != "" {
				title = wf.Name
			}
			items[i] = item{title: title, desc: f}
		}

		return workflowMsg(items)
	}
}

func fetchLocalJobs(path string) tea.Cmd {
	return func() tea.Msg {
		wf, err := parser.ParseWorkflow(path)
		if err != nil {
			return errorMsg(err)
		}

		items := make([]list.Item, 0, len(wf.Jobs))
		for id, job := range wf.Jobs {
			title := id
			if job.Name != "" {
				title = job.Name
			}
			items = append(items, item{title: title, desc: id})
		}

		return jobMsg(items)
	}
}
