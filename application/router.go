package application

func (app *Application) router() {
	app.r.HandleFunc("/health", app.HealthHandler)
	app.r.HandleFunc("/balance/get", app.GetBalanceHandler).Queries("user_id", "{user_id:[0-9]+}").Methods("GET")
	app.r.HandleFunc("/balance/add", app.AddBalanceHandler).Methods("POST")
	app.r.HandleFunc("/balance/sub", app.SubBalanceHandler).Methods("POST")
	app.r.HandleFunc("/transfer", app.TransferHandler).Methods("POST")
}
