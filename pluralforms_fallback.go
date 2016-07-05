package gogettext

import (
	"fmt"
	"strings"
	"errors"
)

type Expression struct {
	Raw string
	Expr func(n uint32) int
}

func one_form(n uint32) int {
	return 0
}

func two_forms_singular_one(n uint32) int {
	if n != 1 {
		return 1
	} else {
		return 0
	}
}

func two_forms_singular_one_and_zero(n uint32) int {
	if n > 1 {
		return 1
	} else {
		return 0
	}
}

func three_forms_special_zero(n uint32) int {
	if n % 10 == 1 && n % 100 != 11 {
		return 0
	} else if n != 0 {
		return 1
	} else {
		return 2
	}
}

func three_forms_special_one_and_two(n uint32) int {
	if n == 1 {
		return 0
	} else if n == 2{
		return 1
	} else {
		return 2
	}
}

func three_forms_romanian(n uint32) int {
	if n == 1 {
		return 0
	} else if n == 0 || (n % 100 > 0 && n % 100 < 20) {
		return 1
	} else {
		return 2
	}
}

func three_forms_lithuanian(n uint32) int {
	if n % 10 == 1 && n % 100 != 11 {
		return 0
	} else if n % 10 >= 2 && (n % 100 < 10 || n % 100 >= 20) {
		return 1
	} else {
		return 2
	}
}

func three_forms_russian(n uint32) int {
	if n % 10 == 1 && n % 100 != 11 {
		return 0
	} else if n % 10 >= 2 && n % 10 <= 4 && (n % 100 < 10 || n % 100 >= 20) {
		return 1
	} else {
		return 2
	}
}

func three_forms_special_one_two_three_four(n uint32) int {
	if n == 1 {
		return 0
	} else if n >= 2 && n <= 4 {
		return 1
	} else {
		return 2
	}
}

func three_forms_polish(n uint32) int {
	if n == 1 {
		return 0
	} else if n % 10 >= 2 && n % 10 <= 4 && (n % 100 < 10 || n % 100 >= 20) {
		return 1
	} else {
		return 2
	}
}

func four_forms_slovenian(n uint32) int {
	if n % 100 == 1 {
		return 0
	} else if n % 100 == 2 {
		return 1
	} else if n % 100 == 3 || n % 100 == 4 {
		return 2
	} else {
		return 3
	}
}

func arabic(n uint32) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else if n == 2 {
		return 2
	} else if n % 100 >= 3 && n % 100 <= 10 {
		return 3
	} else if n % 100 >= 11 {
		return 4
	} else {
		return 5
	}
}

var forms = map[string]func(n uint32)int{
	"0": one_form,
	"n!=1": two_forms_singular_one,
	"n>1": two_forms_singular_one_and_zero,
	"n%10==1&&n%100!=11?0:n!=0?1:2": three_forms_special_zero,
	"n==1?0:n==2?1:2": three_forms_special_one_and_two,
	"n==1?0:(n==0||(n%100>0&&n%100<20))?1:2": three_forms_romanian,
	"n%10==1&&n%100!=11?0:n%10>=2&&(n%100<10||n%100>=20)?1:2": three_forms_lithuanian,
	"n%10==1&&n%100!=11?0:n%10>=2&&n%10<=4&&(n%100<10||n%100>=20)?1:2": three_forms_russian,
	"(n==1)?0:(n>=2&&n<=4)?1:2": three_forms_special_one_two_three_four,
	"n==1?0:n%10>=2&&n%10<=4&&(n%100<10||n%100>=20)?1:2": three_forms_polish,
	"n%100==1?0:n%100==2?1:n%100==3||n%100==4?2:3": four_forms_slovenian,
	"n==0?0:n==1?1:n==2?2:n%100>=3&&n%100<=10?3:n%100>=11?4:5": arabic,
}


func NewExpression(raw string) (Expression, error){
	nospc := strings.Trim(strings.Replace(raw, " ", "", -1), "()")
	expr := Expression{Raw:raw}
	f, ok := forms[nospc]
	if !ok {
		return expr, errors.New(fmt.Sprintf("Unsupported plural forms expression: %s", raw))
	}
	expr.Expr = f
	return expr, nil
}


func (expr Expression) Eval (n uint32) int {
	return expr.Expr(n)
}

func (expr Expression) String () string {
	return fmt.Sprintf("<%s>", expr.Raw)
}
