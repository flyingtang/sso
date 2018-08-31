package configs

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"sso-v2/cmds"
	"sso-v2/constants"
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions"
)

type MySQL struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type Config struct {
	MySQL    MySQL
	HttpAddr string
	HttpPort string
	SessionSecret string
	SessionName string
}

var config *Config

// 获取配置文件
func NewConfig(env string) *Config {
	//env := cmd.GetEnv()

	cmd, err := cmds.GetCmd()
	if err != nil {
		log.Fatal("NewConfig error ", err.Error())
	}
	configpath := cmd.GetConfigPath(constants.CONFIGDIR)
	data, err := ioutil.ReadFile(configpath)
	if err != nil && err != io.EOF {
		log.Fatal("read config file error ", err.Error)
	}

	config = new(Config)
	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatal("transform config filr error ", err.Error())
	}
	return config
}

// 获取http监听地址和端口URL
func (config *Config) GetHttpAddrPort() (ap string) {

	if config != nil && config.HttpAddr != "" {
		ap += config.HttpAddr
	} else {
		ap += "0.0.0.0"
	}

	ap += ":"
	if config != nil && config.HttpPort != "" {
		ap += config.HttpPort
	} else {
		ap += "3000"
	}
	return ap
}

// 获取数据库的链接URL 主要用于常规增删改查
func (config Config) GetDatabaseURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.MySQL.Username, config.MySQL.Password, config.MySQL.Host, config.MySQL.Port, config.MySQL.Database)
}

// 获取MySQL的URL，主要用于创建数据库
func (config Config) GetMySqlURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/", config.MySQL.Username, config.MySQL.Password, config.MySQL.Host, config.MySQL.Port)
}

// 获取config对象
func GetConfig()(  *Config){
	if config == nil {
		return NewConfig(cmds.GetEnv())
	}else {
		return config
	}
}

// 设置Session
func (config Config)SetSession(engine *gin.Engine){

	store := cookie.NewStore([]byte(config.SessionName))
	engine.Use(sessions.Sessions(config.SessionName, store))
	return
}