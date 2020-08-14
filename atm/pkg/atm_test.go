package pkg_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rickducott/techproblems/atm/pkg"
)

var _ = Describe("Atm", func() {

	const (
		id  = "12345"
		pin = "1234"
	)

	var (
		account pkg.Account
		atm     pkg.Atm
		done    chan bool

		amount = pkg.Dollars(20000)
	)

	BeforeEach(func() {
		account = pkg.NewAccount(id, pin, amount)
		atm, done = pkg.NewAtm(1, account)
	})

	AfterEach(func() {
		// prevent goroutine leak
		done <- true
	})

	authorize := func() {
		err := atm.Authorize(id, pin)
		Expect(err).To(BeNil())
	}

	expectBalance := func(expected pkg.Amount, expectedErr error) {
		balance, err := atm.Balance()
		if expectedErr == nil {
			Expect(err).To(BeNil())
			Expect(balance).To(Equal(expected))
		} else {
			Expect(err).To(Equal(expectedErr))
		}
	}

	expectDeposit := func(amount pkg.Amount, expected error) {
		err := atm.Deposit(amount)
		if expected == nil {
			Expect(err).To(BeNil())
		} else {
			Expect(err).To(Equal(expected))
		}
	}

	expectWithdraw := func(amount, expectedAmt pkg.Amount, expectedErr error) {
		withdrawTxn, err := atm.Withdraw(amount)
		if expectedErr == nil {
			Expect(err).To(BeNil())
			Expect(withdrawTxn.Amount).To(Equal(expectedAmt))
		} else {
			Expect(err).To(Equal(expectedErr))
		}
	}

	It("works", func() {
		authorize()
		expectBalance(amount, nil)
		expectDeposit(pkg.Dollars(1000), nil)
		expectWithdraw(pkg.Dollars(2000), pkg.Dollars(-2000), nil)
		expectBalance(pkg.Dollars(19000), nil)
		atm.Logout()
	})

	It("handles atm running out of money", func() {
		authorize()
		expectDeposit(pkg.Dollars(20000), nil)
		expectWithdraw(pkg.Dollars(20000), pkg.Dollars(-10000), nil)
		expectBalance(pkg.Dollars(30000), nil)
		atm.Logout()
	})

	It("properly errors when not authorized", func() {
		expectDeposit(pkg.Dollars(20000), pkg.AuthorizationRequiredError)
		expectWithdraw(pkg.Dollars(20000), pkg.ZeroAmount, pkg.AuthorizationRequiredError)
		expectBalance(pkg.ZeroAmount, pkg.AuthorizationRequiredError)
		_, err := atm.History()
		Expect(err).To(Equal(pkg.AuthorizationRequiredError))
	})

	It("properly logs out after logout seconds", func() {
		authorize()
		time.Sleep(1500 * time.Millisecond)
		expectBalance(pkg.ZeroAmount, pkg.AuthorizationRequiredError)
	})

	It("handles invalid auth", func() {
		err := atm.Authorize(id, "2345")
		Expect(err).To(Equal(pkg.AuthorizationFailedError))
	})

})
