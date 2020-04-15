//+build race

package bio

import (
	"bytes"
	"testing"
	"time"
)

func TestBio(t *testing.T) {
	buf := new(bytes.Buffer)
	log := NewLogger(nil, buf)
	sec := time.After(time.Second)
	for {
		select {
		case <-sec:
			log.Close()
			n := buf.Len()
			if n > 128 {
				n = 128
			}
			s := buf.String()
			for s[n] != '\n' {
				n++
			}
			t.Log(s[:n])
			return
		default:
		}
		log.Println("hi")
		log.Println("0123456789abcdefg")
		log.Println("abcdefghijklmopqrstuvwxyz")
		log.Println("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}
}
