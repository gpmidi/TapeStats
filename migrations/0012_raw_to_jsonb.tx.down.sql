ALTER TABLE submissions
    DROP COLUMN raw;
ALTER TABLE submissions
    ADD COLUMN raw TYPE bytea;
