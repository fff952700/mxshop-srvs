package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"mxshop-srvs/user-srv/global"
)

func InitConfig() {
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user-srv/%s_pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("user-srv/%s_debug.yaml", configFilePrefix)
	}
	v := viper.New()
	// 设置文件路径
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	// 实例化配置结构体
	serverConfig := global.ServerConf
	if err := v.Unmarshal(serverConfig); err != nil {
		panic(err)
	}
	// viper 动态监听变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		// 重新读取
		_ = v.ReadInConfig()
		// 重新解析
		_ = v.Unmarshal(serverConfig)
		fmt.Printf("%+v\n", serverConfig)

	})
}

// 获取环境变量
func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
