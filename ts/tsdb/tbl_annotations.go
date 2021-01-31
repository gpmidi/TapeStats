package tsdb

type AnnotationData struct {
}

type Annotation struct {
	*Ident
	tableName struct{}        `pg:"annotations,discard_unknown_columns"`
	OnID      int64           `pg:"on_id,notnull"`
	On        *Ident          `pg:"rel:has-one,fk:on_id"`
	AType     string          `pg:"atype,notnull"`
	Data      *AnnotationData `pg:"data,notnull"`
}
