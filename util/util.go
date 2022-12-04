// @Title  util
// @Description  收集各种需要使用的工具函数
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:47
package util

import (
	"MGA_OJ/common"
	"MGA_OJ/model"
	"context"
	"fmt"
	"math/rand"
	"net/smtp"
	"regexp"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/jordan-wright/email"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var Units = map[string]uint{
	"mb": 1024,
	"kb": 1,
	"gb": 1024 * 1024,
	"ms": 1,
	"s":  1000,
}

var ctx context.Context = context.Background()

// @title    GetH
// @description   在redis中的一个哈希中获取值
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, H string, key string        k表示选用第几个库，H为哈希，key为在H中的key
// @return   string     	  返回对应的value
func GetH(k int, H string, key string) string {
	client := common.GetRedisClient(k)
	level, _ := client.HGet(ctx, H, key).Result()
	return level
}

// @title    SETH
// @description   在redis中的一个哈希中设置值
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, H string, key string, value string       k表示选用第几个库，H为哈希，key为在H中的key，value为要设置的对应值
// @return   void
func SetH(k int, H string, key string, value string) {
	client := common.GetRedisClient(k)
	client.HSet(ctx, H, key, value)
}

// @title    SetS
// @description   在redis中的一个集合中设置值
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, S string, value string       k表示选用的第几个库，S为集合，value为要设置的对应值
// @return   void
func SetS(k int, S string, value string) {
	client := common.GetRedisClient(k)
	client.SAdd(ctx, S, value)
}

// @title    RemS
// @description   在redis中的一个集合中删除值
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, S string, value string       k表示选用的第几个库，S为集合，value为要删除的对应值
// @return   void
func RemS(k int, S string, value string) {
	client := common.GetRedisClient(k)
	client.SRem(ctx, S, value)
}

// @title    IsS
// @description   在redis中的一个集合中查找某个元素是否存在
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, S string, value string       k表示选用的第几个库，S为集合，value为要查找的对应值
// @return   bool 表示value是否在S中
func IsS(k int, S string, value string) bool {
	client := common.GetRedisClient(k)
	flag, _ := client.SIsMember(ctx, S, value).Result()
	return flag
}

// @title    MembersS
// @description   在redis中的一个集合中查找所有元素
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, S string      k表示选用的第几个库，S为集合
// @return   []string		表示该集合的所有元素
func MembersS(k int, S string) []string {
	client := common.GetRedisClient(k)
	es, _ := client.SMembers(ctx, S).Result()
	return es
}

// @title    InterS
// @description   在redis中集合的交集
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, S string      k表示选用的第几个库，keys为集合
// @return   []string		表示交集中的所有元素
func InterS(k int, keys ...string) []string {
	client := common.GetRedisClient(k)
	es, _ := client.SInter(ctx, keys...).Result()
	return es
}

// @title    UnionS
// @description   在redis中集合的并集
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, S string      k表示选用的第几个库，keys为集合
// @return   []string		表示并集中的所有元素
func UnionS(k int, keys ...string) []string {
	client := common.GetRedisClient(k)
	es, _ := client.SUnion(ctx, keys...).Result()
	return es
}

// @title    CardS
// @description  查看redis中一个集合中元素的个数
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, Key string      k表示选用的第几个库，Key为要查看的集合
// @return   int	表示该集合元素的个数
func CardS(k int, Key string) int {
	client := common.GetRedisClient(k)
	cnt, _ := client.SCard(ctx, Key).Result()
	return int(cnt)
}

// @title    AddZ
// @description  设置redis中一个有序集合中一个元素
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, key string, value string, cost float64      k表示选用的第几个库，key为要设置的有序集合，value表示值，cost表示权重
// @return   int	表示该集合元素的个数
func AddZ(k int, key string, value string, cost float64) {
	client := common.GetRedisClient(k)
	client.ZAdd(ctx, key, redis.Z{Score: cost, Member: value})
}

