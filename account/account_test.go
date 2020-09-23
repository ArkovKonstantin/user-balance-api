package account

import (
	"github.com/lib/pq"
	"log"
	"os"
	"testing"
	"user-balance-api/models"
	"user-balance-api/provider"
)

var (
	config models.Config
	p      provider.Provider
)

func init() {
	env := os.Getenv("ENV")

	var path string
	if env == "" || env == "dev" {
		path = "../config/config.dev.toml"
	} else if env == "prod" {
		path = "../config/config.prod.toml"
	}

	err := config.LoadConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	p = provider.New(&config.SQLDataBase)
	err = p.Open()
	if err != nil {
		log.Fatal(err)
	}
}

func TestAddSubBalance(t *testing.T) {
	amount := 10
	rep := NewAccountRepository(p)
	userID, balance1, err := rep.AddBalance(1, amount) // + 10
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("userID = %d, balance + %d = %d\n", userID, amount, balance1)

	userID, balance2, err := rep.AddBalance(1, amount) // +10
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("userID = %d, balance + %d = %d\n", userID, amount, balance2)

	if diff := balance2 - balance1; diff != amount {
		t.Errorf("diff = %d; the balance should increase by %d\n", diff, amount)
	}

	userID, balance3, err := rep.SubBalance(1, amount) // -10
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("userID = %d, balance - %d = %d\n", userID, amount, balance3)

	if diff := balance2 - balance3; diff != amount {
		t.Errorf("diff = %d; the balance should decrease by %d\n", diff, amount)
	}

}

func TestCheckPositiveBalance(t *testing.T) {
	amount := 1
	rep := NewAccountRepository(p)
	_, balance1, err := rep.AddBalance(1, amount) // + 1
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = rep.SubBalance(1, balance1+amount) // new_balance = balance - (balance + amount)
	if err, ok := err.(*pq.Error); ok {
		if err.Code.Name() != "violation_check" {
			t.Fatal(err, "db must check balance on the negativity")
		}
	}
}
