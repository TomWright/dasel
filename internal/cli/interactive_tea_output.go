package cli

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tomwright/dasel/v3/parsing"
)

type interactiveOutputModel struct {
	sharedData          *interactiveSharedData
	hasUpdatedBefore    bool
	lastSeenSelector    string
	lastSeenFormatIn    parsing.Format
	lastSeenFormatOut   parsing.Format
	lastSeenInput       string
	root                bool
	run                 interactiveDaselExecutor
	output              string
	outputViewport      viewport.Model
	outputViewportReady bool
}

func newInteractiveOutputModel(sharedData *interactiveSharedData, root bool, run interactiveDaselExecutor) *interactiveOutputModel {
	m := &interactiveOutputModel{
		sharedData: sharedData,
		root:       root,
		run:        run,
	}
	m.outputViewport = viewport.New(10, 10)
	return m
}

func (m *interactiveOutputModel) Init() tea.Cmd {
	return nil
}

func (m *interactiveOutputModel) setOutput(output string) {
	m.output = output
	if m.outputViewportReady {
		m.outputViewport.SetContent(m.output)
	}
}

func (m *interactiveOutputModel) setSize(width int, height int) {
	if !m.outputViewportReady {
		m.outputViewportReady = true
	}

	m.outputViewport.Width = width
	m.outputViewport.Height = height
	m.outputViewport.SetContent(m.output)
}

func (m *interactiveOutputModel) setVerticalPosition(pos int) {
	m.outputViewport.YPosition = pos
}

func (m *interactiveOutputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	defer func() {
		m.lastSeenSelector = m.sharedData.selector
		m.lastSeenFormatIn = m.sharedData.formatIn
		m.lastSeenFormatOut = m.sharedData.formatOut
		m.lastSeenInput = m.sharedData.input
	}()
	firstUpdate := !m.hasUpdatedBefore
	m.hasUpdatedBefore = true

	queryChanged := m.lastSeenSelector != m.sharedData.selector ||
		m.lastSeenFormatIn != m.sharedData.formatIn ||
		m.lastSeenFormatOut != m.sharedData.formatOut ||
		m.lastSeenInput != m.sharedData.input

	// Take care of dasel execution + output.
	if firstUpdate || queryChanged {
		m.setOutput("Executing...")
		out, err := m.run(m.sharedData.selector, m.root, m.sharedData.formatIn, m.sharedData.formatOut, m.sharedData.input)
		if err != nil {
			m.setOutput(err.Error())
		} else {
			m.setOutput(out)
		}
	}

	m.outputViewport, cmd = m.outputViewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *interactiveOutputModel) View() string {
	title := "Result"
	if m.root {
		title = "Root"
	}

	content := "Initializing..."
	if m.outputViewportReady {
		content = m.outputViewport.View()
	}

	return lipgloss.JoinVertical(lipgloss.Left, outputHeaderStyle.Render(title), outputContentStyle.Render(content))
}
