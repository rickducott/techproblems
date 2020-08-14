package pkg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	CentsPerDollar = 100
	ZeroAmount     = Amount{}

	AmountParseError = func(msg string) error {
		return errors.New("Error parsing amount: " + msg)
	}
)

// Create a new amount from a combination of dollars and cents.
// For negative amounts, make sure to use the negative sign on both the dollars and cents value
// (otherwise you may get unintuitive behavior, for example, `NewAmount(-1, 15) == NewAmount(0, -85)`).
func NewAmount(dollars, cents int) Amount {
	return Amount{
		cents: dollars*CentsPerDollar + cents,
	}
}

func Cents(cents int) Amount {
	return NewAmount(0, cents)
}

func Dollars(dollars int) Amount {
	return NewAmount(dollars, 0)
}

func ParseAmount(amount string) (Amount, error) {
	negative := false
	if strings.HasPrefix(amount, "-") {
		negative = true
		amount = strings.TrimPrefix(amount, "-")
	}

	parts := strings.Split(amount, ".")
	if parts[0] == "" {
		parts[0] = "0"
	}
	dollars, err := strconv.Atoi(parts[0])
	if err != nil {
		return ZeroAmount, AmountParseError("error parsing dollars: " + err.Error())
	}
	amt := Dollars(dollars)
	if len(parts) > 1 {
		if len(parts[1]) > 2 {
			return ZeroAmount, AmountParseError("too many digits in the cents")
		} else if len(parts[1]) <= 1 {
			parts[1] = parts[1] + "0"
		}
		cents, err := strconv.Atoi(parts[1])
		if err != nil {
			return ZeroAmount, AmountParseError("error parsing cents: " + err.Error())
		}
		amt = amt.Add(Cents(cents))
	}
	if negative {
		return amt.Negative(), nil
	}
	return amt, nil
}

type Amount struct {
	cents int
}

func (a Amount) String() string {
	dollars := a.cents / CentsPerDollar
	cents := a.cents % CentsPerDollar
	if cents < 0 {
		cents = cents * -1
	}
	centsStr := strconv.Itoa(cents)
	if cents < 10 {
		centsStr = "0" + centsStr
	}
	return fmt.Sprintf("%d.%s", dollars, centsStr)
}

func (a Amount) Add(amount Amount) Amount {
	return Amount{
		cents: a.cents + amount.cents,
	}
}

func (a Amount) Subtract(amount Amount) Amount {
	return Amount{
		cents: a.cents - amount.cents,
	}
}

func (a Amount) GreaterThan(other Amount) bool {
	return a.cents > other.cents
}

func (a Amount) Negative() Amount {
	return Amount{
		cents: a.cents * -1,
	}
}

func (a Amount) Abs() Amount {
	if ZeroAmount.GreaterThan(a) {
		return a.Negative()
	} else {
		return a
	}
}

func (a Amount) MultipleOf(other Amount) bool {
	return a.cents%other.cents == 0
}
