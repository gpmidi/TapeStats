package ts

import (
	"crypto/rand"
	"errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"math/big"
)

type PasswordValidator interface {
	SetPasswordHash(hash string) error
	GetPasswordHash() (string, error)
}

func CreatePassword() (string, error) {
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
		// A-Z || a-z
		if (n >= 65 && n <= 90) || (n >= 97 && n <= 122) {
			result += string(n)
		}
	}
}

func CreateSetPassword(a PasswordValidator) (string, error) {
	password, err := CreatePassword()
	if err != nil {
		return "", err
	}
	if err := SetPassword(a, password); err != nil {
		return password, err // Return password in case the change does somehow happen
	}
	return password, nil
}

func SetPassword(a PasswordValidator, passwd string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(([]byte)(passwd), bcrypt.DefaultCost+viper.GetInt("passwords.extracost"))
	if err != nil {
		return err
	}

	if err := a.SetPasswordHash((string)(hashedPassword)); err != nil {
		return err
	}

	return nil
}

func VerifyPassword(a PasswordValidator, passwd string) (bool, error) {
	passwdHash, err := a.GetPasswordHash()
	if err != nil {
		return false, nil
	}
	if passwd == "" || passwdHash == "" {
		// Shouldn't actually happen as GetPasswordHash "should" return an error on empty passwords
		return false, errors.New("password and/or Hashed can't be empty")
	}

	if err := bcrypt.CompareHashAndPassword(([]byte)(passwdHash), ([]byte)(passwd)); err != nil {
		return false, err
	}
	return true, nil
}

func init() {
	viper.SetDefault("passwords.extracost", 2)
	viper.SetDefault("passwords.length", 32)
}
