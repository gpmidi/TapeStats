ALTER TABLE submissions
    ADD COLUMN raw bytea;
ALTER TABLE submissions
    ADD COLUMN kvs hstore;

DROP INDEX tapes_account_id_alt_uci_uindex;
DROP INDEX tapes_account_id_uci_uindex;
DROP INDEX tapes_account_id_serial_number_index;

ALTER TABLE tapes
    ALTER COLUMN serial_number SET NOT NULL;
create unique index tapes_account_id_serial_number_index
    on tapes (account_id, serial_number);

