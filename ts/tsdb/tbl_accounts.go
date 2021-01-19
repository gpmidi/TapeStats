package tsdb

import (
	"crypto/rand"
	"errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"time"
)

type Account struct {
	tableName struct{}  `pg:"accounts,discard_unknown_columns"`
	Id        string    `pg:"id,pk,type:uuid,default:gen_random_uuid()"`
	Created   time.Time `pg:"created,notnull"`
	Modified  time.Time `pg:"modified,notnull"`
	Hashed    string    `pg:"hashed"`
}

func (a *Account) CreatePassword() (string, error) {
	result := ""
	for {
		if len(result) >= viper.GetInt("passwords.length") {
			return result, nil
		}
		num, err := rand.Int(rand.Reader, big.NewInt(int64(127)))
		if err != nil {
			return "", err
		}
		n := num.Int64()
		// Make sure that the number/byte/letter is inside
		// the range of printable ASCII characters (excluding space and DEL)
		if n > 32 && n < 127 {
			result += string(n)
		}
	}
}
func (a *Account) CreateSetPassword() (string, error) {
	password, err := a.CreatePassword()
	if err != nil {
		return "", err
	}
	if err := a.SetPassword(password); err != nil {
		return password, err // Return password in case the change does somehow happen
	}
	return password, nil
}

func (a *Account) SetPassword(passwd string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(([]byte)(passwd), bcrypt.DefaultCost+viper.GetInt("passwords.extracost"))
	if err != nil {
		return err
	}

	a.Hashed = (string)(hashedPassword)

	return nil
}

func (a *Account) VerifyPassword(passwd string) (bool, error) {
	if passwd == "" || a.Hashed == "" {
		return false, errors.New("password and/or Hashed can't be empty")
	}
	err := bcrypt.CompareHashAndPassword(([]byte)(a.Hashed), ([]byte)(passwd))
	if err != nil {
		return false, err
	}
	return true, nil
}

func init() {
	viper.SetDefault("passwords.extracost", 8)
	viper.SetDefault("passwords.length", 16)
}
