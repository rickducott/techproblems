package pkg_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rickducott/techproblems/atm/pkg"
)

var _ = Describe("Account", func() {
	var (
		account pkg.Account
	)

	BeforeEach(func() {
		account = pkg.NewAccount("12345", "1234", pkg.Cents(10000))
	})

	It("works when depositing", func() {
		txn, err := account.Transaction(pkg.Cents(1000))
		Expect(err).To(BeNil())
		Expect(txn.Balance).To(Equal(pkg.Cents(11000)))
		Expect(txn.Amount).To(Equal(pkg.Cents(1000)))
		Expect(txn.Overdraft).To(BeFalse())
	})

	It("works when withdrawing", func() {
		txn, err := account.Transaction(pkg.Cents(-1000))
		Expect(err).To(BeNil())
		Expect(txn.Balance).To(Equal(pkg.Cents(9000)))
		Expect(txn.Amount).To(Equal(pkg.Cents(-1000)))
		Expect(txn.Overdraft).To(BeFalse())
	})

	It("handles overdrafts", func() {
		txn, err := account.Transaction(pkg.Cents(-12000))
		Expect(err).To(BeNil())
		Expect(txn.Balance).To(Equal(pkg.Cents(-2000).Subtract(pkg.OverdraftFee)))
		Expect(txn.Amount).To(Equal(pkg.Cents(-12000)))
		Expect(txn.Overdraft).To(BeTrue())
	})

	It("returns error if already overdrawn", func() {
		_, err := account.Transaction(pkg.Cents(-12000))
		Expect(err).To(BeNil())
		_, err = account.Transaction(pkg.Cents(-10000))
		Expect(err).To(Equal(pkg.AccountOverdrawnError))
	})

	It("updates history and balance", func() {
		_, _ = account.Transaction(pkg.Cents(-1000))
		_, _ = account.Transaction(pkg.Cents(-2000))
		_, _ = account.Transaction(pkg.Cents(-3000))
		Expect(account.Balance()).To(Equal(pkg.Cents(4000)))
		Expect(account.History()).Should(HaveLen(3))
		Expect(account.History()[0].Amount).To(Equal(pkg.Cents(-1000)))
		Expect(account.History()[2].Amount).To(Equal(pkg.Cents(-3000)))
	})
})
