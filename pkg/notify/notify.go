package notify

import (
	"context"
	"errors"
	"net/smtp"

	"github.com/sirupsen/logrus"

	"github.com/MninaTB/vacadm/pkg/database"
)

type Notifier interface {
	NotifyUser(ctx context.Context, userID, action string) error
	NotifyTeam(ctx context.Context, teamID, action string) error
}

var _ Notifier = (*NoopNotifier)(nil)
var _ Notifier = (*Mailer)(nil)

var (
	ErrEmptyTeam = errors.New("team has no member")
)

func NewNoopNotifier() *NoopNotifier {
	return &NoopNotifier{
		logger: logrus.New().WithFields(logrus.Fields{
			"component": "noop-notifier",
		}),
	}
}

type NoopNotifier struct {
	logger logrus.FieldLogger
}

func (m *NoopNotifier) NotifyUser(ctx context.Context, userID, action string) error {
	m.logger.WithFields(logrus.Fields{
		"notify-user": userID,
		"action":      action,
	}).Info("inform user")
	return nil
}

func (m *NoopNotifier) NotifyTeam(ctx context.Context, teamID, action string) error {
	m.logger.WithFields(logrus.Fields{
		"notify-team": teamID,
		"action":      action,
	}).Info("inform team")
	return nil
}

type Mailer struct {
	address string
	auth    smtp.Auth
	from    string
	db      database.Database
	logger  logrus.FieldLogger
}

func NewMailer(smtpHost, smtpPort, user, password string, db database.Database) *Mailer {
	address := smtpHost + ":" + smtpPort
	return &Mailer{
		logger: logrus.New().WithFields(logrus.Fields{
			"component": "mailer",
			"address":   address,
		}),
		address: address,
		auth:    smtp.PlainAuth("", user, password, smtpHost),
		from:    user,
		db:      db,
	}
}

func (m *Mailer) NotifyUser(ctx context.Context, userID, action string) error {
	m.logger.WithFields(logrus.Fields{
		"notify-user": userID,
		"action":      action,
	}).Info("inform user")
	usr, err := m.db.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	message := []byte(action)
	return smtp.SendMail(m.address, m.auth, m.from, []string{usr.Email}, message)
}

func (m *Mailer) NotifyTeam(ctx context.Context, teamID, action string) error {
	m.logger.WithFields(logrus.Fields{
		"notify-team": teamID,
		"action":      action,
	}).Info("inform team")

	users, err := m.db.ListTeamUsers(ctx, teamID)
	if err != nil {
		return err
	}

	to := []string{}
	for _, u := range users {
		to = append(to, u.Email)
	}

	if len(to) == 0 {
		return ErrEmptyTeam
	}

	message := []byte(action)
	return smtp.SendMail(m.address, m.auth, m.from, to, message)
}
