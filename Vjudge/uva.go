// @Title  uva
// @Description  用于操作uva相关提交
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:47
package Vjudge

import (
	"MGA_OJ/Interface"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// UVA			定义了UVA接口
type UVA struct {
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
func (fp *UVA) Login() bool {

	r, err := fp.Session.Get(fp.MainURL)
	if err != nil {
		return false
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return false
	}

	hiddenData := make(map[string]string)
	doc.Find("input[type='hidden']").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		value, _ := s.Attr("value")
		hiddenData[name] = value
	})

	loginURL := fp.MainURL + "/index.php?option=com_comprofiler&task=login"

	formData := make(map[string]string)
	formData["username"] = fp.UserID
	formData["passwd"] = fp.Password
	formData["remember"] = "yes"
	formData["Submit"] = "Login"
	for k, v := range hiddenData {
		formData[k] = v
	}

	body := strings.NewReader(encodeFormData(formData))
	r, err = fp.Session.Post(loginURL, "application/x-www-form-urlencoded", body)
	if err != nil {
		fp.LoggedIn = false
		return false
	} else {
		fp.LoggedIn = true
	}
	r.Body.Close()

	return fp.LoggedIn
}

// @title    Submit
// @description   提交
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    code, probID, lang string 代码，题目id，语言
// @return   string, error 表示提交id、报错消息
func (fp *UVA) Submit(code, probID, lang string) (string, error) {
	if !fp.LoggedIn {
		// TODO 尝试登录
		if !fp.Login() {
			return "", fmt.Errorf("登录失败")
		}
	}

	MapLanguage := map[string]string{
		"ANSI C": "1",
		"JAVA":   "2",
		"C++":    "3",
		"PASCAL": "4",
		"C++11":  "5",
		"PYTH3":  "6",
	}

	// TODO 构建提交表单数据
	saveSubmissionURL := fp.MainURL + "/index.php?option=com_onlinejudge&Itemid=25&page=save_submission"

	formData := make(map[string]string)
	formData["problemid"] = probID
	formData["category"] = ""
	formData["localid"] = "100"
	formData["language"] = MapLanguage[lang]
	formData["code"] = code

	body := strings.NewReader(encodeFormData(formData))

	// TODO 发送POST请求进行提交
	r, err := fp.Session.Post(saveSubmissionURL, "application/x-www-form-urlencoded", body)
	if err != nil {
		if err != nil {
			// TODO 设置为未登录状态
			fp.LoggedIn = false
			return "", err
		}
	}
	defer r.Body.Close()

	rbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if err != nil {
			// TODO 设置为未登录状态
			fp.LoggedIn = false
			return "", err
		}
	}

	// TODO Extracting Run ID
	runIDRe := regexp.MustCompile(`Submission received with ID (\d+)`)
	runIDs := runIDRe.FindStringSubmatch(string(rbody))

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
func (fp *UVA) GetStatus(RunId string, ProbId string, channel chan map[string]string) {
	page := 0
	results := make(map[string]string)
	// TODO 持续请求提交状态
	// TODO 请求次数
	for i := 0; !ojmatchesRegex(results["Result"]) && i < 200; i++ {
		// TODO 逐步扩大请求间隔
		time.Sleep(time.Duration(i*i) * time.Second)
		statusURL := fmt.Sprintf("%s/index.php?option=com_onlinejudge&Itemid=9&limit=50&limitstart=%d", fp.MainURL, page)
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
		results, find := UVAextractLatestSubmission(RunId, string(body))
		// TODO 没找到则翻页
		if !find {
			page += 50
			continue
		}
		// TODO 将读出信息写入管道
		channel <- results
	}
	// TODO 关闭管道
	close(channel)
}

// @title    UVAextractLatestSubmission
// @description   分析UVA提交表单
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    code, probID, lang string 代码，题目id，语言
// @return   string, error 表示提交id、报错消息
func UVAextractLatestSubmission(runid, html string) (map[string]string, bool) {
	submission := make(map[string]string)
	submission["Html"] = html

	// TODO Extracting Run ID
	RunIDRe := regexp.MustCompile(`<td>\s*<a href="index.php\?option=.*?Itemid=.*?page=.*?submission=(\d+)">.*?</a>\s*</td>\s*<td>.*?</td>\s*<td>[\d.]+</td>`)
	RunIDs := RunIDRe.FindAllStringSubmatch(html, -1)

	// TODO Extracting Result
	ReultIDRe := regexp.MustCompile(`<td>\s*<a href="index.php\?option=.*?Itemid=.*?page=.*?submission=\d+">(.*?)</a>\s*</td>\s*<td>.*?</td>\s*<td>[\d.]+</td>`)
	ResultIDs := ReultIDRe.FindAllStringSubmatch(html, -1)

	// TODO Extracting Time
	TimeRe := regexp.MustCompile(`<td>\s*<a href="index.php\?option=.*?Itemid=.*?page=.*?submission=\d+">.*?</a>\s*</td>\s*<td>.*?</td>\s*<td>([\d.]+)</td>`)
	Times := TimeRe.FindAllStringSubmatch(html, -1)

	for i := range RunIDs {
		if RunIDs[i][1] == runid {
			submission["Result"] = ResultIDs[i][1]
			submission["Time"] = Times[i][1]
			return submission, true
		}
	}

	return submission, false
}

func encodeFormData(data map[string]string) string {
	var encodedData []string
	for k, v := range data {
		encodedData = append(encodedData, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(encodedData, "&")
}

func NewUVA(userID string, password string) Interface.VjudgeInterface {
	jar, _ := cookiejar.New(nil)
	return &UVA{
		MainURL:  "https://onlinejudge.org",
		Session:  &http.Client{Jar: jar},
		UserID:   userID,
		Password: password,
		LoggedIn: false,
	}
}
