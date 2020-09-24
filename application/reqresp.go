package application

import (
	"encoding/json"
	"errors"
	"net/http"
	"user-balance-api/models"
)

type (
	ChangeBalanceReq struct {
		UserID int `json:"user_id"`
		Amount int `json:"amount"`
	}

	TransferReq struct {
		SenderID    int `json:"sender_id"`
		RecipientID int `json:"recipient_id"`
		Amount      int `json:"amount"`
	}

	TransferResp struct {
		Sender    models.User `json:"sender"`
		Recipient models.User `json:"recipient"`
	}
)

func decodeChangeBalanceReq(r *http.Request) (req ChangeBalanceReq, err error) {
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return
	}
	if req.Amount < 0 {
		return req, errors.New("amount must be positive")
	}

	return
}

func decodeTransferReq(r *http.Request) (req TransferReq, err error) {
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return
	}
	if req.Amount < 0 {
		return req, errors.New("amount must be positive")
	}

	return
}

func encodeResponse(w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
