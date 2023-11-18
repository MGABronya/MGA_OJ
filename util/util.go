// @Title  util
// @Description  收集各种需要使用的工具函数
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:47
package util

import (
	"MGA_OJ/Interface"
	Handle "MGA_OJ/Language"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/vo"
	"archive/zip"
	"context"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/extrame/xls"
	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
	"github.com/parnurzeal/gorequest"
	uuid "github.com/satori/go.uuid"
	"github.com/tealeg/xlsx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Units		定义了单位换算
var Units = map[string]uint{
	"mb":  1024,
	"mib": 1024,
	"kb":  1,
	"kib": 1,
	"gb":  1024 * 1024,
	"ms":  1,
	"s":   1000,
}

var searchIndex int = 0
var translateIndex int = 0

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
	"C++20":      Handle.NewCppPlusPlus20(),
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
	"组合数学", "排列组合", "容斥原理", "抽屉原理", "置换群", "Polya", "母函数", "MoBius", "偏序关系理论",
	"素数", "整除", "进制", "模运算", "高斯消元", "概率", "欧几里得", "扩展欧几里得",
	"博弈论", "Nim", "极大过程", "极小过程",
	"拓扑排序", "最小生成树", "最短路", "二分图", "匈牙利算法", "KM算法", "仙人掌",
	"网络流", "最小费用最大流", "最小费用流", "最小割", "网络流规约", "差分约束", "双连通分量", "强连通分支", "割边", "割点",
	"背包问题", "01背包", "完全背包", "多维背包", "多重背包", "区间dp", "环形dp", "判定性dp", "棋盘分割", "最长公共子序列", "最长上升子序列",
	"二分判定型dp", "树型动态规划", "最大独立集", "状态压缩dp", "哈密顿路径", "四边形不等式", "单调队列", "单调栈",
	"串", "KMP", "排序", "快排", "快速排序", "归并排序", "逆序数", "堆排序",
	"哈希表", "二分", "并查集", "霍夫曼树", "哈夫曼树", "堆", "线段树", "二叉树", "树状数组", "RMQ", "阿朵莉树",
	"社招", "校招", "面经",
}

// BasicAlgorithmMap			定义基础算法映射表
var BasicAlgorithmMap map[string]bool = map[string]bool{
	"基本算法": true, "基础": true, "简单": true, "练习": true, "编程": true, "枚举": true, "贪心": true, "递归": true, "分治": true, "递推": true, "构造": true, "模拟": true,
}

// SearchMap			定义搜索映射表
var SearchMap map[string]bool = map[string]bool{
	"搜索": true, "深度优先搜索": true, "宽度优先搜索": true, "广度优先搜索": true, "双向搜索": true, "启发式搜索": true, "记忆化搜索": true,
}

// ComputationalGeometryMap			定义计算几何映射表
var ComputationalGeometryMap map[string]bool = map[string]bool{
	"计算几何": true, "几何公式": true, "叉积": true, "点积": true, "多边形": true, "凸包": true, "扫描线": true, "内核": true,
	"几何工具": true, "平面交线": true, "可视图": true, "点集最小圆覆盖": true, "对踵点": true,
}

// NumberTheoryMap			定义数论映射表
var NumberTheoryMap map[string]bool = map[string]bool{
	"数学": true, "数论": true, "组合数学": true, "排列组合": true, "容斥原理": true, "抽屉原理": true, "置换群": true, "Polya": true, "母函数": true, "MoBius": true, "偏序关系理论": true,
	"素数": true, "整除": true, "进制": true, "模运算": true, "高斯消元": true, "概率": true, "欧几里得": true, "扩展欧几里得": true,
	"博弈论": true, "Nim": true, "极大过程": true, "极小过程": true,
}

// GraphTheoryMap			定义图论映射表
var GraphTheoryMap map[string]bool = map[string]bool{
	"图论": true, "拓扑排序": true, "最小生成树": true, "最短路": true, "二分图": true, "匈牙利算法": true, "KM算法": true, "仙人掌": true,
	"网络流": true, "最小费用最大流": true, "最小费用流": true, "最小割": true, "网络流规约": true, "差分约束": true, "双连通分量": true, "强连通分支": true, "割边": true, "割点": true,
}

