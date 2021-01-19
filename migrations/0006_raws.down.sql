DROP INDEX tapes_account_id_serial_number_index;
ALTER TABLE tapes
    ALTER COLUMN serial_number DROP NOT NULL;

create index tapes_account_id_serial_number_index
    on tapes (account_id, serial_number);
create unique index tapes_account_id_uci_uindex
    on tapes (account_id, uci);
create unique index tapes_account_id_alt_uci_uindex
    on tapes (account_id, alt_uci);

ALTER TABLE submissions
    DROP COLUMN kvs;
ALTER TABLE submissions
    DROP COLUMN raw;
