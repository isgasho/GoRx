package rx

//Subject Represents an object that is both an observable sequence as well as an observer.
//Each notification is broadcasted to all subscribed observers.
func Subject(source Observable, input <-chan interface{}) Observable {
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
