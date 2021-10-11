package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/g1ibby/daily-tracker/internal/app"
	"github.com/g1ibby/daily-tracker/internal/config"
	"github.com/g1ibby/daily-tracker/internal/sheet"
	"github.com/g1ibby/daily-tracker/internal/telegram"
	"github.com/powerman/structlog"
)

var (
	log = structlog.New()
)

func main() {
	cfg, err := config.GetApp()
	if err != nil {
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
	lambda.Start(func(request events.APIGatewayProxyRequest) (err error) {
		if cfg.ApiSecret == "" {
			log.Warn("ApiSecret is empty", "secret", cfg.ApiSecret)
			return nil
		}
		secret, ok := request.QueryStringParameters["secret"]
		if !ok {
			return nil
		}
		if cfg.ApiSecret != secret {
			log.Warn("wrong secret", "secret", secret)
			return
		}

		a.Shedule(tg.Ask)
		return nil
	})
}
