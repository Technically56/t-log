package parser

import (
	"regexp"

	"github.com/Technically56/bubble-log/internal/engine/rules"
)

type Parser struct {
	ruleset *rules.Ruleset
}

type LogEntry struct {
	Timestamp    string
	Source       string
	Message      string
	Level        string
	MatchedRules []rules.Rule
}

func (p *Parser) ParseLine(line string) *LogEntry {
	for _, rule := range p.ruleset.Rules {
		regexp.MustCompile(rule.Regex)
	}
}
