package rx

import (
	"container/list"
	"sync"
	"time"
)

//Interval Creates an Observable that emits sequential numbers every specified interval of time
//period: The interval size in milliseconds
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

//Delay Delays the emission of items from the source Observable by a given timeout.
//delay: The delay duration in milliseconds
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
					}
					buffer.PushBack(bufferData{time.Now(), d})
					if buffer.Len() == 1 {
						wg.Add(1)
						go pop(delay)
					}
				case <-stop:
					return
				}
			}
		}
	}
}

//Timer Creates an Observable that starts emitting after an initialDelay and emits ever increasing numbers after each period of time thereafter.
//Its like Interval, but you can specify when should the emissions start.
func Timer(dueTime time.Duration, period time.Duration) Observable {
	interval := Interval(period)
	return func(n Next, s Stop) {
		select {
		case <-s:
			return
		case <-time.After(dueTime * time.Millisecond):
			interval(n, s)
		}
	}
}
