// @Title  Translation
// @Description  定义翻译
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package vo

type TranslationRequest struct {
	Text string `json:"Text"`
}

type TranslationResponse struct {
	Translations []struct {
		Text string `json:"text"`
	} `json:"translations"`
}

type TextRequest struct {
	Text string `json:"Text"`
}
