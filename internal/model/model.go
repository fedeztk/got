package model

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fedeztk/got/pkg/translator"
	"github.com/mattn/go-runewidth"
)

const (
	// application states
	TYPING      = iota // input tab
	CHOOSING           // language list tab
	TRANSLATING        // translation tab
	LOADING            // loading inside input tab
	// pager
	headerHeight = 6 // 3 + 3 (padding of tabs)
	footerHeight = 3
	// list
	listHeight = 34
	listWidth  = 14
)

var (
	states = []int{TYPING, CHOOSING, TRANSLATING, LOADING}
	// prompt
	promptStyleIndicator = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	promptStyleUpperText = lipgloss.NewStyle().Background(lipgloss.Color("6")).Bold(true).MarginLeft(2).Padding(0, 1).Foreground(lipgloss.Color("0"))
	promptStyleSelLang   = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).MarginLeft(2)
	// spinner
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	// list
	titleStyle         = lipgloss.NewStyle().Background(lipgloss.Color("11")).Bold(true).Padding(0, 1).Foreground(lipgloss.Color("0"))
	paginationStyle    = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle          = list.DefaultStyles().HelpStyle.PaddingLeft(4)
	statusMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).MarginLeft(2).Bold(true)
	// tabs
	docStyle        = lipgloss.NewStyle().Align(lipgloss.Center)
	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}
	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}
	highlight = lipgloss.Color("4")
	tab       = lipgloss.NewStyle().
			Border(tabBorder, true).
			BorderForeground(highlight).
			Padding(0, 2)
	activeTab = tab.Copy().Border(activeTabBorder, true).Bold(true)
	tabGap    = tab.Copy().
			BorderTop(false).
			BorderLeft(false).
			BorderRight(false)
	// pager
	footerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	footerTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
)

type model struct {
	textInput    textinput.Model
	spinner      spinner.Model
	viewport     viewport.Model
	langList     list.Model
	langListKeys keyMapList

	result string
	source string
	target string

	termInfoReady bool
	state         int
	err           error
}

type gotTrans struct {
	Err    error
	result string
}

type keyMapList struct {
	sourceLangKey, targetLangKey key.Binding
}

type Config interface {
	Langs() []list.Item
	Source() string
	Target() string
	RememberLastLangs(source, target string)
}

var conf Config

func newModel() *model {
	t := textinput.NewModel()
	t.Placeholder = "your text here"
	t.PromptStyle = promptStyleIndicator
	t.Focus()

	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	keys := getKeyMapLangList()
	confLangs := conf.Langs()
	l := list.NewModel(confLangs, list.NewDefaultDelegate(), listWidth, listHeight)
	l.Title = "Available languages"
	l.AdditionalFullHelpKeys = func() []key.Binding { return []key.Binding{keys.sourceLangKey, keys.targetLangKey} }
	l.Help.ShowAll = true
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return &model{
		langList:     l,
		textInput:    t,
		spinner:      s,
		langListKeys: keys,
		state:        TYPING,
		source:       conf.Source(),
		target:       conf.Target(),
	}
}

