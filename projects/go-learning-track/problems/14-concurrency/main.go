package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrNoOffers means every lender failed.
var ErrNoOffers = errors.New("no lender returned an offer")

// Offer is one lender's rate quote.
type Offer struct {
	Lender string
	Rate   float64
}

// Lender quotes a rate. Implementations must stop work and return ctx.Err()
// promptly once ctx is done.
type Lender interface {
	Quote(ctx context.Context, amount float64) (Offer, error)
}

// FakeLender simulates a slow remote lender.
type FakeLender struct {
	Name  string
	Rate  float64
	Delay time.Duration // how long the "network call" takes
	Err   error         // returned after Delay instead of an offer, when set
}

// Quote waits Delay, then returns Offer{Name, Rate} (or Err when set). If ctx
// is done before Delay has passed, it returns ctx.Err() immediately.
func (f *FakeLender) Quote(ctx context.Context, amount float64) (Offer, error) {
	// TODO
	return Offer{}, nil
}

// AllQuotes asks every lender at the same time and waits for all of them.
// offers holds the successful offers in the same order as lenders. err joins
// every failure, or is nil when none failed.
func AllQuotes(ctx context.Context, lenders []Lender, amount float64) (offers []Offer, err error) {
	// TODO
	return nil, nil
}

// FastestQuote asks every lender at the same time and returns the first
// successful offer, cancelling the lenders still working. When every lender
// fails, the error matches ErrNoOffers and each lender's error. When ctx ends
// first, the error matches ctx.Err().
func FastestQuote(ctx context.Context, lenders []Lender, amount float64) (Offer, error) {
	// TODO
	return Offer{}, nil
}

// QuoteWithin is FastestQuote with a deadline of timeout from now.
func QuoteWithin(ctx context.Context, lenders []Lender, amount float64, timeout time.Duration) (Offer, error) {
	// TODO
	return Offer{}, nil
}

func main() {
	lenders := []Lender{
		&FakeLender{Name: "Northwind", Rate: 0.0640, Delay: 300 * time.Millisecond},
		&FakeLender{Name: "Contoso", Rate: 0.0615, Delay: 120 * time.Millisecond},
		&FakeLender{Name: "Fabrikam", Rate: 0.0599, Delay: 2 * time.Second},
	}
	ctx := context.Background()

	start := time.Now()
	offers, err := AllQuotes(ctx, lenders, 250000)
	fmt.Printf("all quotes after %v: %v (err: %v)\n", time.Since(start).Round(time.Millisecond), offers, err)

	start = time.Now()
	offer, err := QuoteWithin(ctx, lenders, 250000, 500*time.Millisecond)
	fmt.Printf("fastest after %v: %+v (err: %v)\n", time.Since(start).Round(time.Millisecond), offer, err)

	start = time.Now()
	offer, err = QuoteWithin(ctx, lenders, 250000, 50*time.Millisecond)
	fmt.Printf("50ms budget after %v: %+v (err: %v)\n", time.Since(start).Round(time.Millisecond), offer, err)
}
