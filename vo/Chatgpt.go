// @Title  chatgpt
// @Description  定义chatgpt消息
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package vo

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
	Finish  string  `json:"finish_reason"`
	Index   int     `json:"index"`
}
