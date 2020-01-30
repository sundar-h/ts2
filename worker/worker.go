package worker

import (
	"sync"

	"github.com/golang/glog"

	"github.com/could-be/tools/tls2mp4/models"
	"github.com/could-be/tools/tls2mp4/util"
)

type Jobs struct {
	VideoType   models.VideoType
	Vid         string
	TsFiles     []string
	OutDir      string
	VideoName   string
	EpisodeName string
}

type Worker struct {
	bufferChanWg *sync.WaitGroup
	jobs         chan *Jobs
	poolSize     int
	All          int
}

func NewWorker(bufferSize, jobSize int) *Worker {
	w := &Worker{
		bufferChanWg: &sync.WaitGroup{},
		jobs:         make(chan *Jobs, bufferSize),
		poolSize:     jobSize,
	}

	// 协程数
	w.bufferChanWg.Add(w.poolSize)
	for i := 0; i < w.poolSize; i++ {
		go w.Run()
	}

	return w
}

func (w *Worker) Run() {
	for {
		select {
		case job, ok := <-w.jobs:
			if ok {
				if err := util.Transform(job.VideoName, job.EpisodeName, job.VideoType, job.TsFiles, job.OutDir); err != nil {
					glog.Errorf("Transform %s failed: %v", job.Vid, err)
					continue
				}
				glog.Infof("<<%s>> %s.%s SUCCESS!!!", job.VideoName, job.EpisodeName, job.VideoType)
			} else {
				w.bufferChanWg.Done()
				return
			}
		}
	}
}

func (w *Worker) Add(job *Jobs) {
	w.All++
	w.jobs <- job
}

func (w *Worker) CloseWait() {
	close(w.jobs)
	// 等待所有任务队列消费完毕
	w.bufferChanWg.Wait()
}
