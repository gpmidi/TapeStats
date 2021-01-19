ALTER TABLE submissions
    ADD COLUMN requester    hstore,
    ADD COLUMN requester_ip inet,
    ADD COLUMN request_id   VARCHAR(255);
