// @Title  vijos
// @Description  用于操作vijos相关提交
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:47
package Vjudge

import (
	"MGA_OJ/Interface"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"time"
)

// VIJOS			定义了VIJOS接口
type VIJOS struct {
	Interface.VjudgeInterface              // 包含Vjudge方法
	MainURL                   string       // 主url
	Session                   *http.Client // 会话
	UserID                    string       // 用户id
	Password                  string       // 密码
	LoggedIn                  bool         // 是否登录成功
	CsrfTokens                string       // token
}

// @title    Login
// @description   获得登录状态
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    password string       接收一个密码
// @return   bool	返回是否登录成功
func (fp *VIJOS) Login() bool {

	if fp.LoggedIn {
		fmt.Println("You Have Logged in!")
		return false
	}

	loginData := url.Values{}
	loginData.Set("uname", fp.UserID)
	loginData.Set("password", fp.Password)

	resp, err := fp.Session.PostForm(fp.MainURL+"/login", loginData)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}

	resp, err = fp.Session.Get(fp.MainURL + "/login")
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}

	body, _ := ioutil.ReadAll(resp.Body)

	// TODO Extracting csrfToken
	Re := regexp.MustCompile(`"csrf_token":"(.*?)"`)
	csrfTokenMatches := Re.FindStringSubmatch(string(body))
	if len(csrfTokenMatches) > 1 {
		fp.CsrfTokens = csrfTokenMatches[1]
	} else {
		return false
	}

	fp.LoggedIn = true
	return true

}

// @title    Submit
// @description   提交
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    code, probID, lang string 代码，题目id，语言
// @return   string, error 表示提交id、报错消息
func (fp *VIJOS) Submit(code, probID, lang string) (string, error) {
	if !fp.LoggedIn {
		// TODO 尝试登录
		if !fp.Login() {
			return "", fmt.Errorf("登录失败")
		}
	}

	MapLanguage := map[string]string{
		"C":          "c",
		"C++":        "cc",
		"C#":         "cs",
		"Pascal":     "pas",
		"Java":       "java",
		"Python":     "py",
		"Python 3":   "py3",
		"PHP":        "php",
		"Rust":       "rs",
		"Haskell":    "hs",
		"JavaScript": "js",
		"Go":         "go",
		"Ruby":       "rb",
	}

	submitData := url.Values{}
	submitData.Set("lang", MapLanguage[lang])
	submitData.Set("code", code)
	submitData.Set("csrf_token", fp.CsrfTokens)

	submitAddr := fmt.Sprintf("%s/p/%s/submit", fp.MainURL, probID)
	resp, err := fp.Session.PostForm(submitAddr, submitData)
	if err != nil || resp.StatusCode != http.StatusOK {
		// TODO 设置为未登录状态
		fp.LoggedIn = false
		return "", err
	}

	return resp.Request.URL.String(), nil
}

// @title    GetStatus
// @description   跟踪提交状态
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    RunId, channel 提交id, 管道
// @return   string, error 表示提交id、报错消息
func (fp *VIJOS) GetStatus(RunId string, ProbId string, channel chan map[string]string) {

	results := make(map[string]string)
	// TODO 持续请求提交状态
	// TODO 请求次数
	end := false
	for i := 0; i < 200 && !end; i++ {
		// TODO 逐步扩大请求间隔
		time.Sleep(time.Duration(i*i) * time.Second)
		resp, err := fp.Session.Get(RunId)
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
		results, end = VIJOSextractLatestSubmission(string(body))
		// TODO 将读出信息写入管道
		channel <- results
	}
	// TODO 关闭管道
	close(channel)
}

// @title    VIJOSextractLatestSubmission
// @description   分析VIJOS提交表单
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    html, runid string html以及运行id
// @return   map[string]string 表单信息
func VIJOSextractLatestSubmission(html string) (submission map[string]string, end bool) {
	submission = make(map[string]string)
	submission["Html"] = html
	end = false

	// TODO Extracting Result
	resultRe := regexp.MustCompile(`(?s)<span class="record-status--text [a-z]+">\s*([a-z A-Z]+)\s*</span>`)
	resultMatches := resultRe.FindStringSubmatch(html)
	if len(resultMatches) > 1 {
		submission["Result"] = resultMatches[1]
	}

	// TODO Extracting End
	endRe := regexp.MustCompile(`(?s)<span class="record-status--text ([a-z]+)">\s*[a-z A-Z]+\s*</span>`)
	endMatches := endRe.FindStringSubmatch(html)
	if len(endMatches) > 1 {
		if endMatches[1] == "pass" || endMatches[1] == "fail" {
			end = true
		}
	}

	// TODO Extracting Memory
	memoryRe := regexp.MustCompile(`<dt>峰值内存</dt>\s*<dd>([\d.]+) KiB</dd>`)
	memoryMatches := memoryRe.FindStringSubmatch(html)
	if len(memoryMatches) > 1 {
		submission["Memory"] = memoryMatches[1]
	}

	// TODO Extracting Time
	timeRe := regexp.MustCompile(`<dt>总耗时</dt>\s*<dd>(\d*?)ms</dd>`)

	timeMatches := timeRe.FindStringSubmatch(html)
	if len(timeMatches) > 1 {
		submission["Time"] = timeMatches[1]
	}

	return submission, end
}

func NewVIJOS(userID string, password string) Interface.VjudgeInterface {
	jar, _ := cookiejar.New(nil)
	return &VIJOS{
		MainURL:  "https://vijos.org",
		Session:  &http.Client{Jar: jar},
		UserID:   userID,
		Password: password,
		LoggedIn: false,
	}
}