// DynamicProgrammingMap			定义动态规划映射表
var DynamicProgrammingMap map[string]bool = map[string]bool{
	"动态规划": true, "背包问题": true, "01背包": true, "完全背包": true, "多维背包": true, "多重背包": true, "区间dp": true, "环形dp": true, "判定性dp": true, "棋盘分割": true,
	"最长公共子序列": true, "最长上升子序列": true,
	"二分判定型dp": true, "树型动态规划": true, "最大独立集": true, "状态压缩dp": true, "哈密顿路径": true, "四边形不等式": true, "单调队列": true, "单调栈": true,
}

// DataStructureMap			定义数据结构映射表
var DataStructureMap map[string]bool = map[string]bool{
	"数据结构": true, "串": true, "KMP": true, "排序": true, "快排": true, "快速排序": true, "归并排序": true, "逆序数": true, "堆排序": true,
	"哈希表": true, "二分": true, "并查集": true, "霍夫曼树": true, "哈夫曼树": true, "堆": true, "线段树": true, "二叉树": true, "树状数组": true, "RMQ": true, "阿朵莉树": true,
}

// OJMap			支持的oj
var OJMap map[string]string = map[string]string{
	"POJ":        "00000001",
	"HDU":        "00000002",
	"SPOJ":       "00000003",
	"VIJOS":      "00000004",
	"CF":         "00000005",
	"UVA":        "00000006",
	"UOJ":        "00000007",
	"URAL":       "00000008",
	"HACKERRANK": "00000009",
	"ATCODER":    "0000000a",
}

// JOMap			支持的oj，但反向映射
var JOMap map[string]string = map[string]string{
	"00000001": "POJ",
	"00000002": "HDU",
	"00000003": "SPOJ",
	"00000004": "VIJOS",
	"00000005": "CF",
	"00000006": "UVA",
	"00000007": "UOJ",
	"00000008": "URAL",
	"00000009": "HACKERRANK",
	"0000000a": "ATCODER",
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
// @auth      MGAronya             2022-9-16 10:29
// @param     无
// @return    无
func MgaronyaPrint() {
	log.Println(MgaronyaString[rand.New(rand.NewSource(time.Now().UnixNano())).Int()%10])
}

// @title    FileExit
// @description   查看某一文件是否存在
// @auth      MGAronya             2022-9-16 10:29
// @param     path string		文件以及路径
// @return    bool				表示是否存在文件
func FileExit(path string) bool {
	finfo, err := os.Stat(path)
	return err == nil && !finfo.IsDir()
}

// @title    RandomString
// @description   生成一段随机的字符串
// @auth      MGAronya             2022-9-16 10:29
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
// @auth      MGAronya             2022-9-16 10:29
// @param     email string		一串字符串，表示邮箱
// @return    bool    返回是否合法
func VerifyEmailFormat(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// @title    VerifyMobileFormat
// @description   用于验证手机号格式是否正确的工具函数
// @auth      MGAronya             2022-9-16 10:29
// @param     mobileNum string		一串字符串，表示手机号
// @return    bool    返回是否合法
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

// @title    VerifyQQFormat
// @description   用于验证QQ号格式是否正确的工具函数
// @auth      MGAronya             2022-9-16 10:29
// @param     QQNum string		一串字符串，表示QQ
// @return    bool    返回是否合法
func VerifyQQFormat(QQNum string) bool {
	regular := "[1-9][0-9]{4,10}"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(QQNum)
}

// @title    VerifyQQFormat
// @description  用于验证Icon是否为默认图片的工具函数
// @auth      MGAronya             2022-9-16 10:29
// @param     Icon string		一串字符串，表示图像名称
// @return    bool    返回是否合法
func VerifyIconFormat(Icon string) bool {
	regular := "MGA[1-9].jpg"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(Icon)
}

// @title    isEmailExist
// @description   查看email是否在数据库中存在
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func IsEmailExist(db *gorm.DB, email string) bool {
	var user model.User
	db.Where("email = (?)", email).First(&user)
	return user.ID != uuid.UUID{}
}

// @title    isNameExist
// @description   查看name是否在数据库中存在
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func IsNameExist(db *gorm.DB, name string) bool {
	var user model.User
	db.Where("name = (?)", name).First(&user)
	return user.ID != uuid.UUID{}
}

// @title    SendEmailValidate
// @description   发送验证邮件
// @auth      MGAronya       2022-9-16 12:15
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

// @title    SendEmail
// @description   发送验证邮件
// @auth      MGAronya       2022-9-16 12:15
// @param    em []string, ex       接收一个邮箱字符串以及邮件内容
// @return   string, error     返回验证码和error值
func SendEmail(em []string, tx string) error {
	mod := `
	尊敬的%s，您好！

	这是一封来自DOJ的来信，如有打扰，还望海涵，以下为来信内容。

	%s
`
	e := email.NewEmail()
	e.From = "mgAronya <2829214609@qq.com>"
	e.To = em
	// TODO 设置文件发送的内容
	content := fmt.Sprintf(mod, em[0], tx)
	e.Text = []byte(content)
	// TODO 设置服务器相关的配置
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "2829214609@qq.com", "qzinvmfkpsmvdgig", "smtp.qq.com"))
	return err
}

