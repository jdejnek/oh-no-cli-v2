package main

import (
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
			Foreground(lipgloss.Color("#ffffff")).
			BorderForeground(lipgloss.Color("#9c8922")).
			Align(lipgloss.Left)
)

type Model struct {
	tableDefault table.Model
	rowCount     int
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
		}).WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Padding())
		rows = append(rows, row)
	}
	return rows
}

func genTable(columnCount int, rowCount int) table.Model {
	columns := []table.Column{
		table.NewColumn(columnKeyID, "id", 12),
		table.NewColumn(columnKeyStatus, "online", 8),
		table.NewColumn(columnKeyType, "softsim", 8),
		table.NewColumn(columnKeyLabel, "label", 14),
		table.NewColumn(columnKeyICCID, "iccid", 24),
		table.NewColumn(columnKeyConnector, "connector", 16),
		table.NewColumn(columnKeyIPV4, "ipv4", 12),
		table.NewColumn(columnKeyMCC, "mcc", 8),
		table.NewColumn(columnKeyMNC, "mnc", 8),
	}

	rows := genRows(columnCount, rowCount)

	return table.New(columns).WithRows(rows).HeaderStyle(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#c4ac21")))
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
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		case "u":
			m.tableDefault = m.tableDefault.WithPageSize(m.tableDefault.PageSize() - 1)

		case "i":
			m.tableDefault = m.tableDefault.WithPageSize(m.tableDefault.PageSize() + 1)

		case "r":
			m.tableDefault = m.tableDefault.WithCurrentPage(rand.Intn(m.tableDefault.MaxPages()) + 1)

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
	pad := lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color("#43e06d"))

	tables := []string{
		lipgloss.JoinVertical(lipgloss.Center, "Sims", pad.Render(m.tableDefault.View())),
	}

	body.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, tables...))

	return body.String()
}
