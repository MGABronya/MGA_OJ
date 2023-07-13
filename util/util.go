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
	"MGA_OJ/vo"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/extrame/xls"
	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"github.com/tealeg/xlsx"
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

// Tags			定义自动标签
var Tags []string = []string{
	"基本算法", "基础", "简单", "练习", "编程",
	"搜索", "计算几何", "数学", "数论", "图论", "动态规划", "数据结构",
	"枚举", "贪心", "递归", "分治", "递推", "构造", "模拟",
	"深度优先搜索", "宽度优先搜索", "广度优先搜索", "双向搜索", "启发式搜索", "记忆化搜索",
	"几何公式", "叉积", "点积", "多边形", "凸包", "扫描线", "内核", "几何工具", "平面交线", "可视图", "点集最小圆覆盖", "对踵点",
	"组合数学", "排列组合", "容斥原理", "抽屉原理", "置换群", "Polya定理", "母函数", "MoBius反演", "偏序关系理论",
	"素数", "整除", "进制", "模运算", "高斯消元", "概率", "欧几里得", "扩展欧几里得",
	"博弈论", "Nim", "极大过程", "极小过程",
	"拓扑排序", "最小生成树", "最短路", "二分图", "匈牙利算法", "KM算法",
	"网络流", "最小费用最大流", "最小费用流", "最小割", "网络流规约", "差分约束", "双连通分量", "强连通分支", "割边", "割点",
	"背包问题", "01背包", "完全背包", "多维背包", "多重背包", "区间dp", "环形dp", "判定性dp", "棋盘分割", "最长公共子序列", "最长上升子序列",
	"二分判定型dp", "树型动态规划", "最大独立集", "状态压缩dp", "哈密顿路径", "四边形不等式", "单调队列", "单调栈",
	"串", "KMP", "排序", "快排", "快速排序", "归并排序", "逆序数", "堆排序",
	"哈希表", "二分", "并查集", "霍夫曼树", "哈夫曼树", "堆", "线段树", "二叉树", "树状数组", "RMQ",
	"社招", "校招", "面经",
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
	db.Where("email = (?)", email).First(&user)
	return user.ID != uuid.UUID{}
}

// @title    isNameExist
// @description   查看name是否在数据库中存在
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func IsNameExist(db *gorm.DB, name string) bool {
	var user model.User
	db.Where("name = (?)", name).First(&user)
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
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "2829214609@qq.com", "qzinvmfkpsmvdgig", "smtp.qq.com"))
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
	err = db.Model(&model.User{}).Where("email = (?)", em[0]).Updates(model.User{
		Password: string(hasedPassword),
	}).Error

	if err != nil {
		return "密码更新失败"
	}

	// TODO 设置文件发送的内容
	content := fmt.Sprintf(mod, em[0], t, password)
	e.Text = []byte(content)
	// TODO 设置服务器相关的配置
	err = e.Send("smtp.qq.com:25", smtp.PlainAuth("", "2829214609@qq.com", "qzinvmfkpsmvdgig", "smtp.qq.com"))

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

// @title    Read
// @description   读取文件内容
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     file_path string		文件位置
// @return    res [][]string, err error		res为读出的内容，err为可能出现的错误
func Read(file_path string) (res [][]string, err error) {

	extName := path.Ext(file_path)

	if extName == ".csv" {
		return ReadCsv(file_path)
	} else if extName == ".xls" {
		return ReadXls(file_path)
	} else if extName == ".xlsx" {
		return ReadXlsx(file_path)
	}
	return nil, nil
}

