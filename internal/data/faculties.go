package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"arnur.second.try/internal/validator"
)

type Faculty struct {
	ID        int64     `json:"id"`                // Unique integer ID for the faculty
	CreatedAt time.Time `json:"created_at"`        // Timestamp for when the faculty is added to our database
	Founder   string    `json:"founder,omitempty"` // Founder of a faculty
	Title     string    `json:"title"`             // Faculty title
	Year      int32     `json:"year,omitempty"`    // Faculty release year
	Runtime   int32     `json:"runtime"`           // Faculty runtime (in minutes)
	Version   int32     `json:"version"`           // The version number starts at 1 and will be incremented each
	// time the faculty information is updated
}

func ValidateFaculty(v *validator.Validator, faculty *Faculty) {

	v.Check(faculty.Title != "", "title", "must be provided")
	v.Check(len(faculty.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(faculty.Year != 0, "year", "must be provided")
	v.Check(faculty.Year >= 1, "year", "must be greater than 1")
	v.Check(faculty.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(faculty.Runtime != 0, "runtime", "must be provided")
	v.Check(faculty.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(faculty.Founder != "", "founder", "must be provided")
	v.Check(len(faculty.Founder) <= 500, "founder", "must not be more than 500 bytes long")
}

type FacultyModel struct {
	DB *sql.DB
}

func (f FacultyModel) Insert(faculty *Faculty) error {
	query := `
	INSERT INTO faculties (title, year, runtime, founder)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version`

	args := []interface{}{faculty.Title, faculty.Year, faculty.Runtime, faculty.Founder}

	return f.DB.QueryRow(query, args...).Scan(&faculty.ID, &faculty.CreatedAt, &faculty.Version)
}

func (f FacultyModel) Get(id int64) (*Faculty, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT id, created_at, title, year, runtime, founder, version
		FROM faculties
		WHERE id = $1`

	var faculty Faculty

	err := f.DB.QueryRow(query, id).Scan(
		&faculty.ID,
		&faculty.CreatedAt,
		&faculty.Title,
		&faculty.Year,
		&faculty.Runtime,
		&faculty.Founder,
		&faculty.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the struct.
	return &faculty, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (f FacultyModel) Update(faculty *Faculty) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `
UPDATE faculties
SET title = $1, year = $2, runtime = $3, founder = $4, version = version + 1
WHERE id = $5
RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		faculty.Title,
		faculty.Year,
		faculty.Runtime,
		faculty.Founder,
		faculty.ID,
	}
	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version value into the movie struct.
	return f.DB.QueryRow(query, args...).Scan(&faculty.Version)

}

// Add a placeholder method for deleting a specific record from the movies table.
func (f FacultyModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
		DELETE FROM faculties
		WHERE id = $1`
	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	result, err := f.DB.Exec(query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil

}

func (f FacultyModel) GetAll(founder string, title string, year int, filters Filters) ([]*Faculty, error) {

	query := fmt.Sprintf(`
	SELECT id, founder, title, year, runtime,version
	FROM faculties
	WHERE  (to_tsvector('simple', founder) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND  (to_tsvector('simple', title) @@ plainto_tsquery('simple', $2) OR $2 = '')
	ORDER BY %s %s, id ASC
	LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := f.DB.QueryContext(ctx, query, founder, title, filters.limit(), filters.offset())
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	faculties := []*Faculty{}

	for rows.Next() {

		var faculty Faculty

		err := rows.Scan(
			&faculty.ID,
			&faculty.Founder,
			&faculty.Title,
			&faculty.Year,
			&faculty.Runtime,
			&faculty.Version,
		)
		if err != nil {
			return nil, err
		}

		faculties = append(faculties, &faculty)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return faculties, nil
}
