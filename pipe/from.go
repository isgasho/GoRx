package rx

//FromArray 从一个切片创建事件流
func FromArray(array []interface{}) Observable {
	return func(n Next, s Stop) {
		for _, item := range array {
			select {
			case <-s:
				return
			default:
				n <- item
			}
		}
		close(n)
	}
}

//Of 发送一系列值
func Of(array ...interface{}) Observable {
	return FromArray(array)
}

//FromChan 把一个chan转换成事件流
func FromChan(source chan interface{}) Observable {
	return func(n Next, s Stop) {
		for {
			select {
			case <-s:
				return
			case item, ok := <-source:
				if ok {
					n <- item
				} else {
					close(n)
					return
				}
			default:
			}
		}
	}
}

//From 通用数据源
func From(source interface{}) Observable {
	switch source.(type) {
	case []interface{}:
		return FromArray(source.([]interface{}))
	case chan interface{}:
		return FromChan(source.(chan interface{}))
	}
	return Of(source)
}
func empty(n Next, s Stop) {
	close(n)
}

//Empty Creates an Observable that emits no items to the Observer and immediately emits a complete notification.
func Empty() Observable {
	return empty
}

func never(n Next, s Stop) {
}

//Never Creates an Observable that do nothing forever and ever
func Never() Observable {
	return never
}

//Throw Creates an Observable that emits no items to the Observer and immediately emits an error notification.
func Throw(e error) Observable {
	return func(n Next, s Stop) {
		n <- e
		close(n)
	}
}

//Defer Creates an Observable that, on subscribe, calls an Observable factory to make an Observable for each new Observer.
//Creates the Observable lazily, that is, only when it is subscribed.
func Defer(f func() Observable) Observable {
	return func(n Next, s Stop) {
		f()(n, s)
	}
}
