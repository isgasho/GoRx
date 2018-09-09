# GoRx
*ReactiveX* is a new, alternative way of asynchronous programming to callbacks, promises and deferred. It is about processing streams of events or items, with events being any occurrences or changes within the system.

[doc](https://godoc.org/github.com/langhuihui/GoRx/pipe)

# get GoRx
```bash
go get github.com/langhuihui/gorx
```

# build chain.go
```bash
node buildChain.js
```

# usage

GoRx's usage is very like RxJs. However implemention by Golang is very different with Javascript

## traditional usage
in tradition we call `Chain programming`
just like rxjs 5.0
```go
import "github.com/langhuihui/gorx"
rx.Interval(1000).SkipUntil(rx.Of(1).Delay(3000)).Subscribe(func(x interface{}, dispose func()) {
		fmt.Print(x)
	}, nil, nil)
```

## pipe usage

just like rxjs 6.0,but there are still some difference.

```go
import . "github.com/langhuihui/gorx/pipe"
Subscribe(func(x interface{}, dispose func()) {
		fmt.Print(x)
	}, nil, nil)(Interval(1000),SkipUntil(Of(1),Delay(3000)))
```

## Concept
An `Observable` is a synchronous stream of "emitted" values which can be either an empty interface{} or error. Below is how an `Observable` can be visualized:

```bash

                                time -->

(*)-------------(o)--------------(o)---------------(x)----------------|>
 |               |                |                 |                 |
Start          value            value             error              Done

```

## How it works

 - An `Observable` is a function looks like `func(Next,Stop)`.
 - `type Next chan interface{}` used to receive values.
 - `type Stop chan bool` used to unsubscribe streams.
 - Closing a `Next` means stream complete.
 - Calling an `Observable` means subscribe this Obserable.
 - Closing a `Stop` means unsubscribe a stream.

 ### example

 ```go
 func simpleObservable(n Next,s Stop){
	 n<-1//emit a value '1'
	 close(n)//complete
 }
 //subscribe simpleObservable
 next:=make(Next)
 stop:=make(Stop)
 go simpleObservable(next,stop)//now we subscribe
 for d:= range next{
	 //we will get values here
 }
 //we can do something here when complete
 ```

