package rx

import "errors"

//Take 获取最多count数量的事件，然后完成
func Take(count int) Deliver {
	return func(source Observable) Observable {
		return func(next Next, stop Stop) {
			_count := count
			_deliver(func(d Any, n Next, s Stop) (goon bool) {
				_count--
				if goon = _count > 0; goon {
					n <- d
				}
				return
			}, source, next, stop)
		}
	}
}

//TakeWhile 如果测试函数返回false则完成
func TakeWhile(f func(interface{}) bool) Deliver {
	return deliver(func(d Any, n Next, s Stop) (goon bool) {
		if goon = f(d); goon {
			n <- d
		}
		return
	})
}

//TakeUntil 直到开关事件流发出事件前一直接受事件
func TakeUntil(sSrc Observable, delivers ...Deliver) Deliver {
	sSrc = Pipe(sSrc, delivers...)
	return func(source Observable) Observable {
		return func(next Next, s Stop) {
			ssNext := make(Next)
			go sSrc(ssNext, s)
			go source(next, s)
			for {
				select {
				case <-s:
					return
				case d, ok := <-ssNext:
					if _, isError := d.(error); ok && !isError {
						close(s)
						close(next)
						return
					}
				}
			}
		}
	}
}

//TakeLast 取最后几个,输出的是一个channel，里面缓存着数据
func TakeLast(count int) Deliver {
	return func(source Observable) Observable {
		return func(next Next, s Stop) {
			sNext := make(Next)
			buffer := make(Next, count)
			go source(sNext, s)
			for {
				if d, ok := <-sNext; ok {
					if isError(d, next) {
						return
					}
					if len(buffer) == count {
						<-buffer
					}
					buffer <- d
				} else {
					next <- buffer
					close(next)
					return
				}
			}
		}
	}
}

//Skip 跳过最多count数量的事件，然后开始传送
func Skip(count int) Deliver {
	return func(source Observable) Observable {
		return func(next Next, stop Stop) {
			_count := count
			_deliver(func(d Any, n Next, s Stop) bool {
				if _count == 0 {
					n <- d
				} else {
					_count--
				}
				return true
			}, source, next, stop)
		}
	}
}

//SkipWhile 如果测试函数返回false则开始传送
func SkipWhile(f func(interface{}) bool) Deliver {
	return func(source Observable) Observable {
		return func(next Next, stop Stop) {
			activate := false
			_deliver(func(d Any, n Next, s Stop) bool {
				if activate {
					n <- d
				} else {
					activate = !f(d)
				}
				return true
			}, source, next, stop)
		}
	}
}

//SkipUntil 直到开关事件流发出事件前一直跳过事件
func SkipUntil(sSrc Observable, delivers ...Deliver) Deliver {
	sSrc = Pipe(sSrc, delivers...)
	return func(source Observable) Observable {
		return func(next Next, s Stop) {
			activate := false
			ssNext := make(Next)
			sNext := make(Next)
			stopS := make(Stop)
			go sSrc(ssNext, stopS)
			go source(sNext, s)
			for {
				select {
				case data, ok := <-sNext:
					if !ok {
						close(next)
						return
					}
					if activate {
						next <- data
					}
				case <-s:
					close(stopS)
					return
				case d, ok := <-ssNext:
					if _, isError := d.(error); ok && !isError {
						activate = true
						close(stopS)
					}
				}
			}
		}
	}
}

//IgnoreElements 忽略所有元素
func IgnoreElements(source Observable) Observable {
	return func(next Next, s Stop) {
		sNext := make(Next)
		go source(sNext, s)
		for d := range sNext {
			if isError(d, next) {
				return
			}
		}
		close(next)
	}
}

//Filter Filter items emitted by the source Observable by only emitting those that satisfy a specified predicate.
func Filter(f func(interface{}) bool) Deliver {
	return deliver(func(d Any, n Next, s Stop) bool {
		if f(d) {
			n <- d
		}
		return true
	})
}

