package tsdb

import (
	"errors"
	"github.com/gpmidi/TapeStats/ts"
)

type User struct {
	*Ident
	tableName  struct{}   `pg:"users,discard_unknown_columns"`
	OrgID      int64      `pg:"org_id,notnull"`
	Org        *Org       `pg:"fk:org_id"`
	Username   string     `pg:"username,notnull"`
	FullName   string     `pg:"full_name,notnull"`
	RemoteID   string     `pg:"remote_id"`
	Comments   []*Comment `pg:"join_fk:by_id"`
	PasswdHash string     `pg:"passwd,notnull"`
	Active     bool       `pg:"active,notnull"`
}

func (u *User) SetPasswordHash(hash string) error {
	u.PasswdHash = hash
	return nil
}

func (u *User) GetPasswordHash() (string, error) {
	if u.PasswdHash == "" {
		return "", errors.New("empty password hash")
	}
	return u.PasswdHash, nil
}

func (u *User) CreateSetPassword() (string, error) {
	return ts.CreateSetPassword(u)
}

func (u *User) SetPassword(passwd string) error {
	return ts.SetPassword(u, passwd)
}

func (u *User) VerifyPassword(passwd string) (bool, error) {
	return ts.VerifyPassword(u, passwd)
}
