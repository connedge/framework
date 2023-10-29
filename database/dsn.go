package database

import (
	"fmt"
	"log"
)

func dsnStrategy(driver string, cfg Config) string {
	switch driver {
	case MYSQL:
		return mysqlDsn(cfg)
	case POSTGRESQL:
		return pgDsn(cfg)
	default:
		log.Fatalf("driver not implemented %s", driver)
	}

	return ""
}

func mysqlDsn(cfg Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
}

func pgDsn(cfg Config) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
	)
}
