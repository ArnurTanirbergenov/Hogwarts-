package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"arnur.second.try/internal/validator"
)

type Student struct {
	ID        int64     `json:"id"`                   // Unique integer ID for the faculty
	CreatedAt time.Time `json:"-"`                    // Timestamp for when the faculty is added to our database
	Name      string    `json:"name"`                 // Name
	Surname   string    `json:"surname"`              // Surname
	StudyYear int32     `json:"study_year,omitempty"` // year of study
	Age       int32     `json:"age,omitempty"`        // age of a student
	FacultyId int32     `json:"faculty_id,omitempty"` // id of his faculty
	Runtime   int32     `json:"runtime,omitempty"`    // Student runtime (in minutes)
	Version   int32     `json:"version"`              // The version number starts at 1 and will be incremented each
	// time the student information is updated
}

func ValidateStudent(v *validator.Validator, student *Student) {
	v.Check(student.Name != "", "name", "must be provided")
	v.Check(len(student.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(student.Surname != "", "surname", "must be provided")
	v.Check(len(student.Surname) <= 500, "surname", "must not be more than 500 bytes long")
	v.Check(student.StudyYear != 0, "study_year", "must be provided")
	v.Check(student.StudyYear >= 1, "study_year", "must be greater than 1")
	v.Check(student.StudyYear <= 7, "study_year", "must be less than 7")
	v.Check(student.Age != 0, "age", "must be provided")
	v.Check(student.Age >= 11, "age", "must be greater than 11")
	v.Check(student.Age <= 19, "age", "must be less than 19")
	v.Check(student.FacultyId != 0, "faculty_id", "must be provided")
	v.Check(student.Runtime != 0, "runtime", "must be provided")
	v.Check(student.Runtime > 0, "runtime", "must be a positive integer")
}

// Define a StudentModel struct type which wraps a sql.DB connection pool.
type StudentModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the student table.
func (s StudentModel) Insert(student *Student) error {
	query := `
	INSERT INTO students (name, surname, study_year, age, faculty_id, runtime)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, version`

	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []interface{}{student.Name, student.Surname, student.StudyYear, student.Age, student.FacultyId, student.Runtime}

	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the movie struct.
	return s.DB.QueryRow(query, args...).Scan(&student.ID, &student.CreatedAt, &student.Version)

}

// Add a placeholder method for fetching a specific record from the student table.
func (s StudentModel) Get(id int64) (*Student, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, created_at, name, surname, study_year, age, faculty_id,
	runtime, version
	FROM students
	WHERE id = $1`

	var student Student
	err := s.DB.QueryRow(query, id).Scan(
		&student.ID,
		&student.CreatedAt,
		&student.Name,
		&student.Surname,
		&student.StudyYear,
		&student.Age,
		&student.FacultyId,
		&student.Runtime,
		&student.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &student, nil
}

// Add a placeholder method for updating a specific record in the student table.
func (s StudentModel) Update(student *Student) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `
UPDATE students
SET name = $1, surname=$2, study_year=$3, age=$4, faculty_id=$5, runtime = $6, version = version + 1
WHERE id = $7
RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		student.Name,
		student.Surname,
		student.StudyYear,
		student.Age,
		student.FacultyId,
		student.Runtime,
		student.ID,
	}
	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version value into the movie struct.
	return s.DB.QueryRow(query, args...).Scan(&student.Version)
}

// Add a placeholder method for deleting a specific record from the student table.
func (s StudentModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
		DELETE FROM students
		WHERE id = $1`
	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	result, err := s.DB.Exec(query, id)
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

func (s StudentModel) GetAll(name string, surname string, filters Filters, study_year int, age int, faculty_id int) ([]*Student, Metadata, error) {
	// Construct the SQL query to retrieve all movie records.
	query := fmt.Sprintf(`
	SELECT  count(*) OVER(), id, created_at, name, surname, study_year, age, faculty_id,
	runtime, version
	FROM students
	WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND (to_tsvector('simple', surname) @@ plainto_tsquery('simple', $2) OR $2 = '')
	AND (study_year = $3 OR $3 = 0)
	AND (age = $4 OR $4 = 0)
	AND (faculty_id = $5 OR $5 = 0)
	ORDER BY %s %s, id ASC
	LIMIT $6 OFFSET $7`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, query, name, surname, study_year, age, faculty_id, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	students := []*Student{}
	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var student Student
		// Scan the values from the row into the struct.
		err := rows.Scan(
			&totalRecords,
			&student.ID,
			&student.CreatedAt,
			&student.Name,
			&student.Surname,
			&student.StudyYear,
			&student.Age,
			&student.FacultyId,
			&student.Runtime,
			&student.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		students = append(students, &student)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// If everything went OK, then return the slice of movies.
	return students, metadata, nil
}
