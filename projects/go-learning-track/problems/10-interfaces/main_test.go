package main

import (
	"errors"
	"slices"
	"strings"
	"testing"
)

// Compile-time checks: these lines fail to build if a type stops satisfying
// Notifier. Nothing in the types themselves says "implements Notifier".
var (
	_ Notifier = (*EmailNotifier)(nil)
	_ Notifier = (*SMSNotifier)(nil)
	_ Notifier = (*LoggingNotifier)(nil)
)

func TestEmailNotify(t *testing.T) {
	e := &EmailNotifier{To: "chris@example.com"}

	if err := e.Notify("hello"); err != nil {
		t.Fatalf("Notify error = %v, want nil", err)
	}
	if want := []string{"hello"}; !slices.Equal(e.Sent, want) {
		t.Errorf("Sent = %v, want %v", e.Sent, want)
	}
}

func TestEmailInvalidAddress(t *testing.T) {
	e := &EmailNotifier{To: "chris.example.com"}

	if err := e.Notify("hello"); !errors.Is(err, ErrInvalidAddress) {
		t.Errorf("Notify error = %v, want ErrInvalidAddress", err)
	}
	if len(e.Sent) != 0 {
		t.Errorf("Sent = %v, want nothing recorded on failure", e.Sent)
	}
}

func TestSMSNotify(t *testing.T) {
	s := &SMSNotifier{To: "+15555550100"}

	if err := s.Notify("hello"); err != nil {
		t.Fatalf("Notify error = %v, want nil", err)
	}
	if err := s.Notify(strings.Repeat("x", SMSMaxLen)); err != nil {
		t.Fatalf("Notify(%d bytes) error = %v, want nil (the limit is inclusive)", SMSMaxLen, err)
	}
	if len(s.Sent) != 2 || s.Sent[0] != "hello" {
		t.Errorf("Sent = %v, want 2 messages starting with %q", s.Sent, "hello")
	}
}

func TestSMSInvalidAddress(t *testing.T) {
	s := &SMSNotifier{To: "5555550100"}

	if err := s.Notify("hello"); !errors.Is(err, ErrInvalidAddress) {
		t.Errorf("Notify error = %v, want ErrInvalidAddress", err)
	}
	if len(s.Sent) != 0 {
		t.Errorf("Sent = %v, want nothing recorded on failure", s.Sent)
	}
}

func TestSMSTooLong(t *testing.T) {
	s := &SMSNotifier{To: "+15555550100"}

	if err := s.Notify(strings.Repeat("x", SMSMaxLen+1)); !errors.Is(err, ErrMessageTooLong) {
		t.Errorf("Notify(%d bytes) error = %v, want ErrMessageTooLong", SMSMaxLen+1, err)
	}
	if len(s.Sent) != 0 {
		t.Errorf("Sent = %v, want nothing recorded on failure", s.Sent)
	}
}

func TestLoggingNotifierLogsAndDelegates(t *testing.T) {
	inner := &EmailNotifier{To: "chris@example.com"}
	l := &LoggingNotifier{Notifier: inner}

	if err := l.Notify("one"); err != nil {
		t.Fatalf("Notify error = %v, want nil", err)
	}
	_ = l.Notify("two")

	if want := []string{"one", "two"}; !slices.Equal(l.Log, want) {
		t.Errorf("Log = %v, want %v", l.Log, want)
	}
	if want := []string{"one", "two"}; !slices.Equal(inner.Sent, want) {
		t.Errorf("inner Sent = %v, want %v (LoggingNotifier must delegate)", inner.Sent, want)
	}
}

func TestLoggingNotifierLogsFailures(t *testing.T) {
	l := &LoggingNotifier{Notifier: &SMSNotifier{To: "no-plus"}}

	if err := l.Notify("one"); !errors.Is(err, ErrInvalidAddress) {
		t.Errorf("Notify error = %v, want the inner ErrInvalidAddress passed through", err)
	}
	if want := []string{"one"}; !slices.Equal(l.Log, want) {
		t.Errorf("Log = %v, want %v (log even when sending fails)", l.Log, want)
	}
}

func TestLoggingNotifierNoInner(t *testing.T) {
	l := &LoggingNotifier{}

	if err := l.Notify("one"); !errors.Is(err, ErrNoNotifier) {
		t.Errorf("Notify error = %v, want ErrNoNotifier", err)
	}
}

func TestNotifyAllSendsToEveryone(t *testing.T) {
	e := &EmailNotifier{To: "chris@example.com"}
	s := &SMSNotifier{To: "+15555550100"}

	if err := NotifyAll([]Notifier{e, s}, "ready"); err != nil {
		t.Fatalf("NotifyAll error = %v, want nil", err)
	}
	if len(e.Sent) != 1 || len(s.Sent) != 1 {
		t.Errorf("email Sent = %v, sms Sent = %v; want one message each", e.Sent, s.Sent)
	}
}

func TestNotifyAllContinuesAfterFailure(t *testing.T) {
	bad := &EmailNotifier{To: "nope"}
	good := &SMSNotifier{To: "+15555550100"}

	err := NotifyAll([]Notifier{bad, good}, "ready")
	if !errors.Is(err, ErrInvalidAddress) {
		t.Errorf("NotifyAll error = %v, want it to include ErrInvalidAddress", err)
	}
	if len(good.Sent) != 1 {
		t.Errorf("second notifier Sent = %v, want 1 message (keep going after a failure)", good.Sent)
	}
}

func TestNotifyAllJoinsErrors(t *testing.T) {
	badEmail := &EmailNotifier{To: "nope"}
	sms := &SMSNotifier{To: "+15555550100"}
	long := strings.Repeat("x", SMSMaxLen+1)

	err := NotifyAll([]Notifier{badEmail, sms}, long)
	if !errors.Is(err, ErrInvalidAddress) {
		t.Errorf("NotifyAll error = %v, want it to include ErrInvalidAddress", err)
	}
	if !errors.Is(err, ErrMessageTooLong) {
		t.Errorf("NotifyAll error = %v, want it to include ErrMessageTooLong", err)
	}
}

func TestNotifyAllEmpty(t *testing.T) {
	if err := NotifyAll(nil, "ready"); err != nil {
		t.Errorf("NotifyAll(nil) = %v, want nil", err)
	}
}

// countingNotifier is declared only in this test file, and never mentions
// Notifier. Having a Notify(string) error method is enough.
type countingNotifier struct {
	calls int
}

func (c *countingNotifier) Notify(string) error {
	c.calls++
	return nil
}

func TestImplicitSatisfaction(t *testing.T) {
	c := &countingNotifier{}

	if err := NotifyAll([]Notifier{c, c}, "ready"); err != nil {
		t.Fatalf("NotifyAll error = %v, want nil", err)
	}
	if c.calls != 2 {
		t.Errorf("calls = %d, want 2", c.calls)
	}
}
