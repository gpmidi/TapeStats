package tsdb

import "time"

type Ident struct {
	tableName    struct{}          `pg:"ident,discard_unknown_columns"`
	Id           int64             `pg:"id,pk"`
	Guid         string            `pg:"guid,unique,notnull,type:uuid,default:gen_random_uuid()"`
	Created      time.Time         `pg:"created,notnull"`
	Modified     time.Time         `pg:"modified,notnull"`
	RType        string            `pg:"rtype,notnull"`
	InternalName string            `pg:"internal_name,notnull"`
	InternalDesc string            `pg:"internal_description,notnull"`
	Attr         map[string]string `pg:"attr,hstore,notnull"`
	Tags         []string          `pg:"tags,array,notnull"`
	MyIdent      *Ident            `pg:"fk:id,join_fk:id"`
	Annotations  []*Annotation     `pg:"join_fk:on_id"`
	Comments     []*Comment        `pg:"join_fk:on_id"`
	URLs         []*URLs           `pg:"join_fk:on_id"`
}
