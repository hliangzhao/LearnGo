package models

type Video struct {
	Id          int         `uri:"id"`
	Title       string      `json:"title"`              // 使用tag告知数据转换的方式
	Description string      `json:"description"`
}
