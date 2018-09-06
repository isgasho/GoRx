package rx

//Subject 可以编程发送数据的Observable,input 为用于发送数据的外部数据源
func Subject(source Observable, input <-chan Any) Observable {
	var _next Next
	go func() {
		for d := range input {
			if isError(d, _next) {
				return
			}
			_next <- d
		}
		close(_next)
	}()
	return Share()(func(next Next, stop Stop) {
		_next = next
		if source != nil {
			source(next, stop)
		}
	})
}
