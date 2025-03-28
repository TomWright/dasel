package cli

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tomwright/dasel/v3/internal"
	"github.com/tomwright/dasel/v3/parsing"
)

var (
	interactiveKeyQuit       = tea.KeyCtrlC
	interactiveKeyCycleRead  = tea.KeyCtrlE
	interactiveKeyCycleWrite = tea.KeyCtrlD

	headingStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Padding(0, 1, 1, 1)
	}()
	shortcutStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Padding(0, 1).Align(lipgloss.Left)
	}()
	headerStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Padding(1).Align(lipgloss.Center)
	}()
	inputStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Margin(0, 0, 1, 0)
	}()
	inputContentStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.RoundedBorder())
	}()
	inputHeaderStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Padding(0, 2).Margin(0, 0, 1, 0).Underline(true)
	}()
	outputContentStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.RoundedBorder())
	}()
	outputHeaderStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Padding(0, 2).Margin(0, 0, 1, 0).Underline(true)
	}()
)

type interactiveDaselExecutor func(selector string, root bool, formatIn parsing.Format, formatOut parsing.Format, in string) (res string, err error)

func newInteractiveTeaProgram(initialInput string, initialSelector string, formatIn parsing.Format, formatOut parsing.Format, run interactiveDaselExecutor) (*tea.Program, func() string) {
	m := newInteractiveRootModel(initialInput, initialSelector, formatIn, formatOut, run)
	return tea.NewProgram(m, tea.WithAltScreen()), func() string {
		return m.sharedData.selector
	}
}

type interactiveSharedData struct {
	formatIn  parsing.Format
	formatOut parsing.Format
	selector  string
	input     string
}

type interactiveRootModel struct {
	sharedData   *interactiveSharedData
	inputModel   *interactiveInputModel
	outputModels []*interactiveOutputModel
}

func newInteractiveRootModel(initialInput string, initialSelector string, formatIn parsing.Format, formatOut parsing.Format, run interactiveDaselExecutor) *interactiveRootModel {
	res := &interactiveRootModel{
		sharedData: &interactiveSharedData{
			formatIn:  formatIn,
			formatOut: formatOut,
			selector:  initialSelector,
			input:     initialInput,
		},
		outputModels: make([]*interactiveOutputModel, 0),
	}

	res.inputModel = newInteractiveInputModel(res.sharedData)

	outputRootModel := newInteractiveOutputModel(res.sharedData, true, run)
	outputResultModel := newInteractiveOutputModel(res.sharedData, false, run)

	res.outputModels = append(res.outputModels, outputRootModel, outputResultModel)

	return res
}

func (m *interactiveRootModel) Init() tea.Cmd {
	return nil
}

func cycleFormats(all []parsing.Format, current parsing.Format) parsing.Format {
	slices.SortFunc(all, func(i, j parsing.Format) int {
		return strings.Compare(string(i), string(j))
	})
	cur := -1
	for i, format := range all {
		if format == current {
			cur = i
			break
		}
	}
	next := cur + 1
	if next > len(all)-1 {
		next = 0
	}
	return all[next]
}

func (m *interactiveRootModel) cycleReader() {
	m.sharedData.formatIn = cycleFormats(parsing.RegisteredReaders(), m.sharedData.formatIn)
}

func (m *interactiveRootModel) cycleWriter() {
	m.sharedData.formatOut = cycleFormats(parsing.RegisteredWriters(), m.sharedData.formatOut)
}

func (m *interactiveRootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case interactiveKeyQuit:
			return m, tea.Quit
		case interactiveKeyCycleRead:
			m.cycleReader()
		case interactiveKeyCycleWrite:
			m.cycleWriter()
		default:
		}

	case tea.WindowSizeMsg:
		headerStyle = headerStyle.Width(msg.Width).MaxWidth(msg.Width)

		var headerHeight int
		{
			headerHeight += lipgloss.Height(m.headerView())
			headerHeight += lipgloss.Height(m.inputView())
		}
		verticalMarginHeight := headerHeight

		numCols := len(m.outputModels)

		viewportHeight := msg.Height - verticalMarginHeight - (2 * numCols)
		viewportWidth := (msg.Width / numCols) - (2 * numCols)

		for _, outputModel := range m.outputModels {
			outputModel.setSize(viewportWidth, viewportHeight)
			outputModel.setVerticalPosition(verticalMarginHeight)
		}
	}

	{
		var model tea.Model
		model, cmd = m.inputModel.Update(msg)
		m.inputModel = model.(*interactiveInputModel)
		cmds = append(cmds, cmd)
	}

	for i, outputModel := range m.outputModels {
		var model tea.Model
		model, cmd = outputModel.Update(msg)
		m.outputModels[i] = model.(*interactiveOutputModel)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *interactiveRootModel) headerView() string {
	header := headingStyle.Render("Dasel interactive mode - " + internal.Version)

	shortcuts := "\n"
	shortcuts += fmt.Sprintf("%s: %s\n", interactiveKeyQuit, "Quit")
	shortcuts += fmt.Sprintf("%s: %s\n", interactiveKeyCycleRead, "Cycle reader")
	shortcuts += fmt.Sprintf("%s: %s\n", interactiveKeyCycleWrite, "Cycle writer")

	out := append([]string{header}, shortcutStyle.Render(shortcuts))

	out = append(out, fmt.Sprintf("\nReader: %s | Writer: %s", m.sharedData.formatIn, m.sharedData.formatOut))

	return headerStyle.Render(out...)
}

func (m *interactiveRootModel) inputView() string {
	return inputStyle.Render(m.inputModel.View())
}

func (m *interactiveRootModel) View() string {
	rows := make([]string, 0)

	rows = append(rows, m.headerView())

	rows = append(rows, m.inputView())

	{
		cols := make([]string, 0)
		for _, outputModel := range m.outputModels {
			cols = append(cols, outputModel.View())
		}
		if len(cols) > 0 {
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cols...))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
