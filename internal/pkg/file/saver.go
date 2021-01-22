package file

import (
	"io"
	"os"
	"path"

	"github.com/pkg/errors"
)

//Saver saves file
type Saver struct {
	dir string
}

//NewSaver creates temporary file saver dir
func NewSaver(dir string) (*Saver, error) {
	res := Saver{}
	res.dir = dir
	if dir == "" {
		return nil, errors.New("No temp dir")
	}
	err := os.MkdirAll(dir, 0700)
	return &res, err
}

//Save saves file to temp dir
func (s *Saver) Save(name string, reader io.Reader) (string, error) {
	fn := path.Join(s.dir, name)

	file, err := os.Create(fn)
	if err != nil {
		return "", errors.Wrap(err, "Can't create file "+fn)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", errors.Wrap(err, "Can't write file "+fn)
	}

	return fn, nil
}
