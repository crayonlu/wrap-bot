package bot

import (
	"log"
	"runtime/debug"
	"time"
)

func Recovery() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v\n%s", err, debug.Stack())
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}

func Logger() HandlerFunc {
	return func(ctx *Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)

		switch ctx.Event.PostType {
		case EventTypeMessage:
			log.Printf("[%s] Type=%s UserID=%d GroupID=%d Message=%s Duration=%v",
				ctx.Event.PostType,
				ctx.Event.MessageType,
				ctx.Event.UserID,
				ctx.Event.GroupID,
				ctx.Event.RawMessage,
				duration,
			)
		case EventTypeNotice:
			log.Printf("[%s] NoticeType=%s UserID=%d GroupID=%d Duration=%v",
				ctx.Event.PostType,
				ctx.Event.NoticeType,
				ctx.Event.UserID,
				ctx.Event.GroupID,
				duration,
			)
		case EventTypeRequest:
			log.Printf("[%s] RequestType=%s UserID=%d Comment=%s Duration=%v",
				ctx.Event.PostType,
				ctx.Event.RequestType,
				ctx.Event.UserID,
				ctx.Event.Comment,
				duration,
			)
		}
	}
}

func FilterEventType(eventTypes ...EventType) HandlerFunc {
	typeMap := make(map[EventType]bool)
	for _, t := range eventTypes {
		typeMap[t] = true
	}

	return func(ctx *Context) {
		if !typeMap[ctx.Event.PostType] {
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func FilterMessageType(messageTypes ...MessageType) HandlerFunc {
	typeMap := make(map[MessageType]bool)
	for _, t := range messageTypes {
		typeMap[t] = true
	}

	return func(ctx *Context) {
		if !typeMap[ctx.Event.MessageType] {
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func FilterGroup(groupIDs ...int64) HandlerFunc {
	groupMap := make(map[int64]bool)
	for _, id := range groupIDs {
		groupMap[id] = true
	}

	return func(ctx *Context) {
		if !groupMap[ctx.Event.GroupID] {
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func FilterUser(userIDs ...int64) HandlerFunc {
	userMap := make(map[int64]bool)
	for _, id := range userIDs {
		userMap[id] = true
	}

	return func(ctx *Context) {
		if !userMap[ctx.Event.UserID] {
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func RateLimit(maxRequests int, duration time.Duration) HandlerFunc {
	type entry struct {
		count     int
		resetTime time.Time
	}

	cache := make(map[int64]*entry)

	return func(ctx *Context) {
		userID := ctx.Event.UserID
		now := time.Now()

		e, exists := cache[userID]
		if !exists || now.After(e.resetTime) {
			cache[userID] = &entry{
				count:     1,
				resetTime: now.Add(duration),
			}
			ctx.Next()
			return
		}

		if e.count >= maxRequests {
			ctx.Abort()
			return
		}

		e.count++
		ctx.Next()
	}
}

func AdminOnly(adminIDs ...int64) HandlerFunc {
	adminMap := make(map[int64]bool)
	for _, id := range adminIDs {
		adminMap[id] = true
	}

	return func(ctx *Context) {
		if !adminMap[ctx.Event.UserID] {
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func GroupAdminOnly() HandlerFunc {
	return func(ctx *Context) {
		if ctx.Event.Sender == nil {
			ctx.Abort()
			return
		}

		role := ctx.Event.Sender.Role
		if role != "admin" && role != "owner" {
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func Authentication(allowedUsers []int64, allowedGroups []int64) HandlerFunc {
	userMap := make(map[int64]bool)
	for _, id := range allowedUsers {
		userMap[id] = true
	}

	groupMap := make(map[int64]bool)
	for _, id := range allowedGroups {
		groupMap[id] = true
	}

	return func(ctx *Context) {
		if len(userMap) == 0 && len(groupMap) == 0 {
			ctx.Next()
			return
		}

		if ctx.Event.MessageType == MessageTypePrivate {
			if len(userMap) > 0 && !userMap[ctx.Event.UserID] {
				ctx.Abort()
				return
			}
		} else if ctx.Event.MessageType == MessageTypeGroup {
			if len(groupMap) > 0 && !groupMap[ctx.Event.GroupID] {
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