// @title    SendEmailPass
// @description   发送密码邮件
// @auth      MGAronya       2022-9-16 12:15
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
// @auth      MGAronya       2022-9-16 12:15
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
// @auth      MGAronya       2022-9-16 12:15
// @param    email string, v string       接收一个邮箱和一个验证码
// @return   void
func SetRedisEmail(ctx *gin.Context, email string, v string) {
	client := common.GetRedisClient(0)

	client.Set(ctx, email, v, 300*time.Second)
}

// @title    ScoreChange
// @description   用于计算分数变化
// @auth      MGAronya       2022-9-16 12:15
// @param    email string, v string       接收一个邮箱和一个验证码
// @return   void
func ScoreChange(fre float64, sum float64, del float64, total float64) float64 {
	return (0.07/(fre+1) + 0.04) * sum * (math.Pow(2, 10*del-0.5)) / (math.Pow(2, 10*del-0.5) + 1) * (math.Pow(2, 0.1*total-5)) / (math.Pow(2, 0.1*total-5) + 1) / total
}

// @title    StringMerge
// @description   用于字符串的合并
// @auth      MGAronya       2022-9-16 12:15
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
// @auth      MGAronya             2022-9-16 10:29
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
// @auth      MGAronya             2022-9-16 10:29
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
// @auth      MGAronya             2022-9-16 10:29
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
// @auth      MGAronya             2022-9-16 10:29
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
// @auth      MGAronya             2022-9-16 10:29
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
// @auth      MGAronya             2022-9-16 10:29
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
	req.Header.Set("Ocp-Apim-Subscription-Key", common.SearchSubscriptionKey[searchIndex%len(common.SearchSubscriptionKey)])
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

