package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	flag "github.com/spf13/pflag"

	"github.com/slavka/another-less/internal/app"
	"github.com/slavka/another-less/internal/ingest"
)

var (
	version   = "dev"
	commit    = "none"
	sourceURL = "https://github.com/slavka/another-less"
	goVersion = "unknown"
)

func main() {
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("aless %s (commit: %s, go: %s)\n%s\n", version, commit, goVersion, sourceURL)
		os.Exit(0)
	}

	var reader *ingest.Reader

	args := flag.Args()
	if len(args) > 0 {
		// Read from file(s) - for now just first file
		r, err := ingest.NewFileReader(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		reader = r
	} else {
		// Read from stdin
		fi, _ := os.Stdin.Stat()
		if fi.Mode()&os.ModeCharDevice != 0 {
			fmt.Fprintln(os.Stderr, "usage: aless [file] or pipe data to aless")
			os.Exit(1)
		}
		reader = ingest.NewReader(os.Stdin)
	}

	reader.Start()

	// Detect if stdout is piped (for Q pipe mode)
	outFi, _ := os.Stdout.Stat()
	pipeOutput := outFi.Mode()&os.ModeCharDevice == 0

	model := app.New(app.Config{
		LinesCh:    reader.Chan(),
		PipeOutput: pipeOutput,
	})

	// When stdout is piped, use stderr for the TUI
	var opts []tea.ProgramOption
	opts = append(opts, tea.WithAltScreen())
	if pipeOutput {
		opts = append(opts, tea.WithOutput(os.Stderr))
	}

	p := tea.NewProgram(model, opts...)

	finalModel, err := p.Run()
	reader.Stop()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// If pipe quit was used, write the buffer to stdout
	if m, ok := finalModel.(app.Model); ok {
		if buf := m.PipeBuffer(); buf != "" {
			fmt.Print(buf)
		}
	}
}
