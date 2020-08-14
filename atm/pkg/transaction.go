package pkg

import (
	"fmt"
	"time"
)

type Transaction struct {
	Date      time.Time
	Amount    Amount
	Balance   Amount
	Overdraft bool
}

func NewTransaction(amount, balance Amount) Transaction {
	overdraft := ZeroAmount.GreaterThan(amount) && ZeroAmount.GreaterThan(balance)
	return Transaction{
		Date:      time.Now(),
		Amount:    amount,
		Balance:   balance,
		Overdraft: overdraft,
	}
}

func (t Transaction) String() string {
	return fmt.Sprintf("%v %v %v", t.Date.Format("2006-01-02 15:04:05"), t.Amount, t.Balance)
}
