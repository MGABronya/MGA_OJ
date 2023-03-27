// @Title  MassageInterface
// @Description  该文件用于封装信息交流方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package Interface

import "github.com/gin-gonic/gin"

// MassageInterface		定义了信息交流方法
type MassageInterface interface {
	Send(ctx *gin.Context)        // 发送私信
	LinkList(ctx *gin.Context)    // 列出连接列表
	ChatList(ctx *gin.Context)    // 列出聊天列表
	Receive(ctx *gin.Context)     // 建立实时接收
	ReceiveLink(ctx *gin.Context) // 建立连接实时接收
	RemoveLink(ctx *gin.Context)  // 移除某个连接
}
