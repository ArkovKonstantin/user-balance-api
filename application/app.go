package application

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"user-balance-api/repository"
)

type Application struct {
	servicePort int
	r           *mux.Router
	rep         repository.Account
}

func New(rep repository.Account) Application {
	return Application{servicePort: 9000, r: mux.NewRouter(), rep: rep}
}

func (app *Application) Start() {
	app.router()
	fmt.Printf("start and listening on :%d ...\n", app.servicePort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(`:%d`, app.servicePort), app.r))
}
