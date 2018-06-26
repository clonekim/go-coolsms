package coolsms

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	. "github.com/ohomango/config" //이것은 환경설정 로딩한 구조체
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func checkHangulLength(text string) int {

	var buf bytes.Buffer
	wr := transform.NewWriter(&buf, korean.EUCKR.NewEncoder())
	wr.Write([]byte(text))
	wr.Close()

	return len(buf.String())
}

func Send(msg *CoolMessage) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	salt := RandBytes(5)
	mac := hmac.New(md5.New, []byte(Conf.Sms.Apisecret))
	mac.Write([]byte(timestamp + salt))
	signature := hex.EncodeToString(mac.Sum(nil))

	var mType string
	if mType = "SMS"; checkHangulLength(msg.Message) > 80 {
		mType = "LMS"
	}

	msg.Type = mType

	resp, err := http.PostForm("https://api.coolsms.co.kr/sms/2/send", url.Values{
		"api_key":   {Conf.Sms.Apikey},
		"timestamp": {timestamp},
		"salt":      {salt},
		"to":        {msg.To},
		"from":      {"Your Phone number"},
		"text":      {msg.Message},
		"type":      {msg.Type},
		"signature": {signature},
	})

	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
    
  //TODO body처리  
	fmt.Println(string(body))

}

type CoolMessage struct {
	To      string
	Message string
	Type    string
}
