package main

import (
	"fmt"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/jdejnek/oh-no-cui/http_client"
)

const (
	columnKeyID        = "id"
	columnKeyStatus    = "status"
	columnKeyType      = "type"
	columnKeyICCID     = "iccid"
	columnKeyLabel     = "label"
	columnKeyIPV4      = "ipv4"
	columnKeyConnector = "connector"
	columnKeyMCC       = "mcc"
	columnKeyMNC       = "mnc"
)

var (
	styleSubtle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))

	styleBase = lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("240"))
)

type Model struct {
	tableDefault table.Model
	rowCount     int

	columnSortKey string
	sortDirection string
}

func genRows(columnCount int, rowCount int) []table.Row {
	queryParams := http_client.QueryParams{
		Params: map[string]interface{}{
			"offset": "0",
			"limit":  "101",
		},
	}
	sims := http_client.CallApiWithParams("GET", "/sims", queryParams)
	var rows []table.Row
	for _, sim := range sims {
		row := table.NewRow(table.RowData{
			columnKeyID:        sim.Id,
			columnKeyStatus:    sim.Online,
			columnKeyType:      sim.Softsim,
			columnKeyLabel:     sim.Label,
			columnKeyICCID:     sim.Iccid,
			columnKeyConnector: sim.Connector,
			columnKeyIPV4:      sim.Ipv4,
			columnKeyMCC:       sim.Mcc,
			columnKeyMNC:       sim.Mnc,
		})
		rows = append(rows, row)
	}
	return rows
}

func genTable(columnCount int, rowCount int) table.Model {
	columns := []table.Column{
		table.NewColumn(columnKeyID, "id", 12).WithFiltered(true),
		table.NewColumn(columnKeyStatus, "online", 8),
		table.NewColumn(columnKeyType, "softsim", 8),
		table.NewColumn(columnKeyLabel, "label", 14).WithFiltered(true),
		table.NewColumn(columnKeyICCID, "iccid", 24).WithFiltered(true),
		table.NewColumn(columnKeyConnector, "connector", 16).WithFiltered(true),
		table.NewColumn(columnKeyIPV4, "ipv4", 14).WithFiltered(true),
		table.NewColumn(columnKeyMCC, "mcc", 8).WithFiltered(true),
		table.NewColumn(columnKeyMNC, "mnc", 8).WithFiltered(true),
	}

	rows := genRows(columnCount, rowCount)

	return table.New(columns).WithRows(rows).BorderRounded().Filtered(true)
}

func NewModel() Model {
	const startingRowCount = 105

	m := Model{
		rowCount: startingRowCount,
		tableDefault: genTable(9, startingRowCount).WithPageSize(20).
			WithBaseStyle(styleBase).
			Focused(true),
	}

	m.regenTableRows()

	return m
}

func (m *Model) regenTableRows() {
	m.tableDefault = m.tableDefault.WithRows(genRows(9, m.rowCount))
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmds = append(cmds, tea.Quit)
		case "u":
			m.tableDefault = m.tableDefault.WithPageSize(m.tableDefault.PageSize() - 1)

		case "o":
			m.tableDefault = m.tableDefault.WithPageSize(m.tableDefault.PageSize() + 1)

		case "r":
			m.tableDefault = m.tableDefault.WithCurrentPage(rand.Intn(m.tableDefault.MaxPages()) + 1)

		case "i":
			m.columnSortKey = columnKeyID
			m.tableDefault = m.tableDefault.SortByAsc(m.columnSortKey)
		case "I":
			m.columnSortKey = columnKeyID
			m.tableDefault = m.tableDefault.SortByDesc(m.columnSortKey)

		case "enter":
			selectedId := m.tableDefault.HighlightedRow().Data[columnKeyID]
			selectedSim := http_client.CallApiWithPath("GET", "/sims/", fmt.Sprintf("%v", selectedId))
			fmt.Sprintln(selectedSim)
		case "z":
			if m.rowCount < 10 {
				break
			}

			m.rowCount -= 10
			m.regenTableRows()

		case "x":
			m.rowCount += 10
			m.regenTableRows()
		}
	}

	m.tableDefault, cmd = m.tableDefault.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	ansiHeader := fmt.Sprintln(coloredText.Render(`
   _____    _
  / ___/   (_)   ____ ___    _____
  \__ \   / /   / __ '__ \  / ___/
 ___/ /  / /   / / / / / / (__  )
/____/  /_/   /_/ /_/ /_/ /____/
`))

	helptext := fmt.Sprintln("\n\n'/' to search table\n'I' to sort by id (asc)\n'i' to sort by id (desc)\n'q' to quit")

	body.WriteString(ansiHeader)
	body.WriteString(m.tableDefault.View())
	body.WriteString(helptext)
	return body.String()
}
