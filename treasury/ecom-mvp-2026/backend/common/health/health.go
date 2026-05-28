package health

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

// Liveness returns basic service info and liveness status.
var osHostname = os.Hostname

func Liveness(version, commit string) gin.HandlerFunc {
	hostname, err := osHostname()
	if err != nil {
		hostname = fmt.Sprintf("unknown host err: %s", err.Error())
	}

	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"hostname": hostname,
			"version":  strings.ReplaceAll(version, "\n", ""),
			"commit":   commit,
		})
	}
}

// Readiness indicates the server is ready to receive traffic.
func Readiness() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

// Metrics returns a snapshot of runtime memory stats.
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		c.JSON(http.StatusOK, gin.H{
			"memory": gin.H{
				"alloc":        toMB(mem.Alloc),
				"totalAlloc":   toMB(mem.TotalAlloc),
				"sysAlloc":     toMB(mem.Sys),
				"heapInuse":    toMB(mem.HeapInuse),
				"heapIdle":     toMB(mem.HeapIdle),
				"heapReleased": toMB(mem.HeapReleased),
				"stackInuse":   toMB(mem.StackInuse),
				"stackSys":     toMB(mem.StackSys),
			},
		})
	}
}

func megabytes(b uint64) float64 {
	const mb = 1 << 20
	return float64(b) / float64(mb)
}

func toMB(b uint64) string {
	return fmt.Sprintf("%.2f MB", megabytes(b))
}
