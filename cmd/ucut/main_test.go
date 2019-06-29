package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxStringWidth(t *testing.T) {
	type TestData struct {
		desc      string
		line      string
		delimiter string
		fields    string
		expect    []string
	}
	tds := []TestData{
		{desc: "フィールドを直接指定", line: "あ,い,う", delimiter: ",", fields: "1,3", expect: []string{"あ", "う"}},
		{desc: "フィールドがマルチバイト文字", line: "あんいんう", delimiter: "ん", fields: "1,3", expect: []string{"あ", "う"}},
		{desc: "省略記法 2-4", line: "1,2,3,4,5", delimiter: ",", fields: "2-4", expect: []string{"2", "3", "4"}},
		{desc: "省略記法 2-", line: "1,2,3,4,5", delimiter: ",", fields: "2-", expect: []string{"2", "3", "4", "5"}},
		{desc: "省略記法 1,3-", line: "1,2,3,4,5", delimiter: ",", fields: "1,3-", expect: []string{"1", "3", "4", "5"}},
		{desc: "省略記法 -3", line: "1,2,3,4,5", delimiter: ",", fields: "-3", expect: []string{"1", "2", "3"}},
		{desc: "省略記法 -3,5", line: "1,2,3,4,5", delimiter: ",", fields: "-3,5", expect: []string{"1", "2", "3", "5"}},
		{desc: "順序変更 3,2,1", line: "1,2,3,4,5", delimiter: ",", fields: "3,2,1", expect: []string{"3", "2", "1"}},
		{desc: "順序変更 3,2,3", line: "1,2,3,4,5", delimiter: ",", fields: "3,2,3", expect: []string{"3", "2", "3"}},
	}
	for _, v := range tds {
		got, err := ucut(v.line, v.delimiter, v.fields)
		assert.Equal(t, v.expect, got, v.desc)
		assert.Nil(t, err, v.desc)
	}
}
