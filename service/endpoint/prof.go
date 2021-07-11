package endpoint

import (
	"net/http"
	"runtime/pprof"
	"sync"
)

const memProfURI = "/v1/api/debug/memprof"

cons