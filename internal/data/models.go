package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.
type Models struct {
	Faculties   FacultyModel
	Students    StudentModel
	Users       UserModel
	Tokens      TokenModel
	Permissions PermissionModel
}

// For ease of use, we also add a New() method which returns a Models struct containing
// the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Faculties:   FacultyModel{DB: db},
		Students:    StudentModel{DB: db},
		Users:       UserModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Permissions: PermissionModel{DB: db},
	}
}
