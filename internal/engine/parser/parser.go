package parser

import "github.com/Technically56/bubble-log/internal/engine/rules"

type Parser struct {
	filepath string
	ruleset  *rules.Ruleset
}

type LogEntry struct {
	Timestamp    string
	Source       string
	Message      string
	Level        string
	MatchedRules []rules.Rule
}
