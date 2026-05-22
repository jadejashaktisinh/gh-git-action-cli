package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	stateWorkflowSelection state = iota
	stateJobSelection
	stateRunning
)

type Model struct {
	state      state
	list       list.Model
	Workflow   string
	Job        string
	err        error
}

func (m Model) Init() tea.Cmd {
	return fetchLocalWorkflows()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.state == stateWorkflowSelection {
				if i, ok := m.list.SelectedItem().(item); ok {
					m.Workflow = i.desc // path
					m.state = stateJobSelection
					m.list.Title = "Select Job"
					return m, fetchLocalJobs(m.Workflow)
				}
			} else if m.state == stateJobSelection {
				if i, ok := m.list.SelectedItem().(item); ok {
					m.Job = i.title
					m.state = stateRunning
					return m, tea.Quit
				}
			}
		case "backspace", "esc":
			if m.state == stateJobSelection {
				m.state = stateWorkflowSelection
				m.list.Title = "Select Workflow"
				return m, fetchLocalWorkflows()
			}
		}

	case workflowMsg:
		m.list.SetItems(msg)
		return m, nil

	case jobMsg:
		m.list.SetItems(msg)
		return m, nil

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
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select Local Workflow"
	
	return Model{
		state: stateWorkflowSelection,
		list:  l,
	}
}
