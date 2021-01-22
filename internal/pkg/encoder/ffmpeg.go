package encoder

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

//FFMpeg wrapper for file conversion
type FFMpeg struct {
	convertFunc func([]string) error
}

//NewFFMpeg return new ffmpeg converter wrapper
func NewFFMpeg() (*FFMpeg, error) {
	res := FFMpeg{}
	res.convertFunc = runCmd
	return &res, nil
}

//Convert returns name of new converted file
func (e *FFMpeg) Convert(nameIn string, format string, metadata []string) (string, error) {
	resName := getNewFile(nameIn, format)
	params := []string{"ffmpeg", "-i", nameIn}
	params = append(params, getMetadataParams(metadata)...)
	params = append(params, resName)
	err := e.convertFunc(params)
	if err != nil {
		return "", err
	}
	return resName, nil
}

func runCmd(cmdArr []string) error {
	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
	var outputBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &outputBuffer
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "Output: "+string(outputBuffer.Bytes()))
	}
	return nil
}

func getNewFile(file string, format string) string {
	f := filepath.Base(file)
	d := filepath.Dir(file)
	return filepath.Join(d, fmt.Sprintf("%s.%s", f, format))
}

func getMetadataParams(prm []string) []string {
	res := []string{}
	for _, p := range prm {
		pt := strings.TrimSpace(p)
		if pt != "" {
			res = append(res, "-metadata")
			res = append(res, pt)
		}
	}
	return res
}
