package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"time"
)

type Msg struct {
	UserIndex string `json:"userIndex"`
	Result    string `json:"result"`
	Message   string `json:"message"`
}

func main() {
	current, err := user.Current()
	if err != nil {
		fmt.Println("获取用户目录失败")
		return
	}
	cp := path.Join(current.HomeDir, ".config\\configstore\\login.json")
	fmt.Println(cp)
	bytes, err := ioutil.ReadFile(path.Join(current.HomeDir, ".config\\configstore\\login.json"))
	if err != nil {
		err := os.MkdirAll(path.Join(current.HomeDir, ".config\\configstore"), os.ModePerm)
		if err != nil {
			fmt.Println("文件夹已存在或者创建文件夹失败")
		}
		OpenDir(strings.ReplaceAll(cp, "/", "\\"))
		fmt.Println("读取配置文件出错:", err)
		fmt.Println("请编辑配置文件后重试,请打开login.json")
		file, err := os.Create(path.Join(current.HomeDir, ".config\\configstore\\login.json"))
		if err != nil {
			fmt.Println("创建文件失败")
			return
		}
		file.WriteString(`{
  "userId": "********", 
  "password": "*******",
  "localURL": "http://223.2.10.172"
}`)
		err = file.Close()
		if err != nil {
			fmt.Println("关闭文件错误")
			return
		}
		_, err2 := fmt.Scan()
		if err2 != nil {
			fmt.Println("阻塞io失败")
			return
		}
	}
	var loginer Loginer
	err = json.Unmarshal(bytes, &loginer)
	if err != nil {
		return
	}
	fmt.Println("读取账户信息成功,登陆的学号为:", loginer.UserId)
	res := loginer.PostInfo()
	var msg Msg
	err = json.Unmarshal([]byte(res), &msg)
	if err != nil {
		fmt.Println("unmarshal过程出错")
		return
	}
	if msg.UserIndex == "" {
		fmt.Printf("返回的结果为:%s,消息为:%s\n", msg.Result, msg.Message)
		OpenDir(strings.ReplaceAll(cp, "/", "\\"))
		loginer.ShowDialog()
	} else {
		fmt.Printf("登陆成功,用户的编号是%s\n", msg.UserIndex)
		fmt.Println("窗口将在2s后关闭")
	}
	time.Sleep(time.Second * 2)
}

const (
	queryString = "wlanuserip%253Df02d3ceebe91a66d1cc6e302e4732a44%2526wlanacname%253Db58fcda622d8bb5b%2526ssid%253D52eefd2d44d14e03%2526nasip%253De54b2b351c575112839c042a61804a4c%2526snmpagentip%253D%2526mac%253D83f50810620ed332f4567e3d08199054%2526t%253Dwireless-v2%2526url%253D2c0328164651e2b4f13b933ddf36628bea622dedcc302b30%2526apmac%253D%2526nasid%253Db58fcda622d8bb5b%2526vid%253D2758adca22277fb2%2526port%253Da059c1f642b68d9b%2526nasportid%253Da4ea9de28cbfe9fc659fd94ce1242135eceea7b545b082b366971ebb1790ab061b3ffa5a426253b6"
	loginRouter = "/eportal/InterFace.do?method=login"
	unicom      = "%25E8%2581%2594%25E9%2580%259A%25E6%259C%258D%25E5%258A%25A1"
	tele        = "%E7%94%B5%E4%BF%A1%E6%9C%8D%E5%8A%A1"
)

type Loginer struct {
	LocalURL string `json:"localURL"`
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

func (loginer *Loginer) GetLocationHref() string {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	request, err := http.NewRequest("GET", loginer.LocalURL, nil)
	if err != nil {
		fmt.Println("创建GET请求失败", err)
		return ""
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("访问认证界面失败:", err)
		return ""
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("读取网页内容失败")
		return ""
	}
	body := string(bytes)
	f := strings.IndexByte(body, '\'')
	e := strings.LastIndexByte(body, '\'')
	href := body[f+1 : e]

	fmt.Printf("网页内容为:\n %s\n\n", href)
	return href
}

func (loginer *Loginer) PostInfo() string {
	client := http.Client{
		Timeout: time.Second * 5,
	}
	body := strings.NewReader(fmt.Sprintf("userId=%s&password=%s&service=%s&queryString=%s&operatorPwd=&operatorUserId=&validcode=&passwordEncrypt=false", loginer.UserId, loginer.Password, unicom, queryString))
	request, err := http.NewRequest("POST", loginer.LocalURL+loginRouter, body)
	if err != nil {
		fmt.Println("创建GET请求失败", err)
		return ""
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("Accept-Encoding", "gzip, deflate")
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Content-Length", "697")
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Post账户和密码出错")
		return ""
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	res := string(bytes)

	return res
}
func (loginer *Loginer) ShowDialog() {
	fmt.Println("登陆过程发生错误，请检查用户名和密码是否正确")
	fmt.Println("如果无法解决问题，请联系开发者")
	fmt.Println("Gmail : ********")
}

func OpenDir(path string) {
	cmd := exec.Command("C:/Windows/explorer", path)
	err := cmd.Start()
	if err != nil {
		fmt.Println("打开文件资源管理器错误")
		return
	}
}
