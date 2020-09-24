package application


func (app *Application) router() {
	app.r.HandleFunc("/health", app.HealthHandler)
}
