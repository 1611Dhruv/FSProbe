package tui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mode int

const (
	modeMenu mode = iota
	modeCreateFS
	modeRunFS
)

type model struct {
	list      list.Model
	textInput textinput.Model
	mode      mode

	width, height int
	inodeBlocks   int
	dataBlocks    int
	choice        string
	quitting      bool
}

var (
	titleStyle        = lipgloss.NewStyle().Margin(1, 2).Bold(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	headerStyle       = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#0000FF")).
				Padding(1, 2).
				Bold(true)

	paneStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#FFA500")).
			Margin(1, 1).
			Padding(1, 1)
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func initialModel() model {
	items := []list.Item{
		item{title: "Select Disk Image", desc: "Choose an existing disk image to mount"},
		item{title: "Create Disk Image", desc: "Create a new disk image"},
	}

	const defaultWidth = 20

	l := list.New(items, list.NewDefaultDelegate(), defaultWidth, 14)
	l.Title = "Welcome to MyFuseApp"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle

	ti := textinput.New()
	ti.Placeholder = "Number of inode blocks"
	ti.Focus()

	return model{list: l, textInput: ti, mode: modeMenu}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if m.mode == modeMenu {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					if i.title == "Create Disk Image" {
						m.mode = modeCreateFS
						m.textInput.Focus()
					} else if i.title == "Select Disk Image" {
						// Handle selecting a disk image
					}
				}
			} else if m.mode == modeCreateFS {
				m.inodeBlocks, _ = strconv.Atoi(m.textInput.Value())
				m.textInput.SetValue("")
				m.textInput.Placeholder = "Number of data blocks"
				m.textInput.Focus()
				m.mode = modeRunFS
			} else if m.mode == modeRunFS {
				m.dataBlocks, _ = strconv.Atoi(m.textInput.Value())
				m.createFS()
				m.runFS()
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	if m.mode == modeCreateFS || m.mode == modeRunFS {
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return fmt.Sprintf("You chose: %s\n", m.choice)
	}
	if m.quitting {
		return "Goodbye!\n"
	}
	switch m.mode {
	case modeMenu:
		header := headerStyle.Render("Header")
		headerHeight := lipgloss.Height(header)

		paneHeight := (m.height - headerHeight) / 2
		paneWidth := m.width / 2

		pane := func(content string) string {
			return paneStyle.
				Width(paneWidth).
				Height(paneHeight).
				Render(content)
		}

		grid := lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.JoinVertical(lipgloss.Top, pane("Pane 1"), pane("Pane 3")),
			lipgloss.JoinVertical(lipgloss.Top, pane("Pane 2"), pane("Pane 4")),
		)

		return lipgloss.JoinVertical(lipgloss.Top, header, grid)
	case modeCreateFS, modeRunFS:
		return m.textInput.View()
	default:
		return ""
	}
}

func (m model) createFS() {
	// Round inodeBlocks and dataBlocks to multiples of 2
	m.inodeBlocks = (m.inodeBlocks + 1) &^ 1
	m.dataBlocks = (m.dataBlocks + 1) &^ 1

	// Initialize the filesystem
	fs := initializeFS(m.inodeBlocks, m.dataBlocks)

	// Save the filesystem to a disk image file
	file, err := os.Create("disk.img")
	if err != nil {
		log.Fatalf("failed to create disk image: %v", err)
	}
	defer file.Close()

	_, err = file.Write(fs)
	if err != nil {
		log.Fatalf("failed to write to disk image: %v", err)
	}
}

func initializeFS(inodeBlocks, dataBlocks int) []byte {
	totalBlocks := inodeBlocks + dataBlocks
	blockSize := 4096
	fs := make([]byte, totalBlocks*blockSize)

	// Simple FS initialization (metadata, inode table, data blocks, etc.)
	// This is just an example, you can implement your actual FS initialization logic here

	return fs
}

func (m model) runFS() {
	cmd := exec.Command("./myfuseapp", "disk.img", "mnt")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run fuse application: %v", err)
	}
}

func Start() error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
