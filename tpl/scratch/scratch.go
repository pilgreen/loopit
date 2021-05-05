// Copyright 2018 The Hugo Authors. All rights reserved.
// Simplified copy of Hugo's scratch.go file

package scratch

import (
  _os "os"
	"sync"
)

// Scratch is a writable context used for stateful operations in Page/Node rendering.
type Scratch struct {
	values map[string]interface{}
	mu     sync.RWMutex
}

// Set stores a value with the given key in the Node context.
// This value can later be retrieved with Get.
func (c *Scratch) Set(key string, value interface{}) string {
	c.mu.Lock()
	c.values[key] = value
	c.mu.Unlock()
	return ""
}

// Reset deletes the given key
func (c *Scratch) Delete(key string) string {
	c.mu.Lock()
	delete(c.values, key)
	c.mu.Unlock()
	return ""
}

// Get returns a value previously set by Add or Set
func (c *Scratch) Get(key string) interface{} {
	c.mu.RLock()
	val := c.values[key]
	c.mu.RUnlock()

	return val
}

func NewScratch() *Scratch {
	return &Scratch{values: make(map[string]interface{})}
}

// Retrieve an environment variable
func Getenv(key string) (string) {
  return _os.Getenv(key)
}
