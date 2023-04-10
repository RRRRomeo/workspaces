package internal

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

var LOG_RELEASE = "0"

type LOG_LEVEL uint8
type TIME_SPEC uint8

// log level
const (
	LOG_LEVEL_DBG LOG_LEVEL = iota
	LOG_LEVEL_INF
	LOG_LEVEL_WRN
	LOG_LEVEL_ERR
	LOG_LEVEL_FTL
)

// time spec
const (
	TIME_MIN TIME_SPEC = iota
	TIME_MIL
	TIME_MIC
	TIME_NAN
)

// max outter name lens
const (
	MAX_OUTTER_NAME_LEN uint8 = 16
	MAX_OUTTER_PATH_LEN uint8 = 32
)

const (
	COLOR_RED = uint8(iota + 31)
	COLOR_GREEN
	COLOR_YELLOW
	COLOR_BLUE
	COLOR_MAGENTA
)

// type Logger interface {
// 	Typ()
// }

type outter struct {
	out_fd *os.File
	wg     *sync.WaitGroup
}

type Logger struct {
	sync.Mutex           // if set to async
	level      LOG_LEVEL // log level
	isout      uint8     // is outer to file?
	depth      uint8
	out_param  *outter
	std        *os.File
	// ...
}

func NewLogger(l LOG_LEVEL, depth uint8) *Logger {
	if l < LOG_LEVEL_DBG || l > LOG_LEVEL_FTL {
		return nil
	}

	return &Logger{
		level:     l,
		isout:     0,
		depth:     depth,
		out_param: &outter{nil, &sync.WaitGroup{}},
		std:       os.Stderr,
	}
}

func (l *Logger) SetOutter(isout uint8, out_ffp string) error {
	if isout > 1 {
		return errors.New("the outter param is invalid")
	}

	l.isout = isout
	// tmp_out_param := &outter{}
	fd, err := os.OpenFile(out_ffp, os.O_CREATE|os.O_WRONLY, 0751)
	if err != nil {
		return err
	}
	l.out_param.out_fd = fd
	// tmp_out_param.wg = wg
	return nil
}

func NewLoggerWithOutter(l LOG_LEVEL, isout uint8, out_ffp string) *Logger {
	logger := NewLogger(l, 3)
	err := logger.SetOutter(isout, out_ffp)
	if err != nil {
		return nil
	}
	return logger
}

func (l *Logger) Typ() string {
	return "logger"
}

func red(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", COLOR_RED, s)
}

func green(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", COLOR_GREEN, s)
}

func yellow(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", COLOR_YELLOW, s)
}

func blue(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", COLOR_BLUE, s)
}

func getLevelPrefix(l LOG_LEVEL) string {
	if l > LOG_LEVEL_FTL || l < LOG_LEVEL_DBG {
		return ""
	}

	switch l {
	case LOG_LEVEL_DBG:
		return blue("DBG")
	case LOG_LEVEL_INF:
		return green("INF")
	case LOG_LEVEL_WRN:
		return yellow("WRN")
	case LOG_LEVEL_ERR:
		return red("ERR")
	case LOG_LEVEL_FTL:
		return red("FTL")
	default:
		return ""
	}

}

func fmtOutputHeader(l *Logger, level LOG_LEVEL, s string) []byte {
	if s == "" {
		return nil
	}

	_, fn, ln, ok := runtime.Caller(int(l.depth))
	if !ok {
		fn = "???"
		ln = 0
	}

	pre_t := fmt.Sprintf("[%s] %s", getLevelPrefix(level), time.Now().Format("2006-01-02 15:04:05.000000000"))
	output := fmt.Sprintf("%s %s:%d %s", pre_t, fn, ln, s)
	return []byte(output)

}

func printf(l *Logger, level LOG_LEVEL, format string, v ...any) {
	if LOG_RELEASE == "1" {
		return
	}
	// ... format the output string
	output := fmtOutputHeader(l, level, fmt.Sprintf(format, v...))

	if l.isout == 1 {
		// ... output into file;
		l.out_param.wg.Add(1)
		go write(l.out_param.wg, l.out_param.out_fd, output) // ... async func need waitgroup it;
	}

	l.std.Write(output)
	l.out_param.wg.Wait()
}

func (l *Logger) Debugf(format string, v ...any) {
	l.Lock()
	defer l.Unlock()

	if l.level > LOG_LEVEL_DBG {
		return
	}
	printf(l, LOG_LEVEL_DBG, format, v...)
}

func (l *Logger) Infof(format string, v ...any) {
	l.Lock()
	defer l.Unlock()

	if l.level > LOG_LEVEL_INF {
		return
	}
	printf(l, LOG_LEVEL_INF, format, v...)
}

func (l *Logger) Warnf(format string, v ...any) {
	l.Lock()
	defer l.Unlock()

	if l.level > LOG_LEVEL_WRN {
		return
	}
	printf(l, LOG_LEVEL_WRN, format, v...)
}

func (l *Logger) Errf(format string, v ...any) {
	l.Lock()
	defer l.Unlock()

	if l.level > LOG_LEVEL_ERR {
		return
	}
	printf(l, LOG_LEVEL_ERR, format, v...)
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.Lock()
	defer l.Unlock()

	printf(l, LOG_LEVEL_FTL, format, v...)
	os.Exit(1)
}

func DelLogger(l *Logger) {
	l.isout = 0
	l.level = LOG_LEVEL_DBG
	l.out_param.out_fd.Close()
	l.out_param.wg = nil
	l.out_param = nil
	l.std = nil
	runtime.GC()
}

func write(wg *sync.WaitGroup, fd *os.File, buf []byte) (int16, error) {
	defer wg.Done()
	if fd == nil {
		return -1, errors.New("inner fd is nil")
	}
	// val := fmt.Sprintf(format, v...)
	n, e := fd.Write(buf)
	if e != nil {
		return -1, e
	}
	return int16(n), nil

	// ... this station needn't close the file fd; close into logger deinit;
}
