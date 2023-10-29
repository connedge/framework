package database

import (
	"gorm.io/gorm"
	"testing"
)

func TestOpen(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		mockDialector := &mockDialector{}
		cfg := Config{Driver: "test"}

		conn, err := Open(cfg)
		if err != nil {
			t.Fatal(err)
		}

		if conn.DB() != mockDialector.db {
			t.Errorf("unexpected db instance")
		}
	})

	t.Run("error", func(t *testing.T) {
		cfg := Config{Driver: "test"}

		conn, err := Open(cfg)
		if err == nil {
			t.Errorf("expected error but got none")
		}

		if conn != nil {
			t.Errorf("expected nil conn on error")
		}
	})

}

type mockDialector struct {
	db  *gorm.DB
	err error
}

func (m *mockDialector) Open(dsn string) (*gorm.DB, error) {
	return m.db, m.err
}

func (m *mockDialector) Close() {
}
