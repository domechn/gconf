// Author : dmc
//
// Date: 2018/9/4 上午10:28
//
// Description:
package gconf

import (
	"fmt"
	"testing"
)

type sRead struct {
	User     string                 `json:"user" default:"hello"`
	Password string                 `json:"password"`
	Sex      int                    `json:"sex"`
	Married  bool                   `json:"married"`
	Ts       map[string]interface{} `json:"ts"`
	Tlist    []string               `json:"tlist"`
	Tll      struct {
		Tlll int `json:"tlll"`
	} `json:"tll"`
}

func init() {
	// err := Register("ymlRead", "./test_file/test.yml")
	// if err != nil {
	// 	panic(err)
	// }
	err := Register("jsonRead", "./test_file/test.json")
	if err != nil {
		panic(err)
	}
}

var testCase = []struct {
	want *sRead
}{
	{
		want: &sRead{
			User:     "domgoer",
			Password: "qq123456",
			Sex:      1,
			Married:  false,
			Ts:       map[string]interface{}{"tss": "dd"},
			Tlist:    []string{"a", "b", "c"},
			Tll: struct {
				Tlll int `json:"tlll"`
			}{Tlll: 123},
		},
	},
}

func TestRead2Struct(t *testing.T) {
	yRead := &sRead{}
	jRead := &sRead{}
	Read2StructByDefault("./test_file/test.yml", yRead)
	Read2Struct("./test_file/test.json", jRead)
	for _, v := range testCase {
		flag := compare(yRead, v.want)
		if !flag {
			fmt.Printf("123,%+v\n", yRead)
			t.Errorf("Read2Struct() appear error")
		}
	}
}

func compare(s1, s2 *sRead) bool {
	if s1.User != s2.User {
		return false
	}
	if s1.Password != s2.Password {
		return false
	}
	if s1.Sex != s2.Sex {
		return false
	}
	if s1.Married != s2.Married {
		return false
	}
	if s1.Ts["tss"] != s2.Ts["tss"] {
		return false
	}
	if s1.Tlist[0] != s2.Tlist[0] || s1.Tlist[1] != s2.Tlist[1] || s1.Tlist[2] != s2.Tlist[2] {
		return false
	}
	return true
}
