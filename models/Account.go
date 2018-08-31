package models

import (
	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	BaseModel
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Sex      uint8  `json:"sex"`
	Email    string `json:"email"`
}

func GetAccount(sql string, args ...interface{}) (*Account, error) {
	if stmt, err := GetMySQLDB().Prepare(sql); err != nil {
		return nil, err
	} else {
		res := stmt.QueryRow(args...)
		var account Account

		if err := res.Scan(&account.Username, &account.Password); err != nil {
			return nil, err
		} else {
			return &account, nil
		}
	}
}

// 比较密码
func (account Account) HasPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(plain))
	if err != nil {
		return false
	} else {
		return true
	}
}
