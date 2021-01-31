CREATE EXTENSION hstore;
CREATE EXTENSION "uuid-ossp";

create type identrowtypes as enum ('accounts', 'ident', 'submissions', 'tapes', 'orgs', 'remote_systems', 'parser_versions', 'submission_parsing', 'comments', 'users', 'annotations', 'parsers', 'urls');

create type parsing_status as enum ('created', 'pending', 'running', 'completed', 'failed', 'new');

create type comment_status as enum ('draft', 'posted', 'removed');

create type annotation_type as enum ('unknown', 'human', 'manual');

create type parsers_status as enum ('active', 'banned', 'inactive');

create type url_type as enum ('unknown', 'rest-api', 'user-provided', 'documentation', 'code-repo', 'home');

create table if not exists ident
(
    id                   bigserial                                                      not null
        constraint ident_pk
            primary key,
    guid                 uuid                     default uuid_generate_v4()            not null,
    created              timestamp with time zone default now()                         not null,
    modified             timestamp with time zone default now()                         not null,
    rtype                identrowtypes                                                  not null,
    internal_name        varchar(1024)            default ''::character varying         not null,
    internal_description text                     default ''::text                      not null,
    attr                 hstore                   default ''::hstore                    not null,
    tags                 varchar(1024)[]          default ARRAY []::character varying[] not null
);

comment on table ident is 'GUIDs';

create unique index if not exists ident_guid_uindex
    on ident (guid);

create table if not exists orgs
(
    name        varchar(1024)            not null,
    description text    default ''::text not null,
    active      boolean default true     not null,
    constraint orgs_pk
        primary key (id)
)
    inherits (ident);

comment on table orgs is 'Org for managing users, data, and accounts';

create unique index if not exists org_name_uindex
    on orgs (name);

create unique index if not exists orgs_guid_uindex
    on orgs (guid);

create table if not exists accounts
(
    name        varchar(1024)            not null,
    description text    default ''::text not null,
    active      boolean default true,
    org_id      bigint                   not null
        constraint accounts_orgs_id_fk
            references orgs
            on update cascade on delete cascade,
    passwd      varchar(1024)            not null,
    constraint accounts_pk
        primary key (id)
)
    inherits (ident);

comment on table accounts is 'Accounts can submit data';

comment on column accounts.passwd is 'Hashed & Salted';

create unique index if not exists accounts_name_uindex
    on accounts (name);

create index if not exists accounts_org_id_active_index
    on accounts (org_id, active);

create unique index if not exists accounts_guid_uindex
    on accounts (guid);

create table if not exists tapes
(
    account_id       uuid         not null,
    uci              varchar(255),
    alt_uci          varchar(255),
    serial_number    varchar(255) not null,
    assigning_org    varchar(255),
    manufacturer     varchar(255),
    manufacture_dt   date,
    density_code     varchar(255),
    medium_type      varchar(255),
    medium_type_info varchar(255),
    lto_version      integer,
    constraint tapes_pk
        primary key (id)
)
    inherits (ident);

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

create unique index if not exists tapes_account_lookup
    on tapes (account_id, manufacturer, manufacture_dt, serial_number, density_code, medium_type, lto_version);

create unique index if not exists tapes_guid_uindex
    on tapes (guid);

create table if not exists remote_systems
(
    remote_ip inet,
    headers   hstore default ''::hstore,
    constraint remote_systems_pk
        primary key (id)
)
    inherits (ident);

comment on table remote_systems is 'Something that submitted data';

create unique index if not exists remote_systems_guid_uindex
    on remote_systems (guid);

create unique index if not exists remote_systems_id_uindex
    on remote_systems (id);

create index if not exists remote_systems_remote_ip_index
    on remote_systems (remote_ip);

create table if not exists comments
(
    subject      varchar(1024)                                   not null,
    comment_text text                                            not null,
    by_id        bigint                                          not null,
    on_id        bigint                                          not null
        constraint comments_ident_id_fk
            references ident
            on update cascade on delete cascade,
    status       comment_status default 'posted'::comment_status not null,
    constraint comments_pkey
        primary key (id),
    constraint comments_guid_key
        unique (guid)
)
    inherits (ident);

comment on table comments is 'Comments or notes on an entry';

create index if not exists comments_on_id_created_index
    on comments (on_id asc, created desc);

create index if not exists comments_on_id_status_index
    on comments (on_id, status);

