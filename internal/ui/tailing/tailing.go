package tailing

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/Technically56/t-log/internal/engine/parser"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func DrawLiveView(pages *tview.Pages, stopFunc context.CancelFunc) *tview.Table {
	table := tview.NewTable().
		SetSelectable(true, false).
		SetFixed(1, 1).
		SetSeparator('|')

	table.SetBorder(true).SetTitle(" [ BİRLEŞTİRİLMİŞ CANLI LOG AKIŞI ] ")

	headers := []string{"LEVEL", "TIMESTAMP", "SOURCE", "MESSAGE", "MATCHED RULE"}
	for i, h := range headers {
		table.SetCell(0, i, tview.NewTableCell(h).
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold).
			SetSelectable(false))
	}
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			pages.SwitchToPage("dashboard")
			stopFunc()
		}
		return event
	})

	return table
}

func UpdateTable(table *tview.Table, entry *parser.LogEntry, levelColor, sourceColor tcell.Color) {
	row := table.GetRowCount()

	ruleNames := []string{}
	for _, r := range entry.MatchedRules {
		ruleNames = append(ruleNames, r.Name)
	}
	ruleStr := strings.Join(ruleNames, ", ")
	if ruleStr == "" {
		ruleStr = "-"
	}
	table.SetCell(row, 0, tview.NewTableCell(entry.Level).
		SetTextColor(levelColor).
		SetAttributes(tcell.AttrBold))

	table.SetCell(row, 1, tview.NewTableCell(entry.Timestamp).
		SetTextColor(tcell.ColorWhite))

	table.SetCell(row, 2, tview.NewTableCell(filepath.Base(entry.SourceFile)).
		SetTextColor(sourceColor))

	table.SetCell(row, 3, tview.NewTableCell(entry.LogMessage).
		SetTextColor(tcell.ColorWhite))

	table.SetCell(row, 4, tview.NewTableCell(ruleStr).
		SetTextColor(tcell.ColorMediumPurple))
	table.ScrollToEnd()
}
