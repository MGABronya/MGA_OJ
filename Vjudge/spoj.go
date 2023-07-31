// @Title  spoj
// @Description  用于操作spoj相关提交
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
	"strconv"
	"strings"
	"time"
)

// SPOJ			定义了SPOJ接口
type SPOJ struct {
	Interface.VjudgeInterface              // 包含Vjudge方法
	MainURL                   string       // 主url
	Session                   *http.Client // 会话
	UserID                    string       // 用户id
	Password                  string       // 密码
	LoggedIn                  bool         // 是否登录成功
}

// @title    Login
// @description   获得登录状态
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    password string       接收一个密码
// @return   bool	返回是否登录成功
func (fp *SPOJ) Login() bool {

	if fp.LoggedIn {
		fmt.Println("You Have Logged in!")
		return false
	}

	loginData := url.Values{}
	loginData.Set("login_user", fp.UserID)
	loginData.Set("password", fp.Password)
	loginData.Set("autologin", "1")
	loginData.Set("next_raw", "/")

	response, err := fp.Session.PostForm(fp.MainURL+"/login", loginData)
	if err != nil {
		fmt.Println("Login failed!")
		return false
	}
	defer response.Body.Close()

	fp.LoggedIn = true
	return fp.LoggedIn
}

// @title    Submit
// @description   提交
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    code, probID, lang string 代码，题目id，语言
// @return   string, error 表示提交id、报错消息
func (fp *SPOJ) Submit(code, probID, lang string) (string, error) {
	if !fp.LoggedIn {
		// TODO 尝试登录
		if !fp.Login() {
			return "", fmt.Errorf("登录失败")
		}
	}

	MapLanguage := map[string]string{
		"CPP":        "1",
		"PAS-GPC":    "2",
		"PERL":       "3",
		"PYTHON":     "4",
		"FORTRAN":    "5",
		"WHITESPACE": "6",
		"ADA95":      "7",
		"OCAML":      "8",
		"ICK":        "9",
		"JAVA":       "10",
		"C":          "11",
		"BF":         "12",
		"ASM32":      "13",
		"CLPS":       "14",
		"PRLG-swi":   "15",
		"ICON":       "16",
		"RUBY":       "17",
		"SCM qobi":   "18",
		"PIKE":       "19",
		"D":          "20",
		"HASK":       "21",
		"PAS-FPC":    "22",
		"ST":         "23",
		"JAR":        "24",
		"NICE":       "25",
		"LUA":        "26",
		"CSHARP":     "27",
		"BASH":       "28",
		"PHP":        "29",
		"NEM":        "30",
		"LISP sbcl":  "31",
		"LISP clisp": "32",
		"SCM guile":  "33",
		"C99":        "34",
		"JS-RHINO":   "35",
		"ERL":        "36",
		"TCL":        "38",
		"SCALA":      "39",
		"SQLITE":     "40",
		"C++ 4.3.2":  "41",
		"ASM64":      "42",
		"OBJC":       "43",
		"CPP14":      "44",
		"ASM32-GCC":  "45",
		"SED":        "46",
		"KTLN":       "47",
		"DART":       "48",
		"VB.NET":     "50",
		"PERL6":      "54",
		"NODEJS":     "56",
		"DOC":        "59",
		"PDF":        "60",
	}

	submitData := url.Values{}
	submitData.Set("file", code)
	submitData.Set("lang", MapLanguage[lang])
	submitData.Set("problemcode", probID)
	submitData.Set("submit", "Submit!")

	response, err := fp.Session.PostForm(fp.MainURL+"/submit/complete/", submitData)
	if err != nil {
		// TODO 设置为未登录状态
		fp.LoggedIn = false
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		// TODO 设置为未登录状态
		fp.LoggedIn = false
		return "", err
	}
	body, _ := ioutil.ReadAll(response.Body)

	// TODO 先尝试找到RunId
	resultRe := regexp.MustCompile(`<input type="hidden" name="newSubmissionId" value="(\d*?)"/>`)
	resultMatches := resultRe.FindStringSubmatch(string(body))
	if len(resultMatches) > 1 {
		return resultMatches[1], nil
	}

	return "", fmt.Errorf("runid lose")
}

