package handler

import (
	"net/http"

	"github.com/g1ibby/daily-tracker/internal/config"
	"github.com/g1ibby/daily-tracker/internal/telegram"
)

func HandlerSetHook(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.GetApp()
	if err != nil {
		log.Fatal(err)
	}

  log.Debug("req", "domain", cfg.Domain)
	err = telegram.SetHookDoamin(cfg.TgBotToken, cfg.Domain)
	if err != nil {
		log.Fatal(err)
	}
}
