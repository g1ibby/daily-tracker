package sheet

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/g1ibby/daily-tracker/internal/app"
	"github.com/powerman/structlog"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

var (
	log = structlog.New()
)

type GSheet struct {
	service *spreadsheet.Service
	id      string
}

func NewGSheet(secret string, spreadSheetID string) (*GSheet, error) {
	ctx := context.Background()

	conf, err := google.JWTConfigFromJSON([]byte(secret), spreadsheet.Scope)
	if err != nil {
		return nil, err
	}
	client := conf.Client(ctx)

	service := spreadsheet.NewServiceWithClient(client)

	svc := &GSheet{
		service,
		spreadSheetID,
	}
	return svc, nil
}

func (s *GSheet) getSheet() (*spreadsheet.Sheet, error) {
	ss, err := s.service.FetchSpreadsheet(s.id)
	if err != nil {
		return nil, err
	}
	sheet, err := ss.SheetByIndex(0)
	if err != nil {
		return nil, err
	}
	return sheet, nil
}

func (s *GSheet) Columns() ([]string, error) {
	sheet, err := s.getSheet()
	if err != nil {
		return []string{}, err
	}

	var columns []string
	for _, cl := range sheet.Rows[0][1:] {
		if cl.Value != "" {
			columns = append(columns, cl.Value)
		}
	}
	return columns, nil
}

func (s *GSheet) GetLastDay() (time.Time, error) {
	sheet, err := s.getSheet()
	if err != nil {
		return time.Time{}, err
	}

	rwIndex := uint(0)
	for _, rw := range sheet.Rows {
		if rw[0].Value == "" {
			break
		}
		rwIndex = rw[0].Row
	}

	day, err := time.Parse(app.DayLayout, sheet.Rows[rwIndex][0].Value)
	if err != nil {
		return time.Time{}, err
	}
	return day, nil
}

func (s *GSheet) AddNewDay(day time.Time) error {
	sheet, err := s.getSheet()
	if err != nil {
		return err
	}
	rwIndex := uint(0)
	for _, rw := range sheet.Rows {
		if rw[0].Value == "" {
			break
		}
		rwIndex = rw[0].Row
	}
	if sheet.Rows[rwIndex][0].Value == day.Format(app.DayLayout) {
		return nil
	}
	sheet.Update(int(rwIndex+1), 0, day.Format(app.DayLayout))
	return sheet.Synchronize()
}

func (s *GSheet) SetValue(day time.Time, column string, value bool) error {
	sheet, err := s.getSheet()
	if err != nil {
		return err
	}
	rwIndex := uint(0)
	for _, rw := range sheet.Rows {
		if rw[0].Value == day.Format(app.DayLayout) {
			rwIndex = rw[0].Row
			continue
		}
	}
	clIndex := uint(0)
	for _, cl := range sheet.Rows[0] {
		if cl.Value == column {
			clIndex = cl.Column
			continue
		}
	}
	if rwIndex == 0 || clIndex == 0 {
		return nil
	}

	sheet.Update(int(rwIndex), int(clIndex), strings.ToUpper(strconv.FormatBool(value)))
	return sheet.Synchronize()
}
