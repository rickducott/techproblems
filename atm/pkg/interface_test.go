package pkg_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rickducott/techproblems/atm/pkg"
	"time"
)

var _ = Describe("Text Interface", func() {

	const (
		id1  = "12345"
		pin1 = "1234"
		id2  = "23456"
		pin2 = "2345"
	)

	var (
		account1, account2 pkg.Account
		atm                pkg.Atm
		ui                 pkg.TextInterface
		done               chan bool

		amount1 = pkg.Dollars(20000)
		amount2 = pkg.Dollars(25)
	)

	BeforeEach(func() {
		account1 = pkg.NewAccount(id1, pin1, amount1)
		account2 = pkg.NewAccount(id2, pin2, amount2)
		atm, done = pkg.NewAtm(1, account1, account2)
		ui = pkg.NewInterface(atm)
	})

	AfterEach(func() {
		// prevent goroutine leak
		done <- true
	})

	It("handles unknown command", func() {
		msg := ui.Execute("foo")
		Expect(msg).To(Equal(pkg.HelpMessage))
	})

	It("handles no command", func() {
		msg := ui.Execute("")
		Expect(msg).To(Equal(pkg.HelpMessage))
	})

	It("handles not authorized for command", func() {
		msg := ui.Execute("balance")
		Expect(msg).To(Equal(pkg.AuthorizationRequiredError.Error()))
	})

	It("handles invalid authorization", func() {
		msg := ui.Execute("authorize foo bar")
		Expect(msg).To(Equal(pkg.AuthorizationFailedError.Error()))
	})

	It("handles valid authorization", func() {
		msg := ui.Execute(fmt.Sprintf("authorize %s %s", id1, pin1))
		Expect(msg).To(Equal(pkg.AuthorizedMessage(id1)))
	})

	It("handles valid authorization", func() {
		_ = ui.Execute(fmt.Sprintf("authorize %s %s", id1, pin1))
		msg := ui.Execute("balance")
		Expect(msg).To(Equal(pkg.BalanceMessage(amount1)))
	})

	It("handles deposit", func() {
		_ = ui.Execute(fmt.Sprintf("authorize %s %s", id1, pin1))
		msg := ui.Execute("deposit 500")
		Expect(msg).To(Equal(pkg.BalanceMessage(amount1.Add(pkg.Dollars(500)))))
	})

	It("handles withdraw", func() {
		_ = ui.Execute(fmt.Sprintf("authorize %s %s", id1, pin1))
		msg := ui.Execute("withdraw 500")
		withdrawAmt := pkg.Dollars(500)
		txn := pkg.Transaction{
			Date:      time.Now(),
			Amount:    withdrawAmt.Negative(),
			Balance:   amount1.Subtract(withdrawAmt),
			Overdraft: false,
		}
		Expect(msg).To(Equal(pkg.WithdrawMessage(withdrawAmt, &txn)))
	})

	It("handles withdraw overdraft", func() {
		_ = ui.Execute(fmt.Sprintf("authorize %s %s", id2, pin2))
		msg := ui.Execute("withdraw 40")
		withdrawAmt := pkg.Dollars(40)
		txn := pkg.Transaction{
			Date:      time.Now(),
			Amount:    withdrawAmt.Negative(),
			Balance:   amount2.Subtract(withdrawAmt).Subtract(pkg.OverdraftFee),
			Overdraft: true,
		}
		Expect(msg).To(Equal(pkg.WithdrawMessage(withdrawAmt, &txn)))
	})

	It("handles withdraw run out of money", func() {
		_ = ui.Execute(fmt.Sprintf("authorize %s %s", id1, pin1))
		msg := ui.Execute("withdraw 20000")
		desiredAmt := pkg.Dollars(20000)
		withdrawAmt := pkg.Dollars(10000)
		txn := pkg.Transaction{
			Date:      time.Now(),
			Amount:    withdrawAmt.Negative(),
			Balance:   amount1.Subtract(withdrawAmt),
			Overdraft: false,
		}
		Expect(msg).To(Equal(pkg.WithdrawMessage(desiredAmt, &txn)))
	})

	It("handles run out of money", func() {
		_ = ui.Execute(fmt.Sprintf("authorize %s %s", id1, pin1))
		_ = ui.Execute("withdraw 20000")
		msg := ui.Execute("withdraw 20")
		Expect(msg).To(Equal(pkg.NoMoneyError.Error()))
	})

})
