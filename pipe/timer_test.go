package rx

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {

	t.Run("interval", func(t *testing.T) {
		t.Parallel()
		getValue, current := 0, 0
		dispose := Subscribe(func(d interface{}, dispose func()) {
			getValue = d.(int)
		}, nil, nil)(Interval(1000))
		defer dispose()
		<-time.After(time.Second / 2)
		for current < 4 {
			<-time.After(time.Second)
			if getValue != current {
				t.Errorf("after %d sec we got %d", current, getValue)
				return
			}
			current++
		}
	})
	t.Run("delay", func(t *testing.T) {
		t.Parallel()
		getValue, current := 0, 0
		dispose := Subscribe(func(d interface{}, dispose func()) {
			getValue = d.(int)
		}, nil, nil)(Interval(1000), Delay(1000))
		defer dispose()
		<-time.After(time.Second + time.Second/2)
		for current < 4 {
			<-time.After(time.Second)
			if getValue != current {
				t.Errorf("after %d sec we got %d", current, getValue)
				return
			}
			current++
		}
	})

}
