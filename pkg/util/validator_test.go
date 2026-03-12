package util

import (
	"fmt"
	"testing"

	"github.com/labstack/echo/v4"
)

func Test_DefaultRules(t *testing.T) {
	type Base struct {
		ID    int    `json:"id" check:"min=1,max=1000,request"`
		Key   string `json:"key" check:"uuid"`
		Email string `json:"email" check:"email"`
		Phone string `json:"phone" check:"fixed=11"`
	}
	bases := []Base{
		// just use base_0 to test all condition...
		// Base{ID: 0, Key: "647a0b56-146b-4db3-9b61", Email: "test0@puzzle.com", Phone: "7211627891"},
		Base{ID: 1, Key: "647a0b56-146b-4db3-9b61-c29f9b0094f6", Email: "test1@puzzle.com", Phone: "18211627891"},
		Base{ID: 2, Key: "92e29510-dc3d-48eb-add2-eec39a29a9df", Email: "test2@puzzle.com", Phone: "15211627891"},
		Base{ID: 3, Key: "818d232b-3a53-4acb-8f4c-30e79fa68222", Email: "test3@puzzle.com", Phone: "13211627891"},
		Base{ID: 4, Key: "3788b8b4-c833-4e44-8b3d-adce4b36f961", Email: "test4@puzzle.com", Phone: "16211627891"},
		Base{ID: 5, Key: "90568701-8279-4e58-9bf9-0796a0e64154", Email: "test5@puzzle.com", Phone: "11211627891"},
	}
	e := echo.New()

	checker := NewValidator(e)
	checker.SetupDefaultRules()
	for i, b := range bases {
		b := b
		t.Run(fmt.Sprintf("base_%d", i), func(t *testing.T) {
			err := checker.Check(b)
			if err != nil {
				t.Errorf("%s\n", err.Error())
			}
		})
	}
}
