package cli

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type interactiveInputModel struct {
	sharedData *interactiveSharedData
	inputModel textarea.Model
}

func newInteractiveInputModel(sharedData *interactiveSharedData) *interactiveInputModel {
	ti := textarea.New()
	ti.Placeholder = "Enter a query..."
	ti.SetValue(sharedData.selector)
	ti.Focus()
	ti.SetHeight(5)
	ti.ShowLineNumbers = false

	return &interactiveInputModel{
		sharedData: sharedData,
		inputModel: ti,
	}
}

func (m *interactiveInputModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m *interactiveInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.sharedData.selector = m.inputModel.Value()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.inputModel.SetWidth(msg.Width)
	}

	m.inputModel, cmd = m.inputModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *interactiveInputModel) View() string {
	return m.inputModel.View()
}
