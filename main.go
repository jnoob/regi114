package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type TerminalAndFile struct {
	file *os.File
}

func NewOut() *TerminalAndFile {
	f, _ := os.OpenFile("a.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	return &TerminalAndFile{
		file: f,
	}
}

func (t *TerminalAndFile) Write(p []byte) (n int, err error) {
	_, _ = t.file.Write(p)
	return os.Stdout.Write(p)
}

func main() {
	_7pre := get7Pre3Second()
	t := time.Unix(_7pre, 0)
	log.Printf("WAIT UNTIL:%v\n", t)
	needWait := time.Until(t)
	log.Printf("Need wait:%v\n", needWait.Seconds())
	<-time.After(needWait)

	log.SetOutput(NewOut())
	log.Printf("--------------------------------------------------\n")
	log.Printf("\tTime:\t\t%s\n", getTime())
	log.Printf("\tDutyDate:\t%s\n", getDutyDate())
	log.Printf("--------------------------------------------------\n")
	log.Printf("SELECT PRODUCT:\n")
	dutyDate := getDutyDate()
	candidates := getProductDetail(dutyDate)
	retryTimes := 0
	for len(candidates) == 0 {
		retryTimes++
		time.Sleep(time.Millisecond * 100)
		log.Printf("--------------------------------------------------\n")
		log.Printf("retry:%d", retryTimes)
		candidates = getProductDetail(dutyDate)
	}
	log.Printf("--------------------------------------------------\n")
	if len(candidates) > 0 {
		for _, can := range candidates {
			log.Printf("\tDOCTOR:	%s\n", can.DoctorName)
			log.Printf("\tTITLE	:	%s\n", can.DoctorTitleName)
			if sendVerifyCode(can.UniqProductKey, dutyDate) {
				log.Printf("--------------------------------------------------\n")
				log.Printf("WAIT VERIFY CODE:\n")
				var verifyCode int64
				fmt.Scanln(&verifyCode)
				if verifyCode > 0 {
					log.Printf("YOUR VERIFY CODE:%d\n", verifyCode)
					if !order(verifyCode, dutyDate, can.UniqProductKey) {
						log.Printf("PRODUCT ORDERED!")
						break
					} else {
						log.Printf("PRODUCT ORDER FAILED!")
					}
				}
			} else {
				log.Printf("--------------------------------------------------\n")
			}
		}
	}
	log.Printf("EXIST...\n")
}