//ElementAt Emits the single value at the specified index in a sequence of emissions from the source Observable.
func ElementAt(count int, defaultValue interface{}) Deliver {
	return func(source Observable) Observable {
		return func(next Next, stop Stop) {
			_count := count
			sNext := make(Next)
			go source(sNext, stop)
			for {
				select {
				case d, ok := <-sNext:
					if !ok {
						if defaultValue != nil {
							next <- defaultValue
						}
						close(next)
						return
					}
					if isError(d, next) {
						return
					} else if _count--; _count == 0 {
						next <- d
						close(stop)
						close(next)
						return
					}
				case <-stop:
					return
				}
			}
		}
	}
}

//Throttle Emits a value from the source Observable, then ignores subsequent source values for a duration determined by another Observable, then repeats this process.
func Throttle(durationSelector func(interface{}) Observable, leading bool, trailing bool) Deliver {
	return func(source Observable) Observable {
		return func(next Next, stop Stop) {
			var tNext Next
			var tStop Stop
			sNext, cache := make(Next), make(chan Any)
			throttle := func(value Any) {
				tNext, tStop = make(Next), make(Stop)
				go durationSelector(value)(tNext, tStop)
			}
			send := func() {
				select {
				case d := <-cache:
					next <- d
					throttle(d)
				default:
				}
			}
			throttleDone := func() {
				if tStop != nil {
					close(tStop)
					tStop = nil
				}
				if trailing {
					send()
				}
			}
			go source(sNext, stop)
			for {
				select {
				case <-stop:
					if tStop != nil {
						close(tStop)
					}
					return
				case <-tNext:
					throttleDone()
					return
				case d, ok := <-sNext:
					if ok {
						if isError(d, next) {
							return
						}
						select {
						case <-cache:
						default:
						}
						cache <- d
						if tStop == nil {
							if leading {
								send()
							} else {
								throttle(d)
							}
						}
					} else {
						throttleDone()
						close(next)
					}
				}
			}
		}
	}
}

//Audit Ignores source values for a duration determined by another Observable, then emits the most recent value from the source Observable, then repeats this process.
func Audit(durationSelector func(interface{}) Observable) Deliver {
	return Throttle(durationSelector, false, true)
}

//FindIndex Emits only the index of the first value emitted by the source Observable that meets some condition.
func FindIndex(f func(interface{}) bool) Deliver {
	return func(source Observable) Observable {
		return func(next Next, stop Stop) {
			i := 0
			_deliver(func(d Any, n Next, s Stop) bool {
				if f(d) {
					n <- i
					return false
				}
				i++
				return true
			}, source, next, stop)
		}
	}
}

//First Emits only the first value (or the first value that meets some condition) emitted by the source Observable.
func First(f func(interface{}) bool, defaultValue interface{}) Deliver {
	return func(source Observable) Observable {
		return func(next Next, stop Stop) {
			sNext := make(Next)
			value := defaultValue
			go source(sNext, stop)
			for {
				select {
				case d, ok := <-sNext:
					if !ok {
						if value != nil {
							next <- value
						} else {
							next <- errors.New("no elements in sequence")
						}
						close(next)
						return
					}
					if isError(d, next) {
						return
					} else if f == nil || f(d) {
						next <- d
						close(stop)
						close(next)
						return
					}
				case <-stop:
					return
				}
			}
		}
	}
}

//Last Returns an Observable that emits only the last item emitted by the source Observable. It optionally takes a predicate function as a parameter, in which case, rather than emitting the last item from the source Observable, the resulting Observable will emit the last item from the source Observable that satisfies the predicate.
func Last(f func(interface{}) bool, defaultValue interface{}) Deliver {
	return func(source Observable) Observable {
		return func(next Next, stop Stop) {
			sNext := make(Next)
			value := defaultValue
			go source(sNext, stop)
			for {
				select {
				case d, ok := <-sNext:
					if !ok {
						if value != nil {
							next <- value
						} else {
							next <- errors.New("no elements in sequence")
						}
						close(next)
						return
					}
					if isError(d, next) {
						return
					} else if f == nil || f(d) {
						value = d
						return
					}
				case <-stop:
					return
				}
			}
		}
	}
}