create table if not exists users
(
    org_id    bigint               not null
        constraint users_orgs_id_fk
            references orgs
            on update cascade on delete cascade,
    username  varchar(255)         not null,
    full_name varchar(1024)        not null,
    remote_id varchar(1024),
    passwd    varchar(1024)        not null,
    active    boolean default true not null,
    constraint users_pkey
        primary key (id),
    constraint users_guid_key
        unique (guid)
)
    inherits (ident);

comment on table users is 'A human-bot-person-thing-login';

comment on column users.full_name is 'Full name';

comment on column users.remote_id is 'For distributed auth systems';

comment on column users.passwd is 'User''s hashed & salted password';

create unique index if not exists users_org_id_username_uindex
    on users (org_id, username);

create table if not exists annotations
(
    on_id bigint          not null
        constraint annotations_ident_id_fk
            references ident
            on update cascade on delete cascade,
    atype annotation_type not null,
    data  jsonb           not null,
    constraint annotations_pkey
        primary key (id),
    constraint annotations_guid_key
        unique (guid)
)
    inherits (ident);

comment on table annotations is 'Automation driver notes on a row of data';

create index if not exists annotations_atype_data_index
    on annotations (atype, data);

create index if not exists annotations_on_id_atype_created_index
    on annotations (on_id, atype, created);

create table if not exists parsers
(
    name        varchar(1024)                                   not null,
    description text           default ''::text                 not null,
    status      parsers_status default 'active'::parsers_status not null,
    urls        hstore         default ''::hstore               not null,
    settings    jsonb          default '{}'::json               not null,
    constraint parsers_pkey
        primary key (id),
    constraint parsers_guid_key
        unique (guid)
)
    inherits (ident);

comment on table parsers is 'Parsers';

create table if not exists parser_versions
(
    tool_id     bigint               not null
        constraint parser_versions_parsers_id_fk
            references parsers
            on update cascade on delete cascade,
    ver         varchar(255)         not null,
    uses        bigint  default 0    not null,
    active      boolean default true not null,
    tool_sha512 varchar(1024),
    tool_path   varchar(8192),
    constraint parser_versions_pkey
        primary key (id),
    constraint parser_versions_guid_key
        unique (guid)
)
    inherits (ident);

comment on table parser_versions is 'A version of the MAM parser ';

create table if not exists submissions
(
    tape_id                bigint not null
        constraint table_name_tapes_id_fk
            references tapes
            on update cascade on delete cascade,
    tape_alert_flags       varchar(255),
    load_count             bigint,
    mam_space_free         bigint,
    assigning_org          varchar(255),
    formatted_density_code bigint,
    init_count             bigint,
    vol_change_ref         bigint,
    ttl_mbytes_life_write  bigint,
    ttl_mbytes_life_read   bigint,
    barcode                varchar(255),
    kvs                    hstore,
    request_id             varchar(255),
    raw                    jsonb,
    submitted_by_id        bigint not null
        constraint submissions_remote_systems_id_fk
            references remote_systems
            on update cascade on delete cascade,
    parser_used_id         bigint
        constraint submissions_parser_versions_id_fk
            references parser_versions
            on update cascade on delete cascade,
    parser_used_run_id     bigint,
    constraint submissions_pk
        primary key (id)
)
    inherits (ident);

comment on table submissions is 'Data submission by someone';

comment on column submissions.tape_alert_flags is '2 TAPEALERT FLAGS (binary, 8 bytes, read-only)';

comment on column submissions.load_count is 'LOAD COUNT (binary, 8 bytes, read-only)';

comment on column submissions.mam_space_free is '4 MAM SPACE REMAINING (binary, 8 bytes, read-only)';

comment on column submissions.assigning_org is '5 ASSIGNING ORGANIZATION (ascii, 8 bytes, read-only)';

comment on column submissions.formatted_density_code is '6 FORMATTED DENSITY CODE (binary, 1 bytes, read-only)';

comment on column submissions.init_count is '7 INITIALIZATION COUNT (binary, 2 bytes, read-only)';

comment on column submissions.vol_change_ref is '9 VOLUME CHANGE REFERENCE (binary, 4 bytes, read-only)';

comment on column submissions.ttl_mbytes_life_write is '220 TOTAL MBYTES WRITTEN IN MEDIUM LIFE (binary, 8 bytes, read-only)';

comment on column submissions.ttl_mbytes_life_read is '221 TOTAL MBYTES READ IN MEDIUM LIFE (binary, 8 bytes, read-only)';

comment on column submissions.barcode is '806 BARCODE (ascii, 32 bytes, read-write)';

comment on column submissions.kvs is 'Submitted Key:Values (string form)';

comment on column submissions.submitted_by_id is 'Info about who submitted the data';

create index if not exists table_name_tape_id_index
    on submissions (tape_id);

