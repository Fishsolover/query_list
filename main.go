package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

const (
	GONGYINGSHANG = "供应商</dt>\n          <dd class=\"information-list__item__definition\">([^<]+)</dd>"
)

type Config struct {
	Name  string
	Regex string
}

type UrlInfo struct {
	Numbers int
	Url     string
	Wg      *sync.WaitGroup
	Config  []*Config
}

func NewUrlInfo(numbers int, url string) UrlInfo {
	return UrlInfo{
		Numbers: numbers,
		Url:     url,
		Wg:      &sync.WaitGroup{},
		Config:  nil,
	}
}

func main() {
	var num int
	var url string
	flag.IntVar(&num, "num", 100, "数量")
	flag.StringVar(&url, "url", "https://apps.apple.com/cn/app/id414478124", "链接")
	flag.Parse()
	info := NewUrlInfo(num, url)
	info.Config = append(info.Config, NewConfig("供应商", GONGYINGSHANG))
	info.run()
	info.Wg.Wait()
}

func (info *UrlInfo) run() {
	fmt.Printf("%v", info.Numbers)
	for i := 0; i < info.Numbers; i++ {
		info.Wg.Add(1)
		//fmt.Printf("ADD" + strconv.Itoa(i))
		go func(a int) {
			defer info.Wg.Done()
			str, err := getContent(info.Url)
			if err != nil {
				fmt.Printf(err.Error())
				return
			}

			//fmt.Printf(str)
			for _, r := range info.Config {
				s, e := r.GetDetails(str)
				if e != nil {
					fmt.Printf(e.Error())
				}
				fmt.Printf(s)
			}

		}(i)
	}
}

func getContent(url string) ([]byte, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	//fmt.Printf("%+v\n",request)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")

	res, err1 := client.Do(request)

	if err1 != nil {
		return nil, err1
	}

	if http.StatusOK != res.StatusCode {
		return nil, err
	}

	data, err2 := ioutil.ReadAll(res.Body)

	if err2 != nil {
		return nil, err2
	}

	return data, nil
}

func NewConfig(name string, regex string) *Config {
	return &Config{
		Name:  name,
		Regex: regex,
	}
}

type Operation interface {
	GetDetails(content []byte) (string, error)
}

func (cf *Config) GetDetails(content []byte) (string, error) {
	match := regexp.MustCompile(cf.Regex)
	matchData := match.FindAllSubmatch(content, -1)
	//fmt.Println(matchData)
	if matchData == nil {
		return "", errors.New("无匹配")
	}
	return fmt.Sprintf("%s:%s", cf.Name, string(bytes.Trim(matchData[0][1], ""))), nil
}
