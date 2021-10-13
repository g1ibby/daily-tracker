
package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/g1ibby/daily-tracker/internal/config"
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
	lambda.Start(func(request events.APIGatewayProxyRequest) (err error) {
    log.Debug("req", "domain", cfg.Domain)
		return telegram.SetHookDoamin(cfg.TgBotToken, cfg.Domain)
	})
}
