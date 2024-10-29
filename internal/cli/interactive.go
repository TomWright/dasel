package cli

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tomwright/dasel/v3/internal"
)

const (
	useHighPerformanceRenderer = false
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.HiddenBorder()
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	resultTitleStyle = func() lipgloss.Style {
		b := lipgloss.HiddenBorder()
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	viewportStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	commandInputStyle = func() lipgloss.Style {
		return lipgloss.NewStyle()
	}()
)

func NewInteractiveCmd(queryCmd *QueryCmd) *InteractiveCmd {
	return &InteractiveCmd{
		Vars:          queryCmd.Vars,
		ExtReadFlags:  queryCmd.ExtReadFlags,
		ExtWriteFlags: queryCmd.ExtWriteFlags,
		InFormat:      queryCmd.InFormat,
		OutFormat:     queryCmd.OutFormat,

		Query: queryCmd.Query,
	}
}

type InteractiveCmd struct {
	Vars          *[]variable         `flag:"" name:"var" help:"Variables to pass to the query. E.g. --var foo=\"bar\" --var baz=json:file:./some/file.json"`
	ExtReadFlags  *[]extReadWriteFlag `flag:"" name:"read-flag" help:"Reader flag to customise parsing. E.g. --read-flag xml-mode=structured"`
	ExtWriteFlags *[]extReadWriteFlag `flag:"" name:"write-flag" help:"Writer flag to customise output"`
	InFormat      string              `flag:"" name:"in" short:"i" help:"The format of the input data."`
	OutFormat     string              `flag:"" name:"out" short:"o" help:"The format of the output data."`

	Query string `arg:"" help:"The query to execute." optional:"" default:""`
}

func (c *InteractiveCmd) Run(ctx *Globals) error {
	var err error
	var stdInBytes []byte = nil

	if ctx.Stdin != nil {
		stdInBytes, err = io.ReadAll(ctx.Stdin)
		if err != nil {
			return err
		}
	}

	m := initialModel(c, c.Query, stdInBytes)

	p := tea.NewProgram(m)

	_, err = p.Run()
	return err
}

func initialModel(it *InteractiveCmd, defaultCommand string, stdInBytes []byte) interactiveTeaModel {
	ti := textarea.New()
	ti.Placeholder = "Enter a query..."
	ti.SetValue(defaultCommand)
	ti.Focus()
	ti.SetHeight(5)
	ti.ShowLineNumbers = false

	return interactiveTeaModel{
		it:           it,
		commandInput: ti,
		err:          nil,
		stdInBytes:   stdInBytes,
		output:       "Loading...",
		firstUpdate:  true,
	}
}

type interactiveTeaModel struct {
	it                     *InteractiveCmd
	err                    error
	originalErr            error
	commandInput           textarea.Model
	currentCommand         string
	output                 string
	outputViewport         viewport.Model
	originalOutput         string
	originalOutputViewport viewport.Model
	outputReady            bool
	previousCommand        string
	stdInBytes             []byte
	firstUpdate            bool
}

func (m interactiveTeaModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m interactiveTeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.previousCommand = m.currentCommand

	m.currentCommand = m.commandInput.Value()
	if m.firstUpdate || m.currentCommand != m.previousCommand {
		m.firstUpdate = false
		// If the command has changed, we need to execute it

		{
			out, err := m.execDasel(m.currentCommand, true)
			if err != nil {
				m.originalErr = err
			} else {
				m.originalErr = nil
				m.originalOutput = out
				m.originalOutputViewport.SetContent(m.originalOutput)
			}
			if m.originalErr != nil {
				m.originalOutput = m.originalErr.Error()
				m.originalOutputViewport.SetContent(m.originalOutput)
			}
		}

		out, err := m.execDasel(m.currentCommand, false)
		if err != nil {
			m.err = err
		} else {
			m.err = nil
			m.output = out
			m.outputViewport.SetContent(m.output)
		}
		if m.err != nil {
			m.output = m.err.Error()
			m.outputViewport.SetContent(m.output)
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Quit
			//if m.commandInput.Focused() {
			//	m.commandInput.Blur()
			//}
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.commandInput.Focused() {
				cmd = m.commandInput.Focus()
				cmds = append(cmds, cmd)
			}
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		commandInputHeight := lipgloss.Height(m.commandInputView())
		resultHeaderHeight := lipgloss.Height(m.resultHeaderView())
		verticalMarginHeight := headerHeight + 2 + commandInputHeight + resultHeaderHeight

		viewportHeight := msg.Height - verticalMarginHeight
		viewportWidth := (msg.Width / 2) - 4

		if !m.outputReady {
			m.outputReady = true

			m.commandInput.SetWidth(msg.Width)

			m.outputViewport = viewport.New(viewportWidth, viewportHeight)
			//m.outputViewport.YPosition = verticalMarginHeight
			m.outputViewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.outputViewport.SetContent(m.output)
			m.commandInput.SetWidth(msg.Width)
			if useHighPerformanceRenderer {
				m.outputViewport.YPosition = headerHeight + 1
			}

			m.originalOutputViewport = viewport.New(viewportWidth, viewportHeight)
			m.originalOutputViewport.YPosition = verticalMarginHeight
			m.originalOutputViewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.originalOutputViewport.SetContent(m.output)

			if useHighPerformanceRenderer {
				m.originalOutputViewport.YPosition = headerHeight + 1
			}
		} else {
			m.commandInput.SetWidth(msg.Width)

			m.outputViewport.Width = viewportWidth
			m.outputViewport.Height = viewportHeight
			m.outputViewport.SetContent(m.output)

			m.originalOutputViewport.Width = viewportWidth
			m.originalOutputViewport.Height = viewportHeight
			m.outputViewport.SetContent(m.originalOutput)
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.outputViewport))
			cmds = append(cmds, viewport.Sync(m.originalOutputViewport))
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.commandInput, cmd = m.commandInput.Update(msg)
	cmds = append(cmds, cmd)

	m.outputViewport, cmd = m.outputViewport.Update(msg)
	cmds = append(cmds, cmd)

	m.originalOutputViewport, cmd = m.originalOutputViewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m interactiveTeaModel) execDasel(selector string, root bool) (res string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	var stdIn *bytes.Reader = nil
	if m.stdInBytes != nil {
		stdIn = bytes.NewReader(m.stdInBytes)
	} else {
		stdIn = bytes.NewReader([]byte{})
	}

	o := runOpts{
		Vars:          m.it.Vars,
		ExtReadFlags:  m.it.ExtReadFlags,
		ExtWriteFlags: m.it.ExtWriteFlags,
		InFormat:      m.it.InFormat,
		OutFormat:     m.it.OutFormat,
		ReturnRoot:    root,
		Unstable:      true,
		Query:         selector,

		Stdin: stdIn,
	}

	outBytes, err := run(o)
	return string(outBytes), err
}

func (m interactiveTeaModel) headerView() string {
	return titleStyle.Render("Dasel Interactive Mode - " + internal.Version + " - ctrl+c or esc to exit")
}

func (m interactiveTeaModel) commandInputView() string {
	return commandInputStyle.Render(m.commandInput.View())
}

func (m interactiveTeaModel) originalHeaderView() string {
	return resultTitleStyle.Render("Root")
}

func (m interactiveTeaModel) resultHeaderView() string {
	return resultTitleStyle.Render("Result")
}

func (m interactiveTeaModel) viewportView() string {
	return viewportStyle.Render(m.outputViewport.View())
}

func (m interactiveTeaModel) originalViewportView() string {
	return viewportStyle.Render(m.originalOutputViewport.View())
}

func (m interactiveTeaModel) View() string {
	res := []string{
		m.headerView(),
		m.commandInputView(),
	}

	var left, right []string

	left = append(left, m.originalHeaderView(), m.originalViewportView())
	right = append(right, m.resultHeaderView(), m.viewportView())

	viewports := lipgloss.JoinHorizontal(lipgloss.Bottom, strings.Join(left, "\n"), strings.Join(right, "\n"))
	res = append(res, viewports)

	return strings.Join(res, "\n")
}
