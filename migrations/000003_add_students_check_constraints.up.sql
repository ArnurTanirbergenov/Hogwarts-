ALTER TABLE students ADD CONSTRAINT students_runtime_check CHECK (runtime >= 0);
ALTER TABLE students ADD CONSTRAINT students_study_year_check CHECK (study_year BETWEEN 0 AND 7);
ALTER TABLE students ADD CONSTRAINT students_age_check CHECK (age BETWEEN 11 AND 19);


