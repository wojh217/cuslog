package cuslog

import (
	"fmt"
	"io"
	"os"
	"sync"
	"unsafe"
)

// std就是包内的Default
var std = New()

type logger struct {
	// logger对象的属性，包含输出位置、输出格式、输出等级等
	// 这里不直接写成字段形式，是为了方便做成选项模式？

	opt *options

	// 这个锁使用来并发控制设置opt的？
	mu sync.Mutex // 这里把锁放在looger中？还是entry中?

	entryPool *sync.Pool
}

// 根据选项创建logger
func New(opts ...Option) *logger {
	logger := &logger{
		opt: initOptions(opts...), // 根据选项，创建真正的opt
	}

	// 创建Pool对象，返回entry对象，后续logger
	logger.entryPool = &sync.Pool{
		New: func() interface{} { return entry(logger) },
	}
	return logger
}

func StdLogger() *logger {
	return std
}

func SetOptions(opts ...Option) {
	std.SetOptions(opts...)
}

func (l *logger) SetOptions(opts ...Option) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, opt := range opts {
		opt(l.opt)
	}
}

func Writer() io.Writer {
	return std
}

func (l *logger) Writer() io.Writer {
	return l
}

// logger.Write调用Entry.Write
// 哪里会直接调用logger.Write呢？
func (l *logger) Write(data []byte) (int, error) {
	l.entry().write(l.opt.stdLevel, FmtEmptySeparate, *(*string)(unsafe.Pointer(&data)))
	return 0, nil
}

// 从池中获取对象
func (l *logger) entry() *Entry {
	return l.entryPool.Get().(*Entry)
}

func (l *logger) Debug(args ...interface{}) {
	l.entry().write(DebugLevel, FmtEmptySeparate, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.entry().write(InfoLevel, FmtEmptySeparate, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.entry().write(WarnLevel, FmtEmptySeparate, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.entry().write(ErrorLevel, FmtEmptySeparate, args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.entry().write(PanicLevel, FmtEmptySeparate, args...)
	panic(fmt.Sprint(args...))
}

func (l *logger) Fatal(args ...interface{}) {
	l.entry().write(FatalLevel, FmtEmptySeparate, args...)
	os.Exit(1)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.entry().write(DebugLevel, format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.entry().write(InfoLevel, format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.entry().write(WarnLevel, format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.entry().write(ErrorLevel, format, args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.entry().write(PanicLevel, format, args...)
	panic(fmt.Sprintf(format, args...))
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.entry().write(FatalLevel, format, args...)
	os.Exit(1)
}

// std logger
func Debug(args ...interface{}) {
	std.entry().write(DebugLevel, FmtEmptySeparate, args...)
}

func Info(args ...interface{}) {
	std.entry().write(InfoLevel, FmtEmptySeparate, args...)
}

func Warn(args ...interface{}) {
	std.entry().write(WarnLevel, FmtEmptySeparate, args...)
}

func Error(args ...interface{}) {
	std.entry().write(ErrorLevel, FmtEmptySeparate, args...)
}

func Panic(args ...interface{}) {
	std.entry().write(PanicLevel, FmtEmptySeparate, args...)
	panic(fmt.Sprint(args...))
}

func Fatal(args ...interface{}) {
	std.entry().write(FatalLevel, FmtEmptySeparate, args...)
	os.Exit(1)
}

func Debugf(format string, args ...interface{}) {
	std.entry().write(DebugLevel, format, args...)
}

func Infof(format string, args ...interface{}) {
	std.entry().write(InfoLevel, format, args...)
}

func Warnf(format string, args ...interface{}) {
	std.entry().write(WarnLevel, format, args...)
}

func Errorf(format string, args ...interface{}) {
	std.entry().write(ErrorLevel, format, args...)
}

func Panicf(format string, args ...interface{}) {
	std.entry().write(PanicLevel, format, args...)
	panic(fmt.Sprintf(format, args...))
}

func Fatalf(format string, args ...interface{}) {
	std.entry().write(FatalLevel, format, args...)
	os.Exit(1)
}

func CurLevel() string {
	return LevelNameMapping[std.opt.level]
}
