// @Title  user
// @Description  定义用户
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package vo

import (
	"MGA_OJ/model"

	uuid "github.com/satori/go.uuid"
)

// user			定义用户
type UserRequest struct {
	Name     string `json:"name"`     // 用户名称
	Email    string `json:"email"`    // 邮箱
	Password string `json:"password"` // 密码
	Verify   string `json:"verify"`   // 验证码
}

// user 		定义用户更新信息
type UserUpdate struct {
	Name    string `json:"name"`    // 用户名称
	Email   string `json:"email"`   // 邮箱
	Blog    string `json:"blog"`    // 博客地址
	Sex     bool   `json:"sex"`     // 性别
	Address string `json:"address"` // 地址
	Verify  string `json:"verify"`  // 验证码
	Icon    string `json:"icon"`    // 头像
}

// UserDto			定义了用户的基本信息
type UserDto struct {
	ID         uuid.UUID `json:"id"`          // 用户id
	Name       string    `json:"name"`        // 用户名称
	Email      string    `json:"email"`       // 邮箱
	Blog       string    `json:"blog"`        // 博客地址
	Sex        bool      `json:"sex"`         // 性别
	Address    string    `json:"address"`     // 地址
	Icon       string    `json:"icon"`        // 这里的Icon存储的是图像文件的地址
	Level      int       `json:"level"`       // 权限等级
	Score      float64   `json:"score"`       // 竞赛分数
	LikeNum    int       `json:"like_num"`    // 收到点赞数量
	UnLikeNum  int       `json:"unlike_num"`  // 收到点踩数量
	CollectNum int       `json:"collect_num"` // 收到收藏数量
	VisitNum   int       `json:"visit_num"`   // 被游览数量
}

// @title    ToUserDto
// @description   用户信息封装
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    user model.User       接收一个用户类
// @return   UserDto			   返回一个用户的基本信息类
func ToUserDto(user model.User) UserDto {
	return UserDto{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Blog:       user.Blog,
		Sex:        user.Sex,
		Address:    user.Address,
		Icon:       user.Icon,
		Level:      user.Level,
		Score:      user.Score,
		LikeNum:    user.LikeNum,
		UnLikeNum:  user.UnLikeNum,
		CollectNum: user.CollectNum,
		VisitNum:   user.VisitNum,
	}
}
