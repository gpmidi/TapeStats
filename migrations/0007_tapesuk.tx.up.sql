DROP INDEX tapes_account_id_serial_number_index;
create unique index tapes_account_lookup
    on tapes (account_id, manufacturer, manufacture_dt, serial_number, density_code, medium_type, lto_version);
