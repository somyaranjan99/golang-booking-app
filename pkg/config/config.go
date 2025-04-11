package config

import (
	"github/somyaranjan99/basic-go-project/pkg/models"
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	Session       *scs.Session
	Infolog       *log.Logger
	ErrorLog      *log.Logger
	MailChan      chan models.MailData
}
