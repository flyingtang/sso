package cmds

import (
	"flag"
	"fmt"
	"strings"

	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type conmand struct {
	Env string
}

var cmd *conmand

func ParseCmd() *conmand {
	var env string
	flag.StringVar(&env, "env", "dev", "-env=pro|dev|test")
	flag.Parse()
	parseEnv(env)
	return cmd
}

// 解析前的环境变量
func parseEnv(env string) *conmand {
	cmd = new(conmand)
	switch env {
	case "dev":
		cmd.Env = "development"
	case "pro":
		cmd.Env = "production"
	case "test":
		cmd.Env = "test"
	}
	log.Info("current environment is ", cmd.Env)
	if cmd.Env == "" {
		log.Fatal("invalid env parameter, please check ...")
	}
	return cmd
}

func GetCmd() (*conmand, error) {
	if cmd != nil {
		return cmd, nil
	} else {
		return nil, errors.New("current cmd is nil")
	}
}

// 获取当前的环境变量
func  GetEnv() string {
	if cmd.Env == "" {
		log.Fatal("no invalid  environment parameter, please check ...")
	}
	return cmd.Env
}

// 获取config 文件路径
func (cmd *conmand) GetConfigPath(dir string) string {
	env := cmd.Env
	fileName := fmt.Sprintf("%s.config.json", env)
	return strings.Join([]string{dir, fileName}, "/")
}

func (cmd *conmand) SetMode() {
	env := cmd.Env
	switch env {
		case "production":
			gin.SetMode(gin.ReleaseMode)
		case "test":
			gin.SetMode(gin.TestMode)
		default:
			gin.SetMode(gin.DebugMode)
	}
}
