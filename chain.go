package rx

import (
	"time"

	p "github.com/langhuihui/gorx/pipe"
)

//Observable emit data
type Observable struct {
	source p.Observable
}

func changeTop(sources []Observable) (pSources []p.Observable) {
	pSources = make([]p.Observable, len(sources))
	for i, source := range sources {
		pSources[i] = source.source
	}
	return
}

//Pipe
func (observable *Observable) Pipe(cbs ...p.Deliver) *Observable {
	return &Observable{p.Pipe(observable.source, cbs...)}
}

//Subject
func Subject(source Observable, input <-chan interface{}) *Observable {
	return &Observable{p.Subject(source.source, input)}
}

//If 
func If(f func() bool, trueOb Observable, falseOb Observable) *Observable {
    return &Observable{p.If(f, trueOb.source, falseOb.source)}
}

//Race 
func Race(sources ...Observable) *Observable {
    return &Observable{p.Race(changeTop(sources)...)}
}

//Concat 
func Concat(sources ...Observable) *Observable {
    return &Observable{p.Concat(changeTop(sources)...)}
}

//Merge 
func Merge(sources ...Observable) *Observable {
    return &Observable{p.Merge(changeTop(sources)...)}
}

//CombineLatest 
func CombineLatest(sources ...Observable) *Observable {
    return &Observable{p.CombineLatest(changeTop(sources)...)}
}

//IgnoreElements 
func IgnoreElements(source Observable) *Observable {
    return &Observable{p.IgnoreElements(source.source)}
}

//FromArray 
func FromArray(array []interface{}) *Observable {
    return &Observable{p.FromArray(array)}
}

//Of 
func Of(array ...interface{}) *Observable {
    return &Observable{p.Of(array...)}
}

//FromChan 
func FromChan(source chan interface{}) *Observable {
    return &Observable{p.FromChan(source)}
}

//From 
func From(source interface{}) *Observable {
    return &Observable{p.From(source)}
}

//Empty 
func Empty() *Observable {
    return &Observable{p.Empty()}
}

//Never 
func Never() *Observable {
    return &Observable{p.Never()}
}

//Throw 
func Throw(e error) *Observable {
    return &Observable{p.Throw(e)}
}

//Defer 
func Defer(f func() Observable) *Observable {
    return &Observable{p.Defer(func() p.Observable {return f().source})}
}

//Interval 
func Interval(period time.Duration) *Observable {
    return &Observable{p.Interval(period)}
}

//Timer 
func Timer(dueTime time.Duration, period time.Duration) *Observable {
    return &Observable{p.Timer(dueTime, period)}
}

//Share 
func (observable *Observable) Share() *Observable {
    return &Observable{p.Share()(observable.source)}
}

//StartWith 
func (observable *Observable) StartWith(xs ...interface{}) *Observable {
    return &Observable{p.StartWith(xs...)(observable.source)}
}

//Take 
func (observable *Observable) Take(count int) *Observable {
    return &Observable{p.Take(count)(observable.source)}
}

//TakeWhile 
func (observable *Observable) TakeWhile(f func(interface{}) bool) *Observable {
    return &Observable{p.TakeWhile(f)(observable.source)}
}

//TakeUntil 
func (observable *Observable) TakeUntil(sSrc Observable, delivers ...p.Deliver) *Observable {
    return &Observable{p.TakeUntil(sSrc.source, delivers...)(observable.source)}
}

//TakeLast 
func (observable *Observable) TakeLast(count int) *Observable {
    return &Observable{p.TakeLast(count)(observable.source)}
}

//Skip 
func (observable *Observable) Skip(count int) *Observable {
    return &Observable{p.Skip(count)(observable.source)}
}

//SkipWhile 
func (observable *Observable) SkipWhile(f func(interface{}) bool) *Observable {
    return &Observable{p.SkipWhile(f)(observable.source)}
}

//SkipUntil 
func (observable *Observable) SkipUntil(sSrc Observable, delivers ...p.Deliver) *Observable {
    return &Observable{p.SkipUntil(sSrc.source, delivers...)(observable.source)}
}

//Reduce 
func (observable *Observable) Reduce(f func(interface{}, interface{}) interface{}, seed ...interface{}) *Observable {
    return &Observable{p.Reduce(f, seed...)(observable.source)}
}

//Count 
func (observable *Observable) Count(f func(interface{}) bool) *Observable {
    return &Observable{p.Count(f)(observable.source)}
}

//Delay 
func (observable *Observable) Delay(delay time.Duration) *Observable {
    return &Observable{p.Delay(delay)(observable.source)}
}

//Scan 
func (observable *Observable) Scan(f func(interface{}, interface{}) interface{}, seed ...interface{}) *Observable {
    return &Observable{p.Scan(f, seed...)(observable.source)}
}

//Map 
func (observable *Observable) Map(f func(interface{}) interface{}) *Observable {
    return &Observable{p.Map(f)(observable.source)}
}

//MapTo 
func (observable *Observable) MapTo(x interface{}) *Observable {
    return &Observable{p.MapTo(x)(observable.source)}
}

//Pairwise 
func (observable *Observable) Pairwise() *Observable {
    return &Observable{p.Pairwise()(observable.source)}
}

//SwitchMap 
func (observable *Observable) SwitchMap(f func(interface{}) Observable, combineResults func(interface{}, interface{}) interface{}) *Observable {
    return &Observable{p.SwitchMap(func(arg0 interface{}) p.Observable {return f(arg0).source}, combineResults)(observable.source)}
}

//SwitchMapTo 
func (observable *Observable) SwitchMapTo(source Observable, combineResults func(interface{}, interface{}) interface{}) *Observable {
    return &Observable{p.SwitchMapTo(source.source, combineResults)(observable.source)}
}

//BufferTime 
func (observable *Observable) BufferTime(period time.Duration, maxBufferSize int) *Observable {
    return &Observable{p.BufferTime(period, maxBufferSize)(observable.source)}
}

//Repeat 
func (observable *Observable) Repeat(count int) *Observable {
    return &Observable{p.Repeat(count)(observable.source)}
}
