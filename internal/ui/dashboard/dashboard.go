package dashboard

import (
	"github.com/Technically56/t-log/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func DrawDashboard(app *tview.Application, pages *tview.Pages, monitors []config.MonitorConfig, onReportSelect func(string), onLiveSelect func()) tview.Primitive {
	list := tview.NewList().ShowSecondaryText(true)
	list.SetBorder(true).SetTitle(" [ SOC Komuta Merkezi ] ")

	for i, monitor := range monitors {
		p := monitor.Path
		list.AddItem(p, "Rapor için R Tuşuna Basın, Canlı İzleme için L Tuşuna Basın", rune('1'+i), nil)
	}

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if list.GetItemCount() == 0 {
			return event
		}

		index := list.GetCurrentItem()
		path, secondary := list.GetItemText(index)

		if secondary == "Uygulamayı kapat" {
			return event
		}

		switch event.Rune() {
		case 'r', 'R':
			if onReportSelect != nil && path != "" {
				onReportSelect(path)
			}
			pages.SwitchToPage("report_page")
			return nil
		case 'l', 'L':
			if onLiveSelect != nil {
				onLiveSelect()
			}
			pages.SwitchToPage("live_view")
			return nil
		}
		return event
	})

	list.AddItem("Çıkış", "Uygulamayı kapat", 'q', func() {
		app.Stop()
	})

	return list
}