// @title    Translator
// @description  调用bing翻译api进行翻译
// @auth      MGAronya             2022-9-16 10:29
// @param     query string			翻译内容
// @return    string, error				翻译结果,报错信息
func Translator(query string) (string, error) {
	requestBody := []vo.TranslationRequest{
		{
			Text: query,
		},
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.cognitive.microsofttranslator.com/translate?api-version=3.0&to=zh-Hans", strings.NewReader(string(jsonData)))
	if err != nil {
		return "", err
	}
	req.Header.Add("Ocp-Apim-Subscription-Key", common.TraslationSubscriptionKey[translateIndex%len(common.TraslationSubscriptionKey)])
	req.Header.Add("Ocp-Apim-Subscription-Region", "eastasia")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	// TODO 解析响应
	var translationResponse []vo.TranslationResponse
	json.Unmarshal(responseBody, &translationResponse)

	if len(translationResponse) == 0 || len(translationResponse[0].Translations) == 0 {
		return "", nil
	}

	// TODO 返回翻译结果
	return translationResponse[0].Translations[0].Text, nil
}

// @title    CountTags
// @description  记录搜索结果中的标签出现次数
// @auth      MGAronya             2022-9-16 10:29
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
// @auth      MGAronya             2022-9-16 10:29
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
// @auth      MGAronya             2022-9-16 10:29
// @param     str string				需要添加前导零的字符串
// @return    string				    添加前导零后的字符串
func PadZero(str string) string {
	str = strings.TrimLeft(str, "0")
	return fmt.Sprintf("%032s", str)
}

// @title    EncodeUUID
// @description  将string类型编码为uuid类型
// @auth      MGAronya             2022-9-16 10:29
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
// @auth      MGAronya             2022-9-16 10:29
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
// @auth      MGAronya             2022-9-16 10:29
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
// @auth      MGAronya             2022-9-16 10:29
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

// @title    SixRandUUID
// @description  用随机数替换uuid的前六位
// @auth      MGAronya             2022-9-16 10:29
// @param     uuidValue uuid.UUID				代替换的uuid
// @return    uuid.UUID,error				结果uuid以及可能的错误
func SixRandUUID(uuidValue uuid.UUID) (uuid.UUID, error) {
	uuidString := uuidValue.String()
	// TODO 设置随机数种子
	rand.Seed(time.Now().UnixNano())

	// TODO 生成6位的随机数
	randomNumber := rand.Intn(0xffffff)

	// TODO 将随机数转换为16进制
	hexNumber := fmt.Sprintf("%06x", randomNumber)

	uuidString = hexNumber + uuidString[6:]

	return uuid.FromString(uuidString)
}

// @title    SixZeroUUID
// @description  用0替换uuid的前六位
// @auth      MGAronya             2022-9-16 10:29
// @param     uuidValue uuid.UUID				代替换的uuid
// @return    uuid.UUID,error				结果uuid以及可能的错误
func SixZeroUUID(uuidValue uuid.UUID) (uuid.UUID, error) {
	uuidString := uuidValue.String()

	uuidString = "000000" + uuidString[6:]

	return uuid.FromString(uuidString)
}

// @title    StateCorrection
// @description  矫正状态
// @auth      MGAronya             2022-9-16 10:29
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
	if regexp.MustCompile(`^ru.*er.*`).MatchString(short) {
		return "Runtime Error"
	}
	if regexp.MustCompile(`^run.*`).MatchString(short) {
		return "Running"
	}
	if regexp.MustCompile(`^co.*er.*`).MatchString(short) {
		return "Compile Error"
	}
	if regexp.MustCompile(`^compil.*`).MatchString(short) {
		return "Compiling"
	}
	if regexp.MustCompile(`^pr.*er.*`).MatchString(short) {
		return "Presentation Error"
	}
	if regexp.MustCompile(`^wr.*an.*`).MatchString(short) {
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
// @auth      MGAronya             2022-9-16 10:29
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

// @title    ProblemCategory
// @description  题目分类
// @auth      MGAronya             2022-9-16 10:29
// @param     ProblemId uuid.UUID				题目id
// @return    map[string]bool				分类结果
func ProblemCategory(ProblemId uuid.UUID) map[string]bool {
	str := make(map[string]bool)
	redix := common.GetRedisClient(0)
	db := common.GetDB()
	ctx := context.Background()

	// TODO 查找数据
	var problemLabels []model.ProblemLabel
	// TODO 先尝试在redis中寻找
	if ok, _ := redix.HExists(ctx, "ProblemLabel", ProblemId.String()).Result(); ok {
		art, _ := redix.HGet(ctx, "ProblemLabel", ProblemId.String()).Result()
		if json.Unmarshal([]byte(art), &problemLabels) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			redix.HDel(ctx, "ProblemLabel", ProblemId.String())
		}
	}

	// TODO 在数据库中查找
	if db.Where("problem_id = (?)", ProblemId).Find(&problemLabels).Error != nil {
		return str
	}
	{
		// TODO 将题目标签存入redis供下次使用
		v, _ := json.Marshal(problemLabels)
		redix.HSet(ctx, "ProblemLabel", ProblemId, v)
	}

leap:

	// TODO 查看标签类型进行分类
	for i := range problemLabels {
		if BasicAlgorithmMap[problemLabels[i].Label] {
			str["BasicAlgorithm"] = true
		} else if ComputationalGeometryMap[problemLabels[i].Label] {
			str["ComputationalGeometry"] = true
		} else if DataStructureMap[problemLabels[i].Label] {
			str["DataStructure"] = true
		} else if DynamicProgrammingMap[problemLabels[i].Label] {
			str["DynamicProgramming"] = true
		} else if GraphTheoryMap[problemLabels[i].Label] {
			str["GraphTheory"] = true
		} else if NumberTheoryMap[problemLabels[i].Label] {
			str["NumberTheory"] = true
		} else if SearchMap[problemLabels[i].Label] {
			str["Search"] = true
		}
	}
	return str
}

// @title    RemoveComments
// @description  移除注释函数
// @auth      MGAronya             2022-9-16 10:29
// @param     text string					代码文本
// @return    string						移除注释后的代码文本
func RemoveComments(text string) string {
	pattern := `\/\/.*|\/\*[\s\S]*?\*\/`
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(text, "")
}

// @title    RemoveComments
// @description  移除头文件
// @auth      MGAronya             2022-9-16 10:29
// @param     text string					代码文本
// @return    string						移除头文件后的代码文本
func RemoveHeads(text string) string {
	pattern := `#.*`
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(text, "")
}

// @title    RemoveSpaces
// @description  移除空格、换行符、制表符等函数
// @auth      MGAronya             2022-9-16 10:29
// @param     text string					代码文本
// @return    string						移除空格、换行符、制表符等函数后的代码文本
func RemoveSpaces(text string) string {
	pattern := `\s+`
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(text, "")
}

// @title    ToLowerCase
// @description  小写化函数
// @auth      MGAronya             2022-9-16 10:29
// @param     text string					代码文本
// @return    string						小写化函数后的代码文本
func ToLowerCase(text string) string {
	return strings.ToLower(text)
}

// @title    ComputeNgramFreq
// @description  计算n-gram词频函数
// @auth      MGAronya             2022-9-16 10:29
// @param     text string					代码文本
// @return    map[string]int				n-gram词频
func ComputeNgramFreq(text string, n int) map[string]int {
	// 移除注释
	text = RemoveComments(text)
	// 移除头文件
	text = RemoveHeads(text)
	// 移除空格、换行符、制表符等
	text = RemoveSpaces(text)
	// 小写化
	text = ToLowerCase(text)
	// 移除变量名
	//text = remove_variables(text)
	// 将文本按照n-gram分割成子串
	sub_strings := make([]string, 0)
	for i := 0; i <= len(text)-n; i++ {
		sub_strings = append(sub_strings, text[i:i+n])
	}
	// 统计每个子串出现的次数
	freq := make(map[string]int)
	for _, subString := range sub_strings {
		freq[subString]++
	}

	return freq
}

// @title    ComputeSimilarity
// @description  计算两个文本的n-gram词频相似度函数
// @auth      MGAronya             2022-9-16 10:29
// @param     text1 string, text2 string, n int					代码文本
// @return    float64											相似度
func ComputeSimilarity(text1 string, text2 string, n int) float64 {
	// 计算n-gram词频
	freq1 := ComputeNgramFreq(text1, n)
	freq2 := ComputeNgramFreq(text2, n)
	// 计算相似度
	common_keys := make(map[string]bool)
	for key := range freq1 {
		common_keys[key] = true
	}
	for key := range freq2 {
		common_keys[key] = true
	}

	numerator := 0
	for key := range common_keys {
		numerator += freq1[key] * freq2[key]
	}

	denominator := 0
	for _, freq := range freq1 {
		denominator += freq * freq
	}
	denominator *= 0
	for _, freq := range freq2 {
		denominator += freq * freq
	}

	return float64(numerator) / float64(denominator)
}

// @title    SimilarityJudge
// @description  计算并给出矩阵图中的连通块
// @auth      MGAronya             2022-9-16 10:29
// @param     arr [][]float64, stValue float64					代码文本
// @return    float64
func SimilarityJudge(arr [][]float64, stValue float64) ([][]int, error) {
	var n int = len(arr)
	p := make([]int, n)

	// TODO 参数校验
	if n == 0 {
		fmt.Printf("arr参数不能为空")
		return nil, fmt.Errorf("arr参数不能为空")
	}
	if stValue < 0 || stValue > 1 {
		return nil, fmt.Errorf("stValue参数超出范围")
	}

	// TODO 寻根函数
	var find func(x int) int
	find = func(x int) int {
		if x != p[x] {
			p[x] = find(p[x])
		}
		return p[x]
	}

	// TODO 相似度判断
	for i := 0; i < n; i++ {
		p[i] = i
	}
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr[i]); j++ {
			if arr[i][j] >= stValue {
				p[find(i)] = find(j)
			}
		}
	}

	// TODO 联通块分组
	groups := make(map[int][]int)
	for i := 0; i < n; i++ {
		groups[find(i)] = append(groups[find(i)], i)
	}

	con := make([][]int, 0)
	for _, group := range groups {
		con = append(con, group)
	}

	return con, nil
}

