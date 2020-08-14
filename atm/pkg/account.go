package pkg

import (
	"errors"
)

var (
	AccountOverdrawnError = errors.New("Your account is overdrawn! You may not make withdrawals at this time.")

	OverdraftFee = Dollars(5)

	_ Account = new(account)
)

type Account interface {
	GetId() string
	Transaction(amount Amount) (*Transaction, error)
	Balance() Amount
	History() []Transaction
	Authorize(pin string) bool
}

func NewAccount(id, pin string, balance Amount) Account {
	return &account{
		id:      id,
		pin:     pin,
		balance: balance,
	}
}

type account struct {
	id           string
	pin          string
	balance      Amount
	transactions []Transaction
}

func (a *account) GetId() string {
	return a.id
}

func (a *account) Transaction(amount Amount) (*Transaction, error) {
	if ZeroAmount.GreaterThan(amount) && ZeroAmount.GreaterThan(a.balance) {
		return nil, AccountOverdrawnError
	}
	a.balance = a.balance.Add(amount)
	transaction := NewTransaction(amount, a.balance)
	if transaction.Overdraft {
		a.balance = a.balance.Subtract(OverdraftFee)
		transaction = NewTransaction(amount, a.balance)
	}
	a.transactions = append(a.transactions, transaction)
	return &transaction, nil
}

func (a *account) Balance() Amount {
	return a.balance
}

func (a *account) History() []Transaction {
	return a.transactions
}

func (a *account) Authorize(pin string) bool {
	return a.pin == pin
}

func Accounts(accounts ...Account) map[string]Account {
	accountMap := make(map[string]Account, len(accounts))
	for _, account := range accounts {
		accountMap[account.GetId()] = account
	}
	return accountMap
}
