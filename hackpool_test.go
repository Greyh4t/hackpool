package hackpool

import (
	"fmt"
	"testing"
)

func TestHackPool(t *testing.T) {

	var hp *HackPool
	numGoroutine := 2
	taskCount := 10000

	hp = New(numGoroutine, func(i ...interface{}) {

		fmt.Println(i[0].(int))

		//time.Sleep(time.Second * 2)
	})

	// 如果同步向chan push数据, 若元素个数超过 chan 定义的长度, 会造成死锁, 所以必须异步push数据
	go func() {

		for i := 0; i < taskCount; i++ {
			hp.Push(i)
		}

		// push任务结束后必须关闭, 否则死锁
		hp.CloseQueue()
	}()

	// 跑起来! 伙计
	hp.Run()
}
