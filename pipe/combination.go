package rx

//Share 共享数据源，当有观察者时，开始观察源，所有观察者取消订阅后，则取消对源的订阅
func Share() Deliver {
	return func(source Observable) Observable {
		children := make(map[Next]Next)
		sendAll := func(data Any) {
			for key := range children {
				key <- data
			}
		}
		closeAll := func() {
			for key := range children {
				close(key)
				delete(children, key)
			}
		}
		return func(next Next, stop Stop) {
			if len(children) == 0 {
				sNext := make(Next)
				sStop := make(Stop)
				go source(sNext, sStop)
				for {
					select {
					case data, ok := <-sNext:
						if ok {
							sendAll(data)
						} else {
							closeAll()
							return
						}
					case <-stop:
						delete(children, next)
						if len(children) == 0 {
							close(sStop)
						}
						return
					}
				}
			} else {
				children[next] = next
			}
		}
	}
}

//If 根据条件选择数据源
func If(f func() bool, trueOb Observable, falseOb Observable) Observable {
	return func(next Next, stop Stop) {
		if f() {
			trueOb(next, stop)
		} else {
			falseOb(next, stop)
		}
	}
}

//Race 谁先到达，发送谁，然后完成
func Race(sources ...Observable) Observable {
	count := len(sources)
	return func(next Next, stop Stop) {
		life := count
		race := func(r Next) {
			for d := range r {
				if _, ok := d.(error); !ok {
					next <- d
					close(stop)
					close(next)
					return
				}
			}
			life--
			if life == 0 {
				close(next)
			}
		}
		for _, source := range sources {
			sNext := make(Next)
			go source(sNext, stop)
			go race(sNext)
		}
	}
}

//Concat 连接多个数据源
func Concat(sources ...Observable) Observable {
	return func(next Next, stop Stop) {
		for _, source := range sources {
			sNext := make(Next)
			go source(sNext, stop)
			for d := range sNext {
				if isError(d, next) {
					return
				}
				next <- d
			}
		}
		close(next)
	}
}

//Merge 合并多个数据源
func Merge(sources ...Observable) Observable {
	count := len(sources)
	return func(next Next, stop Stop) {
		life := count
		merge := func(m Next) {
			for d := range m {
				next <- d
				if _, ok := d.(error); ok {
					break
				}
			}
			life--
			if life == 0 {
				close(next)
			}
		}
		for _, source := range sources {
			sNext := make(Next)
			go source(sNext, stop)
			go merge(sNext)
		}
	}
}

//CombineLatest 合并多个流的最新数据
func CombineLatest(sources ...Observable) Observable {
	count := len(sources)
	return func(next Next, stop Stop) {
		buffer := make([]Any, count)
		nRun := 0
		life := count
		combine := func(n Next, i int) {
			state := 0
			for d := range n {
				if isError(d, next) {
					close(stop)
					return
				}
				switch state {
				case 0:
					nRun++
					buffer[i] = d
					if nRun == count {
						state = 2
						next <- buffer
					} else {
						state = 1
					}
				case 1:
					buffer[i] = d
					if nRun == count {
						state = 2
						next <- buffer
					}
				case 2:
					buffer[i] = d
					next <- buffer
				}
				break
			}
			life--
			if life == 0 {
				close(next)
			}
		}
		for i, source := range sources {
			sNext := make(Next)
			go source(sNext, stop)
			go combine(sNext, i)
		}
	}
}

//StartWith 先发送一些数据
func StartWith(xs ...interface{}) Deliver {
	return func(source Observable) Observable {
		return func(next Next, s Stop) {
			stopped := false
			go func() {
				<-s
				stopped = true
			}()
			for d := range xs {
				if stopped {
					return
				}
				next <- d
			}
			if !stopped {
				source(next, s)
			}
		}
	}
}
