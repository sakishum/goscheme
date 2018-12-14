package goscheme

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Expression interface{}

type Number float64

type String string

type Symbol string

type Boolean bool

type Function func(...Expression) Expression

type NilType struct{}

func (n NilType) String() string {
	return "()"
}

var syntaxes = [...]string{"define", "lambda", "if", "let", "cond", "begin"}

var NilObj = NilType{}

type Undef struct{}

func (u *Undef) String() string {
	return "<UNDEF>"
}

var undefObj = Undef{}

func IsNumber(exp Expression) bool {
	switch v := exp.(type) {
	case string:
		_, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return false
		}
		return true
	case Number:
		return true
	default:
		return false
	}
	return false
}

func IsString(exp Expression) bool {
	switch v := exp.(type) {
	case string:
		ok, err := regexp.MatchString("\".*\"", v)
		if ok && err == nil {
			return true
		}
		return false
	default:
		return false
	}
	return false
}

func IsSpecialSyntaxExpression(exp Expression, name string) bool {
	ops, ok := exp.([]Expression)
	if !ok {
		return false
	}
	operator := ops[0]
	return operator == name
}

func IsSyntaxExpression(exp Expression) bool {
	ops, ok := exp.([]Expression)
	if !ok {
		return false
	}
	operator := ops[0]

	for _, s := range syntaxes {
		if s == operator {
			return true

		}
	}
	return false
}

func IsSymbol(expression Expression) bool {
	_, ok := expression.([]Expression)
	if ok {
		return false
	}
	if IsNumber(expression) || IsString(expression) || IsBoolean(expression) || IsSyntaxExpression(expression) {
		return false
	}
	return true
}

func IsBoolean(exp Expression) bool {
	_, ok := exp.(bool)
	if ok {
		return true
	}
	return exp == "#t" || exp == "#f"
}

func IsTrue(exp Expression) bool {
	if !IsBoolean(exp) {
		return false
	}
	if exp == true {
		return true
	}
	if exp == "#t" {
		return true
	}
	return false
}

func IsNilObj(obj Expression) bool {
	switch obj.(type) {
	case NilType:
		return true
	default:
		return false
	}
}

func IsPair(obj Expression) bool {
	switch obj.(type) {
	case *Pair:
		return true
	default:
		return false
	}
}

type LambdaProcess struct {
	params []Symbol
	body   Expression
	env    *Env
}

func (lambda *LambdaProcess) call(env *Env, args ...Expression) Expression {
	return nil
}

// Should only use with pointer
type Pair struct {
	Car, Cdr Expression
}

func (p *Pair) IsNull() bool {
	return p.Car == nil && p.Cdr == nil
}

func (p *Pair) IsList() bool {
	currentPair := p
	for {
		if currentPair.IsNull() {
			return true
		}
		switch cdr := currentPair.Cdr.(type) {
		case *Pair:
			currentPair = cdr
		case NilType:
			return true
		default:
			return false
		}
	}
}

func (p *Pair) String() string {

	currentPair := p

	var strSlices []string

	for !currentPair.IsNull() {
		if IsPair(currentPair.Car) {
			strSlices = append(strSlices, currentPair.Car.(*Pair).String())
		} else {
			strSlices = append(strSlices, fmt.Sprintf("%v", currentPair.Car))
		}

		if IsPair(currentPair.Cdr) {
			currentPair = currentPair.Cdr.(*Pair)
		} else {
			if IsNilObj(currentPair.Cdr) {
				break
			}
			strSlices = append(strSlices, ".")
			strSlices = append(strSlices, fmt.Sprintf("%v", currentPair.Cdr))
			break
		}
	}

	return "(" + strings.Join(strSlices, " ") + ")"
}