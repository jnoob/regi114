package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func order(verifyCode int64, date string, product string) bool {
	type Payload struct {
		HosCode        string `json:"hosCode"`
		FirstDeptCode  string `json:"firstDeptCode"`
		SecondDeptCode string `json:"secondDeptCode"`
		DutyTime       int    `json:"dutyTime"`
		TreatmentDay   string `json:"treatmentDay"`
		UniqProductKey string `json:"uniqProductKey"`
		CardType       string `json:"cardType"`
		CardNo         string `json:"cardNo"`
		SmsCode        string `json:"smsCode"`
		HospitalCardID string `json:"hospitalCardId"`
		Phone          string `json:"phone"`
		OrderFrom      string `json:"orderFrom"`
	}

	data := Payload{
		CardNo:         "123067684006",
		CardType:       "SOCIAL_SECURITY",
		DutyTime:       0,
		FirstDeptCode:  "75fec1a900e3d4c238cf384556de46de",
		HosCode:        "142",
		OrderFrom:      "HOSP",
		Phone:          "17611269218",
		SecondDeptCode: "200039484",
		SmsCode:        strconv.FormatInt(verifyCode, 10),
		TreatmentDay:   date,
		UniqProductKey: product,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to marshal payload:%+v, err:%v", data, err)
		return false
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://www.114yygh.com/web/order/save?_time=%s", getTime()), body)
	if err != nil {
		log.Printf("failed to new save order request, err:%v", err)
		return false
	}
	refer := fmt.Sprintf("https://www.114yygh.com/hospital/142/submission?hosCode=142&firstDeptCode=75fec1a900e3d4c238cf384556de46de&secondDeptCode=200039484&dutyTime=0&dutyDate=%s&uniqProductKey=%s", date, product)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Sec-Ch-Ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Request-Source", "PC")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Origin", "https://www.114yygh.com")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", refer)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cookie", COOKIE)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("failed to call save order request, err:%v", err)
		return false
	}
	defer resp.Body.Close()
	detail := &BaseResp{}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, detail)
	if err != nil {
		log.Printf("parse product detail:(%s) fail! err:%v\n", string(bodyBytes), err)
		return false
	}
	if detail.ResCode == 0 {
		return true
	}
	return false
}

type BaseResp struct {
	ResCode int32 `json:"resCode,omitempty"`
}
