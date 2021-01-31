package tsdb

import (
	"errors"
	"github.com/gpmidi/TapeStats/ts"
)

type Account struct {
	*Ident
	tableName   struct{} `pg:"accounts,discard_unknown_columns"`
	Name        string   `pg:"name,notnull,unique"`
	Description string   `pg:"description,notnull"`
	Active      bool     `pg:"active"`
	OrgID       int64    `pg:"org_id,notnull"`
	Org         *Org     `pg:"fk:org_id"`
	PasswdHash  string   `pg:"passwd,notnull"`
	Tapes       []*Tape  `pg:"join_fk:account_id"`
}

func (a *Account) SetPasswordHash(hash string) error {
	a.PasswdHash = hash
	return nil
}

func (a *Account) GetPasswordHash() (string, error) {
	if a.PasswdHash == "" {
		return "", errors.New("empty password hash")
	}
	return a.PasswdHash, nil
}

func (a *Account) CreateSetPassword() (string, error) {
	return ts.CreateSetPassword(a)
}

func (a *Account) SetPassword(passwd string) error {
	return ts.SetPassword(a, passwd)
}

func (a *Account) VerifyPassword(passwd string) (bool, error) {
	return ts.VerifyPassword(a, passwd)
}
