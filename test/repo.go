package test

import (
	"os"
	"path"

	"github.com/OpenBazaar/openbazaar-go/repo"
	"github.com/OpenBazaar/openbazaar-go/repo/db"
)

// Repository represents a test (temporary/volitile) repository
type Repository struct {
	Path     string
	Password string
	DB       *db.SQLiteDatastore
}

// NewRepository creates and initializes a new temporary repository for tests
func NewRepository() (*Repository, error) {
	r := &Repository{
		Path:     GetRepoPath(),
		Password: GetPassword(),
	}

	// Create database
	sqliteDB, err := db.Create(r.Path, r.Password, true)
	if err != nil {
		return nil, err
	}

	r.DB = sqliteDB

	// Initialize the IPFS repo if it does not already exist
	err = repo.DoInit(r.Path, 4096, true, r.Password, r.Password, r.DB.Config().Init)
	if err != nil && err != repo.ErrRepoExists {
		return nil, err
	}

	// Reset to blank slate
	r.MustReset()

	return r, nil
}

// RemoveProfile removes the profile json from the repository
func (r *Repository) RemoveProfile() error {
	err := os.Remove(path.Join(r.Path, "root", "profile"))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// RemoveSettings purges settings from the database
func (r *Repository) RemoveSettings() error {
	return r.DB.Settings().Delete()
}

// RemoveRepo removes the test repository
func (r *Repository) RemoveRepo() error {
	return os.RemoveAll(r.Path)
}

// Reset sets the repo state back to a blank slate but retains keys
func (r *Repository) Reset() error {
	err := r.RemoveProfile()
	if err != nil {
		return err
	}

	err = r.RemoveSettings()
	if err != nil {
		return err
	}

	return nil
}

// MustReset sets the repo state back to a blank slate but retains keys
// It panics upon failure instead of allowing tests to continue
func (r *Repository) MustReset() {
	err := r.Reset()
	if err != nil {
		panic(err)
	}
}
