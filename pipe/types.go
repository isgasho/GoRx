package rx

type (
	//Any 任意数据
	Any interface{}
	//Disposable 终止函数
	Disposable func()
	//Next 发送数据的回调函数
	Next chan interface{}
	//Stop 取消订阅
	Stop chan bool
	//Observable 可观察对象，传入SubChan,给SubChan发送Next进行触发
	Observable func(Next, Stop)
	//Deliver 传递者
	Deliver func(source Observable) Observable
	//Observer 观察者
	Observer func(source Observable) Disposable
)
