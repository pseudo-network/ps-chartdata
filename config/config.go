package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

var Conf *Config

// init config function at build
func InitConf() {
	Conf = new(Config)
	Conf.setEnv()
	// if err := Conf.setupViper(); err != nil {
	// 	panic(err)
	// }
	Conf.setupConfig()
}

// service config generate at build
type Config struct {
	Env        string
	ConsulHost string
	ConsulPort string
	MongoRC    []string
}

// set environment
func (c *Config) setEnv() {
	// env := os.Getenv("ENV")
	// if env == DEV {
	// 	c.Env = DEV
	// 	c.ConsulHost = CONSUL_HOST_CLUSTER
	// 	c.ConsulPort = CONSUL_PORT_CLUSTER
	// } else if env == STAGE {
	// 	c.Env = STAGE
	// 	c.ConsulHost = CONSUL_HOST_CLUSTER
	// 	c.ConsulPort = CONSUL_PORT_CLUSTER
	// } else if env == PROD {
	// 	c.Env = PROD
	// 	c.ConsulHost = CONSUL_HOST_CLUSTER
	// 	c.ConsulPort = CONSUL_PORT_CLUSTER
	// } else {
	// 	c.Env = WORKSTATION
	// 	c.ConsulHost = CONSUL_HOST_DEV
	// 	c.ConsulPort = CONSUL_PORT_DEV
	// }
}

// setup viper to use consul as the remote config provider
func (c *Config) setupViper() error {
	consulUrl := fmt.Sprintf("%s:%s", c.ConsulHost, c.ConsulPort)
	if err := viper.AddRemoteProvider("consul", consulUrl, CONSUL_KV); err != nil {
		return err
	}
	viper.SetConfigType("json")
	if err := viper.ReadRemoteConfig(); err != nil {
		return err
	}

	return nil
}

// setup config with premise to either running service on cluster vs locally
func (c *Config) setupConfig() {
	// if c.GetEnv() == WORKSTATION {
	// 	c.MongoRC = MONGOHOSTS_WORKSTATION
	// } else {
	// 	mongoRCI := viper.Get("mongo_rc").([]interface{})
	// 	var mongoRC []string
	// 	for _, mongoHost := range mongoRCI {
	// 		mongoRC = append(mongoRC, mongoHost.(string))
	// 	}

	// 	c.MongoRC = mongoRC
	// }
}

// get environment set at build
func (c *Config) GetEnv() string {
	return c.Env
}

// get jwt secret
func GetJwtSecret() (string, error) {
	if err := viper.ReadRemoteConfig(); err != nil {
		return "", nil
	}
	jwtSecret := viper.GetString("jwt_secret")
	return jwtSecret, nil
}

// get mongo admin user
func GetMongoUser() (*string, error) {
	user := viper.GetString("mongo_user")
	if user == "" {
		return nil, errors.New("no mongo user found")
	}

	return &user, nil
}

// get mongo admin password
func GetMongoPassword() (*string, error) {
	password := viper.GetString("mongo_password")
	if password == "" {
		return nil, errors.New("no mongo password found")
	}

	return &password, nil
}

// get shopmonkey API keys
func GetShopMonkeyAPIKeys() (string, string, error) {
	if err := viper.ReadRemoteConfig(); err != nil {
		return "", "", err
	}
	publicKey := viper.GetString("shop_monkey_public_key")
	if publicKey == "" {
		return "", "", errors.New("Shop Monkey public key is blank")
	}
	privateKey := viper.GetString("shop_monkey_private_key")
	if privateKey == "" {
		return "", "", errors.New("Shop Monkey private key is blank")
	}

	return publicKey, privateKey, nil
}
