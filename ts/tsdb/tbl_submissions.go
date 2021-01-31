package tsdb

import (
	"github.com/gpmidi/TapeStats/ts/mam"
)

type Submission struct {
	*Ident
	tableName            struct{}             `pg:"submissions,discard_unknown_columns"`
	TapeID               int64                `pg:"tape_id,notnull"`
	Tape                 *Tape                `pg:"fk:tape_id"`
	TapeAlertFlags       string               `pg:"tape_alert_flags"`
	LoadCount            int64                `pg:"load_count"`
	MAMSpaceFree         int64                `pg:"mam_space_free"`
	AssigningOrg         string               `pg:"assigning_org"`
	FormattedDensityCode int64                `pg:"formatted_density_code"`
	InitCount            int64                `pg:"init_count"`
	VolChangeRef         int64                `pg:"vol_change_ref"`
	TotalMBytesLifeWrite int64                `pg:"ttl_mbytes_life_write"`
	TotalMBytesLifeRead  int64                `pg:"ttl_mbytes_life_read"`
	Barcode              string               `pg:"barcode"`
	KVS                  map[string]string    `pg:"kvs,hstore"`
	RequestID            string               `pg:"request_id"`
	Raw                  *RawSubmission       `pg:"raw"` // JSONB
	SubmittedByID        int64                `pg:"submitted_by_id,notnull"`
	SubmittedBy          *RemoteSystem        `pg:"fk:submitted_by_id"`
	ParserUsedID         int64                `pg:"parser_used_id"`
	ParserUsed           *ParserVersion       `pg:"fk:parser_used_id"`
	ParserUsedRunID      int64                `pg:"parser_used_run_id"`
	ParserUsedRun        *SubmissionParsing   `pg:"fk:parser_used_run_id"`
	SubmissionParsings   []*SubmissionParsing `pg:"join_fk:submission_id"`
}

type RawSubmission struct {
	// Stuff they sent us
	GETArgs  map[string]string
	POSTArgs map[string]string
	Files    map[string]string
	// Stuff we figured out
	Fields map[string]*mam.Field // Fields we parsed this into
}
