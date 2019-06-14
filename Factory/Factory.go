package Factory

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

//参考模型:工厂流水线->流水线员工->待加工产品
type Payload struct {
	Name     string
	IpMask   string
	CallFunc func(string, string)
}

func (p *Payload) Play() {
	log.Printf("%s任务完成。\n", p.Name)
}

//任务
type Job struct {
	Payload Payload
}

type Worker struct {
	WorkerId        string        //员工ID
	WorkerName      string        //员工名字
	Workbench       chan Job      //员工加工产品的工作台，即来即走(无缓冲)。
	GWorkbenchQueue chan chan Job //等待分配加工产品的员工工作台队列
	Finished        chan bool     //员工结束工作通知通道，无缓冲
}

// 新建一条工厂流水线
func NewWorker(WorkbenchQueue chan chan Job, Id, Name string) *Worker {
	log.Printf("新建流水线:%s \n", Id)
	return &Worker{
		WorkerId:        Id, //员工ID
		WorkerName:      Name,
		Workbench:       make(chan Job),  //员工加工产品的工作台，即来即走(无缓冲)。
		GWorkbenchQueue: WorkbenchQueue,  //等待分配加工产品的员工工作台队列
		Finished:        make(chan bool), //无缓冲
	}
}

// 工人开始工作
func (w *Worker) Start() {
	//开一个新的协程
	go func() {
		for {
			//将当前未分配待加工产品的工作台添加到工作台队列中
			w.GWorkbenchQueue <- w.Workbench
			log.Printf("把[%s:%s]的工作台添加到工作台队列中，当前工作台队列长度：%d\n", w.WorkerId, w.WorkerName, len(w.GWorkbenchQueue))
			select {
			//接收到了新的WorkerJob
			case wJob := <-w.Workbench:
				wJob.Payload.Play()
			case bFinished := <-w.Finished:
				if true == bFinished {
					log.Printf("%s-[%s] 结束工作！\n", w.WorkerId, w.WorkerName)
					return
				}
			}
		}
	}()
}

func (w *Worker) Stop() {
	//w.QuitChannel <- true
	go func() {
		w.Finished <- true
	}()
}

type Dispatcher struct {
	DispatcherId    string         //流水线ID
	MaxWorkers      int            //流水线上的员工(Worker)最大数量
	Workers         []*Worker      //流水线上所有员工(Worker)对象集合
	Closed          chan bool      //流水线工作状态通道
	EndDispatch     chan os.Signal //流水线停止工作信号
	GJobQueue       chan Job       //流水线上的所有代加工产品(Job)队列通道
	GWorkbenchQueue chan chan Job  //流水线上的所有操作台队列通道
}

func NewDispatcher(maxWorkers, maxQueue int) *Dispatcher {
	Closed := make(chan bool)
	EndDispatch := make(chan os.Signal)
	JobQueue := make(chan Job, maxQueue)
	WorkbenchQueue := make(chan chan Job, maxWorkers)
	signal.Notify(EndDispatch, syscall.SIGINT, syscall.SIGTERM)
	return &Dispatcher{
		DispatcherId:    "调度者",
		MaxWorkers:      maxWorkers,
		Closed:          Closed,
		EndDispatch:     EndDispatch,
		GJobQueue:       JobQueue,
		GWorkbenchQueue: WorkbenchQueue,
	}
}

func (d *Dispatcher) Run() {
	// 开始运行
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.GWorkbenchQueue, fmt.Sprintf("work-%s", strconv.Itoa(i)), "bison"+strconv.Itoa(i))
		d.Workers = append(d.Workers, worker)
		//开始工作
		worker.Start()
	}
	//监控
	go d.Dispatch()
}

func (d *Dispatcher) Dispatch() {
	//FLAG:
	for {
		select {
		case endDispatch := <-d.EndDispatch:
			log.Printf("流水线关闭命令[%v]已发出...\n", endDispatch)
			close(d.GJobQueue)
		case wJob, Ok := <-d.GJobQueue:
			if true == Ok {
				log.Println("从流水线获取一个待加工产品(Job)-", wJob.Payload.Name)
				go func(wJob Job) {
					//获取未分配待加工产品的工作台
					Workbench := <-d.GWorkbenchQueue
					//将待加工产品(Job)放入工作台进行加工
					Workbench <- wJob
				}(wJob)
			} else {
				for _, w := range d.Workers {
					w.Stop()
				}
				d.Closed <- true
				return
				///break FLAG
			}
		}
	}
}

type WorkFlow struct {
	GDispatch *Dispatcher
}

func (wf *WorkFlow) StartWorkFlow(maxWorkers, maxQueue int) {
	//初始化一个调度器(流水线)，并指定该流水线上的员工(Worker)和待加工产品(Job)的最大数量
	wf.GDispatch = NewDispatcher(maxWorkers, maxQueue)
	//启动流水线
	wf.GDispatch.Run()
}

func (wf *WorkFlow) AddJob(wJob Job) {
	//向流水线中放入待加工产品(Job)
	wf.GDispatch.GJobQueue <- wJob
}

func (wf *WorkFlow) CloseWorkFlow() {
	closed := <-wf.GDispatch.Closed
	if true == closed {
		log.Println("调度器(流水线)已关闭.")
	}
}
