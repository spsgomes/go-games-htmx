package handlers

import "html/template"

type PageData struct {
	Title string
	Intro string
}

var CustomTemplateFuncs template.FuncMap = map[string]interface{}{
	"isLast": func(index int, len int) bool {
		return index == len-1
	},
	"isNotLast": func(index int, len int) bool {
		return index != len-1
	},
}
