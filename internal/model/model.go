package model

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fedeztk/got/pkg/translator"
)

const (
	// application states
	TYPING      = iota // input tab
	CHOOSING           // language list tab
	TRANSLATING        // translation tab
	LOADING            // loading inside input tab
	// pager
	headerHeight = 6 // 3 + 3 tabs and gaps
	footerHeight = 3
	// list
	defaultListHeight = 34
	defaultListWidth  = 14
)

var (
	states = []int{TYPING, CHOOSING, TRANSLATING, LOADING}
	// prompt
	promptStyleIndicator = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	placeholderStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	promptStyleUpperText = lipgloss.NewStyle().Background(lipgloss.Color("6")).Bold(true).MarginLeft(2).Padding(0, 1).Foreground(lipgloss.Color("0"))
	promptStyleSelLang   = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).MarginLeft(2)
	// spinner
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	// list
	titleStyle         = lipgloss.NewStyle().Background(lipgloss.Color("11")).Bold(true).Padding(0, 1).Foreground(lipgloss.Color("0"))
	paginationStyle    = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle          = list.DefaultStyles().HelpStyle.PaddingLeft(4)
	statusMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).MarginLeft(2).Bold(true)
	// errors
	ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	// tabs
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

	help      help.Model
	genKeyMap genKeyMap

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
	sourceLangKey, targetLangKey, invertLangKey key.Binding
}

type Config interface {
	Source() string
	Target() string
	Engine() string
	RememberLastSettings(source, target string)
}

var conf Config

func newModel() *model {
	t := textinput.NewModel()
	t.Placeholder = "your text here"
	t.PlaceholderStyle = placeholderStyle
	t.PromptStyle = promptStyleIndicator
	t.Focus()

	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	keys := getKeyMapLangList()
	l := list.New(getConfLangs(), list.NewDefaultDelegate(), defaultListWidth, defaultListHeight)
	l.Title = "Available languages"
	l.AdditionalFullHelpKeys = func() []key.Binding { return []key.Binding{keys.sourceLangKey, keys.targetLangKey, keys.invertLangKey} }
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
		help:         help.New(),
		genKeyMap:    getGenKeyMap(),
	}
}

