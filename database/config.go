package database

type Config struct {
	// Host is the hostname or IP address of the database server.
	Host string

	// DBName is the name of the database to connect to.
	DBName string

	// User is the username used to authenticate to the database.
	User string

	// Password is the password used to authenticate to the database.
	Password string

	// Port is the TCP port of the database server.
	Port string

	// Driver specifies the database driver to use, e.g. mysql, postgres, etc.
	Driver string
}
