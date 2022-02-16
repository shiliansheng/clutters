package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"	
	"net/http"
	"strconv"
	"strings"
	"time"
	"ccgpgov/common"
)

type UserInfo struct {
	Account  string `json:"account"`
	Nickname string `json:"nickname"`
}

/**
* 发送content
*/
func PostData(content string) error{
	url := "http://172.16.1.11:51001/json/msg/SendMsg"
	jsonStr := GetJson(strconv.FormatInt(time.Now().UnixMilli(), 10), content)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(([]byte(jsonStr))))
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	statuscode := resp.StatusCode
	hea := resp.Header
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	fmt.Println(statuscode)
	fmt.Println(hea)
	return nil
}

/**
 * time 	发送时间 
 * content 	发送内容
 * return	创建json格式字符串
 */
func GetJson(time, content string) string {
	jsonStr := `
	{
		"userId": "497",
		"sec": "6d12d8e49e5957378abd17593cfae13a",
		"msgid": 0,
		"type": 3,
		"localseq": 0,
		"sender": "497",
		"senderName": "小维智能助手",
		"receiver": "0",
		"receiverName": "",
		"group": "6043",
		"groupName": "招标监测",
		"time": "{time}",
		"mime": "txt",
		"userdata": "{\"atinfo\":{\"type\":0,\"param\":{\"userinfos\":{atinfo}}}}",
		"isencrypt": 0,
		"content": "{content}",
		"privatemsg": 0,
		"ecid": "xw",
		"terminalType": 1
	}
	`
	atinfo, ret := GetAtInfo()
	jsonStr = strings.Replace(jsonStr, "{time}", time, -1)
	jsonStr = strings.Replace(jsonStr, "{content}", ret + content, -1)
	jsonStr = strings.Replace(jsonStr, "{atinfo}", atinfo, -1)

	return jsonStr
}
/**
 * accountArr	配置文件中的account数组
 * nicknameArr	配置文件中的nickname数组
 * return1		对应account和nickname信息json数组字符串
 * return2		@nickname中所有人的字符串
 */
func GetAtInfo() (string, string) {
	var userinfos []UserInfo
	var atContent string = ""
	for i := range common.Config.Account {
		userinfos = append(userinfos, UserInfo{common.Config.Account[i], common.Config.Nickname[i]})
		atContent += "@" + common.Config.Nickname[i] + " "
	}
	ret, err := json.Marshal(userinfos)
	if err != nil {
		fmt.Println("UserInfo数组转换失败！")
	}
	return strings.ReplaceAll(string(ret), "\"", "\\\""), atContent
}