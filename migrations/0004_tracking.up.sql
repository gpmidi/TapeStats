-- Tape Tracking
create table tapes
(
    id               uuid        default uuid_generate_v4() not null
        constraint tapes_pk
            primary key,
    account_id       uuid                                   not null
        constraint tapes_accounts_id_fk
            references accounts
            on update cascade on delete cascade,
    modified         timestamptz default NOW()              not null,
    created          timestamptz default NOW()              not null,
    uci              varchar(28)                            not null,
    alt_uci          VARCHAR(24)                            not null,
    serial_number    varchar(32),
    assigning_org    varchar(8),
    manufacturer     varchar(8),
    manufacture_dt   date,
    density_code     char(1),
    medium_type      char(1),
    medium_type_info char(2),
    lto_version      int
);

comment on table tapes is 'One row per unique tape';

comment on column tapes.uci is '1000 UNIQUE CARTRIDGE IDENTITY (binary, 28 bytes, read-only)';

comment on column tapes.alt_uci is '1001 ALTERNATIVE UNIQUE CARTRIDGE IDENTITY (binary, 24 bytes, read-only)';

comment on column tapes.serial_number is '401 MEDIUM SERIAL NUMBER (ascii, 32 bytes, read-only)';

comment on column tapes.assigning_org is '404 ASSIGNING ORGANIZATION (ascii, 8 bytes, read-only)';

comment on column tapes.manufacturer is '400 MEDIUM MANUFACTURER (ascii, 8 bytes, read-only)';

comment on column tapes.manufacture_dt is 'MEDIUM MANUFACTURE DATE (ascii, 8 bytes, read-only)';

comment on column tapes.density_code is '405 MEDIUM DENSITY CODE (binary, 1 bytes, read-only)';

comment on column tapes.medium_type is '408 MEDIUM TYPE (binary, 1 bytes, read-only)';

comment on column tapes.medium_type_info is '409 MEDIUM TYPE INFORMATION (binary, 2 bytes, read-only)';

create unique index tapes_account_id_alt_uci_uindex
    on tapes (account_id, alt_uci);

create index tapes_account_id_serial_number_index
    on tapes (account_id, serial_number);

create unique index tapes_account_id_uci_uindex
    on tapes (account_id, uci);

-- Tie modified trigger to tapes
CREATE TRIGGER trigger_tapes_set_modified
    BEFORE UPDATE
    ON tapes
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- Submissions
create table submission
(
    id                     bigserial                              not null
        constraint table_name_pk
            primary key,
    tape_id                uuid                                   not null
        constraint table_name_tapes_id_fk
            references tapes
            on update cascade on delete cascade,
    created                timestamp with time zone default now() not null,
    modified               timestamp with time zone default now() not null,
    tape_alert_flags       varchar(8),
    load_count             bigint,
    mam_space_free         bigint,
    assigning_org          varchar(8),
    formatted_density_code bigint,
    init_count             bigint,
    vol_change_ref         bigint,
    ttl_mbytes_life_write  bigint,
    ttl_mbytes_life_read   bigint,
    barcode                varchar(32)
);

comment on column submission.tape_alert_flags is '2 TAPEALERT FLAGS (binary, 8 bytes, read-only)';

comment on column submission.load_count is 'LOAD COUNT (binary, 8 bytes, read-only)';

comment on column submission.mam_space_free is '4 MAM SPACE REMAINING (binary, 8 bytes, read-only)';

comment on column submission.assigning_org is '5 ASSIGNING ORGANIZATION (ascii, 8 bytes, read-only)';

comment on column submission.formatted_density_code is '6 FORMATTED DENSITY CODE (binary, 1 bytes, read-only)';

comment on column submission.init_count is '7 INITIALIZATION COUNT (binary, 2 bytes, read-only)';

comment on column submission.vol_change_ref is '9 VOLUME CHANGE REFERENCE (binary, 4 bytes, read-only)';

comment on column submission.ttl_mbytes_life_write is '220 TOTAL MBYTES WRITTEN IN MEDIUM LIFE (binary, 8 bytes, read-only)';

comment on column submission.ttl_mbytes_life_read is '221 TOTAL MBYTES READ IN MEDIUM LIFE (binary, 8 bytes, read-only)';

comment on column submission.barcode is '806 BARCODE (ascii, 32 bytes, read-write)';

alter table submission
    owner to postgres;

create index table_name_tape_id_created_index
    on submission (tape_id, created);

create index table_name_tape_id_index
    on submission (tape_id);

-- Tie modified trigger to submission
CREATE TRIGGER trigger_submission_set_modified
    BEFORE UPDATE
    ON submission
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
