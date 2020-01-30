package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/golang/glog"

	"github.com/could-be/tools/tls2mp4/models"
)

func Transform(videoName, episodeName string, targetType models.VideoType, tsFiles []string, outDir string) error {
	out := filepath.Join(outDir, videoName)

	if err := os.MkdirAll(out, 0777); err != nil {
		return err
	}

	tmp, cmds := getCmd(out, targetType, tsFiles, episodeName)
	glog.V(4).Info("cmd: ", cmds)

	for _, cmd := range cmds {
		if err := RunCmd(cmd); err != nil {
			return err
		}
	}

	if err := os.RemoveAll(tmp); err != nil {
		glog.Error("os.RemoveAll(%s) failed %s", tmp, err)
	}
	return nil
}

// `cat 0.ts 1.ts > /tmp/all.ts`
// ffmpeg -i all.ts -bsf:a aac_adtstoasc -c copy -vcodec copy 1.mp4
func getCmd(out string, targetType models.VideoType, ts []string, episodeName string) (tmpFile string, cmds []string) {

	cmds = make([]string, 0, len(ts)+1)
	tmpFile = filepath.Join(os.TempDir(), episodeName+".ts")

	for _, v := range ts {
		// 校验文件是否存在
		if _, err := os.Stat(v); err == nil {
			cmds = append(cmds, fmt.Sprintf(`cat '%s' >> '%s'`, v, tmpFile))
		}
	}
	if len(cmds) > 0 {
		out = filepath.Join(out, fmt.Sprintf("%s.%s", episodeName, targetType))
		switch targetType {
		default:
			glog.Fatal("unknown type %s", targetType)
		case models.MP4:
			cmds = append(cmds, fmt.Sprintf(`ffmpeg -i '%s' -bsf:a aac_adtstoasc -c copy -vcodec copy '%s'`, tmpFile, out))
		case models.MP3:
			// ffmpeg -i infile.ts -f mp3 -acodec mp3 -aq 2 -vn outfile.mp3
			// cmds = append(cmds, fmt.Sprintf(`ffmpeg -i '%s' -f mp3 -acodec mp3 -aq 2 -vn '%s'`, tmpFile, out))
			cmds = append(cmds, fmt.Sprintf(`ffmpeg -i '%s' '%s'`, tmpFile, out))
		}
	}
	return tmpFile, cmds
}

func RunCmd(cmd string) error {
	command := exec.Command("/bin/sh", "-c", cmd)
	var out bytes.Buffer
	command.Stderr = &out
	if err := command.Run(); err != nil {
		return fmt.Errorf("%s failed: %v\n, details: %s", cmd, err, &out)
	}

	return nil
}

func Fatal(msg string, err interface{}) {

	switch errMsg := err.(type) {
	case error:
		if errMsg != nil {
			glog.FatalDepth(1, msg, errMsg)
		}
	case string:
		glog.FatalDepth(1, msg, errMsg)
	}

}

func DoWithName(s string) string {
	s = strings.ReplaceAll(s, `，`, "")
	s = strings.ReplaceAll(s, `！`, "")
	s = strings.ReplaceAll(s, ` `, "")
	s = strings.ReplaceAll(s, `\t`, "")
	s = strings.ReplaceAll(s, `“`, "")
	s = strings.ReplaceAll(s, `”`, "")
	s = strings.ReplaceAll(s, `\n`, "")
	s = strings.ReplaceAll(s, `\r`, "")
	s = strings.ReplaceAll(s, `\n\r`, "")
	s = strings.ReplaceAll(s, `【`, "")
	s = strings.ReplaceAll(s, `】`, "")
	s = strings.ReplaceAll(s, `␣`, "")
	s = strings.ReplaceAll(s, `/`, "")
	s = strings.ReplaceAll(s, `\`, "")
	s = strings.ReplaceAll(s, `《`, "")
	s = strings.ReplaceAll(s, `》`, "")
	s = strings.ReplaceAll(s, `（`, "")
	s = strings.ReplaceAll(s, `）`, "")
	s = strings.ReplaceAll(s, `␣␣`, "")

	return s
}
