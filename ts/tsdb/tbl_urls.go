package tsdb

type URLData struct {
}

type URLs struct {
	*Ident
	tableName struct{} `pg:"urls,discard_unknown_columns"`
	OnID      int64    `pg:"on_id,notnull"`
	On        *Ident   `pg:"rel:has-one,fk:on_id"`
	URLType   string   `pg:"utype,notnull"`
	URL       string   `pg:"url,notnull"`
	Data      *URLData `pg:"data,notnull"`
	Active    bool     `pg:"active,notnull"`
	External  bool     `pg:"external,notnull"`
}
