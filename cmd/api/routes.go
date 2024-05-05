package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/faculties", app.requirePermission("faculties:write", app.createFacultyHandler))
	router.HandlerFunc(http.MethodGet, "/v1/faculties/:id", app.requirePermission("faculties:read", app.showFacultyHandler))
	router.HandlerFunc(http.MethodPut, "/v1/faculties/:id", app.requirePermission("faculties:write", app.updateFacultyHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/faculties/:id", app.requirePermission("faculties:write", app.deleteFacultyHandler))
	router.HandlerFunc(http.MethodGet, "/v1/faculties/:id/students", app.requirePermission("faculties:read", app.showFacultyStudentHandler))

	router.HandlerFunc(http.MethodGet, "/v1/students", app.requirePermission("students:read", app.listStudentsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/students", app.requirePermission("students:write", app.createStudentHandler))
	router.HandlerFunc(http.MethodGet, "/v1/students/:id", app.requirePermission("students:read", app.showStudentHandler))
	router.HandlerFunc(http.MethodPut, "/v1/students/:id", app.requirePermission("students:write", app.updateStudentHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/students/:id", app.requirePermission("students:write", app.deleteStudentHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic((app.authenticate(router)))
}
