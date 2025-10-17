package tech_push

import (
	"fmt"
	"reflect"
)

func formatContentForAI(sourceName string, items interface{}, limit int) string {
	val := reflect.ValueOf(items)
	if val.Kind() != reflect.Slice {
		return ""
	}

	maxItems := val.Len()
	if limit > 0 && limit < maxItems {
		maxItems = limit
	}

	content := fmt.Sprintf("\n【%s】\n", sourceName)
	for i := 0; i < maxItems; i++ {
		item := val.Index(i)
		title := extractTitle(item)
		if title != "" {
			content += fmt.Sprintf("%d. %s\n", i+1, title)
		}
	}

	return content
}

func extractTitle(val reflect.Value) string {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return ""
	}

	titleField := val.FieldByName("Title")
	if titleField.IsValid() && titleField.Kind() == reflect.String {
		return titleField.String()
	}
	return ""
}
