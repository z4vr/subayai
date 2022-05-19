package database

// Database is the interface for a database driver.
type Database interface {

	// GENERAL
	Connect(credentials ...interface{}) error
	Close()
}
