package initizlize

import (
	"gitlab.example.com/zhangweijie/tool-sdk/config"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/logger"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"strings"
)

// 解密 config 文件中的加密数据
func decryptConfig(globalConfig *config.Cfg) error {
	secretKey := global.Config.Server.SecretKey
	if secretKey == "" {
		secretKey = os.Getenv("PATHP")
	}

	configValue := reflect.ValueOf(globalConfig).Elem()
	for i := 0; i < configValue.NumField(); i++ {
		field := configValue.Field(i)
		// 检查字段是否是结构体
		if field.Kind() == reflect.Struct {
			// 获取结构体字段的值
			structValue := reflect.Indirect(field)
			//	 使用反射遍历结构体字段
			for j := 0; j < structValue.NumField(); j++ {
				structField := structValue.Field(j)
				structFileInterface := structField.Interface()
				if structField.Kind() == reflect.String &&
					(strings.Contains(strings.ToLower(structValue.Type().Field(j).Name), "password") ||
						strings.Contains(strings.ToLower(structValue.Type().Field(j).Name), "apikey")) &&
					strings.HasPrefix(structFileInterface.(string), "ENC~") {
					data, err := config.DecryptString(structFileInterface.(string)[4:], secretKey)
					if err == nil {
						structField.SetString(data)
					} else {
						logger.Panic(schemas.DecryptConfigErr, err)
						return err
					}
				}
			}
		}

	}

	return nil
}

// LoadConfig 加载配置文件
func LoadConfig(config string) (err error) {
	file, err := os.ReadFile(config)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(file, &global.Config)
	if err != nil {
		return err
	}
	err = decryptConfig(global.Config)
	return err
}

func InitWorker(workerNum int) (err error) {
	global.ValidExecutorChan.WorkExecute = make(chan bool, workerNum)
	for i := 0; i < workerNum; i++ {
		global.ValidExecutorChan.WorkExecute <- true
	}

	return nil
}
