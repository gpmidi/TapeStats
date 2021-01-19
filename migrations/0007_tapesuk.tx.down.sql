DROP INDEX tapes_account_lookup;
create unique index tapes_account_id_serial_number_index
    on tapes (account_id, serial_number);
