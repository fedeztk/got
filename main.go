package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

const (
	headerHeight = 3
	footerHeight = 3
)

type model struct {
	textInput textinput.Model
	spinner   spinner.Model
	viewport  viewport.Model
	result    string

	ready   bool
	typing  bool
	loading bool
	err     error
}

type gotTrans struct {
	Err    error
	result string
}

func main() {
	t := textinput.NewModel()
	t.Placeholder = "your text here"
	t.Focus()
	t.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))

	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	initialModel := model{
		textInput: t,
		spinner:   s,
		typing:    true,
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if err := p.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

			// not typing, not loading, we are in the pager and need to go back to typing
		case "esc", "q":
			if !m.typing && !m.loading {
				m.typing = true
				m.err = nil
				return m, nil
			}

		case "enter":
			if m.typing {
				query := strings.TrimSpace(m.textInput.Value())
				if query != "" {
					m.typing = false
					m.loading = true
					cmds = append(cmds, spinner.Tick)
					// // TODO: remove mock data
					cmds = append(cmds, m.fetchTranslation(query, "it", "en"))
				}
			}
		}

		// called on terminal resize
	case tea.WindowSizeMsg:
		verticalMargins := headerHeight + footerHeight

		// first time receiving terminal size, we don't have a viewport yet
		if !m.ready {
			m.viewport = viewport.Model{Width: msg.Width, Height: msg.Height - verticalMargins}
			m.ready = true
			m.viewport.YPosition = headerHeight + 1
		} else { // resize according to the new terminal size
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMargins
		}

		// translation fetched
	case gotTrans:
		m.loading = false

		if err := msg.Err; err == nil {
			m.err = err
			m.result = msg.result
			m.viewport.SetContent(m.result)
		}
	}

	// we are typing the sentence
	if m.typing {
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	// we are loading the translation
	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	// default, pager msg:
	if !m.typing && !m.loading {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.typing {
		return fmt.Sprintf("Enter sentence (exit with ctrl-c):\n%s", m.textInput.View())
	}

	if m.loading {
		return fmt.Sprintf("%s fetching results... please wait.", m.spinner.View())
	}

	if m.err != nil {
		return fmt.Sprintf("Could not fetch sentence: %v", m.err)
	}

	// wait for terminal size info
	if !m.ready {
		return "Initializing..."
	}

	return m.formatTranslation()
}

func (m model) fetchTranslation(query, source, target string) tea.Cmd {
	return func() tea.Msg {
		text, err := exec.Command("trans", "-t", target, "-s", source, query).Output()
		if err != nil {
			return gotTrans{Err: err}
		}
		return gotTrans{result: string(text)}
	}
}

func (m model) formatTranslation() string {
	headerTop := "╭─────────────╮"
	headerMid := "│ Translation ├"
	headerBot := "╰─────────────╯"
	headerMid += strings.Repeat("─", m.viewport.Width-runewidth.StringWidth(headerMid))
	header := fmt.Sprintf("%s\n%s\n%s", headerTop, headerMid, headerBot)

	footerTop := "╭──────╮"
	footerMid := fmt.Sprintf("┤ %3.f%% │", m.viewport.ScrollPercent()*100)
	footerBot := "╰──────╯"
	gapSize := m.viewport.Width - runewidth.StringWidth(footerMid)
	footerTop = strings.Repeat(" ", gapSize) + footerTop
	footerMid = strings.Repeat("─", gapSize) + footerMid
	footerBot = strings.Repeat(" ", gapSize) + footerBot
	footer := fmt.Sprintf("%s\n%s\n%s", footerTop, footerMid, footerBot)

	return fmt.Sprintf("%s\n%s\n%s", header, m.viewport.View(), footer)
}
