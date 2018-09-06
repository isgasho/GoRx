package rx

func isError(d Any, n Next) (ok bool) {
	if _, ok = d.(error); ok {
		n <- d
		close(n)
	}
	return
}

func _deliver(nn func(Any, Next, Stop) bool, source Observable, next Next, stop Stop) {
	sNext := make(Next)
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
			} else if !nn(d, next, stop) {
				close(stop)
				close(next)
				return
			}
		case <-stop:
			return
		default:

		}
	}
}
func deliver(nn func(Any, Next, Stop) bool) Deliver {
	return func(source Observable) Observable {
		return func(next Next, stop Stop) {
			_deliver(nn, source, next, stop)
		}
	}
}

//Pipe 管道操作，可以将一组操作传入管道中，最后一个参数如果是Observer类型的参数的话就会激活事件流
func Pipe(source Observable, cbs ...Deliver) Observable {
	for _, cb := range cbs {
		source = cb(source)
	}
	return source
}
