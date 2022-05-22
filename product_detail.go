package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func getProductDetail(sd *SessionData, dutyDate string) (candidates []*ProductDetail) {
	var logMsg string
	defer func() {
		if len(logMsg) > 0 && !sd.CandidatesGot() {
			log.Print(logMsg)
		}
	}()

	type Payload struct {
		FirstDeptCode  string `json:"firstDeptCode"`
		SecondDeptCode string `json:"secondDeptCode"`
		HosCode        string `json:"hosCode"`
		Target         string `json:"target"`
	}

	data := Payload{
		FirstDeptCode:  "75fec1a900e3d4c238cf384556de46de",
		SecondDeptCode: "200039484",
		HosCode:        "142",
		Target:         dutyDate,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		logMsg = fmt.Sprintf("failed to parse payload:%+v, err:%v\n", data, err)
		return
	}
	body := bytes.NewReader(payloadBytes)

	reqUrl := fmt.Sprintf("https://www.114yygh.com/web/product/detail?_time=%s", getTime())
	req, err := http.NewRequest("POST", reqUrl, body)
	if err != nil {
		logMsg = fmt.Sprintf("failed to new product detail url, err:%v\n", err)
		return
	}
	req.Host = "www.114yygh.com"
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
	req.Header.Set("Referer", "https://www.114yygh.com/hospital/142/75fec1a900e3d4c238cf384556de46de/200039484/source")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cookie", COOKIE)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logMsg = fmt.Sprintf("request for product detail fail, err:%v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logMsg = fmt.Sprintf("Request not ok! status:%v\n", resp.StatusCode)
		return
	}
	detail := &ProductResp{}
	err = json.NewDecoder(resp.Body).Decode(detail)
	if err != nil {
		logMsg = fmt.Sprintf("parse product detail fail! err:%v\n", err)
		return
	}
	logMsg = fmt.Sprintf("\nrespData:%+v\n\n", detail)
	if len(detail.Data) == 0 {
		logMsg += "\nNO PRODUCT DATA!"
		return
	}

	var (
		morningDetails, afternoonDetails []*ProductDetail
		otherDetails                     = make([][]*ProductDetail, 0)
	)
	for _, d := range detail.Data {
		switch d.DutyCode {
		case "MORNING":
			morningDetails = d.Detail
		case "AFTERNOON":
			afternoonDetails = d.Detail
		default:
			otherDetails = append(otherDetails, d.Detail)
		}
	}

	candidates = []*ProductDetail{}
	addCandidates := func(details []*ProductDetail) {
		reverseDetailSlice(details)
		for _, detail := range details {
			if detail.DoctorTitleName != "普通门诊" && detail.Wnumber%2 != 0 {
				// if detail.Wnumber%2 != 0 {
				// 非普通门诊，并且不是貌似满了的情况，作为一个候选
				candidates = append(candidates, detail)
			}
		}
	}
	if len(morningDetails) > 0 {
		addCandidates(morningDetails)
	}
	if len(afternoonDetails) > 0 {
		addCandidates(afternoonDetails)
	}
	for _, details := range otherDetails {
		addCandidates(details)
	}
	return
}

func reverseDetailSlice(details []*ProductDetail) {
	for i, j := 0, len(details)-1; i < j; i, j = i+1, j-1 {
		details[i], details[j] = details[j], details[i]
	}
}

type ProductResp struct {
	ResCode int64           `json:"resCode,omitempty"`
	Msg     interface{}     `json:"msg",omitempty`
	Data    []*NoonPruducts `json:"data,omitempty"`
}

type ProductDetail struct {
	UniqProductKey  string `json:"uniqProductKey,omitempty"`
	DoctorName      string `json:"doctorName,omitempty"`
	DoctorTitleName string `json:"doctorTitleName,omitempty"`
	Fcode           string `json:"fcode,omitempty"`
	Ncode           string `json:"ncode,omitempty"`
	Wnumber         int64  `json:"wnumber,omitempty"`
	Znumber         int64  `json:"znumber,omitempty"`
}

func (d *ProductDetail) String() string {
	return fmt.Sprintf("(key:%s, name:%s, title:%s, wnumber:%d)", d.UniqProductKey, d.DoctorName, d.DoctorTitleName, d.Wnumber)
}

type NoonPruducts struct {
	DutyCode string           `json:"dutyCode,omitempty"`
	Detail   []*ProductDetail `json:"detail,omitempty"`
}

func (p *NoonPruducts) String() string {
	return fmt.Sprintf("[%s]%v", p.DutyCode, p.Detail)
}