// @title    ChatGPT
// @description  调用ChatGPT的api输入text后回应
// @auth      MGAronya             2022-9-16 10:29
// @param     text					chatgpt的输入
// @return    vo.Response, error					chatgpt应答，以及可能发生的错误
func ChatGPT(text []string, model string) (vo.Response, error) {
	authHeader := fmt.Sprintf("Bearer %s", common.ChatGPTKey)

	var messages []vo.Message

	for i := range text {
		messages = append(messages, vo.Message{Role: "system", Content: text[i]})
	}
	// 构建请求数据
	reqData := vo.Request{
		Model:    model,
		Messages: messages,
	}

	// 发送POST请求
	request := gorequest.New()
	resp, _, errs := request.Post(common.ChatGPTUrl).Set("Authorization", authHeader).Send(reqData).End()
	if errs != nil {
		fmt.Println("请求出错:", errs[0])
		return vo.Response{}, errs[0]
	}

	// 读取响应数据
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("响应数据读取错误:", err)
		return vo.Response{}, err
	}

	// 响应解析
	if resp.StatusCode == 200 {
		// 解析响应数据
		var respData vo.Response
		err = json.Unmarshal(respBody, &respData)
		if err != nil {
			fmt.Println("响应数据解析错误:", err)
			return vo.Response{}, err
		}
		return respData, nil
	} else {
		return vo.Response{}, fmt.Errorf("请求失败:" + resp.Status)
	}
}

