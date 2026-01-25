package rules

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Rule struct {
	Name        string   `yaml:"name"`
	Regex       string   `yaml:"regex"`
	Severity    string   `yaml:"severity"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
}
type Ruleset struct {
	Rules []Rule `yaml:"rules"`
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
