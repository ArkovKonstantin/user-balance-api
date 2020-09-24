package application

import (
	"database/sql"
	"github.com/lib/pq"
	"net/http"
	"strconv"
)

const (
	ErrNegativeBalance = "not enough money"
	ErrInvalidUser     = "user does not exists"
)

func (app *Application) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (app *Application) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.FormValue("user_id"))
	user, err := app.rep.GetBalance(userID)
	switch {
	case err == sql.ErrNoRows:
		http.Error(w, ErrInvalidUser, http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := encodeResponse(w, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *Application) AddBalanceHandler(w http.ResponseWriter, r *http.Request) {
	req, err := decodeChangeBalanceReq(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := app.rep.AddBalance(req.UserID, req.Amount)
	if err, ok := err.(*pq.Error); ok {
		if err.Code.Name() == "foreign_key_violation" {
			http.Error(w, ErrInvalidUser, http.StatusBadRequest)
			return
		}
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := encodeResponse(w, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *Application) SubBalanceHandler(w http.ResponseWriter, r *http.Request) {
	req, err := decodeChangeBalanceReq(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := app.rep.SubBalance(req.UserID, req.Amount)
	if err, ok := err.(*pq.Error); ok {
		switch err.Code.Name() {
		case "check_violation":
			http.Error(w, ErrNegativeBalance, http.StatusBadRequest)
		case "foreign_key_violation":
			http.Error(w, ErrInvalidUser, http.StatusBadRequest)
		}
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := encodeResponse(w, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *Application) TransferHandler(w http.ResponseWriter, r *http.Request) {
	req, err := decodeTransferReq(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sender, recipient, err := app.rep.Transfer(req.SenderID, req.RecipientID, req.Amount)
	if err, ok := err.(*pq.Error); ok {
		switch err.Code.Name() {
		case "check_violation":
			http.Error(w, ErrNegativeBalance, http.StatusBadRequest)
		case "foreign_key_violation":
			http.Error(w, ErrInvalidUser, http.StatusBadRequest)
		}
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := encodeResponse(w, TransferResp{sender, recipient}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}