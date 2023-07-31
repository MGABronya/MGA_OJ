// @Title  hackerrank
// @Description  用于操作hackerrank相关提交
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
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// HACKERRANK			定义了HACKERRANK接口
type HACKERRANK struct {
	Interface.VjudgeInterface              // 包含Vjudge方法
	MainURL                   string       // 主url
	Session                   *http.Client // 会话
	UserID                    string       // 用户id
	Password                  string       // 密码
	LoggedIn                  bool         // 是否登录成功
	UserAgent                 string       // header
	CsrfToken                 string       //token
}

// @title    Login
// @description   获得登录状态
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    password string       接收一个密码
// @return   bool	返回是否登录成功
func (fp *HACKERRANK) Login() bool {
	loginURL := fp.MainURL + "/auth/login"
	req, err := http.NewRequest("GET", loginURL, nil)
	if err != nil {
		return false
	}
	req.Header.Set("User-Agent", fp.UserAgent)
	res, err := fp.Session.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()

	// Extract CSRF token from response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return false
	}
	csrfToken, ok := doc.Find("meta[name=csrf-token]").Attr("content")
	if !ok {
		return false
	}

	// Send POST request to login endpoint
	authURL := "https://www.hackerrank.com/rest/auth/login"
	formData := url.Values{
		"fallback":    {"true"},
		"login":       {fp.UserID},
		"password":    {fp.Password},
		"remember_me": {"true"},
	}
	req, err = http.NewRequest("POST", authURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-CSRF-Token", csrfToken)
	req.Header.Set("User-Agent", fp.UserAgent)
	res, err = fp.Session.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fp.LoggedIn = false
	} else {
		fp.LoggedIn = true
	}

	return fp.LoggedIn
}

// @title    Submit
// @description   提交
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    code, probID, lang string 代码，题目id，语言
// @return   string, error 表示提交id、报错消息
func (fp *HACKERRANK) Submit(code, probID, lang string) (string, error) {
	if !fp.LoggedIn {
		// TODO 尝试登录
		if !fp.Login() {
			return "", fmt.Errorf("登录失败")
		}
	}

	MapLanguage := map[string]string{
		"C":          "c",
		"Clojure":    "clojure",
		"C++11":      "cpp11",
		"C++14":      "cpp14",
		"C++20":      "cpp20",
		"Erlang":     "erlang",
		"Go":         "go",
		"Haskell":    "haskell",
		"Java7":      "java7",
		"Java8":      "java8",
		"Java15":     "java15",
		"Julia":      "julia",
		"Kotlin":     "kotlin",
		"Lua":        "lua",
		"Perl":       "perl",
		"PHP":        "php",
		"Pypy3":      "pypy3",
		"Python3":    "python3",
		"R":          "r",
		"Ruby":       "ruby",
		"Rust":       "rust",
		"Scala":      "scala",
		"Swift":      "swift",
		"TypeScript": "typescript",
	}

	problemURL := fmt.Sprintf("%s/challenges/%s/problem", fp.MainURL, probID)
	req, err := http.NewRequest("GET", problemURL, nil)
	if err != nil {
		// TODO 设置为未登录状态
		fp.LoggedIn = false
		return "", err
	}
	req.Header.Set("User-Agent", fp.UserAgent)
	res, err := fp.Session.Do(req)
	if err != nil {
		// TODO 设置为未登录状态
		fp.LoggedIn = false
		return "", err
	}
	defer res.Body.Close()

	// Extract CSRF token from response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// TODO 设置为未登录状态
		fp.LoggedIn = false
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		// TODO 设置为未登录状态
		fp.LoggedIn = false
		return "", err
	}
	csrfToken, ok := doc.Find("meta[name=csrf-token]").Attr("content")
	if !ok {
		// TODO 设置为未登录状态
		fp.LoggedIn = false
		return "", err
	}
	fp.CsrfToken = csrfToken

	// TODO 构建提交表单数据
	submissionURL := fmt.Sprintf("%s/rest/contests/master/challenges/%s/submissions", fp.MainURL, probID)

	formData := url.Values{
		"code":          {code},
		"contest_slug":  {"master"},
		"language":      {MapLanguage[lang]},
		"playlist_slug": {""},
	}

	// TODO 发送POST请求进行提交
	req, err = http.NewRequest("POST", submissionURL, strings.NewReader(formData.Encode()))
	if err != nil {
		// TODO 设置为未登录状态
		fp.LoggedIn = false
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-CSRF-Token", csrfToken)
	req.Header.Set("User-Agent", fp.UserAgent)
	res, err = fp.Session.Do(req)
	if err != nil {
		// TODO 设置为未登录状态
		fp.LoggedIn = false
		return "", err
	}
	defer res.Body.Close()

	// TODO 先尝试找到RunId
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		// TODO 设置为未登录状态
		fp.LoggedIn = false
		return "", err
	}

	runIDs := regexp.MustCompile(`{"model":{"id":(\d+)`).FindStringSubmatch(string(body))

	if len(runIDs) < 2 {
		// TODO 设置为未登录状态
		fp.LoggedIn = false
		return "", fmt.Errorf("runid lose")
	}

	return runIDs[1], nil
}

// @title    GetStatus
// @description   跟踪提交状态
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    RunId, channel 提交id, 管道
// @return   string, error 表示提交id、报错消息
func (fp *HACKERRANK) GetStatus(RunId string, ProbId string, channel chan map[string]string) {
	statusURL := fmt.Sprintf("%s/challenges/%s/submissions/code/%s", fp.MainURL, ProbId, RunId)
	results := make(map[string]string)
	// TODO 持续请求提交状态
	// TODO 请求次数
	for i := 0; !ojmatchesRegex(results["Result"]) && i < 200; i++ {
		// TODO 逐步扩大请求间隔
		time.Sleep(time.Duration(i*i) * time.Second)
		req, err := http.NewRequest("GET", statusURL, nil)
		if err != nil {
			fmt.Println("Get Status Failed:", err)
			continue
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("X-CSRF-Token", fp.CsrfToken)
		req.Header.Set("User-Agent", fp.UserAgent)
		res, err := fp.Session.Do(req)
		if err != nil {
			fmt.Println("Get Status Failed:", err)
			continue
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Get Status Failed:", err)
			continue
		}
		res.Body.Close()
		results = HACKERRANKextractLatestSubmission(string(body))
		// TODO 将读出信息写入管道
		channel <- results
	}
	// TODO 关闭管道
	close(channel)
}

// @title    HACKERRANKextractLatestSubmission
// @description   分析HACKERRANK提交表单
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    code, probID, lang string 代码，题目id，语言
// @return   string, error 表示提交id、报错消息
func HACKERRANKextractLatestSubmission(html string) map[string]string {
	submission := make(map[string]string)
	submission["Html"] = html

	// TODO Extracting Result
	resultRe := regexp.MustCompile(`<label class="label">Status:</label><span class=".*?">(.*?)</span>`)
	resultMatches := resultRe.FindStringSubmatch(html)
	if len(resultMatches) > 2 {
		submission["Result"] = resultMatches[2]
	}

	return submission
}

func NewHACKERRANK(userID string, password string) Interface.VjudgeInterface {
	jar, _ := cookiejar.New(nil)
	return &HACKERRANK{
		MainURL:   "https://www.hackerrank.com",
		Session:   &http.Client{Jar: jar},
		UserID:    userID,
		Password:  password,
		LoggedIn:  false,
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.104 Safari/537.36",
	}
}
