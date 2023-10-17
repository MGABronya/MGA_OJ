// @Title  atcoder
// @Description  用于操作atcoder相关提交
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:47
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

// ATCODER			定义了ATCODER接口
type ATCODER struct {
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
// @auth      MGAronya       2022-9-16 12:15
// @param    password string       接收一个密码
// @return   bool	返回是否登录成功
func (fp *ATCODER) Login() bool {

	if fp.LoggedIn {
		fmt.Println("You Have Logged in!")
		return false
	}

	resp, err := fp.Session.Get(fp.MainURL + "/login")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return false
	}

	csrfToken, exists := doc.Find("input[name='csrf_token']").Attr("value")
	if !exists {
		return false
	}

	formData := url.Values{
		"csrf_token": {csrfToken},
		"username":   {fp.UserID},
		"password":   {fp.Password},
	}

	resp, err = fp.Session.PostForm(fp.MainURL+"/login", formData)
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
// @auth      MGAronya       2022-9-16 12:15
// @param    code, probID, lang string 代码，题目id，语言
// @return   string, error 表示提交id、报错消息
func (fp *ATCODER) Submit(code, probID, lang string) (string, error) {
	if !fp.LoggedIn {
		// TODO 尝试登录
		if !fp.Login() {
			return "", fmt.Errorf("登录失败")
		}
	}

	parts := strings.Split(probID, "_")

	if len(parts) != 2 {
		return "", fmt.Errorf("proId error")
	}

	resp, err := fp.Session.Get(fmt.Sprintf("%s/contests/%s/submit", fp.MainURL, parts[0]))
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	csrfToken, exists := doc.Find("input[name='csrf_token']").Attr("value")
	if !exists {
		return "", fmt.Errorf("CSRF token not found")
	}

	MapLanguage := map[string]string{
		"C (GCC 9.2.1)":                "4001",
		"C++ (GCC 9.2.1)":              "4003",
		"Python (3.8.2)":               "4006",
		"Haskell (GHC 8.8.3)":          "4027",
		"Haxe (4.0.3); Java":           "4029",
		"Julia (1.4.0)":                "4031",
		"Lua (Lua 5.3.5)":              "4033",
		"Dash (0.5.8)":                 "4035",
		"Ruby (2.7.1)":                 "4049",
		"Standard ML (MLton 20130715)": "4054",
		"Text (cat 8.28)":              "4056",
		"Unlambda (2.0.0)":             "4064",
		"Sed (4.4)":                    "4066",
	}

	if _, ok := MapLanguage[lang]; !ok {
		return "", fmt.Errorf("language error")
	}

	formData := url.Values{
		"csrf_token":          {csrfToken},
		"data.TaskScreenName": {probID},
		"data.LanguageId":     {MapLanguage[lang]},
		"sourceCode":          {code},
		"input-open-file":     {""},
	}

	resp, err = fp.Session.PostForm(fmt.Sprintf("%s/contests/%s/submit", fp.MainURL, parts[0]), formData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err != nil || resp.StatusCode != http.StatusOK {
		// TODO 设置为未登录状态
		fp.LoggedIn = false
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// TODO Extracting Run ID
	runIDRe := regexp.MustCompile(`data-id="(\d+)"`)
	runIDs := runIDRe.FindStringSubmatch(string(body))

	if len(runIDs) < 2 {
		return "", fmt.Errorf("lose runid, maybe code duplication")
	}

	runID := runIDs[1]

	return runID, nil
}

// @title    GetStatus
// @description   跟踪提交状态
// @auth      MGAronya       2022-9-16 12:15
// @param    RunId, channel 提交id, 管道
// @return   string, error 表示提交id、报错消息
func (fp *ATCODER) GetStatus(RunId string, ProbId string, channel chan map[string]string) {

	results := make(map[string]string)
	parts := strings.Split(ProbId, "_")

	if len(parts) != 2 {
		results["Result"] = "system error for problemid"
		return
	}
	// TODO 持续请求提交状态
	// TODO 请求次数
	for i := 0; i < 200 && !ojmatchesRegex(results["Result"]); i++ {
		// TODO 逐步扩大请求间隔
		time.Sleep(time.Duration(i*i) * time.Second)

		resp, err := fp.Session.Get(fmt.Sprintf("%s/contests/%s/submissions/%s", fp.MainURL, parts[0], RunId))
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
		results = ATCODERextractLatestSubmission(string(body))
		// TODO 将读出信息写入管道
		channel <- results
	}
	// TODO 关闭管道
	close(channel)
}

// @title    ATCODERextractLatestSubmission
// @description   分析ATCODER提交表单
// @auth      MGAronya       2022-9-16 12:15
// @param    html, runid string html以及运行id
// @return   map[string]string 表单信息
func ATCODERextractLatestSubmission(html string) (submission map[string]string) {
	submission = make(map[string]string)

	submission["Html"] = html

	// TODO Extracting Time
	timeRe := regexp.MustCompile(`<tr>\s*<th>Exec Time</th>\s*<td class="text-center">(\d+) ms</td>\s*</tr>`)
	Times := timeRe.FindStringSubmatch(html)

	if len(Times) < 2 {
		submission["Time"] = "0"
	} else {
		submission["Time"] = Times[1]
	}

	// TODO Extracting Memory
	memoryRe := regexp.MustCompile(`<tr>\s*<th>Memory</th>\s*<td class="text-center">(\d+) KB</td>\s*</tr>`)
	Memorys := memoryRe.FindStringSubmatch(html)
	if len(Memorys) < 2 {
		submission["Memory"] = "0"
	} else {
		submission["Memory"] = Memorys[1]
	}

	resultRe := regexp.MustCompile(`<th>Status</th>\s*<td id="judge-status" class="text-center"><span class='label .*?' data-toggle='tooltip' data-placement='top' title="(.*?)">.*?</span></td>\s*</tr>`)
	Results := resultRe.FindStringSubmatch(html)
	if len(Results) < 2 {
		submission["Result"] = "Waiting"
	} else {
		submission["Result"] = Results[1]
	}
	return
}

func NewATCODER(userID string, password string) Interface.VjudgeInterface {
	jar, _ := cookiejar.New(nil)
	return &ATCODER{
		MainURL:  "https://atcoder.jp",
		Session:  &http.Client{Jar: jar},
		UserID:   userID,
		Password: password,
		LoggedIn: false,
	}
}
