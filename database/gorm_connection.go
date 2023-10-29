package database

import (
	"gorm.io/gorm"
)

type Database interface {
	Instance() *gorm.DB
	Shutdown() error
}

func Open(cfg Config) (Database, error) {
	dsn := dsnStrategy(cfg.Driver, cfg)
	db, err := gorm.Open(driverStrategy(cfg.Driver, dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	return &database{db: db}, nil
}

type database struct {
	db *gorm.DB
}

func NewConnection(db *gorm.DB) Database {
	return &database{db: db}
}

func (c *database) Instance() *gorm.DB {
	return c.db
}

func (c *database) Shutdown() error {
	db, err := c.db.DB()
	if err != nil {
		return err
	}
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}
