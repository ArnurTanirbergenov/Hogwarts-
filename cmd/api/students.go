package main

import (
	"errors"
	"net/http"

	"arnur.second.try/internal/data"
	"arnur.second.try/internal/validator"
)

func (app *application) createStudentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string `json:"name"`
		Surname   string `json:"surname"`
		StudyYear int32  `json:"study_year"`
		Age       int32  `json:"age"`
		FacultyId int32  `json:"faculty_id"`
		Runtime   int32  `json:"runtime"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	student := &data.Student{
		Name:      input.Name,
		Surname:   input.Surname,
		StudyYear: input.StudyYear,
		Age:       input.Age,
		FacultyId: input.FacultyId,
		Runtime:   input.Runtime,
	}
	v := validator.New()
	if data.ValidateStudent(v, student); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Students.Insert(student)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)

	// headers.Set("Location", fmt.Sprintf("v1/students/%d", student.ID))
	// fmt.Fprintf(w, "%+v\n", input)

	err = app.writeJSON(w, http.StatusCreated, envelope{"student": student}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showStudentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Call the Get() method to fetch the data for a specific movie. We also need to
	// use the errors.Is() function to check if it returns a data.ErrRecordNotFound
	// error, in which case we send a 404 Not Found response to the client.
	student, err := app.models.Students.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"student": student}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateStudentHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the movie ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the existing movie record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	student, err := app.models.Students.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Declare an input struct to hold the expected data from the client.
	var input struct {
		Name      string `json:"name"`
		Surname   string `json:"surname"`
		StudyYear int32  `json:"study_year"`
		Age       int32  `json:"age"`
		FacultyId int32  `json:"faculty_id"`
		Runtime   int32  `json:"runtime"`
	}
	// Read the JSON request body data into the input struct.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the values from the request body to the appropriate fields of the movie
	// record.
	student.Name = input.Name
	student.Surname = input.Surname
	student.StudyYear = input.StudyYear
	student.Age = input.Age
	student.FacultyId = input.FacultyId
	student.Runtime = input.Runtime
	// Validate the updated movie record, snding the client a 422 Unprocessable Entity
	// response if any checks fail.
	v := validator.New()
	if data.ValidateStudent(v, student); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Pass the updated movie record to our new Update() method.
	err = app.models.Students.Update(student)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Write the updated movie record in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"student": student}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) deleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the movie ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the movie from the database, sending a 404 Not Found response to the
	// client if there isn't a matching record.
	err = app.models.Students.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "student successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listStudentsHandler(w http.ResponseWriter, r *http.Request) {
	// To keep things consistent with our other handlers, we'll define an input struct
	// to hold the expected values from the request query string.
	var input struct {
		Name      string
		Surname   string
		StudyYear int
		Age       int
		FacultyId int
		data.Filters
	}
	// Initialize a new Validator instance.
	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Surname = app.readString(qs, "surname", "")
	input.Age = app.readInt(qs, "age", 0, v)
	input.FacultyId = app.readInt(qs, "faculty_id", 0, v)
	input.StudyYear = app.readInt(qs, "study_year", 0, v)

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply a ascending sort on movie ID).

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "surname", "runtime", "study_year", "age", "faculty_id", "-id", "-name", "-surname", "-runtime", "-study_year", "-age", "-faculty_id"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Call the GetAll() method to retrieve the movies, passing in the various filter
	// parameters.
	students, metadata, err := app.models.Students.GetAll(input.Name, input.Surname, input.Filters, input.StudyYear, input.Age, input.FacultyId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send a JSON response containing the movie data.
	err = app.writeJSON(w, http.StatusOK, envelope{"students": students, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
