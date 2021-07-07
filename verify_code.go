package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type VerifyCodeResp struct {
	ResCode int64 `json:"resCode,omitempty"`
	Data    struct {
		EndMilliseconds int64 `json:"endMilliseconds,omitempty"`
	} `json:"data,omitempty"`
}

func sendVerifyCode(product, dutyDate string) (ok bool) {
	reqUrl := fmt.Sprintf("https://www.114yygh.com/web/common/verify-code/get?_time=%s&mobile=17611269218&smsKey=ORDER_CODE&uniqProductKey=%s", getTime(), product)
	fmt.Printf("send verify code url:%s\n", reqUrl)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.Printf("new get verify code request err:%v\n", err)
		return
	}
	referer := fmt.Sprintf("https://www.114yygh.com/hospital/142/submission?hosCode=142&firstDeptCode=75fec1a900e3d4c238cf384556de46de&secondDeptCode=200039484&dutyTime=0&dutyDate=%s&uniqProductKey=%s", dutyDate, product)
	req.Host = "www.114yygh.com"
	req.Header.Set("Sec-Ch-Ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Request-Source", "PC")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", referer)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cookie", COOKIE)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("failed call get verify code request err:%v\n", err)
		return
	}
	defer resp.Body.Close()

	target := &VerifyCodeResp{}
	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		log.Printf("failed to decode verify code request! err:%v\n", err)
		return
	} else {
		log.Printf("resp:%+v\n", target)
		return true
	}
}
