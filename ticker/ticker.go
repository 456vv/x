package ticker
	
import(
	"time"
)

type Ticker struct{
	C		chan time.Time
	t		*time.Ticker
	done	chan bool
	f		func()
}
func NewTicker(d time.Duration) *Ticker {
	ticker := &Ticker{
		C		: make(chan time.Time),
		t		: time.NewTicker(d),
		done	: make(chan bool, 1),
	}
	
	go func(ticker *Ticker){
		for {
			select {
			case c := <-ticker.t.C:
				if ticker.f != nil {
					ticker.f()
					continue
				}
				ticker.C<-c
			case <-ticker.done:
				return
			}
		}
	}(ticker)
	return ticker
}

func (T *Ticker) Func(f func()){
	T.f = f
}
func (T *Ticker) Reset(d time.Duration) {
	T.t.Reset(d)
}
func (T *Ticker) Stop() {
	T.done<-true
	T.t.Stop()
}