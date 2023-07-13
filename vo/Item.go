// @Title  item
// @Description  xml格式题目
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package vo

type Fps struct {
	Item Item `xml:"item"`
}

type Item struct {
	Title        string     `xml:"title"`
	TimeLimit    string     `xml:"time_limit"`
	MemoryLimit  string     `xml:"memory_limit"`
	Description  string     `xml:"description"`
	Input        string     `xml:"input"`
	Output       string     `xml:"output"`
	SampleInput  []string   `xml:"sample_input"`
	SampleOutput []string   `xml:"sample_output"`
	TestInput    []string   `xml:"test_input"`
	TestOutput   []string   `xml:"test_output"`
	Hint         string     `xml:"hint"`
	Source       string     `xml:"source"`
	Solutions    []Solution `xml:"solution"`
}

type Solution struct {
	Language string `xml:"language,attr"`
	Code     string `xml:",cdata"`
}
