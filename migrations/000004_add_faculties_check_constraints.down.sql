ALTER TABLE faculties ADD CONSTRAINT faculties_runtime_check CHECK (runtime >= 0);
ALTER TABLE faculties ADD CONSTRAINT faculties_year_check CHECK (year BETWEEN 0 AND date_part('year', now()));
