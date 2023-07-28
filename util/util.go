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

var searchIndex int = 0

var Max_run int = 4

// timerMap	    定义了当前使用的定时器
var TimerMap map[uuid.UUID]*time.Timer = make(map[uuid.UUID]*time.Timer)

// LanguageMap			定义语言字典，对应其处理方式
var LanguageMap map[string]Interface.CmdInterface = map[string]Interface.CmdInterface{
	"C":          Handle.NewC(),
	"C#":         Handle.NewCs(),
	"C++":        Handle.NewCppPlusPlus(),
	"C++11":      Handle.NewCppPlusPlus11(),
	"C++14":      Handle.NewCppPlusPlus14(),
	"C++17":      Handle.NewCppPlusPlus17(),
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

// OJMap			支持的oj
var OJMap map[string]string = map[string]string{
	"POJ":   "00000001",
	"HDU":   "00000002",
	"SPOJ":  "00000003",
	"VIJOS": "00000004",
	"CF":    "00000005",
}

// JOMap			支持的oj，但反向映射
var JOMap map[string]string = map[string]string{
	"00000001": "POJ",
	"00000002": "HDU",
	"00000003": "SPOJ",
	"00000004": "VIJOS",
	"00000005": "CF",
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

// @title    Search
// @description  调用bing搜索api进行搜索
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     query string			api-key,终结点,搜索内容
// @return    []SearchResult, error				搜索结果,报错信息
func Search(query string) ([]vo.SearchResult, error) {

	// TODO 构建搜索请求的URI
	params := url.Values{}
	params.Set("q", query)
	requestURL := "https://api.bing.microsoft.com/v7.0/search?" + params.Encode()

	// TODO 发送HTTP请求并获取响应
	client := &http.Client{}
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", common.SubscriptionKey[searchIndex%len(common.SubscriptionKey)])
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

	if parsedJson.WebPages.Value == nil || len(parsedJson.WebPages.Value) == 0 {
		searchIndex++
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

// @title    PadZero
// @description  为字符串添加前导零
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     str string				需要添加前导零的字符串
// @return    string				    添加前导零后的字符串
func PadZero(str string) string {
	str = strings.TrimLeft(str, "0")
	return fmt.Sprintf("%032s", str)
}

// @title    EncodeUUID
// @description  将string类型编码为uuid类型
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     proid, source string				用于编码的题目id字符串以及题目来源平台
// @return    uuid.UUID, error				    编码后的uuid以及可能的报错信息
func EncodeUUID(proid, source string) (uuid.UUID, error) {
	// TODO 尝试将题目转为16进制
	proid, err := SixtyFourToSixteen(proid)
	if err != nil {
		return uuid.Nil, err
	}
	// TODO 填充前导零
	paddedString := PadZero(proid)
	// TODO 添加来源标记
	if s, ok := OJMap[source]; !ok {
		return uuid.Nil, fmt.Errorf("不支持的平台")
	} else {
		paddedString = s + paddedString[8:]
		uuidValue, err := uuid.FromString(strings.ReplaceAll(paddedString, "-", ""))
		return uuidValue, err
	}
}

// @title    DeCodeUUID
// @description  将uuid类型解码为string类型
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     uuidValue uuid.UUID				待解码的uuid
// @return    string, string				    解码后的string
func DeCodeUUID(uuidValue uuid.UUID) (proid string, source string, err error) {
	uuidString := strings.TrimLeft(uuidValue.String(), "{-}")
	source = uuidString[:8]
	// TODO 查看source是否合法
	if s, ok := JOMap[source]; !ok {
		return "", "", fmt.Errorf("前缀不合法")
	} else {
		uuidString = uuidString[8:]
		// TODO 将uuid类型还原为原先的字符串
		proid = strings.TrimLeft(uuidString, "0-")
		// TODO 将proid转化为62进制
		proid, err = SixteenToSixtyFour(proid)
		return proid, s, err
	}
}

// @title    SixtyFourToSixteen
// @description  将64进制转化为16进制
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     str string				62进制字符串
// @return    string,error				16进制以及可能的错误
func SixtyFourToSixteen(str string) (string, error) {
	var res int64
	for i := 0; i < len(str); i++ {
		res *= 62
		if unicode.IsNumber(rune(str[i])) {
			res += int64(rune(str[i]) - '0')
		} else if unicode.IsLower(rune(str[i])) {
			res += int64(rune(str[i]) - 'a' + 10)
		} else if unicode.IsUpper(rune(str[i])) {
			res += int64(rune(str[i]) - 'A' + 36)
		} else if str[i] == '-' {
			res += 62
		} else if str[i] == '_' {
			res += 63
		} else {
			return "0", fmt.Errorf("错误字符", rune(str[i]))
		}
	}
	return strconv.FormatInt(res, 16), nil
}

// @title    SixteenToSixtyFour
// @description  将16进制转化为64进制
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     str string				16进制字符串
// @return    string,error				64进制以及可能的错误
func SixteenToSixtyFour(res string) (string, error) {
	r, err := strconv.ParseUint(res, 16, 64)
	if err != nil {
		return "", err
	}
	str := ""
	var base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
	for r != 0 {
		t := r % 64
		r /= 64
		str = string(base62Chars[t]) + str
	}
	return str, nil
}

// @title    StateCorrection
// @description  矫正状态
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     str string				待矫正的状态
// @return    string				矫正后的状态
func StateCorrection(condition string) string {
	if condition == "" {
		return "Waiting"
	}
	short := strings.ToLower(condition)
	if regexp.MustCompile(`^ac.*`).MatchString(short) {
		return "Accepted"
	}
	if regexp.MustCompile(`^wait.*`).MatchString(short) {
		return "Waiting"
	}
	if regexp.MustCompile(`^r.*e.*`).MatchString(short) {
		return "Runtime Error"
	}
	if regexp.MustCompile(`^run.*`).MatchString(short) {
		return "Running"
	}
	if regexp.MustCompile(`^c.*e.*`).MatchString(short) {
		return "Compile Error"
	}
	if regexp.MustCompile(`^compil.*`).MatchString(short) {
		return "Compiling"
	}
	if regexp.MustCompile(`^p.*e.*`).MatchString(short) {
		return "Presentation Error"
	}
	if regexp.MustCompile(`^w.*a.*`).MatchString(short) {
		return "Wrong Answer"
	}
	if regexp.MustCompile(`^t.*l.*e.*`).MatchString(short) {
		return "Time Limit Exceeded"
	}
	if regexp.MustCompile(`^m.*l.*e.*`).MatchString(short) {
		return "Memory Limit Exceeded"
	}
	return condition
}

// @title    RemoveDuplicates
// @description  字符串去重
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     arr []string				去重前的字符串数组
// @return    []string				去重后的字符串数组
func RemoveDuplicates(arr []string) []string {
	// TODO 创建一个空的 map 用于记录已经出现过的字符串
	m := make(map[string]bool)
	result := []string{}

	// TODO 遍历数组中的每个字符串
	for _, str := range arr {
		// TODO 如果该字符串不在 map 中，说明是第一次出现，将其添加到结果数组中，并在 map 中标记为已出现
		if !m[str] {
			result = append(result, str)
			m[str] = true
		}
	}

	return result
}
