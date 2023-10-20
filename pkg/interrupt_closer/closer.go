package interrupt_closer

import (
	"io"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

var DefaultSignals = []os.Signal{os.Interrupt, syscall.SIGINT, syscall.SIGTERM}

type Closer struct {
	closers []io.Closer
	stop    chan os.Signal
}

func NewCloser(signals ...os.Signal) *Closer {
	closer := &Closer{closers: []io.Closer{}, stop: make(chan os.Signal, 1)}

	go func() {
		stop := make(chan os.Signal, 1)
		if len(signals) == 0 {
			signals = DefaultSignals
		}
		signal.Notify(stop, signals...)

		got := <-stop
		closer.stop <- got

		signal.Stop(stop)
		closer.Close()
	}()

	return closer
}

func (c *Closer) Close() error {
	for _, closer := range c.closers {
		log.Println("closing: ", reflect.TypeOf(closer).String())
		err := closer.Close()
		if err != nil {
			log.Println("closer error: ", err)
		}
	}

	return nil
}

func (c *Closer) Wait() {
	<-c.stop
}

func (c *Closer) Add(closer io.Closer) {
	c.closers = append([]io.Closer{closer}, c.closers...)
}
