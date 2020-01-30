package main

import (
	"flag"

	"github.com/golang/glog"

	"github.com/could-be/tools/tls2mp4/models"
	"github.com/could-be/tools/tls2mp4/service"
)

var (
	inDir      = flag.String("in", models.DefaultVideoPath, "-in=.")
	outDir     = flag.String("out", "./out", "-out=./out")
	dbPath     = flag.String("db", models.DefaultDB, "-db=./something.db")
	targetType = flag.Int("t", models.DefaultTargetType, "-t=1 [0:mp4 1:mp3]")
)

func main() {
	flag.Parse()
	defer glog.Flush()

	s := service.New(*inDir, *outDir, *dbPath, models.VideoType(*targetType))
	s.Run()
}
