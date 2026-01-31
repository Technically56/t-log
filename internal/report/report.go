package report

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/Technically56/t-log/internal/engine/parser"
)

func GenerateCsvReport(report *parser.FileReport, output_path string) error {
	file, err := os.Create(output_path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Timestamp", "Level", "SourceFile", "LogMessage", "MostMatchedRule", "MessageLevels"}
	if err := writer.Write(header); err != nil {
		return err
	}

	messageLevelsStr := ""
	for k, v := range report.MessageLevels {
		if messageLevelsStr != "" {
			messageLevelsStr += ";"
		}
		messageLevelsStr += fmt.Sprintf("%s:%d", k, v)
	}

	for _, entry := range report.Entries {
		if entry.Level == "ERROR" || entry.Level == "CRITICAL" {
			row := []string{
				entry.Timestamp,
				entry.Level,
				entry.SourceFile,
				entry.LogMessage,
				report.MostMatchedRule,
				messageLevelsStr,
			}
			if err := writer.Write(row); err != nil {
				return err
			}
		}
	}
	return nil
}
