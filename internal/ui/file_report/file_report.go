package filereport

import (
	"fmt"

	"github.com/Technically56/t-log/internal/engine/parser"
	reporter "github.com/Technically56/t-log/internal/report"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func DrawReportPage(pages *tview.Pages, report *parser.FileReport) *tview.Flex {
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	mainFlex.SetBorder(true).SetTitle(" [ ANALİZ ÖZETİ ] ").SetTitleAlign(tview.AlignCenter)

	header := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("\nDosya: [yellow]%s[-]\nToplam Kayıt: [white]%d[-]\n", report.FilePath, len(report.Entries)))

	middleFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	statsTable := tview.NewTable().SetBorders(false)
	statsTable.SetTitle(" Seviye Dağılımı ").SetBorder(true)

	row := 0
	for level, count := range report.MessageLevels {
		color := tcell.ColorWhite
		switch level {
		case "CRITICAL", "ERROR":
			color = tcell.ColorRed
		case "WARNING":
			color = tcell.ColorYellow
		case "INFO":
			color = tcell.ColorGreen
		}

		statsTable.SetCell(row, 0, tview.NewTableCell(level).SetTextColor(tcell.ColorWhite))
		statsTable.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf(": [ %d ]", count)).
			SetTextColor(color))
		row++
	}

	description := ""
	if report.MostMatchedRuleObj != nil {
		description = report.MostMatchedRuleObj.Description
	}
	ruleName := report.MostMatchedRule
	if ruleName == "" {
		ruleName = "Eşleşen Kural Bulunamadı"
	}
	ruleView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("\n\nEN ÇOK TETİKLENEN KURAL:\n\n[red]%s[-] \n\n%s", ruleName, description))
	ruleView.SetTitle(" Güvenlik Özeti ").SetBorder(true)

	middleFlex.AddItem(statsTable, 0, 1, false)
	middleFlex.AddItem(ruleView, 0, 1, false)

	footer := tview.NewFlex().SetDirection(tview.FlexColumn)

	backBtn := tview.NewButton("Geri Dön (ESC)").SetSelectedFunc(func() {
		pages.SwitchToPage("dashboard")
	})

	exportFunc := func() {
		err := reporter.GenerateCsvReport(report, "./output/output.csv")
		var modal *tview.Modal
		if err != nil {
			modal = tview.NewModal().
				SetText(fmt.Sprintf("Hata oluştu:\n%v", err)).
				AddButtons([]string{"Tamam"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					pages.RemovePage("alert")
				})
		} else {
			modal = tview.NewModal().
				SetText("Rapor başarıyla kaydedildi:\noutput.csv").
				AddButtons([]string{"Tamam"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					pages.RemovePage("alert")
				})
		}
		pages.AddPage("alert", modal, false, true)
	}

	exportBtn := tview.NewButton("CSV Dışa Aktar (S)").SetSelectedFunc(exportFunc)

	footer.AddItem(tview.NewBox(), 0, 1, false)
	footer.AddItem(backBtn, 20, 1, true)
	footer.AddItem(tview.NewBox(), 5, 1, false)
	footer.AddItem(exportBtn, 20, 1, false)
	footer.AddItem(tview.NewBox(), 0, 1, false)

	mainFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			pages.SwitchToPage("dashboard")
			return nil
		}
		if event.Rune() == 's' || event.Rune() == 'S' {
			exportFunc()
			return nil
		}
		return event
	})

	mainFlex.AddItem(header, 5, 1, false)
	mainFlex.AddItem(middleFlex, 0, 1, false)
	mainFlex.AddItem(footer, 3, 1, true)

	return mainFlex
}
