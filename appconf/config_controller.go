package appconf

import (
	"fmt"
	_DB "jojopoper/NBi/StressTest/db"
	_L "jojopoper/NBi/StressTest/log"
	"strings"

	_c "github.com/astaxie/beego/config"
)

// ConfigInstance 配置文件唯一实例
var ConfigInstance *ConfigController

// ConfigController 配置文件读取控制器
type ConfigController struct {
	appConfig _c.Configer
}

// NewConfigController new ConfigController
func NewConfigController() *ConfigController {
	ret := new(ConfigController)
	return ret.Init()
}

// Init 初始化
func (ths *ConfigController) Init() *ConfigController {
	var err error
	ths.appConfig, err = _c.NewConfig("ini", "conf/app.conf")
	if err != nil {
		_L.LoggerInstance.ErrorPrint("Read 'app.conf' has error : \r\n%v\r\n", err)
		panic(err.Error())
	}
	return ths
}

// GetDateBaseConf 获取数据库相关配置
func (ths *ConfigController) GetDateBaseConf() (ret *_DB.DatabaseInfo) {
	ret = new(_DB.DatabaseInfo)
	ret.DbType = _DB.DatabaseType(ths.appConfig.String("DataBase::dbtype"))
	ret.AliasName = ths.appConfig.String("DataBase::aliasname")
	ret.Host = ths.appConfig.String("DataBase::dbHost")
	ret.Port = ths.appConfig.String("DataBase::dbPort")
	ret.UserName = ths.appConfig.String("DataBase::dbUserName")
	ret.Password = ths.appConfig.String("DataBase::dbPassword")
	ret.IsDebug, _ = ths.appConfig.Bool("DataBase::debug")
	return
}

// GetString 通用读取配置文件，返回字符串
func (ths *ConfigController) GetString(section, key string) string {
	findkey := fmt.Sprintf("%s::%s", section, key)
	if len(section) == 0 {
		findkey = fmt.Sprintf("%s", key)
	}
	return ths.appConfig.String(findkey)
}

// GetBool 通用读取配置文件，返回bool
func (ths *ConfigController) GetBool(section, key string) (bool, error) {
	findkey := fmt.Sprintf("%s::%s", section, key)
	if len(section) == 0 {
		findkey = fmt.Sprintf("%s", key)
	}
	return ths.appConfig.Bool(findkey)
}

// GetNumber 通用读取配置文件，返回INT
func (ths *ConfigController) GetNumber(section, key string) int {
	findkey := fmt.Sprintf("%s::%s", section, key)
	if len(section) == 0 {
		findkey = fmt.Sprintf("%s", key)
	}
	return ths.appConfig.DefaultInt(findkey, 0)
}

// GetFloat 通用读取配置文件，返回float64
func (ths *ConfigController) GetFloat(section, key string) float64 {
	findkey := fmt.Sprintf("%s::%s", section, key)
	if len(section) == 0 {
		findkey = fmt.Sprintf("%s", key)
	}
	return ths.appConfig.DefaultFloat(findkey, 1)
}

// GetAppName 读取APP名称
func (ths *ConfigController) GetAppName() string {
	return ths.appConfig.String("baseinfo::appname")
}

// GetHorizonServer 读取Horizon服务器地址
func (ths *ConfigController) GetHorizonServer() string {
	nw := ths.appConfig.String("network::current_network")
	key := "Network::"
	if strings.ToLower(nw) == "live" {
		key += "horizon_live"
	} else {
		key += "horizon_test"
	}
	return ths.appConfig.String(key)
}

// GetCoreServer 获取core服务器地址
func (ths *ConfigController) GetCoreServer() string {
	return ths.appConfig.String("network::stellar_core")
}

// GetHorizonHeader 读取Horizon连接HTTP头配置
func (ths *ConfigController) GetHorizonHeader() map[string]string {
	headers, err := ths.appConfig.GetSection("horizonheader")
	if err != nil {
		panic(err)
	}
	return headers
}
