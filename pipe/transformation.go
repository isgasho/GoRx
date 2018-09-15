package rx

import (
	"time"
)

//Scan Applies an accumulator function over the source Observable, and returns each intermediate result, with an optional seed value.
func Scan(f func(interface{}, interface{}) interface{}, seed ...interface{}) Deliver {
	return func(source Observable) Observable {
		return func(next Next, stop Stop) {
			sNext := make(Next)
			var aac Any
			hasValue := false
			if hasSeed := len(seed) == 1; hasSeed {
				aac = seed[0]
				hasValue = true
			}
			go source(sNext, stop)
			for {
				select {
				case d, ok := <-sNext:
					if !ok {
						close(next)
						return
					}
					if isError(d, next) {
						return
					}
					if hasValue {
						aac = f(aac, d)
					} else {
						hasValue = true
						aac = d
					}
					next <- aac
				case <-stop:
					return
				}
			}
		}
	}
}

//Map Applies a given project function to each value emitted by the source Observable, and emits the resulting values as an Observable.
func Map(f func(interface{}) interface{}) Deliver {
	return deliver(func(d Any, n Next, s Stop) bool {
		n <- f(d)
		return true
	})
}

//MapTo Emits the given constant value on the output Observable every time the source Observable emits a value.
func MapTo(x interface{}) Deliver {
	return deliver(func(d Any, n Next, s Stop) bool {
		n <- x
		return true
	})
}

//Pairwise Groups pairs of consecutive emissions together and emits them as an array of two values.
func Pairwise() Deliver {
	return func(source Observable) Observable {
		return func(next Next, stop Stop) {
			buffer := make(Next, 2)
			_deliver(func(d Any, next Next, stop Stop) bool {
				if len(buffer) == 2 {
					var result [2]Any
					result[0] = <-buffer
					result[1] = <-buffer
					next <- result
				} else {
					buffer <- d
				}
				return true
			}, source, next, stop)
		}
	}
}

//SwitchMap Projects each source value to an Observable which is merged in the output Observable, emitting values only from the most recently projected Observable.
func SwitchMap(f func(interface{}) Observable, combineResults func(interface{}, interface{}) interface{}) Deliver {
	return func(source Observable) Observable {
		return func(next Next, stop Stop) {
			waitInnerStop := false
			sNext := make(Next)
			var innerNext Next
			var innerStop Stop
			var outValue Any
			go source(sNext, stop)
			for {
				select {
				case <-innerStop:
					if waitInnerStop {
						close(next)
						return
					}
					innerStop = nil
				case d, ok := <-innerNext:
					if ok {
						if _, ok = d.(error); ok {
							innerStop = nil
						} else {
							if combineResults != nil {
								d = combineResults(outValue, d)
							}
							next <- d
						}
					} else {
						innerStop = nil
					}
				case <-stop:
					return
				case d, ok := <-sNext:
					if ok {
						if isError(d, next) {
							return
						}
						outValue = d
						if innerStop != nil {
							close(innerStop)
						}
						innerStop = make(Stop)
						innerNext = make(Next)
						go f(d)(innerNext, innerStop)
					} else {
						if innerStop == nil {
							close(next)
							return
						}
						waitInnerStop = true
					}
				}
			}
		}
	}
}

//SwitchMapTo Projects each source value to the same Observable which is flattened multiple times with switchMap in the output Observable.
func SwitchMapTo(source Observable, combineResults func(interface{}, interface{}) interface{}) Deliver {
	return SwitchMap(func(d interface{}) Observable {
		return source
	}, combineResults)
}

//BufferTime Buffers the source Observable values for a specific time period.
func BufferTime(period time.Duration, maxBufferSize int) Deliver {
	return func(source Observable) Observable {
		return func(next Next, stop Stop) {
			sNext := make(Next)
			buffer := make([]Any, 0)
			timer := time.NewTicker(period * time.Millisecond)
			go source(sNext, stop)
			for {
				select {
				case <-timer.C:
					next <- buffer
					buffer = make([]Any, 0)
				case <-stop:
					return
				case d, ok := <-sNext:
					if ok {
						if isError(d, next) {
							return
						}
						if len(buffer) < maxBufferSize {
							buffer = append(buffer, d)
						} else {
							next <- buffer
							buffer = make([]Any, 0)
						}
					} else {
						close(next)
						return
					}
				}
			}
		}
	}
}

//Repeat Returns an Observable that repeats the stream of items emitted by the source Observable at most count times.
func Repeat(count int) Deliver {
	return func(source Observable) Observable {
		return func(next Next, stop Stop) {
			buffer := make([]Any, 10)
			stopped := false
			go func() {
				<-stop
				stopped = true
			}()
			sNext := make(Next)
			go source(sNext, stop)
			for d := range sNext {
				if stopped || isError(d, next) {
					return
				}
				buffer = append(buffer, d)
			}
			for _count := count; _count > 0; _count-- {
				for d := range buffer {
					if stopped {
						return
					}
					next <- d
				}
			}
			close(next)
		}
	}
}
