package hackpool

import (
	"sync"
)

type HackPool struct {
	state    chan struct{}
	messages chan interface{}
	numGo    int
	function func(interface{})
}

func New(numGoroutine int, function func(interface{})) *HackPool {
	return &HackPool{
		state:    make(chan struct{}, numGoroutine),
		messages: make(chan interface{}, numGoroutine),
		numGo:    numGoroutine,
		function: function,
	}
}

func (c *HackPool) QueueCount() int {
	return len(c.messages)
}

func (c *HackPool) Push(data interface{}) {
	c.messages <- data
}

func (c *HackPool) CloseQueue() {
	close(c.messages)
}

func (c *HackPool) Run() {

	var wg sync.WaitGroup

	// 阻塞,等待message有数据或close
	for v := range c.messages {

		// 增加state的目的是为了控制同时执行的协程数量
		c.state <- struct{}{}

		wg.Add(1)

		go func(value interface{}) {

			c.function(value)

			// 协程执行完毕, 下一个协程才能启动
			<-c.state

			wg.Done()
		}(v)
	}

	wg.Wait()
}
