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
