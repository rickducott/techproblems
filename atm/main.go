package main

import (
	"bufio"
	"fmt"
	"github.com/rickducott/techproblems/atm/pkg"
	"os"
	"strings"
)

var (
	AccountData = []pkg.Account{
		pkg.NewAccount("2859459814", "7386", pkg.NewAmount(10, 24)),
		pkg.NewAccount("1434597300", "4557", pkg.NewAmount(90000, 55)),
		pkg.NewAccount("7089382418", "0075", pkg.ZeroAmount),
		pkg.NewAccount("2001377812", "5950", pkg.NewAmount(60, 0)),
	}

	LogoutSeconds = 120
)

func main() {
	atm, done := pkg.NewAtm(LogoutSeconds, AccountData...)
	textUi := pkg.NewInterface(atm)
	reader := bufio.NewReader(os.Stdin)
	end := false
	for !end {
		fmt.Printf("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			end = true
		} else if strings.TrimSpace(line) == "end" {
			end = true
		} else {
			output := textUi.Execute(line)
			fmt.Printf("%s\n", output)
		}
	}

	done <-true
}