create index if not exists table_name_tape_id_created_index
    on submissions (tape_id, created);

create unique index if not exists submissions_guid_uindex
    on submissions (guid);

create index if not exists submissions_barcode_index
    on submissions (barcode);

create unique index if not exists parser_versions_tool_version_uindex
    on parser_versions (tool_id, ver);

create table if not exists submission_parsing
(
    submission_id     bigint                                       not null
        constraint submission_parsing_submissions_id_fk
            references submissions
            on update cascade on delete cascade,
    parser_version_id bigint                                       not null
        constraint submission_parsing_parser_versions_id_fk
            references parser_versions
            on update cascade on delete cascade,
    run               bigint         default 0                     not null,
    result_ok         boolean,
    result_message    text,
    result_code       integer,
    status            parsing_status default 'new'::parsing_status not null,
    constraint submission_parsing_pkey
        primary key (id),
    constraint submission_parsing_guid_key
        unique (guid)
)
    inherits (ident);

comment on table submission_parsing is 'Parsing history for Submissions';

alter table submissions
    add constraint submissions_submission_parsing_id_fk
        foreign key (parser_used_run_id) references submission_parsing
            on update cascade on delete cascade;

create index if not exists submission_parsing_parser_version_id_index
    on submission_parsing (parser_version_id);

create index if not exists submission_parsing_submission_id_index
    on submission_parsing (submission_id);

create unique index if not exists submission_parsing_submission_id_parser_version_id_run_uindex
    on submission_parsing (submission_id, parser_version_id, run);

create unique index if not exists parsers_name_uindex
    on parsers (name);

create table if not exists urls
(
    on_id    bigint                      not null
        constraint urls_ident_id_fk
            references ident
            on update cascade on delete cascade,
    utype    url_type                    not null,
    url      text                        not null,
    data     jsonb   default '{}'::jsonb not null,
    active   boolean default true        not null,
    external boolean default false       not null,
    constraint urls_pkey
        primary key (id),
    constraint urls_guid_key
        unique (guid)
)
    inherits (ident);

comment on table urls is 'URLs for a row of data';

comment on column urls.data is 'Data to submit with URL';

comment on column urls.active is 'URL should work';

comment on column urls.external is 'Safe for external users to see';

create index if not exists urls_utype_url_index
    on urls (utype, url);

create index if not exists urls_on_id_utype_created_index
    on urls (on_id, utype, created);

create or replace view tags(tag, array_agg) as
SELECT utags.tag,
       array_agg(i.rtype) AS array_agg
FROM (SELECT DISTINCT unnest(ident.tags) AS tag
      FROM ident) utags
         JOIN ident i ON utags.tag::text = ANY (i.tags::text[])
GROUP BY utags.tag;

comment on view tags is 'All known tags';

create or replace function trigger_set_timestamp() returns trigger
    language plpgsql
as
$$
BEGIN
    NEW.modified = NOW();
    RETURN NEW;
END;
$$;

create trigger trigger_ident_set_modified
    before update
    on ident
    for each row
execute procedure trigger_set_timestamp();

create trigger trigger_orgs_set_modified
    before update
    on orgs
    for each row
execute procedure trigger_set_timestamp();

create trigger trigger_accounts_set_modified
    before update
    on accounts
    for each row
execute procedure trigger_set_timestamp();

create trigger trigger_tapes_set_modified
    before update
    on tapes
    for each row
execute procedure trigger_set_timestamp();

create trigger trigger_submissions_set_modified
    before update
    on submissions
    for each row
execute procedure trigger_set_timestamp();

create trigger trigger_remote_systems_set_modified
    before update
    on remote_systems
    for each row
execute procedure trigger_set_timestamp();

create trigger trigger_parser_versions_set_modified
    before update
    on parser_versions
    for each row
execute procedure trigger_set_timestamp();

create trigger trigger_submission_parsing_set_modified
    before update
    on submission_parsing
    for each row
execute procedure trigger_set_timestamp();

create trigger trigger_remote_systems_set_modified
    before update
    on comments
    for each row
execute procedure trigger_set_timestamp();

create trigger trigger_users_set_modified
    before update
    on users
    for each row
execute procedure trigger_set_timestamp();

create trigger trigger_annotations_set_modified
    before update
    on annotations
    for each row
execute procedure trigger_set_timestamp();

create trigger trigger_parsers_set_modified
    before update
    on parsers
    for each row
execute procedure trigger_set_timestamp();

create trigger trigger_urls_set_modified
    before update
    on urls
    for each row
execute procedure trigger_set_timestamp();
