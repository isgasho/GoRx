package rx

import (
	"time"

	p "github.com/langhuihui/gorx/pipe"
)

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
func Subject(source p.Observable, input <-chan p.Any) *Observable {
	return &Observable{p.Subject(source, input)}
}

func From(source interface{}) *Observable {
	return &Observable{p.From(source)}
}
func Of(array ...interface{}) *Observable {
	return &Observable{p.Of(array...)}
}
func If(f func() bool, trueOb p.Observable, falseOb p.Observable) *Observable {
	return &Observable{p.If(f, trueOb, falseOb)}
}
func Race(sources ...Observable) *Observable {
	return &Observable{p.Race(changeTop(sources)...)}
}
func Concat(sources ...Observable) *Observable {
	return &Observable{p.Concat(changeTop(sources)...)}
}
func Merge(sources ...Observable) *Observable {
	return &Observable{p.Merge(changeTop(sources)...)}
}
func CombineLatest(sources ...Observable) *Observable {
	return &Observable{p.CombineLatest(changeTop(sources)...)}
}
func Interval(period time.Duration) *Observable {
	return &Observable{p.Interval(period)}
}
func (this *Observable) Delay(delay time.Duration) *Observable {
	this.source = p.Delay(delay)(this.source)
	return this
}
func (this *Observable) StartWith(xs ...p.Any) *Observable {
	return &Observable{p.StartWith(xs...)(this.source)}
}
func (this *Observable) Share() *Observable {
	this.source = p.Share()(this.source)
	return this
}
func (this *Observable) Pipe(cbs ...p.Deliver) *Observable {
	this.source = p.Pipe(this.source, cbs...)
	return this
}
func (this *Observable) Subscribe(n func(interface{}, func()), e func(error), c func()) func() {
	return p.Subscribe(n, e, c)(this.source)
}
func (this *Observable) ToChan(out p.Next) func() {
	return p.ToChan(out)(this.source)
}
func (this *Observable) Take(count int) *Observable {
	this.source = p.Take(count)(this.source)
	return this
}
func (this *Observable) Skip(count int) *Observable {
	this.source = p.Skip(count)(this.source)
	return this
}
func (this *Observable) TakeUntil(sSrc *Observable) *Observable {
	this.source = p.TakeUntil(sSrc.source)(this.source)
	return this
}
func (this *Observable) SkipUntil(sSrc *Observable) *Observable {
	this.source = p.SkipUntil(sSrc.source)(this.source)
	return this
}
