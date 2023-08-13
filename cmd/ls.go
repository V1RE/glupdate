package cmd

import (
	"fmt"
	"glupdate/updaters"
	"log"
	"sort"
	"sync"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"

	"github.com/spf13/cobra"
)

var managers = []updaters.UpdaterFactory{
	updaters.Pnpm(),
	updaters.Brew(),
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List available updates",
	Run:   listUpdates,
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

func listUpdates(cmd *cobra.Command, args []string) {
	pkgsChan := make(chan []updaters.Package, len(managers))

	var wg sync.WaitGroup

	// Concurrently fetch updates from all managers
	for _, factory := range managers {
		wg.Add(1)
		go func(factory updaters.UpdaterFactory) {
			defer wg.Done()

			manager, err := factory()
			if err != nil {
				log.Println(err)
				return
			}

			pkgs, err := manager.List()
			if err != nil {
				log.Println(err)
				return
			}

			pkgsChan <- pkgs
		}(factory)
	}

	go func() {
		wg.Wait()
		close(pkgsChan)
	}()

	columns := []table.Column{
		{Title: "Package-manager", Width: 15},
		{Title: "Name", Width: 4},
		{Title: "Installed", Width: 9},
		{Title: "Latest", Width: 6},
	}

	var pkgs []updaters.Package
	for pkgBatch := range pkgsChan {
		pkgs = append(pkgs, pkgBatch...)
	}

	sort.Slice(pkgs, func(i, j int) bool {
		row1 := pkgs[i]
		row2 := pkgs[j]

		return row1.Updater.Name()+row1.Name < row2.Updater.Name()+row2.Name
	})

	rows := make([]table.Row, 0, len(pkgs))

	for _, pkg := range pkgs {
		row := table.Row{pkg.Updater.Name(), pkg.Name, pkg.Current, pkg.Latest}
		rows = append(rows, row)

		for i, v := range columns {
			columns[i].Width = max(v.Width, len(row[i]))
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(len(rows)),
	)

	baseStyle := lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = lipgloss.NewStyle()
	t.SetStyles(s)

	fmt.Println(baseStyle.Render(t.View()))
}
