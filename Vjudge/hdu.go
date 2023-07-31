// @Title  hdu
// @Description  用于操作hdu相关提交
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

// HDU			定义了HDU接口
type HDU struct {
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
func (fp *HDU) Login() bool {

	loginData := url.Values{
		"username": {fp.UserID},
		"userpass": {fp.Password},
		"login":    {"Sign In"},
	}
	_, err := fp.Session.PostForm(fp.MainURL+"/userloginex.php?action=login", loginData)
	if fmt.Sprint(err) != `Post "http://acm.hdu.edu.cn/userloginex.php?action=login": 302 response missing Location header` {
		fmt.Println("Login failed!")
		fp.LoggedIn = false
		return false
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
func (fp *HDU) Submit(code, probID, lang string) (string, error) {
	if !fp.LoggedIn {
		// TODO 尝试登录
		if !fp.Login() {
			return "", fmt.Errorf("登录失败")
		}
	}

	MapLanguage := map[string]string{
		"G++":    "0",
		"GCC":    "1",
		"C++":    "2",
		"C":      "3",
		"Pascal": "4",
		"Java":   "5",
		"C#":     "6",
	}

	// TODO 构建提交表单数据
	submitData := url.Values{
		"check":     {"0"},
		"problemid": {probID},
		"language":  {MapLanguage[lang]},
		"usercode":  {code},
	}

	// TODO 发送POST请求进行提交
	resp, err := fp.Session.PostForm(fp.MainURL+"/submit.php?action=submit", submitData)
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
	status := fmt.Sprintf("%s/status.php?first=&pid=%s&user_id=%s&lang=%s&status=0",
		fp.MainURL, probID, fp.UserID, lang)

	// TODO 先尝试找到RunId
	// TODO 先尝试找到RunId
	resp, err = fp.Session.Get(status)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// TODO 如果出现查找运行id失败
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("frace failed")
	}
	resp.Body.Close()
	results := HDUextractLatestSubmission(string(body))

	return results["Run ID"], nil
}

// @title    GetStatus
// @description   跟踪提交状态
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    RunId, channel 提交id, 管道
// @return   string, error 表示提交id、报错消息
func (fp *HDU) GetStatus(RunId string, ProbId string, channel chan map[string]string) {
	statusURL := fmt.Sprintf("%s/status.php?first=%s", fp.MainURL, RunId)
	Map := map[string]bool{
		"Accepted":              true,
		"Presentation Error":    true,
		"Wrong Answer":          true,
		"Runtime Error":         true,
		"Time Limit Exceeded":   true,
		"Memory Limit Exceeded": true,
		"Output Limit Exceeded": true,
		"Compilation Error":     true,
		"System Error":          true,
		"Out Of Contest Time":   true,
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
		results = HDUextractLatestSubmission(string(body))
		// TODO 将读出信息写入管道
		channel <- results
	}
	// TODO 关闭管道
	close(channel)
}

// @title    HDUextractLatestSubmission
// @description   分析HDU提交表单
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    code, probID, lang string 代码，题目id，语言
// @return   string, error 表示提交id、报错消息
func HDUextractLatestSubmission(html string) map[string]string {
	submission := make(map[string]string)
	submission["Html"] = html

	// TODO Extracting Run ID
	runIDRe := regexp.MustCompile(`<td height=22px>(\d*?)</td>`)
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

	return submission
}

func NewHDU(userID string, password string) Interface.VjudgeInterface {
	jar, _ := cookiejar.New(nil)
	return &HDU{
		MainURL:  "http://acm.hdu.edu.cn",
		Session:  &http.Client{Jar: jar},
		UserID:   userID,
		Password: password,
		LoggedIn: false,
	}
}
