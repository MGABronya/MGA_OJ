// @Title  cf
// @Description  用于操作cf相关提交
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
	"unicode"

	"golang.org/x/net/html"
)

// CF			定义了CF接口
type CF struct {
	Interface.VjudgeInterface              // 包含Vjudge方法
	MainURL                   string       // 主url
	Session                   *http.Client // 会话
	UserID                    string       // 用户id
	Password                  string       // 密码
	LoggedIn                  bool         // 是否登录成功
	Tta                       string       // _tta
}

// @title    Login
// @description   获得登录状态
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    password string       接收一个密码
// @return   bool	返回是否登录成功
func (fp *CF) Login() bool {

	if fp.LoggedIn {
		fmt.Println("You Have Logged in!")
		return false
	}

	resp, err := fp.Session.Get(fp.MainURL + "/enter")
	defer resp.Body.Close()

	if err != nil {
		return false
	}
	doc, err := html.Parse(resp.Body)

	if err != nil {
		return false
	}

	csrfToken := findCsrfToken(doc)
	_39ce7 := find39ce7(resp.Cookies())
	fp.Tta = calcTta(_39ce7)

	formData := url.Values{
		"csrf_token":    {csrfToken},
		"_tta":          {fp.Tta},
		"action":        {"enter"},
		"handleOrEmail": {fp.UserID},
		"password":      {fp.Password},
		"remember":      {"on"},
	}

	resp, err = fp.Session.PostForm(fp.MainURL+"/enter", formData)
	defer resp.Body.Close()
	if err != nil {
		return false
	}

	if resp.StatusCode != http.StatusOK {
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
func (fp *CF) Submit(code, probID, lang string) (string, error) {
	if !fp.LoggedIn {
		// TODO 尝试登录
		if !fp.Login() {
			return "", fmt.Errorf("登录失败")
		}
	}

	resp, err := fp.Session.Get(fp.MainURL + "/problemset/submit")
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// TODO Extracting token
	tokenRe := regexp.MustCompile(`<meta name="X-Csrf-Token" content="(.*?)"/>`)
	csrfTokens := tokenRe.FindStringSubmatch(string(body))
	if len(csrfTokens) < 2 {
		return "", fmt.Errorf("lose csrfToken")
	}
	csrfToken := csrfTokens[1]

	MapLanguage := map[string]string{
		"Delphi":                "3",
		"FPC":                   "4",
		"PHP":                   "6",
		"Python 2":              "7",
		"Mono C#":               "9",
		"Haskell":               "12",
		"Perl":                  "13",
		"Ocaml":                 "19",
		"D":                     "28",
		"Python 3":              "31",
		"Go":                    "32",
		"JavaScript":            "34",
		"PyPy 2":                "40",
		"PyPy 3":                "41",
		"GNU C11":               "43",
		"GNU C++14":             "50",
		"PascalABC.NET":         "51",
		"Clang++17 Diagnostics": "52",
		"GNU C++17":             "54",
		"Node.js":               "55",
		"MS C++ 2017":           "59",
		"GNU C++17 (64)":        "61",
		"C# 8":                  "65",
		"Ruby 3":                "67",
		"PyPy 3-64":             "70",
		"GNU C++20 (64)":        "73",
		"Rust 2021":             "75",
		"Kotlin 1.6":            "77",
		"C# 10":                 "79",
		"Clang++20 Diagnostics": "80",
		"Kotlin 1.7":            "83",
	}

	formData := url.Values{
		"csrf_token":            {csrfToken},
		"_tta":                  {fp.Tta},
		"action":                {"submitSolutionFormSubmitted"},
		"submittedProblemCode":  {probID},
		"programTypeId":         {MapLanguage[lang]},
		"source":                {code},
		"sourceFile":            {""},
		"sourceCodeConfirmed":   {"true"},
		"doNotShowWarningAgain": {"on"},
	}

	resp, err = fp.Session.PostForm(fp.MainURL+"/problemset/submit?csrf_token="+csrfToken, formData)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != http.StatusOK {
		// TODO 设置为未登录状态
		fp.LoggedIn = false
		return "", err
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// TODO Extracting Run ID
	runIDRe := regexp.MustCompile(`data-submission-id="(\d+)"`)
	runIDs := runIDRe.FindStringSubmatch(string(body))

	if len(runIDs) < 2 {
		return "", fmt.Errorf("lose runid, maybe code duplication")
	}

	runID := runIDs[1]

	return runID, nil
}

// @title    GetStatus
// @description   跟踪提交状态
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    RunId, channel 提交id, 管道
// @return   string, error 表示提交id、报错消息
func (fp *CF) GetStatus(RunId string, ProbId string, channel chan map[string]string) {

	results := make(map[string]string)
	// TODO 持续请求提交状态
	// TODO 请求次数
	for i := 0; i < 200 && !ojmatchesRegex(results["Result"]); i++ {
		// TODO 逐步扩大请求间隔
		time.Sleep(time.Duration(i*i) * time.Second)
		resp, err := fp.Session.Get(fmt.Sprintf(fp.MainURL+"/problemset/submission/%s/%s", keepDigitsOnly(ProbId), RunId))
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
		results = CFextractLatestSubmission(string(body))
		// TODO 将读出信息写入管道
		channel <- results
	}
	// TODO 关闭管道
	close(channel)
}

// @title    CFextractLatestSubmission
// @description   分析CF提交表单
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    html, runid string html以及运行id
// @return   map[string]string 表单信息
func CFextractLatestSubmission(html string) (submission map[string]string) {
	submission = make(map[string]string)

	// TODO Extracting Time
	timeRe := regexp.MustCompile(`(\d+).*?ms`)
	Times := timeRe.FindStringSubmatch(html)

	if len(Times) < 2 {
		submission["Time"] = "0"
	} else {
		submission["Time"] = Times[1]
	}

	// TODO Extracting Memory
	memoryRe := regexp.MustCompile(`(\d+).*?KB`)
	Memorys := memoryRe.FindStringSubmatch(html)
	if len(Memorys) < 2 {
		submission["Memory"] = "0"
	} else {
		submission["Memory"] = Memorys[1]
	}

	resultRe := regexp.MustCompile(`<td partyMemberIds=";\d+;" class="status-cell status-small status-verdict-cell" submissionId="\d+" waiting=".*?">\s*(.*?)\s*</td>`)
	Results := resultRe.FindStringSubmatch(html)
	if len(Results) < 2 {
		resultRe = regexp.MustCompile(`<span class='.*?'>(.*?)<span class=".*?">`)
		Results = resultRe.FindStringSubmatch(html)
		if len(Results) > 1 {
			submission["Result"] = Results[1]
		}
	} else if len(Results[1]) > 20 {
		resultRe = regexp.MustCompile(`<span class='.*?'>(.*?)<span class=".*?">`)
		Results = resultRe.FindStringSubmatch(Results[1])
		if len(Results) > 1 {
			submission["Result"] = Results[1]
		}
	} else {
		submission["Result"] = Results[1]
	}
	return
}

// @title    ojmatchesRegex
// @description   查看是否达到终态，通用
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    str string 			状态字符串
// @return   bool 					是否达到终态
func ojmatchesRegex(str string) bool {
	re := regexp.MustCompile(`^ac|^wrong|^time|^memory|^compilation|^runtime|^\d+|err`)
	return re.MatchString(strings.ToLower(str))
}

func NewCF(userID string, password string) Interface.VjudgeInterface {
	jar, _ := cookiejar.New(nil)
	return &CF{
		MainURL:  "https://codeforces.com",
		Session:  &http.Client{Jar: jar},
		UserID:   userID,
		Password: password,
		LoggedIn: false,
	}
}

// @title    calcTta
// @description   计算tta
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    _39ce7 string
// @return   string
func calcTta(_39ce7 string) string {
	n := len(_39ce7)
	c := 0
	_tta := 0
	for c < n {
		_tta = (_tta + (c+1)*(c+2)*int(_39ce7[c])) % 1009
		if c%3 == 0 {
			_tta += 1
		}
		if c%2 == 0 {
			_tta *= 2
		}
		if c > 0 {
			_tta -= (int(_39ce7[c/2]) / 2) * (_tta % 5)
		}
		for _tta < 0 {
			_tta += 1009
		}
		for _tta >= 1009 {
			_tta -= 1009
		}
		c++
	}
	return fmt.Sprintf("%d", _tta)
}

func findCsrfToken(doc *html.Node) string {
	var csrfToken string

	var search func(*html.Node)
	search = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			for _, attr := range n.Attr {
				if attr.Key == "name" && attr.Val == "X-Csrf-Token" {
					for _, attr := range n.Attr {
						if attr.Key == "value" {
							csrfToken = attr.Val
							return
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			search(c)
		}
	}

	search(doc)
	return csrfToken
}

func find39ce7(cookies []*http.Cookie) string {
	for _, cookie := range cookies {
		if cookie.Name == "39ce7" {
			return cookie.Value
		}
	}
	return ""
}

func keepDigitsOnly(str string) string {
	var result []rune

	for _, char := range str {
		if unicode.IsDigit(char) {
			result = append(result, char)
		}
	}

	return string(result)
}
