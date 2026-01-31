package parser

import (
	"bufio"
	"os"
	"regexp"
	"time"

	"github.com/Technically56/t-log/internal/engine/rules"
)

type LogEntry struct {
	Timestamp    string
	SourceFile   string
	LogMessage   string
	Level        string
	MatchedParts map[string]MultiPartRegexMatch
	MatchedRules []rules.Rule
}

type MultiPartRegexMatch struct {
	Parts map[string]string
}

type FileReport struct {
	FilePath           string
	Entries            []*LogEntry
	MessageLevels      map[string]int
	MostMatchedRule    string
	MostMatchedRuleObj *rules.Rule
}

type Parser struct {
	Ruleset        *rules.CompiledRuleset
	FilePath       string
	TimestampRegex *regexp.Regexp
}

func NewParser(ruleset *rules.Ruleset, file_path string) *Parser {
	compiledRules := rules.CompiledRuleset{}
	timestamp := regexp.MustCompile(`(?P<Timestamp>\d{4}[-/]\d{2}[-/]\d{2}[:\sT]\d{2}:\d{2}:\d{2}(\.\d+)?(Z|[+-]\d{2}:?\d{2})?|\w{3}\s+\d{1,2}\s\d{2}:\d{2}:\d{2})`)
	for _, rule := range ruleset.Rules {
		if rule.Name == "Timestamp" {
			timestamp = regexp.MustCompile(rule.Regex)
			continue
		}
		compiledRule := rules.CompiledRule{
			Rule:          rule,
			CompiledRegex: regexp.MustCompile(rule.Regex),
		}
		compiledRules.Rules = append(compiledRules.Rules, compiledRule)

	}
	return &Parser{
		Ruleset:        &compiledRules,
		FilePath:       file_path,
		TimestampRegex: timestamp,
	}
}
func (p *Parser) ParseLine(line string) *LogEntry {
	matchedRules := []rules.Rule{}
	matchedParts := make(map[string]MultiPartRegexMatch)
	maxLevel := "DEBUG"
	timestamp := p.extractTimestamp(line)

	for _, rule := range p.Ruleset.Rules {
		rgx := rule.CompiledRegex

		if rgx.MatchString(line) {
			matchedRules = append(matchedRules, rule.Rule)
			submatches := rgx.FindStringSubmatch(line)
			names := rgx.SubexpNames()

			currentParts := make(map[string]string)
			namedGroupFound := false

			for i, name := range names {
				if i > 0 && i < len(submatches) && submatches[i] != "" {
					if name != "" {
						currentParts[name] = submatches[i]
						namedGroupFound = true
					}
				}
			}

			if !namedGroupFound {
				currentParts["default"] = rgx.FindString(line)
			}

			matchedParts[rule.Rule.Name] = MultiPartRegexMatch{Parts: currentParts}
			maxLevel = findGreaterLevel(maxLevel, rule.Rule.Level)
		}
	}

	if maxLevel == "DEBUG" {
		upperLine := line

		if regexp.MustCompile(`(?i)\b(CRITICAL|FATAL|PANIC|EMERG|ALERT)\b`).MatchString(upperLine) {
			maxLevel = "CRITICAL"
		} else if regexp.MustCompile(`(?i)\b(ERROR|ERR|FAIL|FAILURE)\b`).MatchString(upperLine) {
			maxLevel = "ERROR"
		} else if regexp.MustCompile(`(?i)\b(WARNING|WARN)\b`).MatchString(upperLine) {
			maxLevel = "WARNING"
		} else if regexp.MustCompile(`(?i)\b(INFO|NOTICE|Accepted)\b`).MatchString(upperLine) {
			maxLevel = "INFO"
		}
	}

	return &LogEntry{
		Timestamp:    timestamp,
		SourceFile:   p.FilePath,
		LogMessage:   line,
		Level:        maxLevel,
		MatchedParts: matchedParts,
		MatchedRules: matchedRules,
	}
}
func findGreaterLevel(level1, level2 string) string {
	levelPriority := map[string]int{
		"DEBUG":    1,
		"INFO":     2,
		"WARNING":  3,
		"ERROR":    4,
		"CRITICAL": 5,
	}
	if levelPriority[level1] > levelPriority[level2] {
		return level1
	}
	return level2
}
func (p *Parser) ParseFile() (*FileReport, error) {
	file, err := os.Open(p.FilePath)
	if err != nil {
		return &FileReport{}, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	entries := []*LogEntry{}
	for scanner.Scan() {
		line := scanner.Text()
		entry := p.ParseLine(line)
		entries = append(entries, entry)

	}
	if err := scanner.Err(); err != nil {
		return &FileReport{}, err
	}
	mostMatchedRuleName, mostMatchedRule := findMostMatchedRule(entries)
	return &FileReport{
		FilePath:           p.FilePath,
		Entries:            entries,
		MessageLevels:      countMessageLevels(entries),
		MostMatchedRule:    mostMatchedRuleName,
		MostMatchedRuleObj: mostMatchedRule,
	}, nil
}
func countMessageLevels(entries []*LogEntry) map[string]int {
	counts := map[string]int{}
	for _, entry := range entries {
		counts[entry.Level]++
	}
	return counts
}
func findMostMatchedRule(entries []*LogEntry) (string, *rules.Rule) {
	ruleCount := map[string]int{}
	for _, entry := range entries {
		for _, rule := range entry.MatchedRules {
			ruleCount[rule.Name]++
		}
	}
	maxCount := 0
	mostMatchedRule := ""
	for ruleName, count := range ruleCount {
		if count > maxCount {
			maxCount = count
			mostMatchedRule = ruleName
		}
	}
	for _, entry := range entries {
		for _, rule := range entry.MatchedRules {
			if rule.Name == mostMatchedRule {
				return mostMatchedRule, &rule
			}
		}
	}
	return mostMatchedRule, nil
}
func (p *Parser) extractTimestamp(line string) string {
	extractedTime := p.TimestampRegex.FindString(line)
	layout := "2006/01/02 - 15:04:05"
	if extractedTime != "" {
		standardTime := standardizeTimestamp(extractedTime)
		return standardTime.Format(layout)
	}
	return time.Now().Format(layout)
}
func standardizeTimestamp(raw string) time.Time {
	layouts := []string{
		time.RFC3339Nano,            // ISO8601 with fractional seconds
		"2006-01-02T15:04:05Z07:00", // RFC3339 / ISO8601
		"Jan 02 15:04:05",           // Syslog (No year)
		"2006-01-02 15:04:05",       // Standard Database/Simple
		"02/Jan/2006:15:04:05",      // Apache/Nginx
	}

	var t time.Time
	var err error

	for _, layout := range layouts {
		t, err = time.Parse(layout, raw)
		if err == nil {
			if t.Year() == 0 {
				t = t.AddDate(time.Now().Year(), 0, 0)
			}
			return t
		}
	}
	return time.Now()
}
