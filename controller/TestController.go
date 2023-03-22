// @Title  TestController
// @Description  该文件提供关于提交测试的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/response"
	"MGA_OJ/util"
	"MGA_OJ/vo"
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ITestController			定义了测试类接口
type ITestController interface {
	Create(ctx *gin.Context) // 包含创建功能
}

// TestController			定义了测试工具类
type TestController struct {
	rw *sync.RWMutex
}

// @title    Create
// @description   创建一篇测试
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TestController) Create(ctx *gin.Context) {
	// TODO 单核处理，此处上锁
	t.rw.Lock()
	// TODO 此处确保资源归还
	defer t.rw.Unlock()
	var requestTest vo.TestRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestTest); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}
	// TODO 进行测试

	output, condition, memory, time := Test(requestTest)

	// TODO 返回测试状态
	response.Success(ctx, gin.H{"output": output, "condition": condition, "memory": memory, "time": time}, "测试创建成功")
}

// @title    NewTestController
// @description   新建一个ITestController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ITestController		返回一个ITestController用于调用各种函数
func NewTestController() ITestController {
	return TestController{rw: &sync.RWMutex{}}
}

// @title    Test
// @description   查看测试输出情况
// @auth      MGAronya（张健）             2022-9-16 10:49
// @param     record model.Record, cmdI Interface.CmdInterface			提交记录以及判题方法
// @return    void			没有回参
func Test(requestTest vo.TestRequest) (output string, condition string, memory uint64, spand int64) {
	// TODO 找出语言对应运行方法
	// TODO 找到提交记录后，开始判题逻辑
	cmdI, ok := util.LanguageMap[requestTest.Language]
	if !ok {
		condition = "Luanguage Error"
		return
	}
	// TODO 清空当前文件夹
	dir, err := ioutil.ReadDir("/user-code/")
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{"user-code/", d.Name()}...))
	}
	// id		定义文件名
	id := cmdI.Name()
	fp, err := os.Create("./user-code/" + id + "." + cmdI.Suffix())
	// TODO 文件错误
	if err != nil {
		// TODO 创建文件失败的原因有：
		// TODO 1、路径不存在  2、权限不足  3、打开文件数量超过上限  4、磁盘空间不足等
		condition = "System Error 1"
		return
	}

	// TODO defer延迟调用 关闭文件，释放资源
	defer fp.Close()

	// TODO 写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(fp)

	write.WriteString(requestTest.Code)

	// TODO Flush将缓存的文件真正写入到文件中
	write.Flush()

	// TODO 编译
	cmd := cmdI.Compile("./user-code/", id)

	// TODO 系统错误
	if err := cmd.Start(); err != nil {
		condition = "System Error 2"
		return
	}
	// TODO 启动routine等待结束
	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	// 设定超时时间，并select它
	after := time.After(time.Duration(20 * time.Second))
	select {
	// TODO 编译超时
	case <-after:
		cmd.Process.Kill()
		condition = "Compile timeout"
		return
	case err = <-done:
	}

	// TODO 编译出错
	if err != nil {
		condition = "Compile Error"
		return
	}

	var bm runtime.MemStats
	runtime.ReadMemStats(&bm)

	// TODO 运行可执行文件
	cmd = cmdI.Run("./user-code/", id)

	var out, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out
	stdinPipe, err := cmd.StdinPipe()
	// TODO 系统错误
	if err != nil {
		condition = "System Error 3"
		return
	}
	io.WriteString(stdinPipe, requestTest.Input)
	now := time.Now().UnixMilli()
	// TODO 系统错误
	if err := cmd.Start(); err != nil {
		condition = "System Error 4"
		return
	}
	// TODO 启动routine等待结束
	done = make(chan error)
	go func() { done <- cmd.Wait() }()

	// 设定超时时间，并select它
	after = time.After(time.Duration(20000*cmdI.TimeMultiplier()+cmdI.RunUpTime()) * time.Millisecond)
	select {
	// TODO 运行超时
	case <-after:
		cmd.Process.Kill()
		condition = "Time Limit Exceeded"
		return
	case err = <-done:
	}
	end := time.Now().UnixMilli()

	// TODO 运行时错误
	if err != nil {
		condition = "Runtime Error"
		return
	}
	// TODO 记录使用时间
	spand = int64(math.Max(float64(end-now-int64(cmdI.RunUpTime())), 0))
	// TODO 记录使用空间
	var em runtime.MemStats
	runtime.ReadMemStats(&em)
	memory = em.Alloc/1024 - bm.Alloc/1024

	condition = "ok"
	output = out.String()
	return
}
