package tech_push

import (
	"fmt"
	"reflect"

	"github.com/crayon/bot_golang/pkgs/napcat"
)

func buildGenericNodes(sourceName string, items interface{}, limit int, botQQ int64) []napcat.ForwardNode {
	var nodes []napcat.ForwardNode

	val := reflect.ValueOf(items)
	if val.Kind() != reflect.Slice {
		return nodes
	}

	maxItems := val.Len()
	if limit > 0 && limit < maxItems {
		maxItems = limit
	}

	for i := 0; i < maxItems; i++ {
		item := val.Index(i)
		segments := structToSegments(item)

		node := napcat.NewMixedForwardNode(
			sourceName,
			botQQ,
			segments...,
		)
		nodes = append(nodes, node)
	}

	return nodes
}

func structToSegments(val reflect.Value) []napcat.MessageSegment {
	var segments []napcat.MessageSegment

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return segments
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if !fieldValue.CanInterface() {
			continue
		}

		fieldName := field.Name

		if fieldName == "Cover" && fieldValue.Kind() == reflect.String && fieldValue.String() != "" {
			segments = append(segments, napcat.NewImageSegment(fieldValue.String()))
			segments = append(segments, napcat.NewTextSegment("\n"))
			continue
		}

		if fieldName == "MobileUrl" {
			continue
		}

		var valueStr string
		switch fieldValue.Kind() {
		case reflect.String:
			valueStr = fieldValue.String()
		case reflect.Int, reflect.Int64:
			valueStr = fmt.Sprintf("%d", fieldValue.Int())
		case reflect.Float64:
			valueStr = fmt.Sprintf("%f", fieldValue.Float())
		default:
			valueStr = fmt.Sprintf("%v", fieldValue.Interface())
		}

		segments = append(segments, napcat.NewTextSegment(fmt.Sprintf("%s: %s\n", fieldName, valueStr)))
	}

	return segments
}
