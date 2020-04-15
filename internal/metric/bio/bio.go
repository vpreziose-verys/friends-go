package bio

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(defaultLogger)
}

var defaultLogger = NewLogger(nil, os.Stderr)

func Printf(format string, v ...interface{}) { defaultLogger.Printf(format, v...) }
func Println(v ...interface{})               { defaultLogger.Println(v...) }
func Print(v ...interface{})                 { defaultLogger.Print(v...) }

// Logger reads events and outputs them in a line-based JSON format.
type Logger struct {
	src  chan interface{}
	dst  io.Writer
	done chan bool
	once chan bool
	halt chan bool
	_    [64]byte
}

// NewLog returns a Log that reads from src and writes the result
// to dst. If src is nil, it creates its own channel accessible by Log.C.
func NewLogger(src chan interface{}, dst io.Writer) *Logger {
	if src == nil {
		src = make(chan interface{}, iounit)
	}
	l := Logger{
		src:  src,
		dst:  dst,
		done: make(chan bool),
		halt: make(chan bool),
		once: make(chan bool, 1),
	}
	l.once <- true
	go l.start()
	return &l
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.src <- fmt.Sprintf(format, v...)
}
func (l *Logger) Println(v ...interface{}) {
	l.src <- fmt.Sprint(v...)
}
func (l *Logger) Print(v ...interface{}) {
	l.src <- fmt.Sprint(v...)
}

func (l *Logger) Write(p []byte) (int, error) {
	l.src <- string(p)
	return len(p), nil
}

// C returns the ingress channel. Send events to this and
// they will be logged.
func (l *Logger) C() chan<- interface{} {
	return l.src
}

// Close closes the logger
func (l *Logger) Close() error {
	select {
	case first := <-l.once:
		if first {
			close(l.once)
			close(l.done)
		}
	default:
	}
	<-l.halt
	return nil
}

const (
	iounit  = 65536
	lowater = 16
	hiwater = iounit - iounit/3
)

func (l *Logger) start() {
	buf := bufio.NewWriterSize(l.dst, iounit)
	defer close(l.halt)
	defer buf.Flush()

	em := emitter{dst: buf}
Loop:
	for {
		select {
		case <-l.done:
			return
		case e, more := <-l.src:
			if !more {
				return
			}
			em.nq = len(l.src)

			const (
				s = int64(time.Second)
			)
			ns := time.Now().UnixNano()
			hi, lo := ns/s, ns%s

			em.prefix = fmt.Sprintf(`{"level":"debug","ts":%d.%d,`, hi, lo)
			em.emit(e)
			for ; em.nq != 0; em.nq-- {
				select {
				case e = <-l.src:
					em.emit(e)
				default:
					continue Loop
				}
			}
		}
	}
}

type emitter struct {
	dst    *bufio.Writer
	prefix string
	e      string
	nq     int
}

func (z *emitter) emit(ev interface{}) {
	var e string
	switch ev := ev.(type) {
	case fmt.Stringer:
		e = ev.String()
	case string:
		e = ev
	case interface{}:
		e = fmt.Sprint(e)
	}
	if e == "" {
		return
	}
	if i := len(e) - 1; e[i] == '\n' {
		e = e[:i]
	}
	if e[0] != '"' {
		e = `"msg": ` + strconv.Quote(e)
	}
	z.dst.WriteString(z.prefix + e + "}\n")
	if z.nq < lowater || z.dst.Buffered() > hiwater {
		z.dst.Flush()
	}
}