func Run(c Config) {
	conf = c
	initialModel := newModel()

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
		if m.langList.FilterState() == list.Filtering {
			break
		}

		if m.state == CHOOSING {
			switch {
			case key.Matches(msg, m.langListKeys.sourceLangKey):
				abbreviation, title := m.langList.SelectedItem().(item).abbreviation, m.langList.SelectedItem().(item).title
				m.source = abbreviation
				statusCmd := m.langList.NewStatusMessage(statusMessageStyle.Render("Source language: " + title))
				cmds = append(cmds, statusCmd)
			case key.Matches(msg, m.langListKeys.targetLangKey):
				abbreviation, title := m.langList.SelectedItem().(item).abbreviation, m.langList.SelectedItem().(item).title
				m.target = abbreviation
				statusCmd := m.langList.NewStatusMessage(statusMessageStyle.Render("Target language: " + title))
				cmds = append(cmds, statusCmd)
			}
		}

		switch msg.String() {
		case "ctrl+c", "esc":
			conf.RememberLastLangs(m.source, m.target)
			return m, tea.Quit

		case "tab":
			m.switchTab(+1)

		case "shift+tab":
			m.switchTab(-1)

		case "enter":
			if m.state == TYPING {
				query := strings.TrimSpace(m.textInput.Value())
				if query != "" {
					m.state = LOADING
					cmds = append(cmds, spinner.Tick)
					cmds = append(cmds, m.fetchTranslation(query))
				}
			}
		}

		// called on terminal resize
	case tea.WindowSizeMsg:
		verticalMargins := headerHeight + footerHeight

		// update pager
		if !m.termInfoReady { // first time receiving terminal size, we don't have a viewport yet
			m.viewport = viewport.Model{Width: msg.Width, Height: msg.Height - verticalMargins}
			m.termInfoReady = true
		} else { // resize according to the new terminal size
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMargins
		}

		// update language list
		m.langList.SetWidth(msg.Width)

		// translation fetched
	case gotTrans:
		m.state = TRANSLATING
		if err := msg.Err; err == nil {
			m.err = err
			m.result = msg.result
			m.viewport.SetContent(m.result)
		}
	}

	switch m.state {
	case TYPING:
		m.textInput, cmd = m.textInput.Update(msg)
	case LOADING:
		m.spinner, cmd = m.spinner.Update(msg)
	case TRANSLATING:
		m.viewport, cmd = m.viewport.Update(msg)
	case CHOOSING:
		m.langList, cmd = m.langList.Update(msg)
	}
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Could not fetch sentence: %v", m.err)
	}

	// wait for terminal size info
	if !m.termInfoReady {
		return "Initializing..."
	}

	// Tabs
	doc := strings.Builder{}
	var row, content string

	switch m.state {
	case TYPING:
		row = lipgloss.JoinHorizontal(
			lipgloss.Top,
			activeTab.Render("Text input"),
			tab.Render("Language selection"),
			tab.Render("Translation"),
		)
		content = promptStyleUpperText.Render("Enter sentence") +
			promptStyleSelLang.Render(fmt.Sprintf("Translating %s →  %s", m.source, m.target)) +
			fmt.Sprintf("\n\n%s\n\n(exit with ctrl-c)", m.textInput.View())
	case LOADING:
		row = lipgloss.JoinHorizontal(
			lipgloss.Top,
			activeTab.Render("Text input"),
			tab.Render("Language selection"),
			tab.Render("Translation"),
		)
		content = fmt.Sprintf("%s fetching results... please wait.", m.spinner.View())
	case TRANSLATING:
		row = lipgloss.JoinHorizontal(
			lipgloss.Top,
			tab.Render("Text input"),
			tab.Render("Language selection"),
			activeTab.Render("Translation"),
		)
		content = m.formatTranslation()
	case CHOOSING:
		row = lipgloss.JoinHorizontal(
			lipgloss.Top,
			tab.Render("Text input"),
			activeTab.Render("Language selection"),
			tab.Render("Translation"),
		)
		content = m.langList.View()
	}

	// activeLanguages := promptStyleSelLang.Render(fmt.Sprintf("%s →  %s", m.source, m.target))
	gap := tabGap.Render(strings.Repeat(" ", m.viewport.Width))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
	doc.WriteString(row + "\n\n")
	return docStyle.Render(doc.String()) + "\n" + content
}

func (m model) fetchTranslation(query string) tea.Cmd {
	return func() tea.Msg {
		response, err := translator.Translate(query, m.source, m.target)
		if err != nil {
			return gotTrans{Err: err}
		}
		return gotTrans{result: response.PrettyPrint() + "\n" + promptStyleSelLang.Render(fmt.Sprintf("%s →  %s", m.source, m.target))}
	}
}

func (m *model) switchTab(direction int) {
	newState := (m.state + direction) % len(states)
	if newState == LOADING {
		newState = states[0]
	}
	if newState < 0 {
		newState = TRANSLATING
	}
	if newState == TYPING {
		m.textInput.Focus()
	} else {
		m.textInput.Blur()
	}
	m.state = newState
}

func (m model) formatTranslation() string {
	footerTop := "╭──────╮"
	percentStr := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)
	footerMid := "┤ " + footerTextStyle.Render(percentStr) + footerStyle.Render(" │")
	footerBot := "╰──────╯"
	gapSize := m.viewport.Width - (runewidth.StringWidth(percentStr) + 4)
	footerTop = strings.Repeat(" ", gapSize) + footerTop
	footerMid = strings.Repeat("─", gapSize) + footerMid
	footerBot = strings.Repeat(" ", gapSize) + footerBot
	footer := fmt.Sprintf("%s\n%s\n%s", footerTop, footerMid, footerBot)

	return fmt.Sprintf("%s\n%s", m.viewport.View(), footerStyle.Render(footer))
}

type item struct {
	title, abbreviation string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.abbreviation }
func (i item) FilterValue() string { return i.title }

func NewListItem(title, abbreviation string) item {
	return item{title, abbreviation}
}

func getKeyMapLangList() keyMapList {
	return keyMapList{
		sourceLangKey: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "choose source language"),
		),
		targetLangKey: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "choose target language"),
		),
	}
}
