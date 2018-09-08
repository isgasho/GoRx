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
func (this *Observable) Pipe(cbs ...p.Deliver) *Observable {
	return &Observable{p.Pipe(this.source, cbs...)}
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

//Interval 
func Interval(period time.Duration) *Observable {
    return &Observable{p.Interval(period)}
}

//Timer 
func Timer(dueTime time.Duration, period time.Duration) *Observable {
    return &Observable{p.Timer(dueTime, period)}
}

//Share 
func (this *Observable) Share() *Observable {
    return &Observable{p.Share()(this.source)}
}

//StartWith 
func (this *Observable) StartWith(xs ...interface{}) *Observable {
    return &Observable{p.StartWith(xs...)(this.source)}
}

//Take 
func (this *Observable) Take(count int) *Observable {
    return &Observable{p.Take(count)(this.source)}
}

//TakeWhile 
func (this *Observable) TakeWhile(f func(interface{}) bool) *Observable {
    return &Observable{p.TakeWhile(f)(this.source)}
}

//TakeUntil 
func (this *Observable) TakeUntil(sSrc Observable, delivers ...p.Deliver) *Observable {
    return &Observable{p.TakeUntil(sSrc.source, delivers...)(this.source)}
}

//TakeLast 
func (this *Observable) TakeLast(count int) *Observable {
    return &Observable{p.TakeLast(count)(this.source)}
}

//Skip 
func (this *Observable) Skip(count int) *Observable {
    return &Observable{p.Skip(count)(this.source)}
}

//SkipWhile 
func (this *Observable) SkipWhile(f func(interface{}) bool) *Observable {
    return &Observable{p.SkipWhile(f)(this.source)}
}

//SkipUntil 
func (this *Observable) SkipUntil(sSrc Observable, delivers ...p.Deliver) *Observable {
    return &Observable{p.SkipUntil(sSrc.source, delivers...)(this.source)}
}

//Reduce 
func (this *Observable) Reduce(f func(interface{}, interface{}) interface{}, seed ...interface{}) *Observable {
    return &Observable{p.Reduce(f, seed...)(this.source)}
}

//Count 
func (this *Observable) Count(f func(interface{}) bool) *Observable {
    return &Observable{p.Count(f)(this.source)}
}

//Delay 
func (this *Observable) Delay(delay time.Duration) *Observable {
    return &Observable{p.Delay(delay)(this.source)}
}

//Scan 
func (this *Observable) Scan(f func(interface{}, interface{}) interface{}, seed ...interface{}) *Observable {
    return &Observable{p.Scan(f, seed...)(this.source)}
}

//Map 
func (this *Observable) Map(f func(interface{}) interface{}) *Observable {
    return &Observable{p.Map(f)(this.source)}
}

//MapTo 
func (this *Observable) MapTo(x interface{}) *Observable {
    return &Observable{p.MapTo(x)(this.source)}
}

//Pairwise 
func (this *Observable) Pairwise() *Observable {
    return &Observable{p.Pairwise()(this.source)}
}

//SwitchMap 
func (this *Observable) SwitchMap(f func(interface{}) Observable, combineResults func(interface{}, interface{}) interface{}) *Observable {
    return &Observable{p.SwitchMap(func(d interface{}) p.Observable {return f(d).source}, combineResults)(this.source)}
}

//SwitchMapTo 
func (this *Observable) SwitchMapTo(source Observable, combineResults func(interface{}, interface{}) interface{}) *Observable {
    return &Observable{p.SwitchMapTo(source.source, combineResults)(this.source)}
}

//BufferTime 
func (this *Observable) BufferTime(period time.Duration, maxBufferSize int) *Observable {
    return &Observable{p.BufferTime(period, maxBufferSize)(this.source)}
}

//Repeat 
func (this *Observable) Repeat(count int) *Observable {
    return &Observable{p.Repeat(count)(this.source)}
}