func getConfLangs() []list.Item {
	items := make([]list.Item, 0)

	for abbrev, title := range translator.GetAllLanguages() {
		items = append(items, item{title, abbrev})
	}
	return items
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

		// global keys
		switch {
		case key.Matches(msg, m.genKeyMap.exitKey):
			conf.RememberLastSettings(m.source, m.target)
			return m, tea.Quit

		case key.Matches(msg, m.genKeyMap.nextTab):
			m.switchTab(+1)

		case key.Matches(msg, m.genKeyMap.prevTab):
			m.switchTab(-1)
		}

		// keys for langauge list view
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
			case key.Matches(msg, m.langListKeys.invertLangKey):
				m.source, m.target = m.target, m.source
				statusCmd := m.langList.NewStatusMessage(statusMessageStyle.Render("Inverted languages: " + m.source + " → " + m.target))
				cmds = append(cmds, statusCmd)
			}
		}

		// keys for text input view (does not have dedicated keys struct, hence the switch on .String())
		switch msg.String() {
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
		m.langList.SetHeight(msg.Height - verticalMargins)

		m.help.Width = msg.Width

		// translation fetched
	case gotTrans:
		m.state = TRANSLATING
		m.err = msg.Err
		m.result = msg.result
		m.viewport.SetContent(m.result)
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
	// wait for terminal size info
	if !m.termInfoReady {
		return "Initializing..."
	}

	// Tabs:
	// tabsRow renders the tabs header
	// content renders the content based on the current tab
	var tabsRow, content string

	switch m.state {
	case TYPING:
		tabsRow = lipgloss.JoinHorizontal(
			lipgloss.Top,
			activeTab.Render("Text input"),
			tab.Render("Language selection"),
			tab.Render("Translation"),
		)
		content = promptStyleUpperText.Render("Enter sentence") + "\n\n" + m.textInput.View()
	case LOADING:
		tabsRow = lipgloss.JoinHorizontal(
			lipgloss.Top,
			activeTab.Render("Text input"),
			tab.Render("Language selection"),
			tab.Render("Translation"),
		)
		content = fmt.Sprintf("%s fetching results... please wait.", m.spinner.View())
	case TRANSLATING:
		tabsRow = lipgloss.JoinHorizontal(
			lipgloss.Top,
			tab.Render("Text input"),
			tab.Render("Language selection"),
			activeTab.Render("Translation"),
		)
		if m.err != nil {
			content = ErrorStyle.Render(m.err.Error())
		} else {
			content = m.viewport.View()
		}
	case CHOOSING:
		tabsRow = lipgloss.JoinHorizontal(
			lipgloss.Top,
			tab.Render("Text input"),
			activeTab.Render("Language selection"),
			tab.Render("Translation"),
		)
		content = m.langList.View()
	}

	// holds top right translation info
	translationStatus := promptStyleSelLang.Render(fmt.Sprintf("%s → %s (%s engine)", m.source, m.target, conf.Engine()))

	lenTabs := lipgloss.Width(translationStatus) + lipgloss.Width(tabsRow) + 2 // still don't know why 2 cells are missing

	gap := tabGap.Render(strings.Repeat(" ", m.viewport.Width-lenTabs) + translationStatus)
	tabsRow = lipgloss.JoinHorizontal(lipgloss.Bottom, tabsRow, gap)

	view := tabsRow + "\n\n\n" + content

	return view + lipgloss.PlaceVertical(
		(m.viewport.Height+headerHeight+footerHeight)- // total height of terminal as originally received
			lipgloss.Height(view), // height of already utilized space
		lipgloss.Bottom,
		m.getFooter())
}

func (m model) fetchTranslation(query string) tea.Cmd {
	return func() tea.Msg {
		response, err := translator.Translate(query, m.source, m.target, conf.Engine())
		if err != nil {
			return gotTrans{Err: err, result: err.Error()}
		}
		return gotTrans{result: response.PrettyPrint()}
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

func (m model) getFooter() string {
	helpMenu := m.help.View(m.genKeyMap)
	helpLen := lipgloss.Width(helpMenu)
	footerTop, footerMid, footerBot := "", "", ""
	gapSize := m.viewport.Width

	// leave some space for the percentage of the viewport
	if m.state == TRANSLATING {
		percentStr := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)
		footerTop = "╭──────╮"
		footerMid = "┤ " + footerTextStyle.Render(percentStr) + footerStyle.Render(" │")
		footerBot = "╰──────╯"
		gapSize -= lipgloss.Width(percentStr) + 4
	}

	footerTop = strings.Repeat(" ", gapSize) + footerTop
	footerMid = strings.Repeat("─", gapSize) + footerMid
	footerBot = helpMenu + strings.Repeat(" ", gapSize-helpLen) + footerStyle.Render(footerBot)
	footer := fmt.Sprintf("%s\n%s\n%s", footerTop, footerMid, footerBot)

	return footerStyle.Render(footer)
}

type item struct {
	title, abbreviation string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.abbreviation }
func (i item) FilterValue() string { return i.title }

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
		invertLangKey: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "invert languages"),
		),
	}
}

type genKeyMap struct {
	nextTab, prevTab, exitKey key.Binding
}

func (k genKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.nextTab, k.prevTab, k.exitKey}
}

func (k genKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.nextTab, k.prevTab, k.exitKey},
	}
}

func getGenKeyMap() genKeyMap {
	return genKeyMap{
		nextTab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab"),
		),
		prevTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift-tab", "previous tab"),
		),
		exitKey: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc/ctrl+c", "exit"),
		),
	}
}
