package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/Technically56/t-log/config"
	"github.com/Technically56/t-log/internal/engine/analyzer"
	"github.com/Technically56/t-log/internal/engine/parser"
	"github.com/Technically56/t-log/internal/ui/dashboard"
	filereport "github.com/Technically56/t-log/internal/ui/file_report"
	"github.com/Technically56/t-log/internal/ui/tailing"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Application struct {
	tviewApp *tview.Application
	pages    *tview.Pages
	analyzer *analyzer.Analyzer
	config   *config.AppConfig
	table    *tview.Table

	levelColorMap  map[string]tcell.Color
	sourceColorMap map[string]tcell.Color

	logChan  chan *parser.LogEntry
	stopCtx  context.Context
	stopFunc context.CancelFunc
}

func NewApplication(cfg *config.AppConfig, log_folder string) *Application {
	levelColorMap := map[string]tcell.Color{
		"DEBUG":    tcell.ColorBlue,
		"INFO":     tcell.ColorGreen,
		"WARNING":  tcell.ColorYellow,
		"ERROR":    tcell.ColorRed,
		"CRITICAL": tcell.ColorRed,
	}
	sourceColorMap := make(map[string]tcell.Color)
	for _, monitor := range cfg.Monitors {
		sourceColorMap[monitor.Path] = parseColor(monitor.SourceColor)
	}

	// Initialize Analyzer first so we can use it in callbacks
	analyzer, err := analyzer.NewAnalyzer(*cfg, log_folder)
	if err != nil {
		panic(err)
	}

	logChan := make(chan *parser.LogEntry)
	stopCtx, stopFunc := context.WithCancel(context.Background())
	pages := tview.NewPages()
	tviewApp := tview.NewApplication()

	// Setup Report Page Placeholder
	// It will be replaced dynamically when a file is selected
	reportTextView := tview.NewTextView().SetText("Rapor hazırlanıyor...")
	reportPage := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(reportTextView, 0, 1, false)

	// Setup Dashboard
	// Callbacks for navigation
	onReportSelect := func(selectedPath string) {
		// 1. Reset Report Page to Loading State
		loadingTextView := tview.NewTextView().
			SetText(fmt.Sprintf("\n\nRapor hazırlanıyor...\nDosya: %s\n\nLütfen bekleyin.", selectedPath)).
			SetTextAlign(tview.AlignCenter).
			SetDynamicColors(true)

		loadingPage := tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(loadingTextView, 0, 1, false)

		pages.RemovePage("report_page")
		pages.AddPage("report_page", loadingPage, true, true)
		pages.SwitchToPage("report_page")

		analyzer.Analyze(selectedPath, func(report *parser.FileReport, err error) {
			tviewApp.QueueUpdateDraw(func() {
				if err != nil {
					errorView := tview.NewTextView().
						SetDynamicColors(true).
						SetTextAlign(tview.AlignCenter).
						SetText(fmt.Sprintf("\n\n[red]HATA OLUŞTU[-]\n\nDosya analiz edilemedi:\n%v\n\n[yellow]Geri dönmek için ESC tuşuna basınız.", err))

					errorPage := tview.NewFlex().SetDirection(tview.FlexRow).
						AddItem(errorView, 0, 1, false)
					errorPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
						if event.Key() == tcell.KeyEsc {
							pages.SwitchToPage("dashboard")
						}
						return event
					})

					pages.RemovePage("report_page")
					pages.AddPage("report_page", errorPage, true, true)
					return
				}

				fullReportPage := filereport.DrawReportPage(pages, report)
				pages.RemovePage("report_page")
				pages.AddPage("report_page", fullReportPage, true, true)
			})
		})
	}
	onLiveSelect := func() {

	}

	dashboardPage := dashboard.DrawDashboard(tviewApp, pages, cfg.Monitors, onReportSelect, onLiveSelect)
	liveViewTable := tailing.DrawLiveView(pages, stopFunc)

	pages.AddPage("dashboard", dashboardPage, true, true)
	pages.AddPage("live_view", liveViewTable, true, false)
	pages.AddPage("report_page", reportPage, true, false)

	return &Application{
		tviewApp:       tviewApp,
		pages:          pages,
		analyzer:       analyzer,
		config:         cfg,
		levelColorMap:  levelColorMap,
		sourceColorMap: sourceColorMap,
		logChan:        logChan,
		stopCtx:        stopCtx,
		stopFunc:       stopFunc,
		table:          liveViewTable,
	}
}
func parseColor(colorName string) tcell.Color {
	c := strings.ToLower(strings.TrimSpace(colorName))

	colors := map[string]tcell.Color{
		"white":  tcell.ColorWhite,
		"black":  tcell.ColorBlack,
		"gray":   tcell.ColorGray,
		"silver": tcell.ColorSilver,

		"red":       tcell.ColorRed,
		"maroon":    tcell.ColorMaroon,
		"orange":    tcell.ColorOrange,
		"orangered": tcell.ColorOrangeRed,
		"darkred":   tcell.ColorDarkRed,
		"gold":      tcell.ColorGold,
		"yellow":    tcell.ColorYellow,

		"green":       tcell.ColorGreen,
		"lime":        tcell.ColorLime,
		"forestgreen": tcell.ColorForestGreen,
		"springgreen": tcell.ColorSpringGreen,

		"blue":    tcell.ColorBlue,
		"navy":    tcell.ColorNavy,
		"aqua":    tcell.ColorAqua,
		"teal":    tcell.ColorTeal,
		"purple":  tcell.ColorPurple,
		"fuchsia": tcell.ColorFuchsia,
		"pink":    tcell.ColorPink,
		"hotpink": tcell.ColorHotPink,
	}

	if color, ok := colors[c]; ok {
		return color
	}

	if strings.HasPrefix(c, "#") {
		return tcell.GetColor(c)
	}

	return tcell.ColorWhite
}
func (a *Application) Run() {
	a.tviewApp.SetRoot(a.pages, true).EnableMouse(true)
	a.HandleLogEntry()
	a.StartTailing()

	if err := a.tviewApp.Run(); err != nil {
		panic(err)
	}
}
func (a *Application) Stop() {
	a.tviewApp.Stop()
}
func (a *Application) StartTailing() {
	a.StopTailing()
	a.stopCtx, a.stopFunc = context.WithCancel(context.Background())
	go a.analyzer.Tail(a.logChan, a.stopCtx)
}
func (a *Application) StopTailing() {
	if a.stopFunc != nil {
		a.stopFunc()
	}
}
func (a *Application) HandleLogEntry() {
	go func() {
		for entry := range a.logChan {
			levelColor := a.levelColorMap[entry.Level]
			if levelColor == 0 {
				levelColor = tcell.ColorWhite
			}
			sourceColor := a.sourceColorMap[entry.SourceFile]
			a.tviewApp.QueueUpdateDraw(func() {
				tailing.UpdateTable(a.table, entry, levelColor, sourceColor)
			})
		}
	}()
}
