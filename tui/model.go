package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	stateRepoSelection state = iota
	stateWorkflowSelection
	stateInputForm
	statePolling
)

type Model struct {
	state      state
	list       list.Model
	repo       string
	workflow   string
	branch     string
	inputs     map[string]string
	err        error
}

func (m Model) Init() tea.Cmd {
	return fetchRepos()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.state == stateRepoSelection {
				if i, ok := m.list.SelectedItem().(item); ok {
					m.repo = i.title
					m.state = stateWorkflowSelection
					m.list.Title = "Select Workflow"
					return m, fetchWorkflows(m.repo)
				}
			} else if m.state == stateWorkflowSelection {
				if i, ok := m.list.SelectedItem().(item); ok {
					m.workflow = i.desc // path
					m.state = stateInputForm
					return m, fetchWorkflowContent(m.repo, m.workflow)
				}
			} else if m.state == stateInputForm {
				// For now, just trigger with default branch and no inputs
				m.branch = "main"
				return m, triggerWorkflow(m.repo, m.workflow, m.branch, nil)
			}
		case "backspace", "esc":
			if m.state == stateWorkflowSelection {
				m.state = stateRepoSelection
				m.list.Title = "Select Repository"
				return m, fetchRepos()
			} else if m.state == stateInputForm {
				m.state = stateWorkflowSelection
				m.list.Title = "Select Workflow"
				return m, fetchWorkflows(m.repo)
			}
		}

	case repoMsg:
		m.list.SetItems(msg)
		return m, nil

	case workflowMsg:
		m.list.SetItems(msg)
		return m, nil

	case workflowContentMsg:
		// TODO: Parse YAML and show inputs
		// For now, just show that we got it
		m.list.Title = "Trigger Workflow? (Enter to confirm)"
		return m, nil

	case statusMsg:
		m.list.Title = string(msg)
		return m, tea.Quit

	case errorMsg:
		m.err = msg
		return m, nil

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error()
	}
	return m.list.View()
}

func InitialModel() Model {
	// Initialize with empty list for now
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select Repository"
	
	return Model{
		state: stateRepoSelection,
		list:  l,
	}
}
