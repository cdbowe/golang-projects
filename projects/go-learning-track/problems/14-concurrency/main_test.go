package main

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"
)

// Every test runs inside synctest.Test: time is fake, so a 2-second delay
// costs nothing, and elapsed times are exact. If the test function returns
// while a goroutine it started is still blocked, synctest fails with
// "deadlock: main bubble goroutine has exited but blocked goroutines remain" —
// that means a goroutine leaked.

var (
	errDown    = errors.New("lender down")
	errTimeout = errors.New("lender timed out")
)

func lender(name string, rate float64, delay time.Duration) *FakeLender {
	return &FakeLender{Name: name, Rate: rate, Delay: delay}
}

func failing(name string, delay time.Duration, err error) *FakeLender {
	return &FakeLender{Name: name, Delay: delay, Err: err}
}

func assertElapsed(t *testing.T, start time.Time, want time.Duration) {
	t.Helper()
	if got := time.Since(start); got != want {
		t.Errorf("took %v, want exactly %v", got, want)
	}
}

func TestFakeLender(t *testing.T) {
	tests := []struct {
		name      string
		lender    *FakeLender
		ctxLimit  time.Duration // 0 = no deadline
		want      Offer
		wantErr   error
		wantAfter time.Duration
	}{
		{"offer after delay", lender("A", 0.06, 50*time.Millisecond), 0, Offer{"A", 0.06}, nil, 50 * time.Millisecond},
		{"error after delay", failing("B", 30*time.Millisecond, errDown), 0, Offer{}, errDown, 30 * time.Millisecond},
		{"deadline first", lender("C", 0.05, time.Second), 10 * time.Millisecond, Offer{}, context.DeadlineExceeded, 10 * time.Millisecond},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx := t.Context()
				if tc.ctxLimit > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, tc.ctxLimit)
					defer cancel()
				}
				start := time.Now()

				got, err := tc.lender.Quote(ctx, 1000)

				if !errors.Is(err, tc.wantErr) {
					t.Errorf("err = %v, want %v", err, tc.wantErr)
				}
				if got != tc.want {
					t.Errorf("offer = %+v, want %+v", got, tc.want)
				}
				assertElapsed(t, start, tc.wantAfter)
			})
		})
	}
}

func TestFakeLenderAlreadyCancelled(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		start := time.Now()

		_, err := lender("A", 0.06, time.Second).Quote(ctx, 1000)

		if !errors.Is(err, context.Canceled) {
			t.Errorf("err = %v, want context.Canceled", err)
		}
		assertElapsed(t, start, 0)
	})
}

func TestAllQuotesRunsConcurrently(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		lenders := []Lender{
			lender("A", 0.061, 30*time.Millisecond),
			lender("B", 0.062, 10*time.Millisecond),
			lender("C", 0.063, 20*time.Millisecond),
		}
		start := time.Now()

		offers, err := AllQuotes(t.Context(), lenders, 1000)

		if err != nil {
			t.Fatalf("err = %v, want nil", err)
		}
		want := []Offer{{"A", 0.061}, {"B", 0.062}, {"C", 0.063}}
		if len(offers) != len(want) {
			t.Fatalf("offers = %+v, want %+v", offers, want)
		}
		for i := range want {
			if offers[i] != want[i] {
				t.Errorf("offers[%d] = %+v, want %+v (keep the lenders' order, not arrival order)", i, offers[i], want[i])
			}
		}
		assertElapsed(t, start, 30*time.Millisecond) // one after another would take 60ms
	})
}

func TestAllQuotesPartialFailure(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		lenders := []Lender{
			failing("A", 10*time.Millisecond, errDown),
			lender("B", 0.062, 20*time.Millisecond),
			failing("C", 5*time.Millisecond, errTimeout),
		}

		offers, err := AllQuotes(t.Context(), lenders, 1000)

		if len(offers) != 1 || offers[0] != (Offer{"B", 0.062}) {
			t.Errorf("offers = %+v, want only B's", offers)
		}
		if !errors.Is(err, errDown) || !errors.Is(err, errTimeout) {
			t.Errorf("err = %v, want it to match both errDown and errTimeout", err)
		}
	})
}