// @title    ScoreZ
// @description  查看redis中一个有序集合中一个元素的权重
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, key string, value string      k表示选用的第几个库，key为要设置的有序集合，value表示值
// @return   int	表示该集合元素的个数
func ScoreZ(k int, key string, value string) float64 {
	client := common.GetRedisClient(k)
	cost, _ := client.ZScore(ctx, key, value).Result()
	return cost
}

// @title    IncrByZ
// @description  修改redis中一个有序集合中一个元素的权重
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, key string, value string, cost float64      k表示选用的第几个库，key为要设置的有序集合，value表示值，cost表示权重
// @return   int	表示该集合元素的个数
func IncrByZ(k int, key string, value string, cost float64) {
	client := common.GetRedisClient(k)
	client.ZIncrBy(ctx, key, cost, value)
}

// @title    RangeZ
// @description  查询redis中一个有序集合中一个范围重的值
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, key string, left int64, right int64      k表示选用的第几个库，key为要设置的有序集合，查询[left, right]中的值
// @return   []string	表示该集合范围中的元素
func RangeZ(k int, key string, left int64, right int64) []string {
	client := common.GetRedisClient(k)
	res, _ := client.ZRevRange(ctx, key, left, right).Result()
	return res
}

// @title    RangeWithScoreZ
// @description  查询redis中一个有序集合中一个范围重的值和分数
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, key string, left int64, right int64      k表示选用的第几个库，key为要设置的有序集合，查询[left, right]中的值
// @return   []redis.Z	表示该集合范围中的元素和分数
func RangeWithScoreZ(k int, key string, left int64, right int64) []redis.Z {
	client := common.GetRedisClient(k)
	res, _ := client.ZRevRangeWithScores(ctx, key, left, right).Result()
	return res
}

// @title    CardZ
// @description  查看redis中一个有序集合中元素的个数
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, key string      k表示选用的第几个库，key为要cx查询的有序集合
// @return   int	表示该集合元素的个数
func CardZ(k int, key string) int64 {
	client := common.GetRedisClient(k)
	res, _ := client.ZCard(ctx, key).Result()
	return res
}

// @title    RemZ
// @description  移除redis中一个有序集合中一个元素
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, key string, value string      k表示选用的第几个库，key为要设置的有序集合，value表示值
// @return   int	表示该集合元素的个数
func RemZ(k int, key string, value string) {
	client := common.GetRedisClient(k)
	client.ZRem(ctx, key, value)
}

// @title    Del
// @description  删除redis中的一个键
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, Key string      k表示选用的第几个库，Key为要删除的键
// @return   void, void
func Del(k int, Key string) {
	client := common.GetRedisClient(k)
	client.Del(ctx, Key)
}

// @title    DelH
// @description  删除redis中的一个哈希下的键
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    k int, h string, key string      k表示选用的第几个库，h表示哪个哈希表，key为要删除的键
// @return   void, void
func DelH(k int, h string, key string) {
	client := common.GetRedisClient(k)
	client.HDel(ctx, h, key)
}

