package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"gorm.io/gorm/logger"
	"io"
)

type Cfg struct {
	Server   ServerConfig   `yaml:"server"`
	Elastic  ElasticConfig  `yaml:"elastic"`
	Database DatabaseConfig `yaml:"database"`
	Cache    CacheConfig    `yaml:"redis"`
	Tool     interface{}    `yaml:"tool"`
}

type ServerConfig struct {
	RunMode      string `yaml:"run_mode"`
	RunPort      string `yaml:"run_port"`
	Loglevel     string `yaml:"log_level"`
	LogType      string `yaml:"log_type"`
	LogFilePath  string `yaml:"log_file_path"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	Concurrency  int    `yaml:"concurrency"`
	Worker       int    `yaml:"worker"`
	RootDir      string `yaml:"root_dir"`
	SecretKey    string `yaml:"secret_key"`
	Mode         string `yaml:"mode"`
}

type ElasticConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Index    string `yaml:"index"`
	Activate bool   `yaml:"activate"` // 是否激活
}

type DatabaseConfig struct {
	Host                 string          `yaml:"host"`
	Port                 string          `yaml:"port"`
	DbName               string          `yaml:"db_name"`
	SSLMode              string          `yaml:"sslmode"`
	TimeZone             string          `yaml:"timezone"`
	Username             string          `yaml:"username"`
	Password             string          `yaml:"password"`
	PreferSimpleProtocol bool            `yaml:"prefer_simple_protocol"`
	MaxIdleConns         int             `yaml:"max_idle_conns"`
	MaxOpenConns         int             `yaml:"max_open_conns"`
	LogLevel             logger.LogLevel `yaml:"log_level"`
	SlowThreshold        int             `yaml:"slow_threshold"`
	Activate             bool            `yaml:"activate"` // 是否激活
}

type CacheConfig struct {
	Hosts      string `yaml:"hosts"`
	Password   string `yaml:"password"`
	MasterName string `yaml:"master_name"`
	Sentinel   bool   `yaml:"sentinel"`
	Database   int    `yaml:"database"`
	PoolSize   int    `yaml:"pool_size"`
	Activate   bool   `yaml:"activate"` // 是否激活
}

func Encrypt(sourceData []byte, secretKey string) string {
	block, err := aes.NewCipher([]byte(createHash(secretKey)))
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	seal := gcm.Seal(nonce, nonce, sourceData, nil)
	return base64.StdEncoding.EncodeToString(seal)
}

func DecryptString(encipheredData string, secretKey string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encipheredData)
	if err != nil {
		return "", errors.New("Invalid text to decrypt")
	}
	key := []byte(createHash(secretKey))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return string(plaintext), nil
}

func EncryptString(sourceData string, secretKey string) string {
	return "ENC~" + string(Encrypt([]byte(sourceData), secretKey))
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
