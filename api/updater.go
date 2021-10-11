package handler

import (
	"encoding/json"
	"io/ioutil"

	"net/http"

	"github.com/g1ibby/daily-tracker/internal/app"
	"github.com/g1ibby/daily-tracker/internal/sheet"
	"github.com/g1ibby/daily-tracker/internal/telegram"
	"github.com/g1ibby/daily-tracker/internal/config"
	"github.com/powerman/structlog"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	log = structlog.New()
)

func HandlerUpdater(w http.ResponseWriter, r *http.Request) {
  cfg, err := config.GetApp()
	if err != nil {
		log.Fatal(err)
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
    log.Fatal(err)
	}

	var u tb.Update
	if err = json.Unmarshal(body, &u); err == nil {
		tg.ProcessUpdate(u)
	}
}
