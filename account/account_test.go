package account

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/pkg/errors"
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

func TestAccount_AddSubBalance(t *testing.T) {
	amount := 10
	rep := NewAccountRepository(p)
	userStep1, err := rep.AddBalance(1, amount) // + 10
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("userID = %d, balance + %d = %d\n", userStep1.ID, amount, userStep1.Balance)

	userStep2, err := rep.AddBalance(1, amount) // +10
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("userID = %d, balance + %d = %d\n", userStep2.ID, amount, userStep2.Balance)

	if diff := userStep2.Balance - userStep1.Balance; diff != amount {
		t.Errorf("diff = %d; the balance should increase by %d\n", diff, amount)
	}

	userStep3, err := rep.SubBalance(1, amount) // -10
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("userID = %d, balance - %d = %d\n", userStep1.ID, amount, userStep1.Balance)

	if diff := userStep2.Balance - userStep3.Balance; diff != amount {
		t.Errorf("diff = %d; the balance should decrease by %d\n", diff, amount)
	}
}

func TestAccount_CheckPositiveBalance(t *testing.T) {
	amount := 1
	rep := NewAccountRepository(p)
	user, err := rep.AddBalance(1, amount) // + 1
	if err != nil {
		t.Fatal(err)
	}

	_, err = rep.SubBalance(user.ID, user.Balance+amount) // new_balance = balance - (balance + amount)
	if err, ok := err.(*pq.Error); ok {
		if err.Code.Name() != "violation_check" {
			t.Fatal(err, "db must check balance on the negativity")
		}
	}
}

func TestAccount_GetBalance(t *testing.T) {
	existingID := 1
	absentID := 100
	rep := NewAccountRepository(p)
	user, err := rep.GetBalance(existingID)
	if err != nil {
		t.Error(err)
	}
	t.Logf("user = %#v\n", user)

	_, err = rep.GetBalance(absentID)
	if err != sql.ErrNoRows {
		t.Error(errors.Wrap(err, "expected ErrNoRows, affected %s"))
	}
}

func TestAccount_Transfer(t *testing.T) {
	amount := 10
	transferAmount := 5
	senderID, recipientID := 1, 2
	rep := NewAccountRepository(p)
	u1Step1, err := rep.AddBalance(senderID, amount)
	if err != nil {
		t.Fatal(err)
	}

	u2Step1, err := rep.GetBalance(recipientID)
	if err != nil {
		t.Error(err)
	}

	u1Step2, u2Step2, err := rep.Transfer(senderID, recipientID, transferAmount)
	if err != nil {
		t.Error(err)
	}
	if diff := u1Step1.Balance - u1Step2.Balance; diff != transferAmount {
		t.Errorf("expected %d, affected %d\n", transferAmount, diff)
	}
	if diff := u2Step2.Balance - u2Step1.Balance; diff != transferAmount {
		t.Errorf("expected %d, affected %d\n", transferAmount, diff)
	}
	if sum1, sum2 := u1Step1.Balance+u2Step1.Balance, u1Step2.Balance+u2Step2.Balance; sum1 != sum2 {
		t.Error("balance sums must be equal before and after transfer")
	}
}

//TODO
func TestAccount_ConcurrencyTransfer(t *testing.T) {

}
