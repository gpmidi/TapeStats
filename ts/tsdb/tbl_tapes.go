package tsdb

import "time"

type Tape struct {
	tableName      struct{}     `pg:"tapes,discard_unknown_columns"`
	Id             string       `pg:"id,pk,type:uuid,default:gen_random_uuid()"`
	AccountID      string       `pg:"account_id,type:uuid,notnull"`
	Account        *Account     `pg:"rel:has-one"`
	Created        time.Time    `pg:"created,notnull"`
	Modified       time.Time    `pg:"modified,notnull"`
	UCI            string       `pg:"uci"`
	AltUCI         string       `pg:"alt_uci"`
	SerialNumber   string       `pg:"serial_number,notnull"`
	AssignOrg      string       `pg:"assigning_org"`
	Manufacture    string       `pg:"manufacturer"`
	ManufactureDT  time.Time    `pg:"manufacture_dt"`
	DensityCode    string       `pg:"density_code"`
	MediumType     string       `pg:"medium_type"`
	MediumTypeInfo string       `pg:"medium_type_info"`
	LTOVersion     int          `pg:"lto_version"`
	Submissions    []Submission `pg:"rel:has-many,join_fk:tape_id"`
}
