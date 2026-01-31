package rules

import (
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type Rule struct {
	Name        string   `yaml:"name"`
	Regex       string   `yaml:"regex"`
	Level       string   `yaml:"severity"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
}
type Ruleset struct {
	Rules []Rule `yaml:"rules"`
}
type CompiledRule struct {
	Rule          Rule
	CompiledRegex *regexp.Regexp
}
type CompiledRuleset struct {
	Rules []CompiledRule
}

func NewRuleset(file_location string) (*Ruleset, error) {
	file_data, err := os.ReadFile(file_location)
	if err != nil {
		return nil, err
	}

	var ruleset Ruleset
	err = yaml.Unmarshal(file_data, &ruleset)
	if err != nil {
		return nil, err
	}

	return &ruleset, nil
}
