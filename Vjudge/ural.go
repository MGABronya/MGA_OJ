// @Title  ural
// @Description  用于操作ural相关提交
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
	"time"
)

// URAL			定义了URAL接口
type URAL struct {
	Interface.VjudgeInterface              // 包含Vjudge方法
	MainURL                   string       // 主url
	Session                   *http.Client // 会话
	UserID                    string       // 用户id
	Password                  string       // 密码
	LoggedIn                  bool         // 是否登录成功
}

// @title    Login
// @description   获得登录状态
// @auth      MGAronya       2022-9-16 12:15
// @param    password string       接收一个密码
// @return   bool	返回是否登录成功
func (fp *URAL) Login() bool {
	fp.LoggedIn = true

	return fp.LoggedIn
}

// @title    Submit
// @description   提交
// @auth      MGAronya       2022-9-16 12:15
// @param    code, probID, lang string 代码，题目id，语言
// @return   string, error 表示提交id、报错消息
func (fp *URAL) Submit(code, probID, lang string) (string, error) {

	MapLanguage := map[string]string{
		"Ruby 1.9":             "18",
		"Haskell 7.6":          "19",
		"FreePascal 2.6":       "31",
		"Java 1.8":             "32",
		"Scala 2.11":           "33",
		"Python 3.8 x64":       "57",
		"Go 1.14 x64":          "58",
		"Kotlin 1.4.0":         "60",
		"Visual C# 2019":       "61",
		"Visual C 2019":        "63",
		"Visual C++ 2019":      "64",
		"Visual C 2019 x64":    "65",
		"Visual C++ 2019 x64 ": "66",
		"GCC 9.2 x64":          "67",
		"G++ 9.2 x64":          "68",
		"Clang++ 10 x64":       "69",
		"PyPy 3.8 x64":         "71",
		"Rust 1.58 x64":        "72",
	}

	if _, ok := MapLanguage[lang]; !ok {
		return "", fmt.Errorf("language error")
	}

	// TODO 构建提交表单数据
	saveSubmissionURL := fp.MainURL + "/submit.aspx"

	formData := url.Values{
		"Action":     {"submit"},
		"SpaceID":    {"1"},
		"JudgeID":    {fp.Password},
		"Language":   {MapLanguage[lang]},
		"ProblemNum": {probID},
		"Source":     {code},
	}

	// TODO 发送POST请求进行提交
	resp, err := http.PostForm(saveSubmissionURL, formData)
	if err != nil {
		return "", fmt.Errorf("submit failed, try again")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("submit failed, try again")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("submit failed, try again")
	}

	match := fmt.Sprintf(`<TR class=".*?"><TD class="id">(\d+)</TD><TD class="date"><NOBR>.*?</NOBR><BR><NOBR>.*?</NOBR></TD><TD class="coder"><A HREF="author\.aspx\?id=.*?">%s</A></TD><TD class="problem"><A HREF="problem\.aspx\?space=.*?num=%s">%s<SPAN CLASS="problemname">.*?</SPAN></A></TD><TD class="language">.*?</TD><TD class=".*?">.*?</TD><TD class="test"><BR></TD><TD class="runtime">.*?</TD><TD class="memory">.*?</TD></TR>`, fp.UserID, probID, probID)

	runIDs := regexp.MustCompile(match).FindStringSubmatch(string(body))

	if len(runIDs) < 2 {
		return "", fmt.Errorf("runid lose")
	}

	return runIDs[1], nil
}

// @title    GetStatus
// @description   跟踪提交状态
// @auth      MGAronya       2022-9-16 12:15
// @param    RunId, channel 提交id, 管道
// @return   string, error 表示提交id、报错消息
func (fp *URAL) GetStatus(RunId string, ProbId string, channel chan map[string]string) {
	results := make(map[string]string)
	statusURL := fmt.Sprintf("%s/status.aspx?space=1&from=%s", fp.MainURL, RunId)
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
		results := URALextractLatestSubmission(RunId, string(body))
		// TODO 将读出信息写入管道
		channel <- results
	}
	// TODO 关闭管道
	close(channel)
}

// @title    URALextractLatestSubmission
// @description   分析URAL提交表单
// @auth      MGAronya       2022-9-16 12:15
// @param    code, probID, lang string 代码，题目id，语言
// @return   string, error 表示提交id、报错消息
func URALextractLatestSubmission(runid, html string) map[string]string {
	submission := make(map[string]string)
	submission["Html"] = html

	matchResult := fmt.Sprintf(`<TR class=".*?"><TD class="id">%s</TD><TD class="date"><NOBR>.*?</NOBR><BR><NOBR>.*?</NOBR></TD><TD class="coder"><A HREF="author\.aspx\?id=.*?">.*?</A></TD><TD class="problem"><A HREF="problem\.aspx\?space=.*?num=.*?">.*?<SPAN CLASS="problemname">.*?</SPAN></A></TD><TD class="language">.*?</TD><TD class=".*?">(.*?)</TD><TD class=".*?">.*?</TD><TD class="runtime">.*?</TD><TD class="memory">.*?</TD></TR>`, runid)
	matchTime := fmt.Sprintf(`<TR class=".*?"><TD class="id">%s</TD><TD class="date"><NOBR>.*?</NOBR><BR><NOBR>.*?</NOBR></TD><TD class="coder"><A HREF="author\.aspx\?id=.*?">.*?</A></TD><TD class="problem"><A HREF="problem\.aspx\?space=.*?num=.*?">.*?<SPAN CLASS="problemname">.*?</SPAN></A></TD><TD class="language">.*?</TD><TD class=".*?">.*?</TD><TD class=".*?">.*?</TD><TD class="runtime">([\d.]+)</TD><TD class="memory">.*?</TD></TR>`, runid)
	matchMemory := fmt.Sprintf(`<TR class=".*?"><TD class="id">%s</TD><TD class="date"><NOBR>.*?</NOBR><BR><NOBR>.*?</NOBR></TD><TD class="coder"><A HREF="author\.aspx\?id=.*?">.*?</A></TD><TD class="problem"><A HREF="problem\.aspx\?space=.*?num=.*?">.*?<SPAN CLASS="problemname">.*?</SPAN></A></TD><TD class="language">.*?</TD><TD class=".*?">.*?</TD><TD class=".*?">.*?</TD><TD class="runtime">.*?</TD><TD class="memory">(/d+) KB</TD></TR>`, runid)

	// TODO Extracting Result
	Results := regexp.MustCompile(matchResult).FindStringSubmatch(html)

	// TODO Extracting Time
	Times := regexp.MustCompile(matchTime).FindStringSubmatch(html)

	// TODO Extracting Memory
	Memorys := regexp.MustCompile(matchMemory).FindStringSubmatch(html)

	if len(Results) > 1 {
		submission["Result"] = Results[1]
	}

	if len(Times) > 1 {
		submission["Time"] = Times[1]
	} else {
		submission["Time"] = "0"
	}

	if len(Memorys) > 1 {
		submission["Memory"] = Memorys[1]
	} else {
		submission["Memory"] = "0"
	}

	return submission
}

func NewURAL(userID string, password string) Interface.VjudgeInterface {
	jar, _ := cookiejar.New(nil)
	return &URAL{
		MainURL:  "https://acm.timus.ru",
		Session:  &http.Client{Jar: jar},
		UserID:   userID,
		Password: password,
		LoggedIn: true,
	}
}
