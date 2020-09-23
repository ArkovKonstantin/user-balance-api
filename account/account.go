package account

import (
	"github.com/pkg/errors"
	"log"
	"user-balance-api/provider"
)

type Account interface {
	AddBalance(int, int) (int, int, error)
	SubBalance(int, int) (int, int, error)
	Transfer()
}

type account struct {
	p provider.Provider
}

func NewAccountRepository(p provider.Provider) Account {
	return &account{p}
}

func (acc *account) AddBalance(userID, amount int) (int, balance int, err error) {
	db, err := acc.p.GetConn()
	if err != nil {
		return userID, balance, errors.Wrap(err, "get db connection err:")
	}

	err = db.QueryRow(`insert into account (user_id, balance)
			   values ($1, $2)
			   on conflict(user_id) do update set balance = account.balance + $2
			   where account.user_id = $1
			   returning balance;`, userID, amount).Scan(&balance)
	if err != nil {
		log.Println(errors.Wrap(err, "sub balance, db operation err"))
		return userID, balance, errors.Wrap(err, "add balance, db operation err")
	}

	return userID, balance, nil
}

func (acc *account) SubBalance(userID, amount int) (int, balance int, err error) {
	db, err := acc.p.GetConn()
	if err != nil {
		return userID, balance, errors.Wrap(err, "get db connection err:")
	}

	err = db.QueryRow(`insert into account (user_id, balance)
			   values ($1, 0)
			   on conflict(user_id) do update set balance = account.balance - $2
			   where account.user_id = $1
			   returning balance;`, userID, amount).Scan(&balance)

	if err != nil {
		log.Println(errors.Wrap(err, "sub balance, db operation err"))
		return userID, balance, errors.Wrap(err, "sub balance, db operation err")
	}

	return userID, balance, nil
}

func (acc *account) Transfer() {

}
