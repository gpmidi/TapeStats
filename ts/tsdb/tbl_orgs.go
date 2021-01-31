package tsdb

type Org struct {
	*Ident
	tableName   struct{}   `pg:"orgs,discard_unknown_columns"`
	Name        string     `pg:"name,notnull,unique"`
	Description string     `pg:"description,notnull"`
	Active      bool       `pg:"active"`
	Accounts    []*Account `pg:"join_fk:org_id"`
	Users       []*User    `pg:"join_fk:org_id"`
}
