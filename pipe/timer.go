package rx

import (
	"container/list"
	"sync"
	"time"
)

//Interval 定时器 period 毫秒数
func Interval(period time.Duration) Observable {
	return func(n Next, s Stop) {
		timer := time.NewTicker(period * time.Millisecond)
		i := 0
		for {
			select {
			case <-s:
				timer.Stop()
				return
			case <-timer.C:
				n <- i
				i++
			}
		}
	}
}

//Delay 延迟一定毫秒数后激活
func Delay(delay time.Duration) Deliver {
	type bufferData struct {
		t time.Time
		d Any
	}
	delay = delay * time.Millisecond //转换成Duration
	return func(source Observable) Observable {
		return func(next Next, stop Stop) {
			sNext := make(Next)
			buffer := list.New()
			go source(sNext, stop)
			wg := sync.WaitGroup{}
			pop := func(t time.Duration) {
				defer wg.Done()
				for buffer.Len() > 0 {
					select {
					case <-stop:
						return
					case <-time.After(t):
						e := buffer.Front()
						buffer.Remove(e)
						data := e.Value.(bufferData)
						next <- data.d
						if buffer.Len() > 0 {
							t = buffer.Front().Value.(bufferData).t.Sub(data.t)
						} else {
							return
						}
					}
				}
			}
			for {
				select {
				case d, ok := <-sNext:
					if !ok {
						select {
						case <-stop:
							return
						case <-time.After(delay):
							wg.Wait()
							close(next)
						}
						return
					}
					if isError(d, next) {
						return
					} else {
						buffer.PushBack(bufferData{time.Now(), d})
						if buffer.Len() == 1 {
							wg.Add(1)
							go pop(delay)
						}
					}
				case <-stop:
					return
				}
			}
		}
	}
}
