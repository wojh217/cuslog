package cuslog

import (
	"bytes"
	"runtime"
	"strings"
	"time"
)

type Entry struct {
	logger *logger // 指向logger

	Buffer *bytes.Buffer // 写的实体

	Map map[string]interface{}

	Level Level     // 写的等级
	Time  time.Time // 时间
	File  string    // 文件名
	Line  int       // 行数
	Func  string    // 调用函数

	Format string
	Args   []interface{}
}

func entry(logger *logger) *Entry {
	return &Entry{
		logger: logger,
		Buffer: new(bytes.Buffer),
		Map:    make(map[string]interface{}, 5), // 创建map是为了输出json格式时保存字段用的，有必要么？
	}
}

// Write真正的输出日志，接受level和string信息
func (e *Entry) write(level Level, format string, args ...interface{}) {
	// 如果此级别小于设置的级别，则忽略输出
	if e.logger.opt.level > level {
		return
	}

	e.Time = time.Now()
	e.Level = level
	e.Format = format
	e.Args = args

	// 启动调用堆栈显示
	if !e.logger.opt.disableCaller {
		// skip为0表示当前函数即Write
		// skip为1表示往上一层，表示调用Write函数的函数，即func (l *logger) Debug Info Warn等函数
		// skip为2表示往上两层，即调用Debug的地方
		if pc, file, line, ok := runtime.Caller(2); !ok {
			e.File = "???"
			e.Func = "???"
		} else {
			e.File, e.Line, e.Func = file, line, runtime.FuncForPC(pc).Name()
			e.Func = e.Func[strings.LastIndex(e.Func, "/")+1:]
		}
	}

	// 将内容写入buffer中
	e.format()

	// 输出到io.Writer中，要使用锁？
	e.writer()

	// 输出后释放此Entry的buffer，且将e回收
	e.release()

}

func (e *Entry) format() {
	_ = e.logger.opt.formatter.Format(e)
}

func (e *Entry) writer() {
	e.logger.mu.Lock()
	_, _ = e.logger.opt.output.Write(e.Buffer.Bytes())
	e.logger.mu.Unlock()
}

func (e *Entry) release() {

	e.Buffer.Reset()
	e.Format, e.Line, e.File, e.Args, e.Func = "", 0, "", nil, ""

	// 将对象放入池中时，注意要将对象的属性清空，和新创建一个对象时的属性一致，不然可能会有问题
	// 原代码中map未做清空
	e.logger.entryPool.Put(e)
}
