// Copyright (c) 2014 José Carlos Nieto, https://menteslibres.net/xiam
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package cache

import (
	"math/rand"
	"sync"
	"time"
)

const (
	maxCachedObjects    = 1024 * 4
	mapCleanDivisor     = 1000
	mapCleanProbability = 1
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Cache holds a map of volatile key -> values.
type Cache struct {
	cache map[string]string
	mu    *sync.Mutex
}

// NewCache() initializes a new caching space.
func NewCache() (c *Cache) {
	c = &Cache{
		mu: new(sync.Mutex),
	}
	c.Clear()
	return
}

// Read() attempts to retrieve a cached value from memory. If the value does
// not exists returns an empty string and false.
func (self *Cache) Read(i Cacheable) (data string, ok bool) {
	h := i.Hash()

	if h != "" {
		self.mu.Lock()
		data, ok = self.cache[h]
		self.mu.Unlock()
	}

	return data, ok
}

// Write() stores a value in memory. If the value already exists it gets
// overwritten.
func (self *Cache) Write(i Cacheable, s string) {

	if maxCachedObjects > 0 && maxCachedObjects < len(self.cache) {
		self.Clear()
	} else if rand.Intn(mapCleanDivisor) <= mapCleanProbability {
		self.Clear()
	}

	h := i.Hash()

	if h != "" {
		self.mu.Lock()
		self.cache[h] = s
		self.mu.Unlock()
	}
}

// Clear() generates a new memory space, leaving the old memory unreferenced,
// so it can be claimed by the garbage collector.
func (self *Cache) Clear() {
	self.mu.Lock()
	self.cache = make(map[string]string)
	self.mu.Unlock()
}
