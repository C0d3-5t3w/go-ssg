package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/C0d3-5t3w/go-ssg/inc/config"
	"github.com/C0d3-5t3w/go-ssg/inc/editor"
	"github.com/C0d3-5t3w/go-ssg/inc/gen"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	configPath    string
	genContentDir string
	genOutputDir  string
	genSiteTitle  string

	serveOutputDir string
	servePort      string
)

var rootCmd = &cobra.Command{
	Use:   "go-ssg",
	Short: "Go-SSG is a simple static site generator",
	Long: `A Fast and Flexible Static Site Generator built with Go.
Use subcommands to generate, serve, or edit site content.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cp := configPath
		if !cmd.Flags().Changed("config") {
			if cmd.Parent() != nil && cmd.Parent().PersistentFlags().Lookup("config") != nil {
				cp = cmd.Parent().PersistentFlags().Lookup("config").Value.String()
			}
		}

		loadedCfg, err := config.LoadConfig(cp)
		if err != nil {
			return fmt.Errorf("error loading config file %s: %w", cp, err)
		}

		if cmd.Name() == generateCmd.Name() {
			if cmd.Flags().Changed("contentDir") {
				loadedCfg.ContentDir = genContentDir
			}
			if cmd.Flags().Changed("outputDir") {
				loadedCfg.OutputDir = genOutputDir
			}
			if cmd.Flags().Changed("siteTitle") {
				loadedCfg.SiteTitle = genSiteTitle
			}
		}

		if cmd.Name() == serveCmd.Name() {
			if cmd.Flags().Changed("outputDir") {
				loadedCfg.OutputDir = serveOutputDir
			}
			if cmd.Flags().Changed("port") {
				loadedCfg.ServerPort = servePort
			}
		}

		config.AppConfig = *loadedCfg
		return nil
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate static HTML files from Markdown",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &config.AppConfig
		fmt.Printf("Generating site with Configuration:\n")
		fmt.Printf("  Site Title: %s\n", cfg.SiteTitle)
		fmt.Printf("  Content Dir: %s\n", cfg.ContentDir)
		fmt.Printf("  Output Dir: %s\n", cfg.OutputDir)

		if err := gen.GEN(cfg); err != nil {
			log.Fatalf("Failed to generate site: %v", err)
		}
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the generated static files",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &config.AppConfig

		if _, err := os.Stat(cfg.OutputDir); os.IsNotExist(err) {
			log.Fatalf("Output directory '%s' does not exist. Generate the site first using 'go-ssg generate'.", cfg.OutputDir)
		}

		fmt.Printf("Serving files from '%s' on http://localhost:%s\n", cfg.OutputDir, cfg.ServerPort)
		fs := http.FileServer(http.Dir(cfg.OutputDir))
		http.Handle("/", fs)
		err := http.ListenAndServe(":"+cfg.ServerPort, nil)
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	},
}

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

type item struct {
	title, desc, fullPath string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type listModel struct {
	list         list.Model
	choice       string
	quitting     bool
	err          error
	initialItems []list.Item
}

func (m listModel) Init() tea.Cmd { return nil }

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.fullPath
			}
			m.quitting = true
			return m, tea.Quit
		}
	case error:
		m.err = msg
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	if m.choice != "" {
		return docStyle.Render(fmt.Sprintf("Selected: %s. Opening in editor...", m.choice))
	}
	if m.quitting {
		return docStyle.Render("Quitting...")
	}
	if m.err != nil {
		return docStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}
	return docStyle.Render(m.list.View())
}

func newFileListModel(items []list.Item) listModel {
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select a file to edit"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	l.Styles.PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)

	return listModel{list: l, initialItems: items}
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a Markdown or HTML file using a TUI and external editor",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &config.AppConfig

		var items []list.Item

		mdFiles, err := ioutil.ReadDir(cfg.ContentDir)
		if err == nil {
			for _, f := range mdFiles {
				if !f.IsDir() && (strings.HasSuffix(f.Name(), ".md") || strings.HasSuffix(f.Name(), ".markdown")) {
					items = append(items, item{
						title:    f.Name(),
						desc:     "Markdown file in " + cfg.ContentDir,
						fullPath: filepath.Join(cfg.ContentDir, f.Name()),
					})
				}
			}
		} else {
			fmt.Printf("Warning: Could not read content directory %s: %v\n", cfg.ContentDir, err)
		}

		htmlFiles, err := ioutil.ReadDir(cfg.OutputDir)
		if err == nil {
			for _, f := range htmlFiles {
				if !f.IsDir() && strings.HasSuffix(f.Name(), ".html") {
					items = append(items, item{
						title:    f.Name(),
						desc:     "HTML file in " + cfg.OutputDir,
						fullPath: filepath.Join(cfg.OutputDir, f.Name()),
					})
				}
			}
		} else {
			fmt.Printf("Warning: Could not read output directory %s: %v\n", cfg.OutputDir, err)
		}

		if len(items) == 0 {
			fmt.Println("No .md or .html files found in content or output directories.")
			return
		}

		m := newFileListModel(items)
		p := tea.NewProgram(m, tea.WithAltScreen())

		finalModel, err := p.Run()
		if err != nil {
			log.Fatalf("Error running TUI: %v", err)
		}

		if fm, ok := finalModel.(listModel); ok && fm.choice != "" {
			fmt.Printf("Opening %s in editor...\n", fm.choice)
			if err := editor.OpenFileInEditor(fm.choice); err != nil {
				log.Fatalf("Failed to open file in editor: %v", err)
			}
		} else if fm.err != nil {
			log.Fatalf("TUI returned an error: %v", fm.err)
		} else if !fm.quitting && fm.choice == "" {
			fmt.Println("No file selected or TUI exited unexpectedly.")
		} else {
			fmt.Println("Editor selection cancelled.")
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", config.DefaultConfigPath, "Path to configuration file")

	generateCmd.Flags().StringVar(&genContentDir, "contentDir", "content", "Directory containing markdown content files")
	generateCmd.Flags().StringVar(&genOutputDir, "outputDir", "output", "Directory where HTML files will be generated")
	generateCmd.Flags().StringVar(&genSiteTitle, "siteTitle", "My Static Site", "Title for the site")

	serveCmd.Flags().StringVar(&serveOutputDir, "outputDir", "output", "Directory of generated files to serve")
	serveCmd.Flags().StringVarP(&servePort, "port", "p", "8080", "Port to serve the site on")

	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(editCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func CLI() {
	Execute()
}
