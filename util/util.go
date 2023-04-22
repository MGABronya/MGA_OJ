// @Title  util
// @Description  收集各种需要使用的工具函数
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:47
package util

import (
	"MGA_OJ/Interface"
	Handle "MGA_OJ/Language"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/smtp"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Units		定义了单位换算
var Units = map[string]uint{
	"mb": 1024,
	"kb": 1,
	"gb": 1024 * 1024,
	"ms": 1,
	"s":  1000,
}

var Max_run int = 4

// timerMap	    定义了当前使用的定时器
var TimerMap map[uuid.UUID]*time.Timer = make(map[uuid.UUID]*time.Timer)

// LanguageMap			定义语言字典，对应其处理方式
var LanguageMap map[string]Interface.CmdInterface = map[string]Interface.CmdInterface{
	"C":          Handle.NewC(),
	"C#":         Handle.NewCs(),
	"C++":        Handle.NewCppPlusPlus(),
	"C++11":      Handle.NewCppPlusPlus11(),
	"Erlang":     Handle.NewErlang(),
	"Go":         Handle.NewGo(),
	"Java":       Handle.NewJava(),
	"JavaScript": Handle.NewJavaScript(),
	"Kotlin":     Handle.NewKotlin(),
	"Pascal":     Handle.NewPascal(),
	"PHP":        Handle.NewPHP(),
	"Python":     Handle.NewPython(),
	"Racket":     Handle.NewRacket(),
	"Ruby":       Handle.NewRuby(),
	"Rust":       Handle.NewRust(),
	"Swift":      Handle.NewSwift(),
}

// MgaronyaString			mgaronya字符串
var MgaronyaString []string = []string{
	`

	███▄ ▄███▓  ▄████  ▄▄▄       ██▀███   ▒█████   ███▄    █ ▓██   ██▓ ▄▄▄      
	▓██▒▀█▀ ██▒ ██▒ ▀█▒▒████▄    ▓██ ▒ ██▒▒██▒  ██▒ ██ ▀█   █  ▒██  ██▒▒████▄    
	▓██    ▓██░▒██░▄▄▄░▒██  ▀█▄  ▓██ ░▄█ ▒▒██░  ██▒▓██  ▀█ ██▒  ▒██ ██░▒██  ▀█▄  
	▒██    ▒██ ░▓█  ██▓░██▄▄▄▄██ ▒██▀▀█▄  ▒██   ██░▓██▒  ▐▌██▒  ░ ▐██▓░░██▄▄▄▄██ 
	▒██▒   ░██▒░▒▓███▀▒ ▓█   ▓██▒░██▓ ▒██▒░ ████▓▒░▒██░   ▓██░  ░ ██▒▓░ ▓█   ▓██▒
	░ ▒░   ░  ░ ░▒   ▒  ▒▒   ▓▒█░░ ▒▓ ░▒▓░░ ▒░▒░▒░ ░ ▒░   ▒ ▒    ██▒▒▒  ▒▒   ▓▒█░
	░  ░      ░  ░   ░   ▒   ▒▒ ░  ░▒ ░ ▒░  ░ ▒ ▒░ ░ ░░   ░ ▒░ ▓██ ░▒░   ▒   ▒▒ ░
	░      ░   ░ ░   ░   ░   ▒     ░░   ░ ░ ░ ░ ▒     ░   ░ ░  ▒ ▒ ░░    ░   ▒   
		   ░         ░       ░  ░   ░         ░ ░           ░  ░ ░           ░  ░
															   ░ ░               
	
	`,
	`
	
	_____                    _____                    _____                    _____                   _______                   _____                _____                    _____          
	/\    \                  /\    \                  /\    \                  /\    \                 /::\    \                 /\    \              |\    \                  /\    \         
   /::\____\                /::\    \                /::\    \                /::\    \               /::::\    \               /::\____\             |:\____\                /::\    \        
  /::::|   |               /::::\    \              /::::\    \              /::::\    \             /::::::\    \             /::::|   |             |::|   |               /::::\    \       
 /:::::|   |              /::::::\    \            /::::::\    \            /::::::\    \           /::::::::\    \           /:::::|   |             |::|   |              /::::::\    \      
/::::::|   |             /:::/\:::\    \          /:::/\:::\    \          /:::/\:::\    \         /:::/~~\:::\    \         /::::::|   |             |::|   |             /:::/\:::\    \     
/:::/|::|   |            /:::/  \:::\    \        /:::/__\:::\    \        /:::/__\:::\    \       /:::/    \:::\    \       /:::/|::|   |             |::|   |            /:::/__\:::\    \    
/:::/ |::|   |           /:::/    \:::\    \      /::::\   \:::\    \      /::::\   \:::\    \     /:::/    / \:::\    \     /:::/ |::|   |             |::|   |           /::::\   \:::\    \   
/:::/  |::|___|______    /:::/    / \:::\    \    /::::::\   \:::\    \    /::::::\   \:::\    \   /:::/____/   \:::\____\   /:::/  |::|   | _____       |::|___|______    /::::::\   \:::\    \  
/:::/   |::::::::\    \  /:::/    /   \:::\ ___\  /:::/\:::\   \:::\    \  /:::/\:::\   \:::\____\ |:::|    |     |:::|    | /:::/   |::|   |/\    \      /::::::::\    \  /:::/\:::\   \:::\    \ 
/:::/    |:::::::::\____\/:::/____/  ___\:::|    |/:::/  \:::\   \:::\____\/:::/  \:::\   \:::|    ||:::|____|     |:::|    |/:: /    |::|   /::\____\    /::::::::::\____\/:::/  \:::\   \:::\____\
\::/    / ~~~~~/:::/    /\:::\    \ /\  /:::|____|\::/    \:::\  /:::/    /\::/   |::::\  /:::|____| \:::\    \   /:::/    / \::/    /|::|  /:::/    /   /:::/~~~~/~~      \::/    \:::\  /:::/    /
\/____/      /:::/    /  \:::\    /::\ \::/    /  \/____/ \:::\/:::/    /  \/____|:::::\/:::/    /   \:::\    \ /:::/    /   \/____/ |::| /:::/    /   /:::/    /          \/____/ \:::\/:::/    / 
		/:::/    /    \:::\   \:::\ \/____/            \::::::/    /         |:::::::::/    /     \:::\    /:::/    /            |::|/:::/    /   /:::/    /                    \::::::/    /  
	   /:::/    /      \:::\   \:::\____\               \::::/    /          |::|\::::/    /       \:::\__/:::/    /             |::::::/    /   /:::/    /                      \::::/    /   
	  /:::/    /        \:::\  /:::/    /               /:::/    /           |::| \::/____/         \::::::::/    /              |:::::/    /    \::/    /                       /:::/    /    
	 /:::/    /          \:::\/:::/    /               /:::/    /            |::|  ~|                \::::::/    /               |::::/    /      \/____/                       /:::/    /     
	/:::/    /            \::::::/    /               /:::/    /             |::|   |                 \::::/    /                /:::/    /                                    /:::/    /      
   /:::/    /              \::::/    /               /:::/    /              \::|   |                  \::/____/                /:::/    /                                    /:::/    /       
   \::/    /                \::/____/                \::/    /                \:|   |                   ~~                      \::/    /                                     \::/    /        
	\/____/                                           \/____/                  \|___|                                            \/____/                                       \/____/         
																																															   

	`,
	`
	
	.----------------.  .----------------.  .----------------.  .----------------.  .----------------.  .-----------------. .----------------.  .----------------. 
	| .--------------. || .--------------. || .--------------. || .--------------. || .--------------. || .--------------. || .--------------. || .--------------. |
	| | ____    ____ | || |    ______    | || |      __      | || |  _______     | || |     ____     | || | ____  _____  | || |  ____  ____  | || |      __      | |
	| ||_   \  /   _|| || |  .' ___  |   | || |     /  \     | || | |_   __ \    | || |   .'    '.   | || ||_   \|_   _| | || | |_  _||_  _| | || |     /  \     | |
	| |  |   \/   |  | || | / .'   \_|   | || |    / /\ \    | || |   | |__) |   | || |  /  .--.  \  | || |  |   \ | |   | || |   \ \  / /   | || |    / /\ \    | |
	| |  | |\  /| |  | || | | |    ____  | || |   / ____ \   | || |   |  __ /    | || |  | |    | |  | || |  | |\ \| |   | || |    \ \/ /    | || |   / ____ \   | |
	| | _| |_\/_| |_ | || | \ '.___]  _| | || | _/ /    \ \_ | || |  _| |  \ \_  | || |  \  '--'  /  | || | _| |_\   |_  | || |    _|  |_    | || | _/ /    \ \_ | |
	| ||_____||_____|| || |  '._____.'   | || ||____|  |____|| || | |____| |___| | || |   '.____.'   | || ||_____|\____| | || |   |______|   | || ||____|  |____|| |
	| |              | || |              | || |              | || |              | || |              | || |              | || |              | || |              | |
	| '--------------' || '--------------' || '--------------' || '--------------' || '--------------' || '--------------' || '--------------' || '--------------' |
	 '----------------'  '----------------'  '----------------'  '----------------'  '----------------'  '----------------'  '----------------'  '----------------' 
	
	`,
	`
	
.------..------..------..------..------..------..------..------.
|M.--. ||G.--. ||A.--. ||R.--. ||O.--. ||N.--. ||Y.--. ||A.--. |
| (\/) || :/\: || (\/) || :(): || :/\: || :(): || (\/) || (\/) |
| :\/: || :\/: || :\/: || ()() || :\/: || ()() || :\/: || :\/: |
| '--'M|| '--'G|| '--'A|| '--'R|| '--'O|| '--'N|| '--'Y|| '--'A|
'------''------''------''------''------''------''------''------'

	`,
	`
	
                                                                                                                                                                                       
                                                                                                                                                                                       
MMMMMMMM               MMMMMMMM        GGGGGGGGGGGGG               AAA                                                                                                                 
M:::::::M             M:::::::M     GGG::::::::::::G              A:::A                                                                                                                
M::::::::M           M::::::::M   GG:::::::::::::::G             A:::::A                                                                                                               
M:::::::::M         M:::::::::M  G:::::GGGGGGGG::::G            A:::::::A                                                                                                              
M::::::::::M       M::::::::::M G:::::G       GGGGGG           A:::::::::A           rrrrr   rrrrrrrrr      ooooooooooo   nnnn  nnnnnnnn    yyyyyyy           yyyyyyy  aaaaaaaaaaaaa   
M:::::::::::M     M:::::::::::MG:::::G                        A:::::A:::::A          r::::rrr:::::::::r   oo:::::::::::oo n:::nn::::::::nn   y:::::y         y:::::y   a::::::::::::a  
M:::::::M::::M   M::::M:::::::MG:::::G                       A:::::A A:::::A         r:::::::::::::::::r o:::::::::::::::on::::::::::::::nn   y:::::y       y:::::y    aaaaaaaaa:::::a 
M::::::M M::::M M::::M M::::::MG:::::G    GGGGGGGGGG        A:::::A   A:::::A        rr::::::rrrrr::::::ro:::::ooooo:::::onn:::::::::::::::n   y:::::y     y:::::y              a::::a 
M::::::M  M::::M::::M  M::::::MG:::::G    G::::::::G       A:::::A     A:::::A        r:::::r     r:::::ro::::o     o::::o  n:::::nnnn:::::n    y:::::y   y:::::y        aaaaaaa:::::a 
M::::::M   M:::::::M   M::::::MG:::::G    GGGGG::::G      A:::::AAAAAAAAA:::::A       r:::::r     rrrrrrro::::o     o::::o  n::::n    n::::n     y:::::y y:::::y       aa::::::::::::a 
M::::::M    M:::::M    M::::::MG:::::G        G::::G     A:::::::::::::::::::::A      r:::::r            o::::o     o::::o  n::::n    n::::n      y:::::y:::::y       a::::aaaa::::::a 
M::::::M     MMMMM     M::::::M G:::::G       G::::G    A:::::AAAAAAAAAAAAA:::::A     r:::::r            o::::o     o::::o  n::::n    n::::n       y:::::::::y       a::::a    a:::::a 
M::::::M               M::::::M  G:::::GGGGGGGG::::G   A:::::A             A:::::A    r:::::r            o:::::ooooo:::::o  n::::n    n::::n        y:::::::y        a::::a    a:::::a 
M::::::M               M::::::M   GG:::::::::::::::G  A:::::A               A:::::A   r:::::r            o:::::::::::::::o  n::::n    n::::n         y:::::y         a:::::aaaa::::::a 
M::::::M               M::::::M     GGG::::::GGG:::G A:::::A                 A:::::A  r:::::r             oo:::::::::::oo   n::::n    n::::n        y:::::y           a::::::::::aa:::a
MMMMMMMM               MMMMMMMM        GGGGGG   GGGGAAAAAAA                   AAAAAAA rrrrrrr               ooooooooooo     nnnnnn    nnnnnn       y:::::y             aaaaaaaaaa  aaaa
                                                                                                                                                  y:::::y                              
                                                                                                                                                 y:::::y                               
                                                                                                                                                y:::::y                                
                                                                                                                                               y:::::y                                 
                                                                                                                                              yyyyyyy                                  
                                                                                                                                                                                       
                                                                                                                                                                                       

	`,
	`
	
	___           ___           ___           ___           ___           ___                         ___     
	/  /\         /  /\         /  /\         /  /\         /  /\         /  /\          __           /  /\    
   /  /::|       /  /::\       /  /::\       /  /::\       /  /::\       /  /::|        |  |\        /  /::\   
  /  /:|:|      /  /:/\:\     /  /:/\:\     /  /:/\:\     /  /:/\:\     /  /:|:|        |  |:|      /  /:/\:\  
 /  /:/|:|__   /  /:/  \:\   /  /::\ \:\   /  /::\ \:\   /  /:/  \:\   /  /:/|:|__      |  |:|     /  /::\ \:\ 
/__/:/_|::::\ /__/:/_\_ \:\ /__/:/\:\_\:\ /__/:/\:\_\:\ /__/:/ \__\:\ /__/:/ |:| /\     |__|:|__  /__/:/\:\_\:\
\__\/  /~~/:/ \  \:\__/\_\/ \__\/  \:\/:/ \__\/~|::\/:/ \  \:\ /  /:/ \__\/  |:|/:/     /  /::::\ \__\/  \:\/:/
	  /  /:/   \  \:\ \:\        \__\::/     |  |:|::/   \  \:\  /:/      |  |:/:/     /  /:/~~~~      \__\::/ 
	 /  /:/     \  \:\/:/        /  /:/      |  |:|\/     \  \:\/:/       |__|::/     /__/:/           /  /:/  
	/__/:/       \  \::/        /__/:/       |__|:|~       \  \::/        /__/:/      \__\/           /__/:/   
	\__\/         \__\/         \__\/         \__\|         \__\/         \__\/                       \__\/    

	`,
	`
	
	__   __  _______  _______  ______    _______  __    _  __   __  _______ 
	|  |_|  ||       ||   _   ||    _ |  |       ||  |  | ||  | |  ||   _   |
	|       ||    ___||  |_|  ||   | ||  |   _   ||   |_| ||  |_|  ||  |_|  |
	|       ||   | __ |       ||   |_||_ |  | |  ||       ||       ||       |
	|       ||   ||  ||       ||    __  ||  |_|  ||  _    ||_     _||       |
	| ||_|| ||   |_| ||   _   ||   |  | ||       || | |   |  |   |  |   _   |
	|_|   |_||_______||__| |__||___|  |_||_______||_|  |__|  |___|  |__| |__|
	
	`,
	`
	
__/\\\\____________/\\\\_        _____/\\\\\\\\\\\\_        _____/\\\\\\\\\____        _______________        _______________        _______________        _______________        ________________        
_\/\\\\\\________/\\\\\\_        ___/\\\//////////__        ___/\\\\\\\\\\\\\__        _______________        _______________        _______________        _______________        ________________       
 _\/\\\//\\\____/\\\//\\\_        __/\\\_____________        __/\\\/////////\\\_        _______________        _______________        _______________        ____/\\\__/\\\_        ________________      
  _\/\\\\///\\\/\\\/_\/\\\_        _\/\\\____/\\\\\\\_        _\/\\\_______\/\\\_        __/\\/\\\\\\\__        _____/\\\\\____        __/\\/\\\\\\___        ___\//\\\/\\\__        __/\\\\\\\\\____     
   _\/\\\__\///\\\/___\/\\\_        _\/\\\___\/////\\\_        _\/\\\\\\\\\\\\\\\_        _\/\\\/////\\\_        ___/\\\///\\\__        _\/\\\////\\\__        ____\//\\\\\___        _\////////\\\___    
	_\/\\\____\///_____\/\\\_        _\/\\\_______\/\\\_        _\/\\\/////////\\\_        _\/\\\___\///__        __/\\\__\//\\\_        _\/\\\__\//\\\_        _____\//\\\____        ___/\\\\\\\\\\__   
	 _\/\\\_____________\/\\\_        _\/\\\_______\/\\\_        _\/\\\_______\/\\\_        _\/\\\_________        _\//\\\__/\\\__        _\/\\\___\/\\\_        __/\\_/\\\_____        __/\\\/////\\\__  
	  _\/\\\_____________\/\\\_        _\//\\\\\\\\\\\\/__        _\/\\\_______\/\\\_        _\/\\\_________        __\///\\\\\/___        _\/\\\___\/\\\_        _\//\\\\/______        _\//\\\\\\\\/\\_ 
	   _\///______________\///__        __\////////////____        _\///________\///__        _\///__________        ____\/////_____        _\///____\///__        __\////________        __\////////\//__

	`,
	`
	
	___ __ __      _______      ________       ______        ______       ___   __       __  __     ________      
	/__//_//_/\    /______/\    /_______/\     /_____/\      /_____/\     /__/\ /__/\    /_/\/_/\   /_______/\     
	\::\| \| \ \   \::::__\/__  \::: _  \ \    \:::_ \ \     \:::_ \ \    \::\_\\  \ \   \ \ \ \ \  \::: _  \ \    
	 \:.      \ \   \:\ /____/\  \::(_)  \ \    \:(_) ) )_    \:\ \ \ \    \:. \-\  \ \   \:\_\ \ \  \::(_)  \ \   
	  \:.\-/\  \ \   \:\\_  _\/   \:: __  \ \    \: __ \\ \    \:\ \ \ \    \:. _    \ \   \::::_\/   \:: __  \ \  
	   \. \  \  \ \   \:\_\ \ \    \:.\ \  \ \    \ \ \\ \ \    \:\_\ \ \    \. \\-\  \ \    \::\ \    \:.\ \  \ \ 
		\__\/ \__\/    \_____\/     \__\/\__\/     \_\/ \_\/     \_____\/     \__\/ \__\/     \__\/     \__\/\__\/ 
																												   
	
	`,
	`
	
	___ __ __      _______      ________       ______        ______       ___   __       __  __     ________      
	/__//_//_/\    /______/\    /_______/\     /_____/\      /_____/\     /__/\ /__/\    /_/\/_/\   /_______/\     
	\::\| \| \ \   \::::__\/__  \::: _  \ \    \:::_ \ \     \:::_ \ \    \::\_\\  \ \   \ \ \ \ \  \::: _  \ \    
	 \:.      \ \   \:\ /____/\  \::(_)  \ \    \:(_) ) )_    \:\ \ \ \    \:. \-\  \ \   \:\_\ \ \  \::(_)  \ \   
	  \:.\-/\  \ \   \:\\_  _\/   \:: __  \ \    \: __ \\ \    \:\ \ \ \    \:. _    \ \   \::::_\/   \:: __  \ \  
	   \. \  \  \ \   \:\_\ \ \    \:.\ \  \ \    \ \ \\ \ \    \:\_\ \ \    \. \\-\  \ \    \::\ \    \:.\ \  \ \ 
		\__\/ \__\/    \_____\/     \__\/\__\/     \_\/ \_\/     \_____\/     \__\/ \__\/     \__\/     \__\/\__\/ 
																												   
	
	`,
}

// @title    MgaronyaPrint
// @description   打印一段随机的mgaronya字符串
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     无
// @return    无
func MgaronyaPrint() {
	log.Println(MgaronyaString[rand.New(rand.NewSource(time.Now().UnixNano())).Int()%10])
}

// @title    FileExit
// @description   查看某一文件是否存在
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     path string		文件以及路径
// @return    bool				表示是否存在文件
func FileExit(path string) bool {
	finfo, err := os.Stat(path)
	return err == nil && !finfo.IsDir()
}

// @title    RandomString
// @description   生成一段随机的字符串
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     n int		字符串的长度
// @return    string    一串随机的字符串
func RandomString(n int) string {
	var letters = []byte("qwertyuioplkjhgfdsazxcvbnmQWERTYUOIPLKJHGFDSAZXCVBNM")
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
	return user.ID != uuid.UUID{}
}

// @title    isNameExist
// @description   查看name是否在数据库中存在
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func IsNameExist(db *gorm.DB, name string) bool {
	var user model.User
	db.Where("name = ?", name).First(&user)
	return user.ID != uuid.UUID{}
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
func IsEmailPass(ctx *gin.Context, email string, vertify string) bool {
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
func SetRedisEmail(ctx *gin.Context, email string, v string) {
	client := common.GetRedisClient(0)

	client.Set(ctx, email, v, 300*time.Second)
}

// @title    ScoreChange
// @description   用于计算分数变化
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    email string, v string       接收一个邮箱和一个验证码
// @return   void
func ScoreChange(fre float64, sum float64, del float64, total float64) float64 {
	return (0.07/(fre+1) + 0.04) * sum * (math.Pow(2, 10*del-0.5)) / (math.Pow(2, 10*del-0.5) + 1) * (math.Pow(2, 0.1*total-5)) / (math.Pow(2, 0.1*total-5) + 1) / total
}

// @title    StringMerge
// @description   用于字符串的合并
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    a string, b string       接收两个字符串
// @return   string			返回合并结果
func StringMerge(a string, b string) string {
	if a > b {
		return a + b
	} else {
		return b + a
	}
}
