package hackpool

import (
	"sync"
)

type HackPool struct {
	numGo    int
	messages chan []interface{}
	function func(...interface{})
}

func New(numGoroutine int, function func(...interface{})) *HackPool {
	return &HackPool{
		numGo:    numGoroutine,
		messages: make(chan []interface{}),
		function: function,
	}
}

func (c *HackPool) Push(data ...interface{}) {
	c.messages <- data
}

func (c *HackPool) CloseQueue() {
	close(c.messages)
}

func (c *HackPool) Run() {
	var wg sync.WaitGroup

	wg.Add(c.numGo)

	for i := 0; i < c.numGo; i++ {
		go func() {
			for v := range c.messages {
				c.function(v...)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
