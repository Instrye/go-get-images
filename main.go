package main

import (
	"flag"
	"fmt"
	. "github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"os"
	// "regexp"
)

type Conf struct {
	url       string
	regex     string
	redisPort int
}

var confFile = flag.String("f", "config.json", "load config file")
var confJson Conf

func init() {
	flag.Parse()
	config, err := os.Open(*confFile)
	if err != nil {
		fmt.Printf("configFile :%s is load failed.\r\n%s\r\n", *confFile, err)
		os.Exit(1)
	}
	defer config.Close()
	configData, err := ioutil.ReadAll(config)
	if err != nil {
		fmt.Printf("configFile :%s is read failed.\r\n%s\r\n", *confFile, err)
		os.Exit(1)
	}
	configJson, err := NewJson(configData)
	if err != nil {
		fmt.Printf("configFile :%s is not json.\r\n%s\r\n", *confFile, err)
		os.Exit(1)
	}
	confJson.url = configJson.Get("url").MustString()
	confJson.regex = configJson.Get("regex").MustString()
	confJson.redisPort = configJson.Get("redisPort").MustInt()
}
func main() {
	res, err := http.Get(confJson.url)
	if err != nil {
		fmt.Printf("url :%s is not connection.\r\n%s\r\n", confJson.url, err)
		os.Exit(1)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("url :%s is not get body.\r\n%s\r\n", confJson.url, err)
		os.Exit(1)
	}

}
