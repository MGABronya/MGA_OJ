// @Title  VjudgeInterface
// @Description  该文件用于封装Vjude方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package Interface

// VjudgeInterface			定义了Vjudge方法
type VjudgeInterface interface {
	Submit(code string, probID string, lang string) (string, error)        // 提交
	GetStatus(RunId string, ProbId string, channel chan map[string]string) // 追踪提交状态
}
