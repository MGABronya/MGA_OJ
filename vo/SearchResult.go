// @Title  SearchResult
// @Description  定义搜索返回结果
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package vo

type SearchResult struct {
	DateLastCrawled  string `json:"dateLastCrawled"`
	DisplayUrl       string `json:"displayUrl"`
	ID               string `json:"id"`
	IsFamilyFriendly bool   `json:"isFamilyFriendly"`
	IsNavigational   bool   `json:"isNavigational"`
	Language         string `json:"language"`
	Name             string `json:"name"`
	Snippet          string `json:"snippet"`
	Url              string `json:"url"`
}
