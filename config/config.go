package config

import (
	"errors"
	"fmt"
	"github.com/Unknwon/goconfig"
	"log"
	"os"
)

const configFile = "/conf/conf.ini"

var File *goconfig.ConfigFile

func init() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	configPath := currentDir + configFile

	if !fileExist(configPath) {
		panic(errors.New("配置文件不存在"))
	}
	args := len(os.Args)
	if args > 1 {
		dir := os.Args[1]
		if dir != "" {
			configPath = dir + configFile
		}
	}

	File, err = goconfig.LoadConfigFile(configPath)

	if err != nil {
		log.Fatal("读取配置文件出错", err)
	}

	fmt.Println(File)
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
