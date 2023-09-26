package initizlize

import (
	"gitlab.example.com/zhangweijie/tool-sdk/config"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gopkg.in/yaml.v3"
	"os"
)

// 解密 config 文件中的加密数据
func decryptConfig(globalConfig *config.Cfg) {
	secretKey := global.Config.Server.SecretKey
	if secretKey == "" {
		secretKey = os.Getenv("PATHP")
	}
	decryDatas := []string{globalConfig.Database.Password, globalConfig.Elastic.Password}
	for i, decryData := range decryDatas {
		data, err := config.DecryptString(decryData[4:], secretKey)
		if err == nil && i == 0 {
			globalConfig.Database.Password = data
		}
		if err == nil && i == 1 {
			globalConfig.Elastic.Password = data
		}
	}
}

// LoadConfig 加载配置文件
func LoadConfig(config []byte) (err error) {
	err = yaml.Unmarshal(config, &global.Config)
	if err != nil {
		return err
	}
	//decryptConfig(global.Config)
	return err
}

func InitWorker(workerNum int) (err error) {
	global.ValidExecutorChan.WorkExecute = make(chan bool, workerNum)
	for i := 0; i < workerNum; i++ {
		global.ValidExecutorChan.WorkExecute <- true
	}

	return nil
}
