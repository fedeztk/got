package model

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
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
)

var once sync.Once

type model struct {
	textInput textinput.Model
	spinner   spinner.Model
	viewport  viewport.Model
	langList  list.Model
	help      help.Model

	keyMgr keyBindingMgr

	result      string
	shortResult string
	source      string
	target      string

	termInfoReady bool
	state         int
	err           error
	conf          Config
}

type gotTrans struct {
	Err         error
	result      string
	shortResult string
}

type gotTTS struct {
	Err    error
	result []byte
}

type Config interface {
	Source() string
	Target() string
	Engine() string
	RememberLastSettings(source, target string)
}

func newModel(c Config) *model {
	t := textinput.NewModel()
	t.Placeholder = "your text here"
	t.PlaceholderStyle = placeholderStyle
	t.PromptStyle = promptStyleIndicator
	t.Focus()

	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	l := list.New(getConfLangs(), list.NewDefaultDelegate(), 0, 0)
	l.DisableQuitKeybindings()
	l.SetShowHelp(false)
	l.Title = "Available languages"
	l.AdditionalFullHelpKeys = getListAdditionalKeyMap
	l.Styles.Title = titleStyle

	return &model{
		langList:  l,
		textInput: t,
		spinner:   s,
		state:     TYPING,
		source:    c.Source(),
		target:    c.Target(),
		help:      help.New(),
		conf:      c,
		keyMgr:    newKeyBindingMgr(l.FullHelp()),
	}
}

