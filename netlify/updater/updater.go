package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/g1ibby/daily-tracker/internal/app"
	"github.com/g1ibby/daily-tracker/internal/config"
	"github.com/g1ibby/daily-tracker/internal/sheet"
	"github.com/g1ibby/daily-tracker/internal/telegram"
	"github.com/powerman/structlog"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	log = structlog.New()
)

func main() {
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

	lambda.Start(func(request events.APIGatewayProxyRequest) (err error) {
		var u tb.Update
		if err = json.Unmarshal([]byte(request.Body), &u); err == nil {
			tg.ProcessUpdate(u)
		}
		return
	})
}
