package tsdb

import (
	"errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Account struct {
	tableName struct{}  `pg:"accounts,discard_unknown_columns"`
	Id        string    `pg:"id,pk,type:uuid,default:gen_random_uuid()"`
	created   time.Time `pg:",notnull"`
	modified  time.Time `pg:",notnull"`
	hashed    string
}

func (a *Account) SetPassword(passwd string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(([]byte)(passwd), bcrypt.DefaultCost+viper.GetInt("passwords.extracost"))
	if err != nil {
		return err
	}

	a.hashed = (string)(hashedPassword)

	return nil
}

func (a *Account) VerifyPassword(passwd string) (bool, error) {
	if passwd == "" || a.hashed == "" {
		return false, errors.New("password and/or hashed can't be empty")
	}
	err := bcrypt.CompareHashAndPassword(([]byte)(a.hashed), ([]byte)(passwd))
	if err != nil {
		return false, err
	}
	return true, nil
}

func init() {
	viper.SetDefault("passwords.extracost", 8)
}