// @title    RemoveFirstAndLastLine
// @description  除去字符串的第一行和最后一行
// @auth      MGAronya             2022-9-16 10:29
// @param     str string					输入字符串
// @return    string					除去字符串的第一行和最后一行后的字符串
func RemoveFirstAndLastLine(str string) string {
	lines := strings.Split(str, "\n")
	if len(lines) < 3 {
		return ""
	}
	return strings.Join(lines[1:len(lines)-1], "\n")
}

// @title    Unzip
// @description  用于解压缩指定的 ZIP 文件到目标目录
// @auth      MGAronya             2022-9-16 10:29
// @param     zipfile, destDir string					指定的文件以及指定的路径
// @return    error						可能的报错
func Unzip(zipfile, destDir string) error {
	// TODO 打开 ZIP 文件
	r, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer r.Close()

	// TODO 创建目标目录
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		return err
	}

	// TODO  遍历 ZIP 文件中的每个文件
	for _, f := range r.File {
		// TODO 解压缩单个文件
		err := ExtractFile(f, destDir)
		if err != nil {
			return err
		}
	}

	return nil
}

// @title    ExtractFile
// @description  用于解压缩指定的文件到目标目录
// @auth      MGAronya             2022-9-16 10:29
// @param     f *zip.File, destDir string					指定的文件以及指定的路径
// @return    error						可能的报错
func ExtractFile(f *zip.File, destDir string) error {
	// TODO 打开 ZIP 文件中的文件
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// TODO 构建解压缩后的文件路径
	path := filepath.Join(destDir, f.Name)

	// TODO 如果是目录，则创建
	if f.FileInfo().IsDir() {
		os.MkdirAll(path, f.Mode())
	} else {
		// TODO 确保目标目录存在
		os.MkdirAll(filepath.Dir(path), f.Mode())

		// TODO 创建目标文件
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer file.Close()

		// TODO 将 ZIP 文件中的数据复制到目标文件中
		_, err = io.Copy(file, rc)
		if err != nil {
			return err
		}
	}

	return nil
}

