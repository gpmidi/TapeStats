package tsdb

type ParserVersion struct {
	*Ident
	tableName          struct{}             `pg:"parser_versions,discard_unknown_columns"`
	ToolID             int64                `pg:"tool_id,notnull,unique:toolver"`
	Tool               *Parser              `pg:"fk:tool_id"`
	Version            string               `pg:"ver,notnull,unique:toolver"`
	Uses               int64                `pg:"uses,notnull"`
	Active             bool                 `pg:"active,notnull"`
	ToolSHA512         string               `pg:"tool_sha512"`
	ToolPath           string               `pg:"tool_path"`
	SubmissionParsings []*SubmissionParsing `pg:"join_fk:parser_version_id"`
	UsedToParse        []*Submission        `pg:"join_fk:parser_used_id"`
}
