package analyzer

import (
	"context"
	"fmt"
	"io"

	"github.com/Technically56/t-log/config"
	"github.com/Technically56/t-log/internal/engine/parser"
	"github.com/Technically56/t-log/internal/engine/rules"
	"github.com/hpcloud/tail"
)

type Analyzer struct {
	parsers    map[string]*parser.Parser
	log_folder string
}

func NewAnalyzer(cfg config.AppConfig, log_folder string) (*Analyzer, error) {
	parsers := make(map[string]*parser.Parser)
	for _, monitor := range cfg.Monitors {
		ruleset, err := rules.NewRuleset(monitor.RulesPath)
		if err != nil {
			return nil, err
		}
		parser := parser.NewParser(ruleset, monitor.Path)
		parsers[monitor.Path] = parser
	}
	return &Analyzer{
		parsers:    parsers,
		log_folder: log_folder,
	}, nil
}

func (a *Analyzer) Analyze(file_path string, onComplete func(*parser.FileReport, error)) {
	go func() {
		parser := a.parsers[file_path]
		if parser == nil {
			onComplete(nil, fmt.Errorf("parser not found for file: %s", file_path))
			return
		}
		report, err := parser.ParseFile()
		if err != nil {
			onComplete(nil, err)
			return
		}
		onComplete(report, nil)
	}()
}

func (a *Analyzer) Tail(outChan chan *parser.LogEntry, stopCtx context.Context) {
	for _, current_parser := range a.parsers {
		go func(parser *parser.Parser) {
			t, err := tail.TailFile(current_parser.FilePath, tail.Config{
				Follow:    true,
				MustExist: true,
				Poll:      true,
				ReOpen:    true,
				Location:  &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd},
				Logger:    tail.DiscardingLogger,
			})
			if err != nil {
				return
			}
			defer t.Stop()
			for {
				select {
				case line := <-t.Lines:
					entry := current_parser.ParseLine(line.Text)
					outChan <- entry
				case <-stopCtx.Done():
					return
				}
			}
		}(current_parser)
	}
}
