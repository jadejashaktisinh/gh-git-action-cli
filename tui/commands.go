package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/cli/go-gh/v2/pkg/api"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
)

type repoMsg []list.Item
type workflowMsg []list.Item
type statusMsg string
type workflowContentMsg string
type errorMsg error

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func fetchRepos() tea.Cmd {
	return func() tea.Msg {
		client, err := api.DefaultRESTClient()
		if err != nil {
			return errorMsg(err)
		}

		var repos []struct {
			FullName string `json:"full_name"`
			Description string `json:"description"`
		}

		err = client.Get("user/repos?sort=updated&per_page=100", &repos)
		if err != nil {
			return errorMsg(err)
		}

		items := make([]list.Item, len(repos))
		for i, r := range repos {
			items[i] = item{title: r.FullName, desc: r.Description}
		}

		return repoMsg(items)
	}
}

func fetchWorkflows(repo string) tea.Cmd {
	return func() tea.Msg {
		client, err := api.DefaultRESTClient()
		if err != nil {
			return errorMsg(err)
		}

		var resp struct {
			Workflows []struct {
				Name  string `json:"name"`
				Path  string `json:"path"`
				State string `json:"state"`
			} `json:"workflows"`
		}

		path := fmt.Sprintf("repos/%s/actions/workflows", repo)
		err = client.Get(path, &resp)
		if err != nil {
			return errorMsg(err)
		}

		var items []list.Item
		for _, w := range resp.Workflows {
			if w.State == "active" {
				items = append(items, item{title: w.Name, desc: w.Path})
			}
		}

		return workflowMsg(items)
	}
}

func triggerWorkflow(repo, workflow, branch string, inputs map[string]string) tea.Cmd {
	return func() tea.Msg {
		client, err := api.DefaultRESTClient()
		if err != nil {
			return errorMsg(err)
		}

		body := map[string]interface{}{
			"ref": branch,
		}
		if len(inputs) > 0 {
			body["inputs"] = inputs
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			return errorMsg(err)
		}

		workflowFile := filepath.Base(workflow)

		path := fmt.Sprintf("repos/%s/actions/workflows/%s/dispatches", repo, workflowFile)
		err = client.Post(path, bytes.NewReader(jsonBody), nil)
		if err != nil {
			return errorMsg(err)
		}

		return statusMsg("Workflow triggered successfully!")
	}
}

func fetchWorkflowContent(repo, path string) tea.Cmd {
	return func() tea.Msg {
		client, err := api.DefaultRESTClient()
		if err != nil {
			return errorMsg(err)
		}

		var resp struct {
			Content string `json:"content"`
		}

		apiPath := fmt.Sprintf("repos/%s/contents/%s", repo, path)
		err = client.Get(apiPath, &resp)
		if err != nil {
			return errorMsg(err)
		}

		return workflowContentMsg(resp.Content)
	}
}
