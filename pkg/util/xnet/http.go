package xnet

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var (
	lock         *sync.Mutex = &sync.Mutex{}
	curlInstance *HttpSingleton
)

type HttpSingleton struct {
	Client *http.Client
}

func GetHttpInstance() *HttpSingleton {
	if curlInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if curlInstance == nil {
			curlInstance = &HttpSingleton{}
			//设置忽律https 请求的证书
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			//设置请求超时时间
			httpSendTimeOut := time.Second * 5

			//请求执行
			curlInstance.Client = &http.Client{Transport: tr, Timeout: httpSendTimeOut}
		}
	}
	return curlInstance
}

func HttpGet(Url string) ([]byte, error) {
	resp, err := GetHttpInstance().Client.Get(Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("statusCode is " + strconv.Itoa(resp.StatusCode))
	}

	//获取返回数据
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	return body, nil
}

func HttpPost(postUrl string, param map[string]string) ([]byte, error) {
	postValue := url.Values{}
	for k, v := range param {
		postValue.Set(k, v)
	}
	resp, err := GetHttpInstance().Client.PostForm(postUrl, postValue)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(" statusCode is " + strconv.Itoa(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
