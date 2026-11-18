package p2craft

import (
	"io"
	"sync"
)

type ReadWriteCloser interface {
	io.Reader
	io.Writer
	io.Closer
}

func Forward(a, b ReadWriteCloser) {
	var wg sync.WaitGroup
	wg.Add(2)

	copyStream := func(dst, src ReadWriteCloser) {
		defer wg.Done()

		_, _ = io.Copy(dst, src)

		// Propagate EOF to the other side.
		if cw, ok := dst.(interface {
			CloseWrite() error
		}); ok {
			_ = cw.CloseWrite()
		}
	}

	go copyStream(a, b)
	go copyStream(b, a)

	wg.Wait()
}