// @title    GetFiles
// @description   获取一个目录下的所有文件
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     folder string	指定目录
// @return    []string    所有文件的文件名
func GetFiles(folder string) ([]vo.File, error) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	res := make([]vo.File, 0)
	for _, file := range files {
		// TODO 尝试读出所有文件的相关信息
		var f vo.File
		f.Name = file.Name()
		if folder[len(folder)-1] == '/' {
			f.Path = folder + file.Name()
		} else {
			f.Path = folder + "/" + file.Name()
		}
		f.Size = file.Size()
		f.LastWriteTime = file.ModTime()
		if file.IsDir() {
			f.Type = "Dir"
		} else {
			f.Type = path.Ext(file.Name())
		}
		res = append(res, f)
	}
	return res, nil
}

// @title    PathExists
// @description   判断文件夹是否存在
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     path string	指定目录
// @return    bool, error    查看文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// @title    Mkdir
// @description   建立文件夹
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     path string	指定路径
// @return   error    查看是否出错
func Mkdir(dir string) error {
	exist, err := PathExists(dir)
	if err != nil {
		return err
	}

	if !exist {
		// TODO 创建文件夹
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// @title    CopyFile
// @description   复制文件
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     sourceFilePath, destinationDir string	指定路径
// @return   error    查看是否出错
func CopyFile(sourceFilePath, destinationDir string) error {
	// TODO 获取文件名
	fileName := filepath.Base(sourceFilePath)

	// TODO 拼接目标路径
	destinationPath := filepath.Join(destinationDir, fileName)

	// TODO 打开源文件
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// TODO 创建目标目录
	err = os.MkdirAll(destinationDir, os.ModePerm)
	if err != nil {
		return err
	}

	// TODO 创建目标文件
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// TODO 拷贝源文件内容到目标文件
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

// @title    RemoveFile
// @description   函数用于删除指定文件
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     sfilePath	指定路径下的文件
// @return   error    查看是否出错
func RemoveFile(filePath string) error {
	// TODO 检查文件是否存在
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("File does not exist")
		}
		return err
	}

	// TODO 删除文件
	err = os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

// @title    RenameFile
// @description   函数用于将指定文件重命名为新的文件名
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     filePath, newName string		指定的文件和新名字
// @return   error    查看是否出错
func RenameFile(filePath, newName string) error {
	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("File does not exist")
		}
		return err
	}

	// 重命名文件
	err = os.Rename(filePath, newName)
	if err != nil {
		return err
	}

	return nil
}

// @title    CopyDir
// @description   函数用于复制源目录及其内容到目标目录
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     filePath, newName string		指定的文件和复制的目标路径
// @return   error    查看是否出错
func CopyDir(sourceDir, destinationDir string) error {
	// TODO 创建目标目录
	err := os.MkdirAll(destinationDir, os.ModePerm)
	if err != nil {
		return err
	}

	// TODO 打开源目录
	dir, err := os.Open(sourceDir)
	if err != nil {
		return err
	}
	defer dir.Close()

	// TODO 读取源目录中的所有文件和子目录
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	// TODO 遍历源目录中的文件和子目录
	for _, fileInfo := range fileInfos {
		// TODO 构建文件/目录的源路径和目标路径
		sourcePath := filepath.Join(sourceDir, fileInfo.Name())
		destinationPath := filepath.Join(destinationDir, fileInfo.Name())

		if fileInfo.IsDir() {
			// TODO 如果是子目录，则递归调用 copyDir 函数复制子目录
			err = CopyDir(sourcePath, destinationPath)
			if err != nil {
				return err
			}
		} else {
			// TODO 如果是文件，则调用 copyFile 函数复制文件
			err = CopyFile(sourcePath, destinationPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// @title    RemoveDir
// @description   删除指定目录及其子目录和文件
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     path string		需要删除的目录路径
// @return   error    查看是否出错
func RemoveDir(path string) error {
	// TODO 打开目录
	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer dir.Close()

	// TODO 读取目录下的所有文件和子目录
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		// TODO 拼接文件路径
		filePath := filepath.Join(path, fileInfo.Name())

		if fileInfo.IsDir() {
			// TODO 若是子目录，则递归调用删除子目录
			err = RemoveDir(filePath)
			if err != nil {
				return err
			}
		} else {
			// TODO 若是文件，则直接删除
			err = os.Remove(filePath)
			if err != nil {
				return err
			}
		}
	}

	// TODO 删除目录
	err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}
