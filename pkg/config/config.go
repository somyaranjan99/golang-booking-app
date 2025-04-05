package config

import (
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
}
