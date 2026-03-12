/*
some std errors
*/
package errors

import (
	"fmt"
)

// std error
type Err struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e Err) Error() string {
	return fmt.Sprintf("error(code = %d, msg = %s)", e.Code, e.Msg)
}

type IErr struct {
	Err
}

func NewIErr(code int, msg string) error {
	return IErr{
		Err: Err{
			Code: code,
			Msg:  msg,
		},
	}
}

type BErr struct {
	Err
}

func NewBErr(code int, msg string) error {
	return BErr{
		Err: Err{
			Code: code,
			Msg:  msg,
		},
	}
}

type CErr struct {
	Err
}

func NewCErr(code int, msg string) error {
	return CErr{
		Err: Err{
			Code: code,
			Msg:  msg,
		},
	}
}
