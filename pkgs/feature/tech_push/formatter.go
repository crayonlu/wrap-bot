package tech_push

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/crayon/wrap-bot/pkgs/feature/tech_push/handlers"
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

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("\n【%s】\n", sourceName))
	for i := 0; i < maxItems; i++ {
		item := val.Index(i)
		if item.Kind() == reflect.Ptr && !item.IsNil() {
			item = item.Elem()
		}
		if !item.IsValid() {
			continue
		}

		switch v := item.Interface().(type) {
		case handlers.JuejinArticle:
			builder.WriteString(fmt.Sprintf("%d. %s | author:%s | hot:%d | link:%s\n", i+1, v.Title, v.Author, v.Hot, v.URL))
		case *handlers.JuejinArticle:
			if v != nil {
				builder.WriteString(fmt.Sprintf("%d. %s | author:%s | hot:%d | link:%s\n", i+1, v.Title, v.Author, v.Hot, v.URL))
			}
		case handlers.BilibiliVideo:
			builder.WriteString(formatBilibiliLine(i, v))
		case *handlers.BilibiliVideo:
			if v != nil {
				builder.WriteString(formatBilibiliLine(i, *v))
			}
		default:
			title := extractTitle(item)
			if title != "" {
				builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, title))
			}
		}
	}

	return builder.String()
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

func formatBilibiliLine(index int, video handlers.BilibiliVideo) string {
	desc := strings.TrimSpace(video.Desc)
	if len(desc) > 120 {
		desc = desc[:120] + "..."
	}

	return fmt.Sprintf("%d. %s | up:%s | hot:%d | link:%s | desc:%s\n",
		index+1,
		video.Title,
		video.Author,
		video.Hot,
		video.URL,
		desc,
	)
}
