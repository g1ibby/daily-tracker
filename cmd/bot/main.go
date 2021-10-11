package main

import (
	"time"

	"github.com/powerman/structlog"

	"github.com/g1ibby/daily-tracker/internal/app"
	"github.com/g1ibby/daily-tracker/internal/config"
	"github.com/g1ibby/daily-tracker/internal/sheet"
	"github.com/g1ibby/daily-tracker/internal/telegram"
	"github.com/go-co-op/gocron"
)

var (
	log = structlog.New()
)

func main() {
	cfg, err := config.GetApp()
	if err != nil {
		log.Fatal(err)
	}
	s := gocron.NewScheduler(time.UTC)

	var sh app.Repo
	if cfg.SheetSecret != "" && cfg.SheetID != "" {
		sh, err = sheet.NewGSheet(cfg.SheetSecret, cfg.SheetID)
		if err != nil {
			log.Fatal(err)
		}
	}
	a := app.New(sh, cfg.TgUserID, cfg.ApiSecret)
	tg, err := telegram.New(cfg.TgBotToken, a, false)
	if err != nil {
		log.Fatal(err)
	}

	// run every day
	s.Cron("0 10 * * *").Do(func() {
		a.Shedule(tg.Ask)
		time.Sleep(2 * time.Second)
	})
	s.StartAsync()

	tg.Start()
}
