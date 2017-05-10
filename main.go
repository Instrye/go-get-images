package main

import (
	"flag"
	"fmt"
	. "github.com/bitly/go-simplejson"
	// "github.com/garyburd/redigo/redis"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

type Conf struct {
	url       string
	imgRegex  string
	urlRegex  string
	redisPort string
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
	confJson.urlRegex = configJson.Get("urlregex").MustString()
	confJson.imgRegex = configJson.Get("imgregex").MustString()
	confJson.redisPort = configJson.Get("redisPort").MustString()
}
func main() {
	// redisconn, err := redis.Dial("tcp", "127.0.0.1:"+confJson.redisPort)
	// if err != nil {
	// 	fmt.Printf("redis :%s is not connection.\r\n%s\r\n", confJson.redisPort, err)
	// 	os.Exit(1)
	// }
	urlchan := make(chan string, 100)
	urlchan <- confJson.url
	paserUrl(urlchan)
}

func paserUrl(urlchan chan string) {

	for {
		url := <-urlchan
		fmt.Println(url)
		go paserHtml(url, urlchan)
	}

}

func paserHtml(url string, urlchan chan<- string) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("url :%s is not connection.\r\n%s\r\n", url, err)
	} else {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("url :%s is not get body.\r\n%s\r\n", url, err)
	}
	urlregex := regexp.MustCompile(confJson.urlRegex)
	urlarray := urlregex.FindAllString(string(body), -1)
	for _, v := range urlarray {
		urlchan <- v
	}
}
