package bot

import (
	"sync"
)

type Context struct {
	Event    *Event
	handlers []HandlerFunc
	index    int
	mu       sync.RWMutex
	Keys     map[string]interface{}
}

func newContext(event *Event, handlers []HandlerFunc) *Context {
	return &Context{
		Event:    event,
		handlers: handlers,
		index:    -1,
		Keys:     make(map[string]interface{}),
	}
}

func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = len(c.handlers)
}

func (c *Context) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Keys[key] = value
}

func (c *Context) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.Keys[key]
	return value, exists
}

func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

func (c *Context) GetString(key string) string {
	if val, ok := c.Get(key); ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (c *Context) GetInt64(key string) int64 {
	if val, ok := c.Get(key); ok {
		if i, ok := val.(int64); ok {
			return i
		}
	}
	return 0
}

func (c *Context) GetBool(key string) bool {
	if val, ok := c.Get(key); ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}
