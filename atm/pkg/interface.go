package pkg

import (
	"fmt"
	"strings"
)

const (
	HelpMessage          = "Must provide command: authorize, withdraw, deposit, balance, history, logout, or end"
	HelpAuthorizeMessage = "Authorize command requires two arguments: <id> <pin>"
	HelpWithdrawMessage  = "Withdraw command requires one argument: <value>"
	HelpDepositMessage   = "Deposit command requires one argument: <value>"
)

var (
	AuthorizedMessage = func(id string) string {
		return fmt.Sprintf("%s successfully authorized.", id)
	}

	BalanceMessage = func(amount Amount) string {
		return fmt.Sprintf("Current balance: %v", amount)
	}

	WithdrawMessage = func(desiredAmt Amount, txn *Transaction) string {
		msg := ""
		if desiredAmt.GreaterThan(txn.Amount.Abs()) {
			msg += fmt.Sprintf("Unable to dispense full amount requested at this time. ")
		}
		msg += fmt.Sprintf("Amount dispensed: $%v\n", txn.Amount.Abs())
		if txn.Overdraft {
			msg += "You have been charged an overdraft fee of $5. "
		}
		msg += BalanceMessage(txn.Balance)
		return msg
	}

	HistoryMessage = func(history []Transaction) string {
		msg := ""
		for i := len(history) - 1; i >= 0; i-- {
			msg += fmt.Sprintf("%v", history[i])
			if i > 0 {
				msg += "\n"
			}
		}
		return msg
	}

	LogoutMessage = func(accountId string) string {
		return fmt.Sprintf("Account %s logged out.", accountId)
	}
)

func NewInterface(atm Atm) TextInterface {
	return &textInterface{atm: atm}
}

type TextInterface interface {
	Execute(command string) string
}

type textInterface struct {
	atm Atm
}

func (t *textInterface) Execute(command string) string {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return HelpMessage
	}

	switch fields[0] {
	case "authorize":
		if len(fields) != 3 {
			return HelpAuthorizeMessage
		} else {
			if err := t.atm.Authorize(fields[1], fields[2]); err != nil {
				return err.Error()
			} else {
				return AuthorizedMessage(fields[1])
			}
		}
	case "withdraw":
		if len(fields) != 2 {
			return HelpWithdrawMessage
		}
		amount, err := ParseAmount(fields[1])
		if err != nil {
			return err.Error()
		} else {
			txn, err := t.atm.Withdraw(amount)
			if err != nil {
				return err.Error()
			} else {
				return WithdrawMessage(amount, txn)
			}
		}
	case "deposit":
		if len(fields) != 2 {
			return HelpDepositMessage
		}
		amount, err := ParseAmount(fields[1])
		if err != nil {
			return err.Error()
		} else {
			err := t.atm.Deposit(amount)
			if err != nil {
				return err.Error()
			} else {
				return t.balance()
			}
		}
	case "balance":
		return t.balance()
	case "history":
		history, err := t.atm.History()
		if err != nil {
			return err.Error()
		} else {
			return HistoryMessage(history)
		}
	case "logout":
		accountId, err := t.atm.Logout()
		if err != nil {
			return err.Error()
		} else {
			return LogoutMessage(accountId)
		}
	}
	return HelpMessage
}

func (t *textInterface) balance() string {
	balance, err := t.atm.Balance()
	if err != nil {
		return err.Error()
	} else {
		return BalanceMessage(balance)
	}
}
