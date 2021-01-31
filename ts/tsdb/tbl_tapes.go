package tsdb

import "time"

type Tape struct {
	*Ident
	tableName      struct{}     `pg:"tapes,discard_unknown_columns"`
	AccountID      int64        `pg:"account_id,notnull"`
	Account        *Account     `pg:"fk:account_id"`
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
