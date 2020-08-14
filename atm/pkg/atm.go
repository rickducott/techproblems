package pkg

import (
	"errors"
	"sync"
	"time"
)

var (
	AuthorizationFailedError   = errors.New("Authorization failed.")
	AuthorizationRequiredError = errors.New("Authorization required.")
	NotAuthorizedError         = errors.New("No account currently authorized.")
	InvalidAmountError         = errors.New("Invalid amount.")
	NoMoneyError               = errors.New("Unable to process your withdrawal at this time.")
)

type Atm interface {
	Authorize(id, pin string) error
	Withdraw(amount Amount) (*Transaction, error)
	Deposit(amount Amount) error
	Balance() (Amount, error)
	History() ([]Transaction, error)
	Logout() (string, error)
}

type Session struct {
	AccountId string
	Timer     int
}

func NewAtm(logoutSeconds int, accounts ...Account) (Atm, chan bool) {
	atm := &atm{
		money:    Dollars(10000),
		accounts: Accounts(accounts...),
		mutex:    &sync.Mutex{},
	}
	done := make(chan bool)
	go atm.Start(logoutSeconds, done)
	return atm, done
}

type atm struct {
	money    Amount
	accounts map[string]Account
	session  *Session
	mutex    *sync.Mutex
}

func (a *atm) Start(logoutSeconds int, done chan bool) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			a.mutex.Lock()
			if a.session != nil {
				a.session.Timer += 1
				if a.session.Timer >= logoutSeconds {
					a.session = nil
				}
			}
			a.mutex.Unlock()
		}
	}
}

func (a *atm) Authorize(id, pin string) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	account, ok := a.accounts[id]
	if !ok || !account.Authorize(pin) {
		return AuthorizationFailedError
	}
	a.session = &Session{
		AccountId: id,
		Timer:     0,
	}
	return nil
}

func (a *atm) transaction(amount Amount) (*Transaction, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.session == nil {
		return nil, AuthorizationRequiredError
	}
	account := a.accounts[a.session.AccountId]
	a.session.Timer = 0
	return account.Transaction(amount)
}

func (a *atm) Withdraw(amount Amount) (*Transaction, error) {
	if !amount.GreaterThan(ZeroAmount) || !amount.MultipleOf(Dollars(20)) {
		return nil, InvalidAmountError
	}
	if a.money == ZeroAmount {
		return nil, NoMoneyError
	}
	if amount.GreaterThan(a.money) {
		amount = a.money
	}
	txn, err := a.transaction(amount.Negative())
	if err != nil {
		return nil, err
	} else {
		a.money = a.money.Add(txn.Amount)
		return txn, err
	}
}

func (a *atm) Deposit(amount Amount) error {
	if !amount.GreaterThan(ZeroAmount) {
		return InvalidAmountError
	}
	_, err := a.transaction(amount)
	return err
}

func (a *atm) Balance() (Amount, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.session == nil {
		return ZeroAmount, AuthorizationRequiredError
	}
	account := a.accounts[a.session.AccountId]
	a.session.Timer = 0
	return account.Balance(), nil
}

func (a *atm) History() ([]Transaction, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.session == nil {
		return nil, AuthorizationRequiredError
	}
	account := a.accounts[a.session.AccountId]
	a.session.Timer = 0
	return account.History(), nil
}

func (a *atm) Logout() (string, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.session != nil {
		accountId := a.session.AccountId
		a.session = nil
		return accountId, nil
	} else {
		return "", NotAuthorizedError
	}
}
