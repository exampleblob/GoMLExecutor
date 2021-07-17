package endpoint

import (
	"net/http"
	"runtime/pprof"
	"sync"
)

const memProfURI = "/v1/api/debug/memprof"

const cpuProfIndexURI = "/v1/api/debug/pprof/"

const cpuProfCmdlineURI = "/v1/api/debug/pprof/cmdline"
const cpuProfProfileURI = "/v1/api/debug/pprof/profile"
const cpuProfSymbolURI = "/v1/api/debu