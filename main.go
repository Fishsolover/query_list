package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
)

const (
	GONGYINGSHANG = "供应商</dt>\n          <dd class=\"information-list__item__definition\">([^<]+)</dd>"
)

func main() {
	var numbers int
	flag.IntVar(&numbers, "num", 100, "数量")

	var wg sync.WaitGroup
	for i := 0; i < numbers; i++ {
		wg.Add(1)
		fmt.Printf("ADD" + strconv.Itoa(i))
		go func(a int) {
			defer wg.Done()
			str, err := getContent("https://apps.apple.com/cn/app/id414478124")
			if err != nil {
				fmt.Printf(err.Error())
				return
			}
			fmt.Printf("ADDxxA" + strconv.Itoa(a))

			//fmt.Printf(str)
			details, err := getDetails(str, GONGYINGSHANG)
			if err != nil {
				fmt.Printf(err.Error())
				return
			}
			f, _ := os.OpenFile("a.txt", os.O_CREATE, 0666)
			defer f.Close()

			_, err = bufio.NewWriter(f).WriteString(strconv.Itoa(a) + ":" + details)

			if err != nil {
				fmt.Printf(err.Error())
				return
			}

		}(i)
	}

	wg.Wait()

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

func getDetails(content []byte, regex string) (string, error) {
	match := regexp.MustCompile(regex)
	matchData := match.FindAllSubmatch(content, -1)
	//fmt.Println(matchData)
	if matchData == nil {
		return "", errors.New("无匹配")
	}
	return string(bytes.Trim(matchData[0][1], "")), nil
}
