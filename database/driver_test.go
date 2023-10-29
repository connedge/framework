package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"testing"
)

func TestDriverStrategy(t *testing.T) {

	t.Run("mysql", func(t *testing.T) {
		dialector := driverStrategy("MYSQL", "dsn")
		_, ok := dialector.(*mysql.Dialector)
		if !ok {
			t.Error("expected mysql dialector")
		}
	})

	t.Run("postgres", func(t *testing.T) {
		dialector := driverStrategy("POSTGRESQL", "dsn")
		_, ok := dialector.(*postgres.Dialector)
		if !ok {
			t.Error("expected postgres dialector")
		}
	})

	t.Run("unknown", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic for unknown driver")
			}
		}()

		driverStrategy("UNKNOWN", "dsn")
	})

}
