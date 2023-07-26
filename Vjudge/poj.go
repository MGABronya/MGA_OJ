// @Title  poj
// @Description  用于操作poj相关提交
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

// POJ			定义了POJ接口
type POJ struct {
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
func (fp *POJ) Login() bool {

	// TODO 构建登录表单数据
	loginData := url.Values{}
	loginData.Set("user_id1", fp.UserID)
	loginData.Set("password1", fp.Password)
	loginData.Set("B1", "login")
	loginData.Set("url", ".")

	// TODO 发送POST请求进行登录
	_, err := fp.Session.PostForm(fp.MainURL+"/login", loginData)
	if err != nil {
		fmt.Println("Error during login:", err)
		return false
	}

	fp.LoggedIn = true
	fmt.Println("Login OK!")
	return true
}

// @title    Submit
// @description   提交
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    code, probID, lang string 代码，题目id，语言
// @return   string, error 表示提交id、报错消息
func (fp *POJ) Submit(code, probID, lang string) (string, error) {
	if !fp.LoggedIn {
		// TODO 尝试登录
		if !fp.Login() {
			return "", fmt.Errorf("登录失败")
		}
	}

	MapLanguage := map[string]string{
		"G++":     "0",
		"GCC":     "1",
		"Java":    "2",
		"Pascal":  "3",
		"C++":     "4",
		"C":       "5",
		"Fortran": "6",
	}

	// TODO 构建提交表单数据
	submitData := url.Values{}
	submitData.Set("problem_id", probID)
	submitData.Set("language", MapLanguage[lang])
	submitData.Set("source", code)
	submitData.Set("submit", "Submit")
	submitData.Set("encoded", "0")

	// TODO 发送POST请求进行提交
	resp, err := fp.Session.PostForm(fp.MainURL+"/submit", submitData)
	if err != nil {
		// TODO 设置为未登录状态
		fp.LoggedIn = false
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// TODO 设置为未登录状态
		fp.LoggedIn = false
		return "", err
	}
	statusURL := fmt.Sprintf("%s/status?problem_id=%s&user_id=%s&result=&language=%s",
		fp.MainURL, probID, fp.UserID, lang)

	// TODO 先尝试找到RunId
	resp, err = fp.Session.Get(statusURL)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// TODO 如果出现查找运行id失败
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("trace failed")
	}
	resp.Body.Close()
	results := POJextractLatestSubmission(string(body))

	return results["Run ID"], nil
}

// @title    GetStatus
// @description   跟踪提交状态
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    RunId, channel 提交id, 管道
// @return   string, error 表示提交id、报错消息
func (fp *POJ) GetStatus(RunId string, ProbId string, channel chan map[string]string) {
	runid, _ := strconv.Atoi(RunId)
	statusURL := fmt.Sprintf("%s/status?top=%d", fp.MainURL, runid+1)
	Map := map[string]bool{
		"Accepted":              true,
		"Presentation Error":    true,
		"Time Limit Exceeded":   true,
		"Memory Limit Exceeded": true,
		"Wrong Answer":          true,
		"Runtime Error":         true,
		"Output Limit Exceeded": true,
		"Compile Error":         true,
		"System Error":          true,
		"Validator Error":       true,
	}
	results := make(map[string]string)
	// TODO 持续请求提交状态
	// TODO 请求次数
	for i := 0; !Map[results["Result"]] && i < 200; i++ {
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
		results = POJextractLatestSubmission(string(body))
		// TODO 将读出信息写入管道
		channel <- results
	}
	// TODO 关闭管道
	close(channel)
}

// @title    POJextractLatestSubmission
// @description   分析提交表单
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    code, probID, lang string 代码，题目id，语言
// @return   string, error 表示提交id、报错消息
func POJextractLatestSubmission(html string) map[string]string {
	submission := make(map[string]string)

	// TODO Extracting Run ID
	runIDRe := regexp.MustCompile(`<tr align=center><td>(\d*?)</td>`)
	runIDMatches := runIDRe.FindStringSubmatch(html)
	if len(runIDMatches) > 1 {
		submission["Run ID"] = runIDMatches[1]
	}

	// TODO Extracting Result
	resultRe := regexp.MustCompile(`<font color=(.*?)>(.*?)</font></td>`)
	resultMatches := resultRe.FindStringSubmatch(html)
	if len(resultMatches) > 2 {
		submission["Result"] = resultMatches[2]
	}

	// TODO Extracting Memory
	memoryRe := regexp.MustCompile(`<td>(\d*?)K</td>`)
	memoryMatches := memoryRe.FindStringSubmatch(html)
	if len(memoryMatches) > 1 {
		submission["Memory"] = memoryMatches[1]
	}

	// TODO Extracting Time
	timeRe := regexp.MustCompile(`<td>(\d*?)MS</td>`)
	timeMatches := timeRe.FindStringSubmatch(html)
	if len(timeMatches) > 1 {
		submission["Time"] = timeMatches[1]
	}

	// TODO Extracting Language
	languageRe := regexp.MustCompile(`<a href=showsource\?solution_id=(.*?) target=_blank>(.*?)</a></td>`)
	languageMatches := languageRe.FindStringSubmatch(html)
	if len(languageMatches) > 2 {
		submission["Language"] = languageMatches[2]
	}

	// TODO Extracting Code Length
	codeLengthRe := regexp.MustCompile(`<td>(\d*?)B</td>`)
	codeLengthMatches := codeLengthRe.FindStringSubmatch(html)
	if len(codeLengthMatches) > 1 {
		submission["Code Length"] = codeLengthMatches[1]
	}

	// TODO Extracting Submit Time
	submitTimeRe := regexp.MustCompile(`<td>(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})</td></tr>`)
	submitTimeMatches := submitTimeRe.FindStringSubmatch(html)
	if len(submitTimeMatches) > 1 {
		submission["Submit Time"] = strings.TrimSpace(submitTimeMatches[1])
	}

	return submission
}

func NewPOJ(userID string, password string) Interface.VjudgeInterface {
	jar, _ := cookiejar.New(nil)
	return &POJ{
		MainURL:  "http://poj.org",
		Session:  &http.Client{Jar: jar},
		UserID:   userID,
		Password: password,
		LoggedIn: false,
	}
}
