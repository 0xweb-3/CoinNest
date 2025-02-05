package initialize

import (
	"fmt"
	"github.com/0xweb-3/CoinNest/eth_srv/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	//debug := GetEnvInfo("DEBUG")
	v := viper.New()
	//debugName := ""
	//if debug {
	//	debugName = "_debug"
	//}

	//fileName := fmt.Sprintf("./user/config.yaml/conf/config.yaml%s.yaml", debugName)
	fileName := fmt.Sprintf("./eth_srv/config/conf/config.yaml")

	fmt.Println("配置：", fileName)

	v.SetConfigFile(fileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
	//zap.S().Debug("%v", global.ServerConfig)
	//fmt.Println(serverConfig.Name)

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Debug("配置信息发生变化：%v", global.ServerConfig)
		err := v.ReadInConfig()
		if err != nil {
			panic(err)
		}
		if err := v.Unmarshal(&global.ServerConfig); err != nil {
			panic(err)
		}
	})
}
