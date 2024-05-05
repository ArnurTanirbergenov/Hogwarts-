CREATE TABLE IF NOT EXISTS faculties (
    id bigserial PRIMARY KEY,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    founder text NOT NULL,
    title text NOT NULL,
    year integer NOT NULL, 
    Runtime integer NOT NULL,
    version integer NOT NULL DEFAULT 1
);
