package aliutils

import (
	"bytes"
	"encoding/base64"
	"strconv"
	"strings"
	"time"
)

const (
	ACCESS_FROM_USER = 0
	COLON            = ":"
)

func GetUserName(ak string, instanceId string) string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(ACCESS_FROM_USER))
	buffer.WriteString(COLON)
	buffer.WriteString(instanceId)
	buffer.WriteString(COLON)
	buffer.WriteString(ak)
	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}

func GetUserNameBySTSToken(ak string, instanceId string, stsToken string) string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(ACCESS_FROM_USER))
	buffer.WriteString(COLON)
	buffer.WriteString(instanceId)
	buffer.WriteString(COLON)
	buffer.WriteString(ak)
	buffer.WriteString(COLON)
	buffer.WriteString(stsToken)
	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}

func GetPassword(sk string) string {
	now := time.Now()
	currentMillis := strconv.FormatInt(now.UnixNano()/1000000, 10)
	var buffer bytes.Buffer
	buffer.WriteString(strings.ToUpper(HmacSha1(currentMillis, sk)))
	buffer.WriteString(COLON)
	buffer.WriteString(currentMillis)
	//fmt.Println(currentMillis)
	//fmt.Println(HmacSha1(sk, currentMillis))
	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}
