// @Title  RecordController
// @Description  该文件提供关于操作提交的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/common"
	"MGA_OJ/model"
	rabbitMq "MGA_OJ/rabbitMq"
	"MGA_OJ/response"
	"MGA_OJ/vo"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IRecordController			定义了提交类接口
type IRecordController interface {
	Create(ctx *gin.Context)   // 用户进行提交操作
	Show(ctx *gin.Context)     // 用户查看指定提交
	PageList(ctx *gin.Context) // 用户搜索提交列表
}

// RecordController			定义了提交工具类
type RecordController struct {
	DB *gorm.DB // 含有一个数据库指针
}

var rabbitmq *rabbitMq.RabbitMQ = rabbitMq.NewRabbitMQSimple("MGAronya")

// @title    Create
// @description   用户进行提交操作
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RecordController) Create(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)
	var requestRecord vo.RecordRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestRecord); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 查看当前problem状态
	var problem model.Problem
	if r.DB.Where("id = ?", requestRecord.ProblemId).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}

	// TODO 创建提交
	record := model.Record{
		UserId:        user.ID,
		ProblemId:     requestRecord.ProblemId,
		Language:      requestRecord.Language,
		Code:          requestRecord.Code,
		Condition:     "Waiting",
		CompetitionId: problem.CompetitionId,
		Pass:          0,
	}

	// TODO 插入数据
	if err := r.DB.Create(&record).Error; err != nil {
		response.Fail(ctx, nil, "提交上传出错，数据验证有误")
		return
	}

	// TODO 加入消息队列
	if err := rabbitmq.PublishSimple(fmt.Sprint(record.ID)); err != nil {
		response.Fail(ctx, nil, "消息队列出错")
		return
	}

	// TODO 成功
	response.Success(ctx, nil, "提交成功")
}

// @title    Show
// @description   查看一篇提交的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RecordController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var record model.Record

	// TODO 查看提交是否在数据库中存在
	if r.DB.Where("id = ?", id).First(&record).Error != nil {
		response.Fail(ctx, nil, "提交不存在")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否为比赛提交且提交是否可见
	var competition model.Competition
	if r.DB.Where("id = ?", record.CompetitionId).First(&competition) == nil {
		// TODO 如果比赛未结束且已经开始
		if competition.EndTime.String() > time.Now().Format("2006-01-02 15:04:05") && competition.StartTime.String() <= time.Now().Format("2006-01-02 15:04:05") {
			// TODO 查看是否有权限查看代码
			if competition.Type == "Single" {
				if record.UserId != user.ID {
					record.Code = ""
				}
			} else if competition.Type == "Group" {
				var set model.Set
				if r.DB.Where("id = ?", competition.SetId).First(&set).Error != nil {
					response.Fail(ctx, nil, "表单不存在")
					return
				}

				// TODO 查看用户所在组别
				results := make([]map[string]interface{}, 0)
				r.DB.Table("group_lists").Select("group_lists.group_id").Joins("left join user_lists on user_lists.group_id = group_lists.group_id and user_lists.user_id = ? and group_lists.set_id = ?", user.ID, set.ID).Scan(&results)

				group_id := results[0]["group_id"].(uint)
				// TODO 和提交者不在同一组，则无权查看具体代码
				if r.DB.Where("group_id = ? and user_id = ?", group_id, record.UserId).First(&model.UserList{}).Error != nil {
					record.Code = ""
				}
			}
		}
	}

	response.Success(ctx, gin.H{"record": record}, "成功")
}

// @title    PageList
// @description   获取多篇提交
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RecordController) PageList(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取查询条件
	Luanguage := ctx.DefaultQuery("luanguage", "")
	UserId := ctx.DefaultQuery("user_id", "")
	ProblemId := ctx.DefaultQuery("problem_id", "")
	StartTime := ctx.DefaultQuery("start_time", "")
	EndTime := ctx.DefaultQuery("end_time", "")
	Condition := ctx.DefaultQuery("condition", "")
	PassLow := ctx.DefaultQuery("pass_low", "")
	PassTop := ctx.DefaultQuery("pass_top", "")
	CompetitionId := ctx.DefaultQuery("competition_id", "0")

	db := common.GetDB()

	// TODO 根据参数设置where条件
	if Luanguage != "" {
		db = db.Where("luanguage = ?", Luanguage)
	}
	if UserId != "" {
		db = db.Where("user_id = ?", UserId)
	}
	if ProblemId != "" {
		db = db.Where("problem_id = ?", ProblemId)
	}
	if StartTime != "" {
		db = db.Where("created_at >= ?", StartTime)
	}
	if EndTime != "" {
		db = db.Where("created_at <= ?", EndTime)
	}
	if Condition != "" {
		db = db.Where("condition = ?", Condition)
	}
	if PassLow != "" {
		db = db.Where("pass >= ?", PassLow)
	}
	if PassTop != "" {
		db = db.Where("pass <= ?", PassTop)
	}
	if CompetitionId != "" {
		db = db.Where("competition_id = ?", CompetitionId)
	}

	// TODO 分页
	var records []model.Record

	var total int64

	// TODO 查找所有分页中可见的条目
	db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&records).Count(&total)

	if CompetitionId != "-1" {
		// TODO 查看比赛是否还在进行
		var competition model.Competition
		if db.Where("id = ?", CompetitionId).First(&competition).Error != nil {
			response.Fail(ctx, nil, "比赛不存在")
			return
		}
		// TODO 如果比赛未结束，且比赛已经开始
		if competition.EndTime.String() > time.Now().Format("2006-01-02 15:04:05") && competition.StartTime.String() <= time.Now().Format("2006-01-02 15:04:05") {
			// TODO 禁止查看代码
			for i := range records {
				// TODO 查看是否有权限查看代码
				if competition.Type == "Single" {
					if records[i].UserId != user.ID {
						records[i].Code = ""
					}
				} else if competition.Type == "Group" {
					var set model.Set
					if r.DB.Where("id = ?", competition.SetId).First(&set).Error != nil {
						response.Fail(ctx, nil, "表单不存在")
						return
					}

					// TODO 查看用户所在组别
					results := make([]map[string]interface{}, 0)
					r.DB.Table("group_lists").Select("group_lists.group_id").Joins("left join user_lists on user_lists.group_id = group_lists.group_id and user_lists.user_id = ? and group_lists.set_id = ?", user.ID, set.ID).Scan(&results)

					group_id := results[0]["group_id"].(uint)
					// TODO 和提交者不在同一组，则无权查看具体代码
					if r.DB.Where("group_id = ? and user_id = ?", group_id, records[i].UserId).First(&model.UserList{}).Error != nil {
						records[i].Code = ""
					}
				}
			}
		}
	}

	// TODO 返回数据
	response.Success(ctx, gin.H{"records": records, "total": total}, "成功")
}

// @title    NewRecordController
// @description   新建一个IRecordController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IRecordController		返回一个IRecordController用于调用各种函数
func NewRecordController() IRecordController {
	db := common.GetDB()
	db.AutoMigrate(model.Record{})
	db.AutoMigrate(model.Case{})
	return RecordController{DB: db}
}