func Run(c Config) {
	initialModel := newModel(c)

	p := tea.NewProgram(initialModel, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if err := p.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.langList.FilterState() == list.Filtering {
			break
		}

		// global keybindings
		switch msg.String() {
		case "tab":
			m.switchTab(+1)

		case "shift+tab":
			m.switchTab(-1)

		case "ctrl+c", "esc":
			m.conf.RememberLastSettings(m.source, m.target)
			return m, tea.Quit
		}

		// language list keybindings
		if m.state == CHOOSING {
			switch msg.String() {
			case "s":
				abbreviation, title := m.langList.SelectedItem().(item).abbreviation, m.langList.SelectedItem().(item).title
				m.source = abbreviation
				statusCmd := m.langList.NewStatusMessage(statusMessageStyle.Render("Source language: " + title))
				cmds = append(cmds, statusCmd)

			case "t":
				abbreviation, title := m.langList.SelectedItem().(item).abbreviation, m.langList.SelectedItem().(item).title
				m.target = abbreviation
				statusCmd := m.langList.NewStatusMessage(statusMessageStyle.Render("Target language: " + title))
				cmds = append(cmds, statusCmd)

			case "i":
				m.source, m.target = m.target, m.source
				statusCmd := m.langList.NewStatusMessage(statusMessageStyle.Render("Inverted languages: " + m.source + " → " + m.target))
				cmds = append(cmds, statusCmd)

			case "?":
				m.help.ShowAll = !m.help.ShowAll
			}
		}

		// text input keybindings
		if m.state == TYPING {
			switch msg.String() {
			case "enter":
				query := strings.TrimSpace(m.textInput.Value())
				if query != "" {
					m.setState(LOADING)
					cmds = append(cmds, spinner.Tick)
					cmds = append(cmds, m.fetchTranslation(query))
				}
			}
		}

		// translating keybindings
		if m.state == TRANSLATING && m.shortResult != "" {
			switch msg.String() {
			case "y":
				m.yankTranslated()

			case "p":
				m.setState(LOADING)
				cmds = append(cmds, spinner.Tick)
				cmds = append(cmds, m.fetchTextToSpeech(m.shortResult))
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
		m.setState(TRANSLATING)
		m.err = msg.Err
		m.result = msg.result
		m.shortResult = msg.shortResult
		m.viewport.SetContent(m.result)

	// text to speech fetched
	case gotTTS:
		m.setState(TRANSLATING)
		m.err = msg.Err
		go m.playTTS(msg.result)
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

	tabsRow = lipgloss.JoinHorizontal(lipgloss.Top, m.renderTabs()...)

	switch m.state {
	case TYPING:
		content = promptStyleUpperText.Render("Enter sentence") + "\n\n" + m.textInput.View()
	case LOADING:
		content = fmt.Sprintf("%s fetching results... please wait.", m.spinner.View())
	case TRANSLATING:
		if m.err != nil {
			content = ErrorStyle.Render(m.err.Error())
		} else {
			content = m.viewport.View()
		}
	case CHOOSING:
		content = m.langList.View()
	}

	// holds top right translation info
	translationStatus := promptStyleSelLang.Render(fmt.Sprintf("%s → %s (%s engine)", m.source, m.target, m.conf.Engine()))

	lenTabs := lipgloss.Width(translationStatus) + lipgloss.Width(tabsRow) + 2 // still don't know why 2 cells are missing

	gap := tabGap.Render(strings.Repeat(" ", diffOrZero(m.viewport.Width, lenTabs)) + translationStatus)
	tabsRow = lipgloss.JoinHorizontal(lipgloss.Bottom, tabsRow, gap)

	view := tabsRow + "\n\n\n" + content

	return view + lipgloss.PlaceVertical(
		(m.viewport.Height+headerHeight+footerHeight)- // total height of terminal as originally received
			lipgloss.Height(view), // height of already utilized space
		lipgloss.Bottom,
		m.renderFooter())
}

func (m model) fetchTranslation(query string) tea.Cmd {
	return func() tea.Msg {
		response, err := translator.Translate(query, m.source, m.target, m.conf.Engine())
		if err != nil {
			return gotTrans{Err: err, result: err.Error()}
		}
		return gotTrans{result: response.PrettyPrint(), shortResult: response.TranslatedText}
	}
}

func (m model) fetchTextToSpeech(query string) tea.Cmd {
	return func() tea.Msg {
		response, err := translator.TextToSpeech(query, m.target)
		if err != nil {
			return gotTTS{Err: err}
		}
		return gotTTS{result: response}
	}
}

func (m *model) setState(state int) {
	m.state = state
	m.keyMgr.state = state
}

func (m *model) switchTab(direction int) {
	states := []int{TYPING, CHOOSING, TRANSLATING}

	var newState int
	if direction > 0 {
		newState = states[(m.state+1)%len(states)]
	} else {
		newState = states[(m.state-1+len(states))%len(states)]
	}

	if newState == TYPING {
		m.textInput.Focus()
	} else {
		m.textInput.Blur()
	}
	if newState != CHOOSING {
		m.help.ShowAll = false
	}

	m.setState(newState)
}

func (m *model) yankTranslated() {
	clipboard.WriteAll(m.shortResult)
}

func (m *model) renderFooter() string {
	helpMenu := m.help.View(m.keyMgr)
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
	footerBot = helpMenu + strings.Repeat(" ", diffOrZero(gapSize, helpLen)) + footerStyle.Render(footerBot)
	footer := fmt.Sprintf("%s\n%s\n%s", footerTop, footerMid, footerBot)

	return footerStyle.Render(footer)
}

func (m *model) renderTabs() []string {
	stateMaps := []string{
		TYPING:      "Text input",
		CHOOSING:    "Language selection",
		TRANSLATING: "Translation",
	}
	checkActive := func(i int, title string) string {
		if i == m.state {
			return activeTab.Render(title)
		}
		return tab.Render(title)
	}

	s := []string{}
	for state, tab := range stateMaps {
		s = append(s, checkActive(state, tab))
	}
	return s
}

func (m *model) playTTS(audio []byte) {
	streamer, format, _ := mp3.Decode(ioutil.NopCloser(bytes.NewReader(audio)))
	defer streamer.Close()

	once.Do(func() {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	})

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}

type item struct {
	title, abbreviation string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.abbreviation }
func (i item) FilterValue() string { return i.title }

func getConfLangs() []list.Item {
	items := make([]list.Item, 0)

	for abbrev, title := range translator.GetAllLanguages() {
		items = append(items, item{title, abbrev})
	}
	return items
}

func diffOrZero(x, y int) int {
	if x > y {
		return x - y
	}
	return 0
}