func TestAllQuotesEmpty(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		offers, err := AllQuotes(t.Context(), nil, 1000)
		if len(offers) != 0 || err != nil {
			t.Errorf("AllQuotes(nil) = (%v, %v), want (empty, nil)", offers, err)
		}
	})
}

func TestFastestQuote(t *testing.T) {
	tests := []struct {
		name      string
		lenders   []Lender
		want      Offer
		wantAfter time.Duration
	}{
		{
			name:      "fastest wins, not lowest rate",
			lenders:   []Lender{lender("Slow", 0.05, 30*time.Millisecond), lender("Fast", 0.07, 10*time.Millisecond)},
			want:      Offer{"Fast", 0.07},
			wantAfter: 10 * time.Millisecond,
		},
		{
			name:      "failures are skipped",
			lenders:   []Lender{failing("Broken", 5*time.Millisecond, errDown), lender("OK", 0.06, 20*time.Millisecond)},
			want:      Offer{"OK", 0.06},
			wantAfter: 20 * time.Millisecond,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				start := time.Now()

				got, err := FastestQuote(t.Context(), tc.lenders, 1000)

				if err != nil {
					t.Fatalf("err = %v, want nil", err)
				}
				if got != tc.want {
					t.Errorf("offer = %+v, want %+v", got, tc.want)
				}
				assertElapsed(t, start, tc.wantAfter)
			})
		})
	}
}

// spyLender reports how its inner lender's Quote ended.
type spyLender struct {
	Lender
	done chan error
}

func (s *spyLender) Quote(ctx context.Context, amount float64) (Offer, error) {
	o, err := s.Lender.Quote(ctx, amount)
	s.done <- err
	return o, err
}

func TestFastestQuoteCancelsTheRest(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		slow := &spyLender{Lender: lender("Slow", 0.05, 2*time.Second), done: make(chan error, 1)}
		lenders := []Lender{lender("Fast", 0.06, 10*time.Millisecond), slow}
		start := time.Now()

		if _, err := FastestQuote(t.Context(), lenders, 1000); err != nil {
			t.Fatalf("err = %v, want nil", err)
		}
		synctest.Wait() // let every goroutine run until it finishes or blocks

		select {
		case err := <-slow.done:
			if !errors.Is(err, context.Canceled) {
				t.Errorf("slow lender ended with %v, want context.Canceled", err)
			}
		default:
			t.Error("slow lender was still waiting after the winner returned; cancel the others")
		}
		assertElapsed(t, start, 10*time.Millisecond)
	})
}

func TestFastestQuoteAllFail(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		lenders := []Lender{failing("A", 10*time.Millisecond, errDown), failing("B", 20*time.Millisecond, errTimeout)}

		_, err := FastestQuote(t.Context(), lenders, 1000)

		for _, want := range []error{ErrNoOffers, errDown, errTimeout} {
			if !errors.Is(err, want) {
				t.Errorf("err = %v, want it to match %v", err, want)
			}
		}
	})
}

func TestFastestQuoteNoLenders(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		if _, err := FastestQuote(t.Context(), nil, 1000); !errors.Is(err, ErrNoOffers) {
			t.Errorf("err = %v, want ErrNoOffers", err)
		}
	})
}

func TestQuoteWithin(t *testing.T) {
	lenders := func() []Lender {
		return []Lender{lender("A", 0.06, 10*time.Millisecond), lender("B", 0.05, 40*time.Millisecond)}
	}

	t.Run("answer inside the budget", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			got, err := QuoteWithin(t.Context(), lenders(), 1000, 15*time.Millisecond)
			if err != nil || got != (Offer{"A", 0.06}) {
				t.Errorf("got (%+v, %v), want A's offer", got, err)
			}
		})
	})

	t.Run("budget runs out", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			start := time.Now()
			_, err := QuoteWithin(t.Context(), lenders(), 1000, 5*time.Millisecond)
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Errorf("err = %v, want context.DeadlineExceeded", err)
			}
			assertElapsed(t, start, 5*time.Millisecond)
		})
	})
}
