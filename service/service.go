package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/jinzhu/gorm"

	"github.com/could-be/tools/tls2mp4/models"
	"github.com/could-be/tools/tls2mp4/util"
	"github.com/could-be/tools/tls2mp4/worker"
)

type Service interface {
	Run()
}

type service struct {
	in         string
	out        string
	dbPath     string
	targetType models.VideoType
	db         *gorm.DB
	workerPool *worker.Worker
}

func New(in, out, dbPath string, targetType models.VideoType) Service {
	db, err := util.OpenDb(dbPath)
	util.Fatal(fmt.Sprintf("open db %s", dbPath), err)

	w := worker.NewWorker(100, 10)

	return &service{
		in:         in,
		out:        out,
		dbPath:     dbPath,
		targetType: targetType,

		db:         db,
		workerPool: w,
	}
}

func (s *service) Run() {

	start := time.Now()
	var err error

	switch s.in {
	default:
		glog.V(4).Infof("walk directory %s", s.in)
		s.walkVideos(s.in)

	case models.DefaultVideoPath:
		// 第一层, 获取所有视频目录
		tsDirs, err := ioutil.ReadDir(s.in)
		util.Fatal(fmt.Sprintf(`ioutil.ReadDir(%s)`, s.in), err)
		// 遍历视频目录, 合并视频下 ts 片段
		for _, tsDir := range tsDirs {
			if tsDir.IsDir() {
				glog.V(4).Infof("walk directory %s", tsDir.Name())
				s.walkVideos(filepath.Join(s.in, tsDir.Name()))
			}
		}
	}

	s.workerPool.CloseWait()
	util.Fatal(fmt.Sprintf("Walk(%s", s.in), err)
	glog.Infof("TRANSFORM %d FILES, TIME USED: %f", s.workerPool.All, time.Since(start).Seconds())
	if err = s.db.Close(); err != nil {
		glog.Error(err)
	}
}

// 获取第一个字段
func (s *service) walkVideos(tsDir string) {

	// d0023cbbqd0.320091.hls --> d0023cbbqd0
	vid := strings.Split(filepath.Base(tsDir), ".")[0]
	info, err := s.VideoInfo(vid)
	if err == gorm.ErrRecordNotFound {
		glog.Errorf("%s not found", vid)
		return
	}
	util.Fatal(fmt.Sprintf("video info of vid(%s) in(%s)", vid, tsDir), err)

	// 第二层, 获取视频所有 ts 文件
	tsFiles, err := filepath.Glob(filepath.Join(tsDir, "*/*.ts"))
	util.Fatal(`filepath.Glob("*/*.ts")`, err)

	// 所有视频ts 片段排序
	sort.Slice(tsFiles, func(i, j int) bool {
		return tsName(tsFiles[i]) <= tsName(tsFiles[j])
	})

	// 添加任务
	s.workerPool.Add(&worker.Jobs{
		VideoType:   s.targetType,
		Vid:         vid,
		TsFiles:     tsFiles,
		OutDir:      s.out,
		VideoName:   util.DoWithName(info.VideoName),
		EpisodeName: util.DoWithName(info.EpisodeName),
	})
}

// d0023cbbqd0.320091.hls_120_149/120.ts --> 120
func tsName(f string) string {
	if filepath.Ext(f) != ".ts" {
		panic(fmt.Sprintf("unknown extension %s", f))
	}

	base := filepath.Base(f)
	ext := filepath.Ext(f)
	return strings.TrimSuffix(base, ext)
}

func (s *service) VideoInfo(vid string) (*models.VideoInfo, error) {
	glog.V(4).Infof("getVideoInfo(%s) from db", vid)

	var record models.DownloadRecord
	if err := s.db.Table(record.TableName()).First(&record, "vid = ?", vid).Error; err != nil {
		return nil, err
	}

	glog.V(4).Infof("get record %s", record)

	if record.VideoInfo != nil {
		if record.VideoInfo.VideoName == "" {
			record.VideoInfo.VideoName = vid
		}

		if record.VideoInfo.EpisodeName == "" {
			record.VideoInfo.EpisodeName = vid
		}
		return record.VideoInfo, nil
	}

	return nil, errors.New("nil")
}
