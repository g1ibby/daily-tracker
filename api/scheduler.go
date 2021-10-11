package handler

import (
	"net/http"

	"github.com/g1ibby/daily-tracker/internal/app"
	"github.com/g1ibby/daily-tracker/internal/config"
	"github.com/g1ibby/daily-tracker/internal/sheet"
	"github.com/g1ibby/daily-tracker/internal/telegram"
)

func HandlerScheduler(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.GetApp()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.ApiSecret == "" {
		log.Warn("ApiSecret is empty", "secret", cfg.ApiSecret)
		return
	}
	secret := string(r.URL.Query().Get("secret"))
	if cfg.ApiSecret != secret {
		log.Fatal("wrong secret", "secret", secret)
		return
	}

	sh, err := sheet.NewGSheet(cfg.SheetSecret, cfg.SheetID)
	if err != nil {
		log.Fatal(err)
	}
	a := app.New(sh, cfg.TgUserID, cfg.ApiSecret)
	tg, err := telegram.New(cfg.TgBotToken, a, true)
	if err != nil {
		log.Fatal(err)
	}
	a.Shedule(tg.Ask)
}
