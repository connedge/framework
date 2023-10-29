package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

const (
	MYSQL      = "MYSQL"
	POSTGRESQL = "POSTGRESQL"
)

func driverStrategy(driver string, dsn string) gorm.Dialector {
	switch driver {
	case MYSQL:
		return mysqlDriver(dsn)
	case POSTGRESQL:
		return pgDriver(dsn)
	default:
		log.Panicf("driver not implemented %s", driver)
	}

	return nil
}

func mysqlDriver(dsn string) gorm.Dialector {
	return mysql.Open(dsn)
}

func pgDriver(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}
