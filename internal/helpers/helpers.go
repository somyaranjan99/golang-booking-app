package helpers

import (
	"fmt"
	"github/somyaranjan99/basic-go-project/pkg/config"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

func NewErrorLogs(a *config.AppConfig) {
	app = a
}
func ClientError(w http.ResponseWriter, status int) {
	app.Infolog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
