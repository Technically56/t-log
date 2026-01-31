package main

import (
	"github.com/Technically56/t-log/config"
	app "github.com/Technically56/t-log/internal/ui"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	application := app.NewApplication(cfg, "logs")

	application.StartTailing()

	application.Run()
}
