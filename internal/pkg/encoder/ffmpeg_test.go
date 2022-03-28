package encoder

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var encoder *FFMpeg

func initTest(t *testing.T) {
	var err error
	encoder, err = NewFFMpeg()
	assert.Nil(t, err)
}

func TestFile_Mp3(t *testing.T) {
	initTest(t)
	var s []string
	encoder.convertFunc = func(cmd []string) error {
		s = cmd
		return nil
	}
	d, err := encoder.Convert("/dir/file.wav", "mp3", nil)
	assert.Equal(t, []string{"ffmpeg", "-i", "/dir/file.wav", "-q:a", "4", "/dir/file.wav.mp3"}, s)
	assert.Equal(t, "/dir/file.wav.mp3", d)
	assert.Nil(t, err)
}

func TestFile_M4a(t *testing.T) {
	initTest(t)
	var s []string
	encoder.convertFunc = func(cmd []string) error {
		s = cmd
		return nil
	}
	d, err := encoder.Convert("/dir/file.wav", "m4a", nil)
	assert.Equal(t, []string{"ffmpeg", "-i", "/dir/file.wav", "/dir/file.wav.m4a"}, s)
	assert.Equal(t, "/dir/file.wav.m4a", d)
	assert.Nil(t, err)
}

func TestFile_Fail(t *testing.T) {
	initTest(t)
	encoder.convertFunc = func(cmd []string) error {
		return errors.New("olia")
	}
	_, err := encoder.Convert("/dir/file.wav", "m4a", nil)
	assert.NotNil(t, err)
}

func TestMetadata(t *testing.T) {
	initTest(t)
	var s []string
	encoder.convertFunc = func(cmd []string) error {
		s = cmd
		return nil
	}
	encoder.Convert("file.wav", "m4a", []string{"a=aaaa"})
	assert.Equal(t, []string{"ffmpeg", "-i", "file.wav", "-metadata",
		"a=aaaa", "file.wav.m4a"}, s)
}

func TestMetadata_Several(t *testing.T) {
	initTest(t)
	var s []string
	encoder.convertFunc = func(cmd []string) error {
		s = cmd
		return nil
	}
	encoder.Convert("file.wav", "m4a", []string{"a=aaaa", "b=b and c"})
	assert.Equal(t, []string{"ffmpeg", "-i", "file.wav", "-metadata",
		"a=aaaa", "-metadata", "b=b and c", "file.wav.m4a"}, s)
}

func TestMetadata_Mp3(t *testing.T) {
	initTest(t)
	var s []string
	encoder.convertFunc = func(cmd []string) error {
		s = cmd
		return nil
	}
	encoder.Convert("file.wav", "mp3", []string{"a=aaaa"})
	assert.Equal(t, []string{"ffmpeg", "-i", "file.wav",
		"-q:a", "4",
		"-metadata", "a=aaaa", "file.wav.mp3"}, s)
}
