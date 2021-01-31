package tsdb

type SubmissionParsing struct {
	*Ident
	tableName       struct{}       `pg:"submission_parsing,discard_unknown_columns"`
	SubmissionID    int64          `pg:"submission_id,notnull,unique:subverrun"`
	Submission      *Submission    `pg:"fk:submission_id"`
	ParserVersionID int64          `pg:"parser_version_id,notnull,unique:subverrun"`
	ParserVersion   *ParserVersion `pg:"fk:parser_version_id"`
	RunNumber       int64          `pg:"run,notnull,unique:subverrun"`
	ResultOk        bool           `pg:"result_ok"`
	ResultMessage   string         `pg:"result_message"`
	ResultCode      int32          `pg:"result_code"`
	Status          string         `pg:"status,notnull"`
	Submissions     []*Submission  `pg:"join_fk:parser_used_run_id"`
}
