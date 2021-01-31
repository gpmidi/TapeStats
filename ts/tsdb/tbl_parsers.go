package tsdb

type ParserSettings struct {
}

type Parser struct {
	*Ident
	tableName   struct{}          `pg:"parsers,discard_unknown_columns"`
	Name        string            `pg:"name,notnull,unique"`
	Description string            `pg:"description,notnull"`
	Status      string            `pg:"status,notnull"`
	URLs        map[string]string `pg:"urls,hstore,notnull"`
	Settings    *ParserSettings   `pg:"settings,notnull"`
}
