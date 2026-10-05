package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidAddress = errors.New("invalid address")
	ErrMessageTooLong = errors.New("message too long")
	ErrNoNotifier     = errors.New("no notifier configured")
)

// SMSMaxLen is the longest message, in bytes, an SMSNotifier accepts.
const SMSMaxLen = 160

// Notifier sends a message to whoever the implementation is configured for.
type Notifier interface {
	Notify(msg string) error
}

// EmailNotifier is a fake: it records messages in Sent instead of emailing.
type EmailNotifier struct {
	To   string
	Sent []string
}

// Notify records msg. To must contain "@", or it returns ErrInvalidAddress.
func (e *EmailNotifier) Notify(msg string) error {
	if !strings.Contains(e.To, "@") {
		return ErrInvalidAddress
	}

	e.Sent = append(e.Sent, msg)
	return nil
}

// SMSNotifier is a fake: it records messages in Sent instead of texting.
type SMSNotifier struct {
	To   string
	Sent []string
}

// Notify records msg. To must start with "+" (ErrInvalidAddress), and msg
// must be at most SMSMaxLen bytes (ErrMessageTooLong).
func (s *SMSNotifier) Notify(msg string) error {
	if !strings.Contains(s.To, "+") {
		return ErrInvalidAddress
	}
	if len(msg) > SMSMaxLen {
		return ErrMessageTooLong
	}

	s.Sent = append(s.Sent, msg)
	return nil
}

// LoggingNotifier wraps any Notifier and records every message it is asked
// to send, whether or not the send succeeds.
type LoggingNotifier struct {
	Notifier          // embedded: its Notify is promoted unless we define our own
	Log      []string // every message, in order
}

// Notify records msg in Log, then delegates to the embedded Notifier. With no
// embedded Notifier it returns ErrNoNotifier instead of panicking.
func (l *LoggingNotifier) Notify(msg string) error {
	l.Log = append(l.Log, msg)

	if l.Notifier == nil {
		return ErrNoNotifier
	}
	if err := l.Notifier.Notify(msg); err != nil {
		return err
	}
	return nil
}

// NotifyAll sends msg through every notifier, even after a failure, and
// returns all the failures joined into one error (nil if none failed).
func NotifyAll(notifiers []Notifier, msg string) error {
	notifyErrs := make([]error, 0, len(notifiers))

	for _, notifier := range notifiers {
		if err := notifier.Notify(msg); err != nil {
			notifyErrs = append(notifyErrs, err)
		}
	}

	if len(notifyErrs) > 0 {
		return errors.Join(notifyErrs...)
	}

	return nil
}

func main() {
	sms := &SMSNotifier{To: "+15555550100"}
	logged := &LoggingNotifier{Notifier: sms} // embedded the SMSNotifier, effectively inheriting its implementation
	// logged := &LoggingNotifier{} // will cause a nil-ref panic if (l *LoggingNotifier) Notify(msg string) is not present
	email := &EmailNotifier{To: "chris@example.com"}
	broken := &EmailNotifier{To: "not-an-address"}

	err := NotifyAll([]Notifier{email, logged, broken}, "Appraisal scheduled for loan L-100")

	fmt.Println("error:     ", err)
	fmt.Println("email sent:", email.Sent)
	fmt.Println("sms sent:  ", sms.Sent)
	fmt.Println("log:       ", logged.Log)
}
