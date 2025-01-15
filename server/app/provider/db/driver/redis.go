package driver

import (
	"fmt"
	"github.com/go-redis/redis"
	"wnadmin/app/provider/db/model"
)

func RedisByConfig(redisCfg model.Redis) redis.UniversalClient {
	var client redis.UniversalClient
	// 使用集群模式
	if redisCfg.UseCluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    redisCfg.ClusterAddrs,
			Password: redisCfg.Password,
		})
	} else {
		// 使用单例模式
		client = redis.NewClient(&redis.Options{
			Addr:     redisCfg.Addr,
			Password: redisCfg.Password,
			DB:       redisCfg.DB,
		})
	}
	pong, err := client.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("redis链接成功" + pong)
	}
	return client
}
