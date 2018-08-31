package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"sso-v2/cmds"
	"sso-v2/configs"
)

type BaseModel struct {
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type MySql struct {
	DB *sql.DB

}


var globalMysql *MySql

// 初始化表格的sql语句集合
var tables = []string{`
CREATE TABLE IF NOT EXISTS Account(
   	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
   	username VARCHAR(128) not null unique,
	password VARCHAR(128) not null ,
	nickname VARCHAR(128),
	sex TINYINT,
	email VARCHAR(128),
   	created_at bigint  NOT NULL,
   	updated_at bigint  NOT NULL,
	INDEX username_index(username),
   	PRIMARY KEY (id)
);`,
	`
CREATE TABLE IF NOT EXISTS Ticket(
   	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	ticket VARCHAR(128) NOT NULL UNIQUE,
  	created_at bigint  NOT NULL,
   	updated_at bigint  NOT NULL,
	ttl INT,
	user_id INT UNSIGNED NOT NULL,
	is_verify tinyint NOT NULL DEFAULT 0,
	INDEX ticket_index(ticket),
	FOREIGN KEY (user_id) REFERENCES Account(id),
   	PRIMARY KEY (id)
);`}



// 数据库如果不存在就创建数据库
func InitialDatabase() {
	MysqlUrl := configs.GetConfig().GetMySqlURL()
	db, err := sql.Open("mysql", MysqlUrl)
	if err != nil {
		log.Fatal("sql.Open", err.Error())
	}
	defer db.Close()

	createDatabaseSql := fmt.Sprintf("create database if not exists `sso_%s` character set UTF8", cmds.GetEnv())
	if _, err := db.Exec(createDatabaseSql); err != nil {
		log.Fatal("db.Exec(mkdirDatabase)", err.Error())
	} else {
		log.Info("create database successful...")
	}
}

// 初始化表
func InitialTable() {
	MysqlUrl := configs.GetConfig().GetDatabaseURL()
	db, err := sql.Open("mysql", MysqlUrl)
	if err != nil {
		log.Fatal("sql.Open", err.Error())
	}
	defer db.Close()
	for _, sql := range tables {
		if _, err := db.Exec(sql); err != nil {
			log.Fatal("db.Exec(sql)", err.Error())
			continue
		}
		log.Info(sql[:45], ".... successful")
	}
}

// 初始化数据库
func NewMysql() (*sql.DB, error) {
	// 初始化数据库
	InitialDatabase()
	MysqlUrl := configs.GetConfig().GetDatabaseURL()
	log.Info("connect MySQL addr is: ", MysqlUrl)

	db, err := sql.Open("mysql", MysqlUrl)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	// 初始化表
	InitialTable()
	globalMysql = new(MySql)
	globalMysql.DB = db
	return db, nil
}

func GetMySQLDB() *sql.DB{

	if globalMysql != nil {

		return globalMysql.DB
	}else {

		db,err:=  NewMysql()
		if err != nil {
			log.Fatal("invalid MySQL DB ...")
		}
		return db
	}
}