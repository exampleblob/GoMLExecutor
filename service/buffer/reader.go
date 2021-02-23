
package buffer

import (
	"bytes"
	"fmt"
	"io"
	"net/http/httputil"
)

//Read reads data with buffer Pool
func Read(pool httputil.BufferPool, reader io.Reader) ([]byte, int, error) {
	data := pool.Get()
	readTotal := 0
	offset := 0
	for i := 0; i < len(data); i++ {