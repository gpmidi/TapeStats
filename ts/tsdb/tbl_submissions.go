package tsdb

import (
	"github.com/gpmidi/TapeStats/ts/mam"
	"net"
	"time"
)

type Submission struct {
	tableName            struct{}          `pg:"submissions,discard_unknown_columns"`
	Id                   int64             `pg:"id,pk"`
	TapeID               string            `pg:"tape_id,type:uuid,notnull"`
	Tape                 *Tape             `pg:"rel:has-one"`
	Created              time.Time         `pg:"created,notnull"`
	Modified             time.Time         `pg:"modified,notnull"`
	TapeAlertFlags       string            `pg:"tape_alert_flags"`
	LoadCount            int64             `pg:"load_count"`
	MAMSpaceFree         int64             `pg:"mam_space_free"`
	AssigningOrg         string            `pg:"assigning_org"`
	FormattedDensityCode int64             `pg:"formatted_density_code"`
	InitCount            int64             `pg:"init_count"`
	VolChangeRef         int64             `pg:"vol_change_ref"`
	TotalMBytesLifeWrite int64             `pg:"ttl_mbytes_life_write"`
	TotalMBytesLifeRead  int64             `pg:"ttl_mbytes_life_read"`
	Barcode              string            `pg:"barcode"`
	Raw                  *RawSubmission    `pg:"raw"`
	KVS                  map[string]string `pg:"kvs,hstore"`
	Requester            map[string]string `pg:"requester,hstore"`
	RequesterIP          net.IP            `pg:"requester_ip"`
	RequestID            string            `pg:"request_id"`
}

type RawSubmission struct {
	// Stuff they sent us
	GETArgs  map[string]string
	POSTArgs map[string]string
	Files    map[string]string
	// Stuff we figured out
	Fields map[string]*mam.Field // Fields we parsed this into
}
