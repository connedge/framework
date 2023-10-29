package database

import (
	"github.com/connedge/framework/connedge"
	"github.com/connedge/framework/ioc"
	"log"
)

type Provider struct{}

func (r Provider) Register(app *connedge.App) {

	database, err := Open(Config{
		Host:     app.Config.GetString("database.host"),
		Port:     app.Config.GetString("database.port"),
		Password: app.Config.GetString("database.password"),
		User:     app.Config.GetString("database.user"),
		Driver:   app.Config.GetString("database.driver"),
		DBName:   app.Config.GetString("database.name"),
	})

	if err != nil {
		log.Fatalf("database.Provider error: %s", err.Error())
	}

	ioc.Bind(app.Container, func(c *ioc.Container) (Database, error) {
		return NewConnection(database.Instance()), nil
	})

}

func (r Provider) Name() string {
	return "DatabaseProvider"
}

func (r Provider) Handler(app *connedge.App) []connedge.Handler {
	return nil
}
