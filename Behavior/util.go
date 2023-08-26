// @Title  util
// @Description  该文件提供行为映射以及各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Handle

import (
	"MGA_OJ/Interface"
	"fmt"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// behaviors	    定义了用户行为映射表
var Behaviors map[string]Interface.BehaviorInterface = map[string]Interface.BehaviorInterface{}

// @title    checkExpression
// @description  函数接收一个不带括号的表达式，查看表达式是否正确
// @auth      MGAronya             2022-9-16 10:29
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
			if _, ok := Behaviors[string(expr[j:i+1])]; !ok {
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

// @title    CheckExpression
// @description  函数接收一个表达式，查看表达式是否正确
// @auth      MGAronya             2022-9-16 10:29
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
// @auth      MGAronya             2022-9-16 10:29
// @param     expr string		表达式
// @return    int, string		计算表达式的值
func EvaluateExpression(expression string, userId uuid.UUID) (int, error) {
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

			value, err := Behaviors[varName].UserBehavior(userId)
			if err != nil {
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
					return 0, fmt.Errorf("invalid expression")
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
				return 0, fmt.Errorf("invalid expression")
			}
		} else if char == '+' || char == '-' || char == '*' || char == '/' {
			// TODO 处理运算符
			for len(operators) > 0 && (operators[len(operators)-1] == "*" || operators[len(operators)-1] == "/") {
				topOperator := operators[len(operators)-1]
				operators = operators[:len(operators)-1]

				if len(stack) < 2 {
					return 0, fmt.Errorf("invalid expression")
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
						return 0, fmt.Errorf("division by zero")
					}
					result = num1 / num2
				}
				// TODO 将计算结果压入栈中
				stack = append(stack, strconv.Itoa(result))
			}

			operators = append(operators, string(char))
		} else {
			return 0, fmt.Errorf(fmt.Sprintf("invalid character %c", char))
		}
	}

	for len(operators) > 0 {
		topOperator := operators[len(operators)-1]
		operators = operators[:len(operators)-1]

		if len(stack) < 2 {
			return 0, fmt.Errorf("invalid expression")
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
		return 0, fmt.Errorf("invalid expression")
	}
	// TODO 将结果转换为整数
	value, err := strconv.Atoi(stack[0])
	if err != nil {
		return 0, fmt.Errorf("invalid expression")
	}

	return value, nil
}