// @title    ReadCsv
// @description   读取Csv文件内容
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     file_path string		文件位置
// @return    res [][]string, err error		res为读出的内容，err为可能出现的错误
func ReadCsv(file_path string) (res [][]string, err error) {
	file, err := os.Open(file_path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// TODO 初始化csv-reader
	reader := csv.NewReader(file)
	// TODO 设置返回记录中每行数据期望的字段数，-1 表示返回所有字段
	reader.FieldsPerRecord = -1
	// TODO 允许懒引号（忘记遇到哪个问题才加的这行）
	reader.LazyQuotes = true
	// TODO 返回csv中的所有内容
	records, read_err := reader.ReadAll()
	if read_err != nil {
		return nil, read_err
	}
	return records, nil
}

// @title    ReadXls
// @description   读取Xls文件内容
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     file_path string		文件位置
// @return    res [][]string, err error		res为读出的内容，err为可能出现的错误
func ReadXls(file_path string) (res [][]string, err error) {
	if xlFile, err := xls.Open(file_path, "utf-8"); err == nil {
		fmt.Println(xlFile.Author)
		// TODO 第一个sheet
		sheet := xlFile.GetSheet(0)
		if sheet.MaxRow != 0 {
			temp := make([][]string, sheet.MaxRow)
			for i := 0; i < int(sheet.MaxRow); i++ {
				row := sheet.Row(i)
				data := make([]string, 0)
				if row.LastCol() > 0 {
					for j := 0; j < row.LastCol(); j++ {
						col := row.Col(j)
						data = append(data, col)
					}
					temp[i] = data
				}
			}
			res = append(res, temp...)
		}
	} else {
		return nil, err
	}
	return res, nil
}

// @title    ReadXlsx
// @description   读取Xlsx文件内容
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     file_path string		文件位置
// @return    res [][]string, err error		res为读出的内容，err为可能出现的错误
func ReadXlsx(file_path string) (res [][]string, err error) {
	if xlFile, err := xlsx.OpenFile(file_path); err == nil {
		for index, sheet := range xlFile.Sheets {
			// TODO 第一个sheet
			if index == 0 {
				temp := make([][]string, len(sheet.Rows))
				for k, row := range sheet.Rows {
					var data []string
					for _, cell := range row.Cells {
						data = append(data, cell.Value)
					}
					temp[k] = data
				}
				res = append(res, temp...)
			}
		}
	} else {
		return nil, err
	}
	return res, nil
}

// @title    RemoveWhiteSpace
// @description  函数接收一个字符串，返回一个去掉空白字符的新字符串
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     str string		目标字符
// @return    string		去掉空白字符的新字符串
func RemoveWhiteSpace(str string) string {
	// TODO 创建一个空字符串res用于存储处理后的结果
	var res string
	// TODO 遍历传入的字符串中的每个字符
	for _, char := range str {
		// TODO 如果当前字符不是空格、制表符、回车或换行符等空白字符
		if !unicode.IsSpace(char) {
			// TODO 将该字符加入到结果字符串中
			res += string(char)
		}
	}
	// TODO 返回处理后的结果字符串
	return res
}

// behaviors	    定义了用户行为映射表
var behaviors map[string]func(uuid.UUID) (float64, string) = map[string]func(uuid.UUID) (float64, string){
	"aaa": a,
	"bbb": b,
}

func a(uuid.UUID) (float64, string) {
	return 1, ""
}

func b(uuid.UUID) (float64, string) {
	return 2, ""
}

// @title    checkExpression
// @description  函数接收一个不带括号的表达式，查看表达式是否正确
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     expr []byte		简单表达式
// @return    bool, string		表示表达式是否正确，并给出错误原因
func checkExpression(expr []byte) (bool, string) {
	println(string(expr))
	cp := 0
	// TODO 定义一个栈用于匹配括号
	for i := 0; i < len(expr); i++ {
		if expr[i] >= 'a' && expr[i] <= 'z' {
			if (cp & 1) == 1 {
				return false, "变量位置错误"
			}
			j := i
			for i+1 < len(expr) && expr[i+1] >= 'a' && expr[i+1] <= 'z' {
				i++
			}
			if _, ok := behaviors[string(expr[j:i+1])]; !ok {
				return false, "变量" + string(expr[j:i+1]) + "未定义"
			}
		} else if expr[i] == '#' {
			if (cp & 1) == 1 {
				return false, "变量位置错误"
			}
			for i+1 < len(expr) && expr[i+1] == '#' {
				i++
			}
		} else if expr[i] > '0' && expr[i] <= '9' {
			if (cp & 1) == 1 {
				return false, "变量位置错误"
			}
			for i+1 < len(expr) && expr[i] >= '0' && expr[i] <= '9' {
				i++
			}
		} else if expr[i] == '0' {
			if (cp & 1) == 1 {
				return false, "变量位置错误"
			}
		} else if expr[i] != '+' && expr[i] != '-' && expr[i] != '*' && expr[i] != '/' {
			if (cp & 1) == 0 {
				return false, "计算符位置错误"
			}
			// TODO 非法字符
			return false, "非法字符"
		}
		cp++
	}
	if (cp & 1) == 1 {
		return true, ""
	}
	return false, "计算符或变量缺失"
}

// @title    checkExpression
// @description  函数接收一个表达式，查看表达式是否正确
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     expr []byte		表达式
// @return    bool, string		表示表达式是否正确，并给出错误原因
func CheckExpression(expr []byte) (bool, string) {
	// TODO 定义一个栈用于匹配括号
	stack := make([]byte, 0)
	// TODO 记录需要替换的子串
	// TODO 记录当前最外层括号的起始位置和长度
	start := -1

	for i, ch := range expr {
		if ch == '(' {
			// TODO 左括号直接入栈
			if len(stack) == 0 {
				start = i
			}
			stack = append(stack, ch)
		} else if ch == ')' {
			// TODO 右括号需要和栈顶的左括号匹配
			if len(stack) > 0 && stack[len(stack)-1] == '(' {
				stack = stack[:len(stack)-1]
				if len(stack) == 0 {
					if ok, err := CheckExpression(expr[start+1 : i]); !ok {
						return false, err
					}
					// TODO 将最外层括号替换为#
					for j := start; j < i+1; j++ {
						expr[j] = '#'
					}
				}
			} else {
				return false, "右括号不匹配"
			}
		}
	}

	if len(stack) > 0 {
		return false, "右括号不匹配"
	}

	if ok, err := checkExpression(expr); !ok {
		return false, err
	}

	return true, ""
}

// @title    EvaluateExpression
// @description  函数接收一个表达式，计算表达式的值
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     expr string		表达式
// @return    int, string		计算表达式的值
func EvaluateExpression(expression string, userId uuid.UUID) (int, string) {
	// TODO 运算数栈
	stack := make([]string, 0)
	// TODO 运算符栈
	operators := make([]string, 0)

	// TODO 去除表达式中的空格
	expression = strings.ReplaceAll(expression, " ", "")

	for i := 0; i < len(expression); i++ {
		char := expression[i]
		// TODO 处理数字
		if char >= '0' && char <= '9' {
			numStr := string(char)

			// TODO 找到数字结束位置
			for j := i + 1; j < len(expression); j++ {
				nextChar := expression[j]
				if nextChar >= '0' && nextChar <= '9' {
					numStr += string(nextChar)
				} else {
					break
				}
			}

			// TODO 将数字压入栈中
			stack = append(stack, numStr)
			// TODO 更新索引位置
			i += len(numStr) - 1
		} else if char >= 'a' && char <= 'z' {
			// TODO 处理变量
			varName := string(char)

			// TODO 找到变量结束位置
			for j := i + 1; j < len(expression); j++ {
				nextChar := expression[j]
				if nextChar >= 'a' && nextChar <= 'z' {
					varName += string(nextChar)
				} else {
					break
				}
			}

			value, err := behaviors[varName](userId)
			if err != "" {
				return 0, err
			}

			// TODO 将变量值压入栈中
			stack = append(stack, strconv.Itoa(int(value)))
			// TODO 更新索引位置
			i += len(varName) - 1
		} else if char == '(' {
			// TODO 处理左括号
			operators = append(operators, string(char))
		} else if char == ')' {
			// TODO 处理右括号
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				topOperator := operators[len(operators)-1]
				operators = operators[:len(operators)-1]

				if len(stack) < 2 {
					return 0, "invalid expression"
				}

				num2, _ := strconv.Atoi(stack[len(stack)-1])
				stack = stack[:len(stack)-1]
				num1, _ := strconv.Atoi(stack[len(stack)-1])
				stack = stack[:len(stack)-1]

				result := 0
				switch topOperator {
				case "+":
					result = num1 + num2
				case "-":
					result = num1 - num2
				case "*":
					result = num1 * num2
				case "/":
					result = num1 / num2
				}
				// TODO 将计算结果压入栈中
				stack = append(stack, strconv.Itoa(result))
			}

			if len(operators) > 0 && operators[len(operators)-1] == "(" {
				operators = operators[:len(operators)-1]
			} else {
				return 0, "invalid expression"
			}
		} else if char == '+' || char == '-' || char == '*' || char == '/' {
			// TODO 处理运算符
			for len(operators) > 0 && (operators[len(operators)-1] == "*" || operators[len(operators)-1] == "/") {
				topOperator := operators[len(operators)-1]
				operators = operators[:len(operators)-1]

				if len(stack) < 2 {
					return 0, "invalid expression"
				}

				num2, _ := strconv.Atoi(stack[len(stack)-1])
				stack = stack[:len(stack)-1]
				num1, _ := strconv.Atoi(stack[len(stack)-1])
				stack = stack[:len(stack)-1]

				result := 0
				switch topOperator {
				case "*":
					result = num1 * num2
				case "/":
					if num2 == 0 {
						return 0, "division by zero"
					}
					result = num1 / num2
				}
				// TODO 将计算结果压入栈中
				stack = append(stack, strconv.Itoa(result))
			}

			operators = append(operators, string(char))
		} else {
			return 0, fmt.Sprintf("invalid character %c", char)
		}
	}

	for len(operators) > 0 {
		topOperator := operators[len(operators)-1]
		operators = operators[:len(operators)-1]

		if len(stack) < 2 {
			return 0, "invalid expression"
		}

		num2, _ := strconv.Atoi(stack[len(stack)-1])
		stack = stack[:len(stack)-1]
		num1, _ := strconv.Atoi(stack[len(stack)-1])
		stack = stack[:len(stack)-1]

		result := 0
		switch topOperator {
		case "+":
			result = num1 + num2
		case "-":
			result = num1 - num2
		}
		// TODO 将计算结果压入栈中
		stack = append(stack, strconv.Itoa(result))
	}

	if len(stack) != 1 {
		return 0, "invalid expression"
	}
	// TODO 将结果转换为整数
	value, err := strconv.Atoi(stack[0])
	if err != nil {
		return 0, "invalid expression"
	}

	return value, ""
}

