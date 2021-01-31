drop type ghstore;

drop table gopg_migrations;

drop table accounts;

alter table submissions
    drop constraint submissions_submission_parsing_id_fk;

drop table submission_parsing;

drop type parsing_status;

drop table submissions;

drop table tapes;

drop table remote_systems;

drop table parser_versions;

drop table comments;

drop type comment_status;

drop table users;

drop table orgs;

drop table annotations;

drop type annotation_type;

drop table parsers;

drop type parsers_status;

drop table urls;

drop type url_type;

drop table ident;

drop type hstore;

drop type identrowtypes;

drop view tags;

drop function trigger_set_timestamp();

DROP EXTENSION "uuid-ossp";
DROP EXTENSION "hstore";
