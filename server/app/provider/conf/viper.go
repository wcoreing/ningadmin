package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"wnadmin/app/config"
)

func init() {
	configapp()
	configdb()

}

// Viper //
// 优先级: 命令行 > 环境变量 > 默认值
// Author [SliverHorn](https://github.com/SliverHorn)
func Viper(yaml string) *viper.Viper {
	file := fmt.Sprintf("etc/%s.yaml", yaml)
	v := viper.New()
	v.SetConfigFile(file)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	return v

}
func configapp() {
	v := Viper("app")
	if err := v.Unmarshal(&config.App); err != nil {
		panic(err)
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&config.App); err != nil {
			fmt.Println(err)
		}
	})
}
func configdb() {
	v := Viper("db")
	if err := v.Unmarshal(&config.DB); err != nil {
		panic(err)
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&config.DB); err != nil {
			fmt.Println(err)
		}
	})
}
