package models

type VideoType int

const (
	MP4 = iota
	MP3
)

func (vt VideoType) String() string {
	switch vt {
	case MP4:
		return "mp4"
	case MP3:
		return "mp3"
	default:
		return "mp4"
	}
}

var (
	DefaultVideoPath  = "/Users/admin/Library/Containers/com.tencent.tenvideo/Data/Library/Application Support/Download/video"
	DefaultDB         = "/Users/admin/Library/Containers/com.tencent.tenvideo/Data/Library/Application Support/CoreData/Download/downloadTask.db"
	DefaultTargetType = MP4
)
