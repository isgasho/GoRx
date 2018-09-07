package rx

//Reduce 累加器，在事件流完成时把结果发出去
func Reduce(f func(interface{}, interface{}) interface{}, seed ...interface{}) Deliver {
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
						if hasValue {
							next <- aac
						}
						close(next)
						return
					}
					if err, ok := d.(error); ok {
						next <- err
					} else {
						if hasValue {
							aac = f(aac, d)
						} else {
							hasValue = true
							aac = d
						}
					}
				case <-stop:
					return
				default:
				}
			}
		}
	}
}

//Count 计数
func Count(f func(interface{}) bool) Deliver {
	return Reduce(func(aac interface{}, c interface{}) interface{} {
		if f == nil || f(c) {
			return aac.(int) + 1
		}
		return aac
	}, 0)
}
