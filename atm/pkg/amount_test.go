package pkg_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rickducott/techproblems/atm/pkg"
)

var _ = Describe("Amount", func() {
	Context("arithmetic", func() {
		It("works", func() {
			amount := pkg.NewAmount(-1, -15)
			Expect(amount).To(Equal(pkg.Cents(-115)))
			Expect(pkg.ZeroAmount.GreaterThan(amount)).To(BeTrue())
			Expect(amount.Negative()).To(Equal(pkg.Cents(115)))
			Expect(amount.Add(pkg.Cents(115))).To(Equal(pkg.ZeroAmount))
			Expect(amount.Subtract(pkg.Cents(115))).To(Equal(pkg.Cents(-230)))
		})
	})

	Context("parsing", func() {
		It("works", func() {
			amt, err := pkg.ParseAmount("1.15")
			Expect(err).To(BeNil())
			Expect(amt).To(Equal(pkg.Cents(115)))
		})

		It("handles negative amounts", func() {
			amt, err := pkg.ParseAmount("-1.15")
			Expect(err).To(BeNil())
			Expect(amt).To(Equal(pkg.Cents(-115)))
		})

		It("handles dollar shorthands", func() {
			amt, err := pkg.ParseAmount("-.15")
			Expect(err).To(BeNil())
			Expect(amt).To(Equal(pkg.Cents(-15)))
			amt, err = pkg.ParseAmount(".15")
			Expect(err).To(BeNil())
			Expect(amt).To(Equal(pkg.Cents(15)))
		})

		It("handles cents shorthands", func() {
			amt, err := pkg.ParseAmount("-.1")
			Expect(err).To(BeNil())
			Expect(amt).To(Equal(pkg.Cents(-10)))
			amt, err = pkg.ParseAmount("-.01")
			Expect(err).To(BeNil())
			Expect(amt).To(Equal(pkg.Cents(-1)))
			amt, err = pkg.ParseAmount("1.")
			Expect(err).To(BeNil())
			Expect(amt).To(Equal(pkg.Cents(100)))
		})

		It("errors on invalid cents", func() {
			_, err := pkg.ParseAmount("-.123")
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("Error parsing amount: too many digits in the cents"))
		})

	})

	Context("rendering", func() {
		It("works", func() {
			amt := pkg.Cents(1005)
			Expect(amt.String()).To(Equal("10.05"))
			amt = pkg.Cents(-1005)
			Expect(amt.String()).To(Equal("-10.05"))
			amt = pkg.Cents(1000)
			Expect(amt.String()).To(Equal("10.00"))
			amt = pkg.Cents(1099)
			Expect(amt.String()).To(Equal("10.99"))
		})
	})
})
