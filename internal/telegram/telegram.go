package telegram

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/g1ibby/daily-tracker/internal/app"
	"github.com/powerman/structlog"
	tb "gopkg.in/tucnak/telebot.v2"
)

// Ctx is a synonym for convenience.
type Ctx = context.Context

type Svc struct {
	b   *tb.Bot
	log *structlog.Logger
	app app.Appl
}

var (
	selector = &tb.ReplyMarkup{}
	btnYes   = selector.Data("YES", "btnYes")
	btnNo    = selector.Data("NO", "btnNo")
)

func New(botToken string, a app.Appl, serverless bool) (*Svc, error) {
	ctx := context.Background()
	log := structlog.FromContext(ctx, nil)
	settings := tb.Settings{
		Token:       botToken,
		Synchronous: true,
		ParseMode:   tb.ModeMarkdown,
	}
	if !serverless {
		settings.Synchronous = false
		settings.Poller = &tb.LongPoller{Timeout: 10 * time.Second}
	}
	b, err := tb.NewBot(settings)
	if err != nil {
		return nil, err
	}
	svc := &Svc{
		b:   b,
		log: log,
		app: a,
	}

	svc.b.Handle("/start", svc.start)
	svc.b.Handle(&btnYes, svc.btnYes)
	svc.b.Handle(&btnNo, svc.btnNo)

	return svc, nil
}

func (s *Svc) start(m *tb.Message) {
	userID := int64(m.Sender.ID)

	s.b.Send(m.Chat, "Welcome")

	isAuth := s.app.Auth(userID)
	if !isAuth {
		s.b.Send(m.Chat, fmt.Sprintf("Set ENV variable *TG_USER_ID* to _%d_", userID))
		return
	}
	_, isSheet, isSecret := s.app.CheckSettings()
	if !isSheet {
		s.b.Send(m.Chat, "ENV variables *SHEET_SECRET* or *SHEET_ID* are not set")
	}
	if !isSecret {
		s.b.Send(m.Chat, "ENV variable *API_SECRET* is not set (it's necessary if you use serverless functions)")
	}
}

func (s *Svc) btnYes(c *tb.Callback) {
	if !s.app.Auth(int64(c.Sender.ID)) {
		return
	}

	category, day, err := btnDataParse(c.Data)
	if err != nil {
		return
	}
	err = s.app.SetValue(day, category, true)
	if err != nil {
		return
	}
	s.log.Debug("btnYes", "category", category, "day", day)
	text := "DONE " + category + " at " + day.Format(app.DayLayout)
	s.b.Respond(c, &tb.CallbackResponse{})
	s.b.Edit(c.Message, text)
}

func (s *Svc) btnNo(c *tb.Callback) {
	if !s.app.Auth(int64(c.Sender.ID)) {
		return
	}

	category, day, err := btnDataParse(c.Data)
	if err != nil {
		return
	}
	err = s.app.SetValue(day, category, false)
	if err != nil {
		return
	}
	s.log.Debug("btnNo", "category", category, "day", day)
	text := "FAIL " + category + " at " + day.Format(app.DayLayout)
	s.b.Respond(c, &tb.CallbackResponse{})
	s.b.Edit(c.Message, text)
}

func (s *Svc) Ask(user int64, category string, day time.Time) {
	men := &tb.ReplyMarkup{}
	btnYes.Data = btnData(category, day)
	btnNo.Data = btnData(category, day)

	men.Inline(
		men.Row(btnYes, btnNo),
	)
	text := "Have you done " + category + " at " + day.Format(app.DayLayout)
	_, err := s.b.Send(tb.ChatID(user), text, men)
	if err != nil {
		s.log.Warn("Ask", "err", err)
	}
}

func (s *Svc) Start() {
	s.b.Start()
}

func (s *Svc) ProcessUpdate(u tb.Update) {
	s.b.ProcessUpdate(u)
}

func btnData(category string, day time.Time) string {
	return category + "|" + day.Format(time.RFC3339)
}

func btnDataParse(data string) (string, time.Time, error) {
	words := strings.Split(data, "|")
	t, err := time.Parse(time.RFC3339, words[1])
	if err != nil {
		return "", time.Time{}, err
	}
	return words[0], t, nil
}

func SetHookDoamin(botToken, url string) error {
  u := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook?url=https://%s/api/updater", botToken, url)
	_, err := http.Get(u)
	if err != nil {
		return err
	}
	return nil
}
