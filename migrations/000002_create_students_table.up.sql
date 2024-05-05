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

