package Logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

//logging 是一个默认的日志对象，提供全局的Error, Info函数供使用，必须调用InitLogging
//函数进行初始化
var logging *Logger

var DEBUG = 0
var INFO = 3
var ERROR = 5
var LOGDir = "logs/"

//InitLogging 初始化默认的日志对象，初始化后，就能使用Error，Info函数记录日志
func InitLogging(inputfilename string, level int) {
	if !PathExists(LOGDir) {
		os.Mkdir(LOGDir, os.ModePerm)
	}
	logging = New(inputfilename, true, false,
		level, 3)
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//Error 默认日志对象方法，记录一条错误日志，需要先初始化
func Error(format string, v ...interface{}) {
	logging.Error(format, v...)
}

//Errorln 默认日志对象方法，记录一条消息日志，需要先初始化
func Errorln(args ...interface{}) {
	logging.Errorln(args...)
}

//Info 默认日志对象方法，记录一条消息日志，需要先初始化
func Info(format string, v ...interface{}) {
	logging.Info(format, v...)
}

//Infoln 默认日志对象方法，记录一条消息日志，需要先初始化
func Infoln(args ...interface{}) {
	logging.Infoln(args...)
}

//Debug 默认日志对象方法，记录一条消息日志，需要先初始化
func Debug(format string, v ...interface{}) {
	logging.Debug(format, v...)
}

//Debugln 默认日志对象方法，记录一条调试日志，需要先初始化
func Debugln(args ...interface{}) {
	logging.Debugln(args...)
}

type Logger struct {
	level         int // debug 0 info 3 err 5
	innerLogger   *log.Logger
	curFile       *os.File
	todaydate     string
	filename      string
	runtimeCaller int
	logFilePath   bool
	logFunc       bool
	msgQueue      chan string // 所有的日志先到这来
	closed        bool
}

//New 创建一个自己的日志对象。
// filename:在logs文件夹下创建的文件名
// logFilePath: 日志中记录文件路径
// logFunc: 日志中记录调用函数
// level: 打印等级。DEBUG, INFO, ERROR
// runtimeCaller: 文件路径深度，设定适当的值，否则文件路径不正确
func New(filename string, logFilePath bool,
	logFunc bool, level int, runtimeCaller int) *Logger {

	// result := newLogger(logFile, flag)
	result := new(Logger)
	result.msgQueue = make(chan string, 1000)
	result.closed = false

	var multi io.Writer

	if filename != "" {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		logFile, err := os.OpenFile(dir+"/"+LOGDir+filename,
			os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println(err.Error())
		}
		result.curFile = logFile

		fmt.Println("newLogger use MultiWriter")
		multi = io.MultiWriter(logFile, os.Stdout)
	} else {
		result.curFile = nil

		fmt.Println("newLogger use stdout")
		multi = os.Stdout
	}

	result.innerLogger = log.New(multi, "", 0)

	result.filename = filename

	result.runtimeCaller = runtimeCaller
	result.logFilePath = logFilePath
	result.logFunc = logFunc
	result.level = level
	result.todaydate = time.Now().Format("2006-01-02")

	// 启动日志切换
	go result.logworker()

	return result
}

// Close 关闭这一个日志对象
func (logobj *Logger) Close() error {
	logobj.closed = true
	return nil
}

func (logobj *Logger) getFormat(prefix, format string) string {
	var buf bytes.Buffer

	// 增加时间
	buf.WriteString(time.Now().Format("2006-01-02 15:04:05 "))

	buf.WriteString(prefix)

	// 增加文件和行号
	funcName, file, line, ok := runtime.Caller(logobj.runtimeCaller)
	if ok {
		if logobj.logFilePath {
			buf.WriteString(filepath.Base(file))
			buf.WriteString(":")
			buf.WriteString(strconv.Itoa(line))
			buf.WriteString(" ")
		}
		if logobj.logFunc {
			buf.WriteString(runtime.FuncForPC(funcName).Name())
			buf.WriteString(" ")
		}
		buf.WriteString(format)
		format = buf.String()
	}
	return format
}

//Error 记录一条错误日志
func (logobj *Logger) Error(format string, v ...interface{}) {
	if logging.level > 5 {
		return
	}

	format = logobj.getFormat("ERROR ", format)
	logobj.msgQueue <- fmt.Sprintf(format, v...)
}

//Errorln 打印一行错误日志
func (logobj *Logger) Errorln(args ...interface{}) {
	if logging.level > 5 {
		return
	}

	prefix := logobj.getFormat("ERROR ", "")
	logobj.msgQueue <- fmt.Sprintln(append([]interface{}{prefix}, args...)...)
}

//Info 记录一条消息日志
func (logobj *Logger) Info(format string, v ...interface{}) {
	if logging.level > 3 {
		return
	}

	format = logobj.getFormat("INFO ", format)
	logobj.msgQueue <- fmt.Sprintf(format, v...)
}

//Infoln 打印一行消息日志
func (logobj *Logger) Infoln(args ...interface{}) {
	if logging.level > 3 {
		return
	}

	prefix := logobj.getFormat("INFO ", "")
	logobj.msgQueue <- fmt.Sprintln(append([]interface{}{prefix}, args...)...)
}

//Debug 记录一条消息日志
func (logobj *Logger) Debug(format string, v ...interface{}) {
	if logging.level > 0 {
		return
	}

	format = logobj.getFormat("DEBUG ", format)
	logobj.msgQueue <- fmt.Sprintf(format, v...)
}

//Debugln 打印一行调试日志
func (logobj *Logger) Debugln(args ...interface{}) {
	if logging.level > 0 {
		return
	}

	prefix := logobj.getFormat("DEBUG ", "")
	logobj.msgQueue <- fmt.Sprintln(append([]interface{}{prefix}, args...)...)
}

func (logobj *Logger) logworker() {
	for logobj.closed == false {
		msg := <-logobj.msgQueue
		logobj.innerLogger.Println(msg)

		//跨日改时间，后台启动
		nowDate := time.Now().Format("2006-01-02")
		if nowDate != logobj.todaydate {
			logobj.Debug("doRotate run %v %v", nowDate, logging.todaydate)
			logobj.doRotate()
		}
	}
}

func (logobj *Logger) doRotate() {
	// 日志按天切换文件，日志对象记录了程序启动时的时间，当当前时间和程序启动的时间不一致
	// 则会启动到这个函数来改变文件
	// 首先关闭文件句柄，把当前日志改名为昨天，再创建新的文件句柄，将这个文件句柄赋值给log对象
	// 最后尝试删除5天前的日志
	fmt.Println("doRotate run")

	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Printf("doRotate %v", rec)
		}
	}()

	if logobj.curFile == nil {
		fmt.Println("doRotate curfile nil, return")
		return
	}

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	prefile := logobj.curFile

	_, err := prefile.Stat()
	if err == nil {
		filePath := dir + "/" + LOGDir + logobj.filename

		err := prefile.Close()
		fmt.Printf("doRotate close err %v", err)
		nowTime := time.Now()
		time1dAgo := nowTime.Add(-1 * time.Hour * 24)
		err = os.Rename(filePath, filePath+"."+time1dAgo.Format("2006-01-02"))
		fmt.Printf("doRotate rename err %v", err)
	}

	if logobj.filename != "" {
		nextfile, err := os.OpenFile(dir+"/"+LOGDir+logobj.filename,
			os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println(err.Error())
		}
		logobj.curFile = nextfile

		fmt.Println("newLogger use MultiWriter")
		multi := io.MultiWriter(nextfile, os.Stdout)
		logobj.innerLogger.SetOutput(multi)
	}

	fmt.Println("doRotate ending")

	// 更新标记，这个标记决定是否会启动文件切换
	nowDate := time.Now().Format("2006-01-02")
	logobj.todaydate = nowDate
	logobj.deleteHistory()
}

func (logobj *Logger) deleteHistory() {
	// 尝试删除7天前的日志
	fmt.Println("deleteHistory run")
	nowTime := time.Now()
	time7dAgo := nowTime.Add(-1 * time.Hour * 24 * 7)

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	filePath := dir + "/" + LOGDir + logobj.filename + "." + time7dAgo.Format("2006-01-02")

	_, err := os.Stat(filePath)
	if err == nil {
		os.Remove(filePath)
	}
}
