package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// grom 自定义类型, 需要实现这两个接口
// sql.Scanner 		db -> models
// driver.Valuer 	model -> db

type DownloadRecord struct {
	RecordId  string     `gorm:"record_id" `
	Vid       string     `gorm:"vid"`
	VideoInfo *VideoInfo `gorm:"type:text"`
}

func (r *DownloadRecord) String() string {
	v, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err.Error()
	}

	return string(v)
}

func (v DownloadRecord) TableName() string {
	return "download_record"
}

// 遇到问题, 自定义 Scan接口的时候, 及时按照驼峰形式默认 VideoName 的默认 json 是video_name
// 但是如果不写 json tag, 也会失败, 解码不出来
// 只有 VideoName
type VideoInfo struct {
	VideoName        string `json:"video_name"`   // 视频合辑名
	EpisodeName      string `json:"episode_name"` // 分集: 当前视频名
	ImagePath        string `json:"image_path"`
	ImageURL         string `json:"image_url"`
	DownloadPriority int    `json:"download_priority"`
	VideoIndex       int    `json:"video_index"`
	CoverimagePath   string `json:"coverimage_path"`
	VInfoKeyID       string `json:"vInfoKeyId"`
	DurationTime     int    `json:"duration_time"`
	IsMp4            bool   `json:"is_mp4"`
	WatchTime        int    `json:"watch_time"`
	IsHevc           bool   `json:"is_hevc"`
	CoverimageURL    string `json:"coverimage_url"`
	VideoType        int    `json:"video_type"`
	FormatID         string `json:"formatId"`
}

// func (vInfo *VideoInfo) UnmarshalJSON(data []byte) error {
// 	if vInfo == nil {
// 		return errors.New("nil receiver")
// 	}
//
// 	return json.Unmarshal(data, vInfo)
// }

// db -> model
func (vInfo *VideoInfo) Scan(value interface{}) error {
	return scan(vInfo, value)
}

// model -> db
func (vInfo *VideoInfo) Value() (driver.Value, error) {
	return value(vInfo)
}

func (vInfo *VideoInfo) String() string {
	if vInfo == nil {
		return "nil value"
	}

	bytes, err := json.MarshalIndent(vInfo, "", "  ")
	if err != nil {
		return err.Error()
	}

	return string(bytes)
}

// 抽象通用的函数
// db -> model
func scan(data, value interface{}) error {
	if data == nil {
		return errors.New("target is nil")
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, data)
	case string:
		return json.Unmarshal([]byte(v), data)
	default:
		return fmt.Errorf("invalid type, is %+v", value)
	}
}

// model -> db
func value(data interface{}) (interface{}, error) {
	vi := reflect.ValueOf(data)
	if vi.IsZero() {
		return nil, nil
	}

	return json.Marshal(data)
}
