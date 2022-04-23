package notify

import (
	"context"
	"errors"
	"fmt"
	"net/smtp"

	"github.com/sirupsen/logrus"

	"github.com/MninaTB/vacadm/pkg/database"
)

// Notifier implements methods to inform a user or team about a ongoing action.
type Notifier interface {
	NotifyUser(ctx context.Context, userID, action string) error
	NotifyTeam(ctx context.Context, teamID, action string) error
}

var _ Notifier = (*NoopNotifier)(nil)
var _ Notifier = (*Mailer)(nil)

// ErrEmptyTeam is returned if a requested team does not contain users.
var ErrEmptyTeam = errors.New("team has no member")

// NewNoopNotifier returns a new NoopNotifier.
func NewNoopNotifier() *NoopNotifier {
	return &NoopNotifier{
		logger: logrus.New().WithFields(logrus.Fields{
			"component": "noop-notifier",
		}),
	}
}

// NoopNotifier does not fulfill any operation. All actions are simply logged to
// console.
type NoopNotifier struct {
	logger logrus.FieldLogger
}

// NotifyUser logs userID and action to console.
func (m *NoopNotifier) NotifyUser(ctx context.Context, userID, action string) error {
	m.logger.WithFields(logrus.Fields{
		"notify-user": userID,
		"action":      action,
	}).Info("inform user")
	return nil
}

// NotifyUser logs teamID and action to console.
func (m *NoopNotifier) NotifyTeam(ctx context.Context, teamID, action string) error {
	m.logger.WithFields(logrus.Fields{
		"notify-team": teamID,
		"action":      action,
	}).Info("inform team")
	return nil
}

// Mailer contains all information to send mails via smtp.
type Mailer struct {
	address string
	auth    smtp.Auth
	from    string
	db      database.Database
	logger  logrus.FieldLogger
}

// NewMailer returns a new Mailer.
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

// NotifyUser sends an e-Mail a user based in the given userID. Content is provided
// by the given action.
func (m *Mailer) NotifyUser(ctx context.Context, userID, action string) error {
	m.logger.WithFields(logrus.Fields{
		"notify-user": userID,
		"action":      action,
	}).Info("inform user")
	usr, err := m.db.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	message := RFC822(m.from, usr.Email, "Vacation Serivce", action)
	return smtp.SendMail(m.address, m.auth, m.from, []string{usr.Email}, []byte(message))
}

// NotifyTeam sends e-Mails a all users in a Team based on the given teamID.
// Content is provided by the given action.
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

	team, err := m.db.GetTeamByID(ctx, teamID)
	if err != nil {
		return err
	}

	// displayed receiver
	displayedReceiver := fmt.Sprintf("team-%s@inform-software.de", team.Name)
	message := RFC822(m.from, displayedReceiver, "Vacation Serivce", action)
	return smtp.SendMail(m.address, m.auth, m.from, to, []byte(message))
}

// RFC822 returns a well formatted mail.
// From: someone@example.com
// To: someone_else@example.com
// Subject: An RFC 822 formatted message
//
// This is the plain text body of the message. Note the blank line
// between the header information and the body of the message.
func RFC822(from, to, subject, body string) string {
	res := fmt.Sprintf("From: %s\n", from)
	res += fmt.Sprintf("To: %s\n", to)
	res += fmt.Sprintf("Subject: %s\n", subject)
	res += fmt.Sprintf("\n%s\n", body)
	return res
}