// @title    GetStatus
// @description   跟踪提交状态
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    RunId, channel 提交id, 管道
// @return   string, error 表示提交id、报错消息
func (fp *SPOJ) GetStatus(RunId string, ProbId string, channel chan map[string]string) {
	page := 0

	results := make(map[string]string)
	// TODO 持续请求提交状态
	// TODO 请求次数
	for i := 0; !spojmatchesRegex(results["Result"]) && i < 200; i++ {
		// TODO 逐步扩大请求间隔
		time.Sleep(time.Duration(i*i) * time.Second)
		statusURL := fmt.Sprintf("%s/status/%s,%s/all/start=%d", fp.MainURL, ProbId, fp.UserID, page)
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
		results = SPOJextractLatestSubmission(string(body), RunId)
		// TODO 未找到结果，翻到下一页
		if results == nil {
			page += 20
			continue
		}
		// TODO 将读出信息写入管道
		channel <- results
	}
	// TODO 关闭管道
	close(channel)
}

// @title    SPOJextractLatestSubmission
// @description   分析SPOJ提交表单
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    html, runid string html以及运行id
// @return   map[string]string 表单信息
func SPOJextractLatestSubmission(html, runid string) map[string]string {
	result := make(map[string]string)
	result["Html"] = html
	// TODO Extracting Result
	str := fmt.Sprintf(`(?s)<td class="statusres text-center" id="statusres_%s" status="\d*?" final="\d*?" manual="\d*?">\s*(.*?)\s*<span class="small">`, runid)
	re := regexp.MustCompile(str)
	match := re.FindStringSubmatch(html)

	if len(match) < 2 {
		return nil
	}

	if len(match[1]) > 25 {
		re := regexp.MustCompile(`<a href=".*?" title="View details.">\s*(.*?)\s*</a>`)
		match = re.FindStringSubmatch(match[1])
	}
	result["Result"] = match[1]

	// TODO Extracting Memory
	str = fmt.Sprintf(`(?s)<td class="smemory statustext text-center" id="statusmem_%s">\s*([\d.]+)M\s*</td>`, runid)
	re = regexp.MustCompile(str)
	match = re.FindStringSubmatch(html)
	if len(match) >= 2 {
		memory, _ := strconv.ParseFloat(match[1], 64)
		memory *= 1024
		result["Memory"] = strconv.Itoa(int(memory))
	} else {
		result["Memory"] = "0"
	}

	str = fmt.Sprintf(`(?s)<td class="stime statustext text-center" id="statustime_%s">\s*<a href="/ranks/.*?" title="See the best solutions">\s*(\d+\.\d+)\s*</a>\s*</td>`, runid)
	// TODO Extracting Time
	re = regexp.MustCompile(str)
	match = re.FindStringSubmatch(html)
	if len(match) >= 2 {
		time, _ := strconv.ParseFloat(match[1], 64)
		time *= 1000
		result["Time"] = strconv.Itoa(int(time))
	} else {
		result["Time"] = "0"
	}

	return result
}

// @title    matchesRegex
// @description   查看是否达到终态
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    str string 			状态字符串
// @return   bool 					是否达到终态
func spojmatchesRegex(str string) bool {
	re := regexp.MustCompile(`^accepted|^memory|^wrong|^time|^compilation|^runtime|^\d+`)
	return re.MatchString(str)
}

func NewSPOJ(userID string, password string) Interface.VjudgeInterface {
	jar, _ := cookiejar.New(nil)
	return &SPOJ{
		MainURL:  "https://www.spoj.com",
		Session:  &http.Client{Jar: jar},
		UserID:   strings.ToLower(userID),
		Password: password,
		LoggedIn: false,
	}
}
