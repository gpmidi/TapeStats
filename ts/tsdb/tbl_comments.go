package tsdb

type Comment struct {
	*Ident
	tableName struct{} `pg:"comments,discard_unknown_columns"`
	Subject   string   `pg:"subject,notnull"`
	Text      string   `pg:"comment_text,notnull"`
	OnID      int64    `pg:"on_id,notnull"`
	On        *Ident   `pg:"rel:has-one,fk:on_id"`
	ByID      int64    `pg:"by_id,notnull"`
	By        *User    `pg:"fk:by_id"`
}
