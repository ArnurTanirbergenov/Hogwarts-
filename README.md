# Hogwarts
Here is my Golang API for Hogwarts. Users can create, change, and see different faculties and students by it.This project was invented when i decided to rewatch Harry Potter.
I will be writing using Go language, PostgreSQL.

If you will download this project to start it you can simply write
```
  go run ./cmd/api
```
To see my database i use
```
psql --host=localhost --dbname=hogwarts --username=hogwarts
```
Here is my endpoints
```
router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/faculties", app.requirePermission("faculties:write", app.createFacultyHandler))
	router.HandlerFunc(http.MethodGet, "/v1/faculties/:id", app.requirePermission("faculties:read", app.showFacultyHandler))
	router.HandlerFunc(http.MethodPut, "/v1/faculties/:id", app.requirePermission("faculties:write", app.updateFacultyHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/faculties/:id", app.requirePermission("faculties:write", app.deleteFacultyHandler))

	router.HandlerFunc(http.MethodGet, "/v1/students", app.requirePermission("students:read", app.listStudentsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/students", app.requirePermission("students:write", app.createStudentHandler))
	router.HandlerFunc(http.MethodGet, "/v1/students/:id", app.requirePermission("students:read", app.showStudentHandler))
	router.HandlerFunc(http.MethodPut, "/v1/students/:id", app.requirePermission("students:write", app.updateStudentHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/students/:id", app.requirePermission("students:write", app.deleteStudentHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
```
As you can clearly see i have users authorization, so to start you need to first go to  "/v1/users", than use token that will come to your email or you can take it from answer and use it in "/v1/users/activated",
then there is authentication "/v1/tokens/authentication" where you will get your token, you will use it in any other endpoints, but not in "/v1/healthcheck".Take a note that not any user can write

There are my tables 
```
CREATE TABLE IF NOT EXISTS faculties (
    id bigserial PRIMARY KEY,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    founder text NOT NULL,
    title text NOT NULL,
    year integer NOT NULL, 
    Runtime integer NOT NULL,
    version integer NOT NULL DEFAULT 1
);


CREATE TABLE IF NOT EXISTS students (
    id bigserial PRIMARY KEY,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    surname text NOT NULL,
    study_year integer NOT NULL, 
    age integer NOT NULL,
    Faculty_id INTEGER NOT NULL,
    Runtime integer NOT NULL,
    version integer NOT NULL DEFAULT 1,
    FOREIGN KEY (Faculty_id) REFERENCES faculties(id)
);
```
Also tables for storing users
```
CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    code text NOT NULL
);
CREATE TABLE IF NOT EXISTS users_permissions (
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);
-- Add the two permissions to the table.
INSERT INTO permissions (code)
VALUES
('students:read'),
('students:write'),
('faculties:read'),
('faculties:write');
```

