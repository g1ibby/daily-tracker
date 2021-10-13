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

  log.Debug("req", "url", cfg.VercelURL)
	err = telegram.SetHookDoamin(cfg.TgBotToken, "http://"+cfg.VercelURL)
	if err != nil {
		log.Fatal(err)
	}
}
