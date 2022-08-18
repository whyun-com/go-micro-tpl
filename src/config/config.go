package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const DefaultConfigPath = "/etc/micro-config.yaml"

var Config ConfigYaml

type ConfigItem interface {
	Init()
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

func findConfigFile(path string) bool {
	return Exists(path) && IsFile(path)
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Err(err)
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

type ConfigYaml struct {
	Redis struct {
		Hosts       []string
		ClusterMode bool
	}
	Kafka struct {
		Log struct {
			Brokers string
			Topics  struct {
				AccessLog string
			}
		}
		Communication struct {
			Brokers string
			Topic   string
		}
	}
}

func init() {
	CONFIG_PATH := os.Getenv("CONFIG_PATH")
	if CONFIG_PATH != "" && findConfigFile(CONFIG_PATH) {
		// log.Info().Str("config file", CONFIG_PATH)
	} else if findConfigFile(DefaultConfigPath) {
		CONFIG_PATH = DefaultConfigPath
		// log.Info().Str("config file", CONFIG_PATH)
	} else {
		CONFIG_PATH = GetCurrentDirectory() + "/config.yaml"
		// log.Info().Str("config file", CONFIG_PATH)
		if CONFIG_PATH == "" || !findConfigFile((CONFIG_PATH)) {
			log.Fatal().Str("尝试从", CONFIG_PATH).Msg("未找到yaml配置文件")
			return
		}
	}
	yamlFile, err := ioutil.ReadFile(CONFIG_PATH)
	if err != nil {
		log.Fatal().Str("读取配置文件失败", CONFIG_PATH).Err(err)
		return
	}

	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		log.Fatal().Str("解析配置文件失败", CONFIG_PATH).Err(err)
		return
	}

	redisConfig := &RedisConfig{}
	redisConfig.Init()
	KafkaConfig := &KafkaConfig{}
	KafkaConfig.Init()
}
