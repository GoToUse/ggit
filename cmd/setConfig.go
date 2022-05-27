package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type ConfigYaml struct {
	viper *viper.Viper
}

// NewConfig 初始化配置文件
func NewConfig() (*ConfigYaml, error) {
	homePath, _ := os.UserHomeDir()
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")

	currentConfPath, _ := filepath.Abs("configs/")
	ggitDir := filepath.Join(homePath, ".ggit")

	_, err := os.Stat(currentConfPath)
	_, errGgit := os.Stat(ggitDir)
	if os.IsNotExist(err) {
		if os.IsNotExist(errGgit) {
			return nil, errors.New(fmt.Sprintf("%s or %s not exist.", currentConfPath, ggitDir))
		} else {
			vp.AddConfigPath(ggitDir + "/")
			err = vp.ReadInConfig()
			if err != nil {
				return nil, err
			}
			return &ConfigYaml{vp}, nil
		}
	} else {
		vp.AddConfigPath("configs/")
		err = vp.ReadInConfig()
		if err != nil {
			return nil, err
		}
		return &ConfigYaml{vp}, nil
	}
}

// ReadSection 将 key 对应的配置信息写入到 v 中
func (c *ConfigYaml) ReadSection(key string, v interface{}) error {
	err := c.viper.UnmarshalKey(key, v)
	if err != nil {
		return err
	}

	return nil
}
