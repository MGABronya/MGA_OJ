// @Title  vjudge
// @Description  该文件用于初始化vjudge，以及包装一个向外提供vjudge的功能
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package common

import (
	"MGA_OJ/Interface"
	"MGA_OJ/Vjudge"

	"github.com/spf13/viper"
)

// VjudgeMap    定义了vjudge映射
var VjudgeMap map[string]Interface.VjudgeInterface = map[string]Interface.VjudgeInterface{}

// @title    InitVjudge
// @description   从配置文件中读取Vjudge相关信息后，完成Vjudge初始化
// @auth      MGAronya（张健）             2022-9-16 10:07
// @param     void        void         没有入参
// @return    void        void         没有回参
func InitVjudge() {
	VjudgeMap["POJ"] = Vjudge.NewPOJ(viper.GetString("poj.user"), viper.GetString("poj.password"))
	VjudgeMap["HDU"] = Vjudge.NewHDU(viper.GetString("hdu.user"), viper.GetString("hdu.password"))
	VjudgeMap["SPOJ"] = Vjudge.NewSPOJ(viper.GetString("spoj.user"), viper.GetString("spoj.password"))
	VjudgeMap["VIJOS"] = Vjudge.NewVIJOS(viper.GetString("vijos.user"), viper.GetString("vijos.password"))
	VjudgeMap["CF"] = Vjudge.NewCF(viper.GetString("cf.user"), viper.GetString("cf.password"))
	VjudgeMap["ATCODER"] = Vjudge.NewATCODER(viper.GetString("atcoder.user"), viper.GetString("atcoder.password"))
	VjudgeMap["UVA"] = Vjudge.NewUVA(viper.GetString("uva.user"), viper.GetString("uva.password"))
	VjudgeMap["UOJ"] = Vjudge.NewUOJ(viper.GetString("uoj.user"), viper.GetString("uoj.password"))
	VjudgeMap["URAL"] = Vjudge.NewURAL(viper.GetString("ural.user"), viper.GetString("ural.password"))
	VjudgeMap["HACKERRANK"] = Vjudge.NewHACKERRANK(viper.GetString("hackerrank.user"), viper.GetString("hackerrank.password"))
}