// @title    Search
// @description  调用bing搜索api进行搜索
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     query string			api-key,终结点,搜索内容
// @return    []SearchResult, error				搜索结果,报错信息
func Search(query string) ([]vo.SearchResult, error) {

	// TODO 构建搜索请求的URI
	params := url.Values{}
	params.Set("q", query)
	requestURL := common.Endpoint + "?" + params.Encode()

	// TODO 发送HTTP请求并获取响应
	client := &http.Client{}
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", common.SubscriptionKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// TODO 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// TODO 解析JSON响应体到结构体数组
	var parsedJson struct {
		WebPages struct {
			Value []vo.SearchResult `json:"value"`
		} `json:"webPages"`
	}
	err = json.Unmarshal(body, &parsedJson)
	if err != nil {
		return nil, err
	}

	return parsedJson.WebPages.Value, nil
}

// @title    CountTags
// @description  记录搜索结果中的标签出现次数
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     query string			api-key,终结点,搜索内容
// @return    []SearchResult, error				搜索结果,报错信息
func CountTags(searchResults []vo.SearchResult, tags ...string) []vo.TagCount {
	// TODO 创建一个空的切片，用于存储每个标签字符串出现的次数
	counts := make([]vo.TagCount, 0)

	// TODO 遍历每个标签字符串
	for _, tag := range tags {
		// TODO 将标签字符串转换为小写
		lowerTag := strings.ToLower(tag)

		count := 0
		for _, searchResult := range searchResults {
			// TDOO 使用strings.Count函数计算标签字符串在输入字符串中出现的次数
			count += strings.Count(strings.ToLower(searchResult.Name), lowerTag)
			count += strings.Count(strings.ToLower(searchResult.Snippet), lowerTag)
		}

		// TODO 如果次数大于0，则将标签和次数存储到切片中
		if count > 0 {
			tagCount := vo.TagCount{tag, count}
			counts = append(counts, tagCount)
		}
	}

	// TODO 按次数从大到小排序
	sort.SliceStable(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})

	return counts
}

// @title    GetInfoFromXML
// @description  将题目从xml格式中读取出来
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     xmlString string			xml格式题目
// @return    vo.Item, error				题目信息,报错信息
func GetInfoFromXML(xmlString string) (vo.Item, error) {
	var result vo.Fps

	err := xml.Unmarshal([]byte(xmlString), &result)
	if err != nil {
		return result.Item, err
	}
	return result.Item, nil
}