// @title    RandomString
// @description   生成一段随机的字符串
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     n int		字符串的长度
// @return    string    一串随机的字符串
func RandomString(n int) string {
	var letters = []byte("qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	// TODO 不断用随机字母填充字符串
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// @title    VerifyEmailFormat
// @description   用于验证邮箱格式是否正确的工具函数
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     email string		一串字符串，表示邮箱
// @return    bool    返回是否合法
func VerifyEmailFormat(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// @title    VerifyMobileFormat
// @description   用于验证手机号格式是否正确的工具函数
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     mobileNum string		一串字符串，表示手机号
// @return    bool    返回是否合法
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

// @title    VerifyQQFormat
// @description   用于验证QQ号格式是否正确的工具函数
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     QQNum string		一串字符串，表示QQ
// @return    bool    返回是否合法
func VerifyQQFormat(QQNum string) bool {
	regular := "[1-9][0-9]{4,10}"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(QQNum)
}

// @title    VerifyQQFormat
// @description  用于验证Icon是否为默认图片的工具函数
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     Icon string		一串字符串，表示图像名称
// @return    bool    返回是否合法
func VerifyIconFormat(Icon string) bool {
	regular := "MGA[1-9].jpg"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(Icon)
}

// @title    isEmailExist
// @description   查看email是否在数据库中存在
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func IsEmailExist(db *gorm.DB, email string) bool {
	var user model.User
	db.Where("email = ?", email).First(&user)
	return user.ID != 0
}

// @title    isNameExist
// @description   查看name是否在数据库中存在
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func IsNameExist(db *gorm.DB, name string) bool {
	var user model.User
	db.Where("name = ?", name).First(&user)
	return user.ID != 0
}

// @title    SendEmailValidate
// @description   发送验证邮件
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    em []string       接收一个邮箱字符串
// @return   string, error     返回验证码和error值
func SendEmailValidate(em []string) (string, error) {
	mod := `
	尊敬的%s，您好！

	您于 %s 提交的邮箱验证，本次验证码为%s，为了保证账号安全，验证码有效期为5分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。
	此邮箱为系统邮箱，请勿回复。
`
	e := email.NewEmail()
	e.From = "mgAronya <2829214609@qq.com>"
	e.To = em
	// TODO 生成6位随机验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vCode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	t := time.Now().Format("2006-01-02 15:04:05")
	// TODO 设置文件发送的内容
	content := fmt.Sprintf(mod, em[0], t, vCode)
	e.Text = []byte(content)
	// TODO 设置服务器相关的配置
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "2829214609@qq.com", "rmdtxokuuqyrdgii", "smtp.qq.com"))
	return vCode, err
}

// @title    SendEmailPass
// @description   发送密码邮件
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    em []string       接收一个邮箱字符串
// @return   string, error     返回验证码和error值
func SendEmailPass(em []string) string {
	mod := `
	尊敬的%s，您好！

	您于 %s 提交的邮箱验证，已经将密码重置为%s，为了保证账号安全。切勿向他人泄露，并尽快更改密码，感谢您的理解与使用。
	此邮箱为系统邮箱，请勿回复。
`
	e := email.NewEmail()
	e.From = "mgAronya <2829214609@qq.com>"
	e.To = em
	// TODO 生成8位随机密码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	password := fmt.Sprintf("%08v", rnd.Int31n(100000000))
	t := time.Now().Format("2006-01-02 15:04:05")

	db := common.GetDB()

	// TODO 创建密码哈希
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "密码加密失败"
	}

	// TODO 更新密码
	err = db.Model(&model.User{}).Where("email = ?", em[0]).Updates(model.User{
		Password: string(hasedPassword),
	}).Error

	if err != nil {
		return "密码更新失败"
	}

	// TODO 设置文件发送的内容
	content := fmt.Sprintf(mod, em[0], t, password)
	e.Text = []byte(content)
	// TODO 设置服务器相关的配置
	err = e.Send("smtp.qq.com:25", smtp.PlainAuth("", "2829214609@qq.com", "rmdtxokuuqyrdgii", "smtp.qq.com"))

	if err != nil {
		return "邮件发送失败"
	}

	return "密码已重置"
}

// @title    IsEmailPass
// @description   验证邮箱是否通过
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    em []string       接收一个邮箱字符串
// @return   string, error     返回验证码和error值
func IsEmailPass(email string, vertify string) bool {
	client := common.GetRedisClient(0)
	V, err := client.Get(ctx, email).Result()
	if err != nil {
		return false
	}
	return V == vertify
}

// @title    SetRedisEmail
// @description   设置验证码，并令其存活五分钟
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    email string, v string       接收一个邮箱和一个验证码
// @return   void
func SetRedisEmail(email string, v string) {
	client := common.GetRedisClient(0)

	client.Set(ctx, email, v, 300*time.Second)
}
