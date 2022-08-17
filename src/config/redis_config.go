package config

import (
	// "os"
	// "encoding/json"
	// "fmt"
	"os"
	// "regexp"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)
const (
	ENV_REDIS_HOSTS = "REDIS_HOSTS"
	ENV_REDIS_CLUSTER_MODE = "REDIS_CLUSTER_MODE"
)

var RedisClient redis.UniversalClient

type RedisConfig struct {
}
func getRedisOption(hosts []string) *redis.UniversalOptions {
	log.Info().Msg("redis配置:" + strings.Join(hosts[:], ","))
	return &redis.UniversalOptions{
		Addrs: hosts,
	}
}

func (config *RedisConfig) Init() {
	hostsFromEnv := os.Getenv(ENV_REDIS_HOSTS)

	if hostsFromEnv != "" {
		hosts := strings.Split(hostsFromEnv, ",")
		Config.Redis.Hosts = hosts
		// if len(hosts) == 1 {
		// 	if ok, _ := regexp.Match("[A-Za-z]+", []byte(hosts[0])); ok {//使用域名模式，则强制使用cluster
		// 		RedisClient = redis.NewUniversalClient(getRedisOption([]string{hosts[0], hosts[0]}))
		// 	} else {
		// 		RedisClient = redis.NewUniversalClient(getRedisOption(hosts))
		// 	}
		// } else {
		// 	RedisClient = redis.NewUniversalClient(getRedisOption(hosts))
		// }
		// return
	}
	clusterModeFromEnv := strings.ToLower(os.Getenv(ENV_REDIS_CLUSTER_MODE))
	if clusterModeFromEnv == "true" {
		Config.Redis.ClusterMode = true
	}

	if Config.Redis.ClusterMode && len(Config.Redis.Hosts) == 1 {// 阿里云的 redis cluster 只使用单域名节点做代理
		host := Config.Redis.Hosts[0]
		Config.Redis.Hosts = []string{host, host}
	}


	// if len(redisConfig.Cluster) > 0 {
	// 	var hosts []string
	// 	for _, host := range redisConfig.Cluster {
	// 		hosts = append(hosts, fmt.Sprintf("%s:%d",host.Host, host.Port))
	// 	}
	// 	RedisClient = redis.NewUniversalClient(getRedisOption(hosts))
	// 	return
	// }
	if len(Config.Redis.Hosts) == 0 {
		log.Fatal().Msg("读取的redis配置地址为空")
		return
	}

	RedisClient = redis.NewUniversalClient(getRedisOption(Config.Redis.Hosts))
}