package app

import (
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

const (
  DayLayout = "02/01/2006"
)

type Repo interface {
	Columns() ([]string, error)
	GetLastDay() (time.Time, error)
	AddNewDay(day time.Time) error
	SetValue(day time.Time, column string, value bool) error
}

type Telegram interface {
	Start()
	ProcessUpdate(u tb.Update)
	Ask(user int64, category string, day time.Time)
}

type Appl interface {
	Auth(userID int64) bool
	CheckSettings() (isUserID, isSheet, isSecret bool)
	SetValue(day time.Time, category string, vval bool) error
	Shedule(tgAsk func(user int64, category string, day time.Time)) error
}

type App struct {
	repo      Repo
	userID    int64
	apiSecret string
}

func New(repo Repo, userID int64, apiSecret string) *App {
	return &App{
		repo:      repo,
		userID:    userID,
		apiSecret: apiSecret,
	}
}

func (a *App) Auth(userID int64) bool {
	return a.userID == userID && a.userID != 0
}

func (a *App) CheckSettings() (isUserID, isSheet, isSecret bool) {
	isUserID = a.userID != 0
	isSheet = a.repo != nil
	isSecret = a.apiSecret != ""
	return
}

func (a *App) SetValue(day time.Time, category string, val bool) error {
	return a.repo.SetValue(day, category, val)
}

func (a *App) Shedule(tgAsk func(user int64, category string, day time.Time)) error {
  isUserID, isSheet, isSecret := a.CheckSettings()
  if !isUserID || !isSheet || !isSecret {
    return nil
  }

	err := a.repo.AddNewDay(time.Now())
	if err != nil {
		return err
	}
	clms, err := a.repo.Columns()
	if err != nil {
		return err
	}
	for _, cl := range clms {
		tgAsk(a.userID, cl, time.Now())
	}
	return nil
}
