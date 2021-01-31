package tsdb

import "net"

type RemoteSystem struct {
	*Ident
	tableName   struct{}          `pg:"remote_systems,discard_unknown_columns"`
	IP          net.IP            `pg:"remote_ip"`
	Headers     map[string]string `pg:"headers,hstore,notnull"`
	Submissions []*Submission     `pg:"join_fk:submitted_by_id"`
}
