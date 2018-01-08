package utils

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
    "os/exec"
	"reflect"
)

func ToSlice(key string, arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		panic("toslice arr not slice")
	}
	l := v.Len()
	ret := make([]interface{}, l+1)
	ret[0] = key
	for i := 0; i < l; i++ {
		ret[i+1] = v.Index(i).Interface()
	}
	return ret
}

/*
get 获取接口数据
中文参数必须过滤
*/
func Get(urlStr string, query string) (rs []byte, err error) {
	//esponse, err := http.Get(Url.String())
	client := &http.Client{}
	request, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	q := request.URL.Query()
	q.Add("q", string(query))
	request.URL.RawQuery = q.Encode()
	//fmt.Println(request.URL.String())
	request.Header.Set("Content-Type", "application/json; encoding=utf-8")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

/**
直接get请求url
*/
func GetUrl(urlStr string) (rs []byte, err error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json; encoding=utf-8")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

/**
post方式
*/
func Post(apiURL string, body []byte) (rs []byte, err error) {
	postReq, err := http.NewRequest("POST", apiURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	postReq.Header.Set("Content-Type", "application/json; encoding=utf-8")
	client := &http.Client{}
	resp, err := client.Do(postReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

/**
postform方式
*/
func PostForm(apiURL string, params map[string][]string)  (rs []byte, err error) {
    resp, err := http.PostForm(apiURL, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

/**
一次性读取文件
*/
func ReadFile(path string) (rs []byte, err error) {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	return ioutil.ReadAll(fi)
}

/**
逐行读取文件（性能可能不高)
*/
func ReadLine(path string) (rs string, err error) {
	var buffer bytes.Buffer
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		buffer.WriteString(line)
		buffer.WriteString("<br/>")
	}
	return buffer.String(), nil
}

/**
执行shell操作
*/
func ExcuteCmd(cmdStr string) (string, error){
    cmd := exec.Command("/bin/bash", "-c", cmdStr)
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return "", err
    }
    //执行命令
    if err := cmd.Start(); err != nil {
        return "", err
    }
    bytestr, err := ioutil.ReadAll(stdout)
    if err != nil {
       return "", err
    }
    if err := cmd.Wait(); err != nil {
        return "", err
    }
    return string(bytestr), nil
}

