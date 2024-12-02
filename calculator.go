package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

func Errors(expression string) error {
	len := len(expression)
	flag := false
	start := 0
	end := 0

	for i := 0; i < len; i++ {
		curr := expression[i]
		next := byte(0)
		if i < len-1 {
			next = expression[i+1]
		}

		if curr == '(' {
			start++
		}
		if curr == ')' {
			end++
		}
		if 48 <= curr && curr <= 57 && !flag {
			flag = true
		}

		switch {
		case i == 0 && (curr == ')' || curr == '*' || curr == '+' || curr == '-' || curr == '/'):
			return errors.New("the first character is the operator")
		case i == len-1 && (curr == '(' || curr == '*' || curr == '+' || curr == '-' || curr == '/'):
			return errors.New("the last character is the operator")
		case curr == '(' && next == ')':
			return errors.New("empty brackets")
		case curr == ')' && next == '(':
			return errors.New("no symbol between brackets")
		case (curr == '*' || curr == '+' || curr == '-' || curr == '/') && (next == '*' || next == '+' || next == '-' || next == '/'):
			return errors.New("the two operands are next to each other")
		case curr != ' ' && (curr < '(' || curr > '9'):
			return errors.New("the wrong character was found")
		case len <= 2:
			return errors.New("invalid expression")
		}
	}

	if start > end {
		return errors.New("the bracket is not closed")
	} else if end > start {
		return errors.New("the bracket is not open")
	}
	if !flag {
		return errors.New("operands not found")
	}
	return nil
}

func (s *Stack) LineToStacks(expression string) {
	var tmp string
	var len int = len([]rune(expression))

	for index, char := range expression {
		switch {
		case char == ' ':
			continue
		case '0' <= char && char <= '9' || char == '.' || char == ',':
			tmp += string(char)
			if index == len-1 {
				num, _ := strconv.ParseFloat(tmp, 64)
				s.numbers = append(s.numbers, num)
				tmp = ""
			}
		case char == '(' || char == ')' || char == '*' || char == '+' || char == '-' || char == '/':
			if tmp != "" {
				num, _ := strconv.ParseFloat(tmp, 64)
				s.numbers = append(s.numbers, num)
				tmp = ""
			}
			s.operands = append(s.operands, string(char))
		}
	}
}

type Stack struct {
	numbers  []float64
	operands []string
}

type StackOperators interface {
	Push(interface{})
	Pop(string) interface{}
}

func (s *Stack) Push(item interface{}) {
	switch char := item.(type) {
	case float64:
		s.numbers = append(s.numbers, char)
	case string:
		s.operands = append(s.operands, char)
	}
}

func (s *Stack) Pop(StackType string) interface{} {
	switch StackType {
	case "num":
		len := len(s.numbers)
		value := s.numbers[len-1]
		s.numbers = s.numbers[:len-1]
		return value
	case "op":
		len := len(s.operands)
		value := s.operands[len-1]
		s.operands = s.operands[:len-1]
		return value
	}
	return 0
}

func Operations(x, y float64, operand string) (float64, error) {
	switch operand {
	case "+":
		return x + y, nil
	case "-":
		return x - y, nil
	case "*":
		return x * y, nil
	case "/":
		if y == 0 {
			return 0, errors.New("division by zero")
		}
		return x / y, nil
	}
	fmt.Println(operand)
	return 0, errors.New("unknown operator")
}

func Priority(operand string) int {
	switch operand {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}
	return 0
}

func Calc(expression string) (float64, error) {
	err := Errors(expression)
	if err != nil {
		return 0, err
	}

	tempNum := make([]float64, 0)
	tempOp := make([]string, 0)
	flag := false
	bracketNum := 0.0
	bracketOp := ""
	s := Stack{}
	s.LineToStacks(expression)

	for len(s.operands) > 0 {
		op := s.Pop("op").(string)

		if op == ")" {
			flag2 := true
			brackNum := make([]float64, 0)
			brackOp := make([]string, 0)

			for len(s.operands) > 0 && s.operands[len(s.operands)-1] != "(" {
				op = s.Pop("op").(string)

				if Priority(op) == 2 || !flag2 {
					y := s.Pop("num").(float64)
					x := s.Pop("num").(float64)
					result, err := Operations(x, y, op)

					if err != nil {
						return 0, err
					}
					s.Push(result)

				} else {
					switch op {
					case "-":
						brackOp = append(brackOp, "+")
						brackNum = append(brackNum, (s.Pop("num").(float64) * (-1)))
					case "+":
						brackOp = append(brackOp, op)
						brackNum = append(brackNum, s.Pop("num").(float64))
					}
				}

				if s.operands[len(s.operands)-1] == "(" && flag2 {
					s.operands = append(s.operands, brackOp...)
					s.numbers = append(s.numbers, brackNum...)
					flag2 = false
				}
			}

			s.Pop("op")
			if bracketOp != "" {
				s.operands = append(s.operands, bracketOp)
				s.numbers = append(s.numbers, bracketNum)
				bracketOp = ""
				bracketNum = 0.0
			}

		} else {
			if ((Priority(op) == 2 || flag) && len(s.operands) == 0) || ((Priority(op) == 2 || flag) && s.operands[len(s.operands)-1] != ")") {
				y := s.Pop("num").(float64)
				x := s.Pop("num").(float64)
				result, err := Operations(x, y, op)

				if err != nil {
					return 0, err
				}
				s.Push(result)

			} else {
				switch op {
				case "-":
					tempOp = append(tempOp, "+")
					tempNum = append(tempNum, (s.Pop("num").(float64) * (-1)))
				case "+":
					tempOp = append(tempOp, op)
					tempNum = append(tempNum, s.Pop("num").(float64))
				}
			}
		}

		if len(s.operands) != 0 {
			if s.operands[len(s.operands)-1] == ")" {
				bracketNum = s.Pop("num").(float64)
				bracketOp = op

			}
		}

		if len(s.operands) == 0 && !flag {
			s.operands = tempOp
			s.numbers = append(s.numbers, tempNum...)
			flag = true
		}
	}

	return s.numbers[0], nil
}

func main() {
	fmt.Print("Enter the expression: ")
	var (
		expression string
	)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	expression = scanner.Text()

	result, err := Calc(expression)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}

}
