// @Title  uoj
// @Description  用于操作uoj相关提交
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:47
package Vjudge

import (
	"MGA_OJ/Interface"
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"time"
)

// UOJ			定义了UOJ接口
type UOJ struct {
	Interface.VjudgeInterface              // 包含Vjudge方法
	MainURL                   string       // 主url
	Session                   *http.Client // 会话
	UserID                    string       // 用户id
	Password                  string       // 密码
	LoggedIn                  bool         // 是否登录成功
	Token                     string       // token
}

// @title    Login
// @description   获得登录状态
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    password string       接收一个密码
// @return   bool	返回是否登录成功
func (fp *UOJ) Login() bool {

	resp, err := fp.Session.Get(fp.MainURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	fp.Token = extractToken(string(body))

	salt := "uoj233_wahaha!!!!"
	pwd := hmacMD5(salt, fp.Password)

	formData := url.Values{
		"_token":   {fp.Token},
		"login":    {""},
		"username": {fp.UserID},
		"password": {pwd},
	}

	resp, err = fp.Session.PostForm(fmt.Sprintf("%s/login", fp.MainURL), formData)
	if err != nil {
		fp.LoggedIn = false
		return false
	} else {
		fp.LoggedIn = true
	}
	resp.Body.Close()

	return fp.LoggedIn
}

// @title    Submit
// @description   提交
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    code, probID, lang string 代码，题目id，语言
// @return   string, error 表示提交id、报错消息
func (fp *UOJ) Submit(code, probID, lang string) (string, error) {
	if !fp.LoggedIn {
		// TODO 尝试登录
		if !fp.Login() {
			return "", fmt.Errorf("登录失败")
		}
	}

	MapLanguage := map[string]bool{
		"C++":       true,
		"C++03":     true,
		"C++11":     true,
		"C++14":     true,
		"C++17":     true,
		"C++20":     true,
		"C":         true,
		"Python3":   true,
		"Python2.7": true,
		"Java8":     true,
		"Java11":    true,
		"Java17":    true,
		"Pascal":    true,
	}

	if !MapLanguage[lang] {
		return "", fmt.Errorf("language error")
	}

	// TODO 构建提交表单数据
	formData := url.Values{
		"_token":                    {fp.Token},
		"answer_answer_language":    {"C++17"},
		"answer_answer_upload_type": {"editor"},
		"answer_answer_editor":      {code},
		"submit-answer":             {"answer"},
	}

	resp, err := fp.Session.PostForm(fmt.Sprintf("%s/problem/%s", fp.MainURL, probID), formData)
	if err != nil {
		if err != nil {
			// TODO 设置为未登录状态
			fp.LoggedIn = false
			return "", err
		}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if err != nil {
			// TODO 设置为未登录状态
			fp.LoggedIn = false
			return "", err
		}
	}

	// TODO Extracting Run ID
	runIDRe := regexp.MustCompile(`<a href="/submission/\d+">#(\d+)</a>`)
	runIDs := runIDRe.FindStringSubmatch(string(body))

	if len(runIDs) < 2 {
		if err != nil {
			// TODO 设置为未登录状态
			fp.LoggedIn = false
			return "", fmt.Errorf("lose runid, maybe code duplication")
		}
	}

	return runIDs[1], nil
}

// @title    GetStatus
// @description   跟踪提交状态
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    RunId, channel 提交id, 管道
// @return   string, error 表示提交id、报错消息
func (fp *UOJ) GetStatus(RunId string, ProbId string, channel chan map[string]string) {
	results := make(map[string]string)
	statusURL := fmt.Sprintf("https://uoj.ac/submission/%s", RunId)
	// TODO 持续请求提交状态
	// TODO 请求次数
	for i := 0; !ojmatchesRegex(results["Result"]) && i < 200; i++ {
		// TODO 逐步扩大请求间隔
		time.Sleep(time.Duration(i*i) * time.Second)
		resp, err := fp.Session.Get(statusURL)
		if err != nil {
			fmt.Println("Get Status Failed:", err)
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Get Status Failed:", err)
			continue
		}
		resp.Body.Close()
		results := UOJextractLatestSubmission(RunId, string(body))
		// TODO 将读出信息写入管道
		channel <- results
	}
	// TODO 关闭管道
	close(channel)
}

// @title    UOJextractLatestSubmission
// @description   分析UOJ提交表单
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    code, probID, lang string 代码，题目id，语言
// @return   string, error 表示提交id、报错消息
func UOJextractLatestSubmission(runid, html string) map[string]string {
	submission := make(map[string]string)
	submission["Html"] = html

	// TODO Extracting Result
	ReultIDRe := regexp.MustCompile(`<a href="/submission/\d+"" class="small">(.*?)</a></td><td>\d+ms</td><td>\d+kb</td>`)
	ResultIDs := ReultIDRe.FindStringSubmatch(html)

	if len(ResultIDs) > 1 {
		submission["Result"] = ResultIDs[1]
		// TODO Extracting Time
		TimeRe := regexp.MustCompile(`<a href="/submission/\d+"" class="small">.*?</a></td><td>(\d+)ms</td><td>\d+kb</td>`)
		Times := TimeRe.FindStringSubmatch(html)
		if len(Times) > 1 {
			submission["Time"] = Times[1]
		} else {
			submission["Time"] = "0"
		}

		// TODO Extracting Memory
		MemoryRe := regexp.MustCompile(`<a href="/submission/\d+"" class="small">.*?</a></td><td>\d+ms</td><td>(\d+)kb</td>`)
		Memorys := MemoryRe.FindStringSubmatch(html)
		if len(Memorys) > 1 {
			submission["Memory"] = Memorys[1]
		} else {
			submission["Memory"] = "0"
		}
		return submission
	}

	ReultIDRe = regexp.MustCompile(`<div class="col-sm-\d">(.*?)</div><div class="col-sm-\d">time: \d+ms</div><div class="col-sm-\d">memory: \d+kb</div>`)
	ResultIDss := ReultIDRe.FindAllStringSubmatch(html, -1)

	TimeRe := regexp.MustCompile(`<div class="col-sm-\d">.*?</div><div class="col-sm-\d">time: (\d+)ms</div><div class="col-sm-\d">memory: \d+kb</div>`)
	Timess := TimeRe.FindAllStringSubmatch(html, -1)

	MemoryRe := regexp.MustCompile(`<div class="col-sm-\d">.*?</div><div class="col-sm-\d">time: \d+ms</div><div class="col-sm-\d">memory: (\d+)kb</div>`)
	Memoryss := MemoryRe.FindAllStringSubmatch(html, -1)

	if len(ResultIDss) == 0 {
		return submission
	}

	submission["Result"] = ResultIDss[0][1]
	submission["Result"] = Timess[0][1]
	submission["Result"] = Memoryss[0][1]

	for i := range ResultIDss {
		if ResultIDss[i][1] != "Accepted" {
			submission["Result"] = ResultIDss[i][1]
			submission["Result"] = Timess[i][1]
			submission["Result"] = Memoryss[i][1]
		}
	}

	return submission
}

func extractToken(body string) string {
	re := regexp.MustCompile(`_token : "(.*?)"`)
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func hmacMD5(salt, passwd string) string {
	h := hmac.New(md5.New, []byte(salt))
	h.Write([]byte(passwd))
	return hex.EncodeToString(h.Sum(nil))
}

func NewUOJ(userID string, password string) Interface.VjudgeInterface {
	jar, _ := cookiejar.New(nil)
	return &UOJ{
		MainURL:  "https://uoj.ac",
		Session:  &http.Client{Jar: jar},
		UserID:   userID,
		Password: password,
		LoggedIn: false,
	}
}
