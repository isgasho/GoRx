package rx

//Subscribe 观察者
func Subscribe(n func(interface{}, Disposable), e func(error), c func()) func(Observable, ...Deliver) Disposable {
	return func(source Observable, sources ...Deliver) (dispose Disposable) {
		next := make(Next)
		stop := make(Stop)
		disposed := false
		dispose = func() {
			disposed = true
			close(stop)
		}
		go func() {
			for x := range next {
				if err, ok := x.(error); ok {
					if e != nil {
						e(err)
					}
				} else {
					n(x, dispose)
				}
				if disposed {
					return
				}
			}
			if c != nil {
				c()
			}
		}()
		go Pipe(source, sources...)(next, stop)
		return
	}
}

//ToChan 用channel方式订阅事件流
func ToChan(out Next) func(Observable, ...Deliver) Disposable {
	return func(source Observable, sources ...Deliver) Disposable {
		stop := make(Stop)
		go Pipe(source, sources...)(out, stop)
		return func() {
			close(stop)
		}
	}
}
