package initialize

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"mxshop_srvs/inventory_srv/global"
)

// InitConfig 通过先通过viper获取本地nacos配置在获取服务配置
func init() {
	v := viper.New()
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("inventory_srv/%s_pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("inventory_srv/%s_prev.yaml", configFilePrefix)
	}
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		zap.S().Panicw("read config failed", "err", err)
	}
	// 实例化配置
	if err := v.Unmarshal(global.NacosConf); err != nil {
		zap.S().Panicw("unmarshal config failed", "err", err)
	}

	// 创建clientConfig
	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConf.Namespace, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./inventory_srv/consul/log",
		CacheDir:            "./inventory_srv/consul/cache",
		LogLevel:            "debug",
	}
	// 至少一个ServerConfig
	sc := []constant.ServerConfig{
		{
			IpAddr:      global.NacosConf.Host,
			ContextPath: "/consul",
			Port:        uint64(global.NacosConf.Port), // 使用 HTTP 端口
			Scheme:      global.NacosConf.Scheme,       // 强制使用 HTTP 协议
		},
	}

	// 创建动态配置客户端的另一种方式 (推荐)
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		zap.S().Panicw("create config client failed", "err", err)
	}
	// 获取配置

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConf.DataId,
		Group:  global.NacosConf.Group,
	})
	fmt.Println(content)
	if err != nil {
		zap.S().Panicw("get config failed", "err", err)
	}
	// 监听配置变化
	if err = configClient.ListenConfig(vo.ConfigParam{
		DataId: global.NacosConf.DataId,
		Group:  global.NacosConf.Group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	}); err != nil {
		zap.S().Panicw("config listen failed", "err", err)
	}

	// 将json序列化为struct
	// 实例化配置对象
	serverConfig := global.ServerConf
	// 通过toml配置赋值给serverConfig
	if err = toml.Unmarshal([]byte(content), &serverConfig); err != nil {
		zap.S().Panicw("unmarshal config failed", "err", err)
	}
	//if err = json.Unmarshal([]byte(content), &serverConfig); err != nil {
	//	zap.S().Panicw("unmarshal config failed", "err", err)
	//}
}

// 获取环境变量
func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
