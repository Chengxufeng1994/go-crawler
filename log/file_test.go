package log_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Chengxufeng1994/go-crawler/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestFile(t *testing.T) {
	const filePrefix = "test"
	const fileSuffix = ".log"
	const gzipSuffix = ".gz"

	p, c := log.NewFilePlugin(filePrefix+fileSuffix, zapcore.DebugLevel)
	logger := log.NewLogger(p)
	var b = make([]byte, 10000)
	var count = 10000
	for count > 0 {
		count--
		logger.Info(string(b))
	}
	err := c.Close()
	require.NoError(t, err)
	// NOTE: 目前Lumberjack的实现，close不会停掉压缩协程
	// 等待Lumberjack压缩日志文件
	time.Sleep(1 * time.Second)
	fs, err := os.ReadDir(".")
	require.NoError(t, err)
	var gzCount, logCount int
	for _, f := range fs {
		name := f.Name()
		if strings.HasPrefix(name, filePrefix) {
			if strings.HasSuffix(name, fileSuffix) {
				logCount++
				assert.NoError(t, os.Remove(f.Name()))
				continue
			}
			if strings.HasSuffix(name, gzipSuffix) {
				gzCount++
				logCount++
				assert.NoError(t, os.Remove(f.Name()))
				continue
			}
		}
	}

	require.Equal(t, 3, logCount)
	require.Equal(t, 2, gzCount)
}
