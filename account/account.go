package account

import (
	"github.com/pkg/errors"
	"log"
	"user-balance-api/models"
	"user-balance-api/provider"
)

type Account interface {
	AddBalance(userID, amount int) (user models.User, err error)
	SubBalance(userID, amount int) (user models.User, err error)
	Transfer(int, int, int) (models.User, models.User, error)
	GetBalance(userID int) (user models.User, err error)
}

type account struct {
	p provider.Provider
}

func NewAccountRepository(p provider.Provider) Account {
	return &account{p}
}

func (acc *account) AddBalance(userID, amount int) (user models.User, err error) {
	user.ID = userID
	db, err := acc.p.GetConn()
	if err != nil {
		return user, errors.Wrap(err, "get db connection err:")
	}

	err = db.QueryRow(`insert into account (user_id, balance)
			   values ($1, $2)
			   on conflict(user_id) do update set balance = account.balance + $2
			   where account.user_id = $1
			   returning balance;`, userID, amount).Scan(&user.Balance)
	if err != nil {
		log.Println(errors.Wrap(err, "add balance, db operation err"))
		return user, errors.Wrap(err, "add balance, db operation err")
	}

	return user, nil
}

func (acc *account) SubBalance(userID, amount int) (user models.User, err error) {
	user.ID = userID
	db, err := acc.p.GetConn()
	if err != nil {
		return user, errors.Wrap(err, "get db connection err:")
	}

	err = db.QueryRow(`insert into account (user_id, balance)
			   values ($1, 0)
			   on conflict(user_id) do update set balance = account.balance - $2
			   where account.user_id = $1
			   returning balance;`, userID, amount).Scan(&user.Balance)

	if err != nil {
		log.Println(errors.Wrap(err, "sub balance, db operation err"))
		return user, errors.Wrap(err, "sub balance, db operation err")
	}

	return user, nil
}

func (acc *account) Transfer(senderID, recipientID, amount int) (sender, recipient models.User, err error) {
	sender.ID = senderID
	recipient.ID = recipientID
	db, err := acc.p.GetConn()
	if err != nil {
		return sender, recipient, errors.Wrap(err, "get db connection err:")
	}

	tx, err := db.Begin()
	if err != nil {
		return sender, recipient, errors.Wrap(err, "begin tx err:")
	}

	err = db.QueryRow(`insert into account (user_id, balance)
			   values ($1, 0)
			   on conflict(user_id) do update set balance = account.balance - $2
			   where account.user_id = $1
			   returning balance;`, senderID, amount).Scan(&sender.Balance)
	if err != nil {
		tx.Rollback()
		return
	}

	err = db.QueryRow(`insert into account (user_id, balance)
			   values ($1, $2)
			   on conflict(user_id) do update set balance = account.balance + $2
			   where account.user_id = $1
			   returning balance;`, recipientID, amount).Scan(&recipient.Balance)
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	return
}

func (acc *account) GetBalance(userID int) (user models.User, err error) {
	db, err := acc.p.GetConn()
	if err != nil {
		return user, errors.Wrap(err, "get db connection err:")
	}

	err = db.QueryRow(`select user_id, balance from account where user_id = $1`, userID).Scan(&user.ID, &user.Balance)
	if err != nil {
		return
	}

	return
}
