package service

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

type (
	// FileSaver saves the file with the provided name
	FileSaver interface {
		Save(name string, reader io.Reader) (string, error)
	}

	// Encoder encodes file, returns file name
	Encoder interface {
		Convert(nameIn string, format string, metadata []string) (string, error)
	}

	//Data is service operation data
	Data struct {
		Port int

		Saver FileSaver
		Coder Encoder

		readFunc func(string) ([]byte, error)
	}
)

//StartWebServer starts the HTTP service and listens for the convert requests
func StartWebServer(data *Data) error {
	goapp.Log.Infof("Starting HTTP audio convert service at %d", data.Port)
	portStr := strconv.Itoa(data.Port)
	data.readFunc = ioutil.ReadFile
	e := initRoutes(data)

	e.Server.Addr = ":" + portStr
	e.Server.IdleTimeout = 3 * time.Minute
	e.Server.ReadHeaderTimeout = 15 * time.Second
	e.Server.ReadTimeout = 180 * time.Second
	e.Server.WriteTimeout = 270 * time.Second

	w := goapp.Log.Writer()
	defer w.Close()
	gracehttp.SetLogger(log.New(w, "", 0))

	return gracehttp.Serve(e.Server)
}

func initRoutes(data *Data) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	e.POST("/convert", convert(data))
	e.GET("/live", live(data))

	goapp.Log.Info("Routes:")
	for _, r := range e.Routes() {
		goapp.Log.Infof("  %s %s", r.Method, r.Path)
	}
	return e
}

type input struct {
	Data     string   `json:"audio"`
	Format   string   `json:"format"`
	Metadata []string `json:"metadata"`
}

type output struct {
	Data string `json:"audio"`
}

func convert(data *Data) func(echo.Context) error {
	return func(c echo.Context) error {
		defer goapp.Estimate("Service method")()
		r := new(input)
		if err := c.Bind(r); err != nil {
			goapp.Log.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "Can get data")
		}

		if err := validate(r); err != nil {
			goapp.Log.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		id := uuid.New().String()
		fileName := id + ".wav"

		est := goapp.Estimate("Saving")
		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(r.Data))
		fileNameIn, err := data.Saver.Save(fileName, reader)
		if err != nil {
			goapp.Log.Error(err)
			return errors.Wrap(err, "Can not save file")
		}
		defer deleteFile(fileNameIn)
		est()

		est = goapp.Estimate("Convert")
		fileNameOut, err := data.Coder.Convert(fileNameIn, getFormat(r.Format), r.Metadata)
		if err != nil {
			goapp.Log.Error(err)
			return errors.Wrap(err, "Can not encode file")
		}
		defer deleteFile(fileNameOut)
		est()

		est = goapp.Estimate("Read")
		fd, err := data.readFunc(fileNameOut)
		if err != nil {
			goapp.Log.Error(err)
			return errors.Wrap(err, "Can not read file")
		}
		est()

		est = goapp.Estimate("Decode")
		res := &output{}
		res.Data = base64.StdEncoding.EncodeToString(fd)
		est()

		return c.JSON(http.StatusOK, res)
	}
}

func deleteFile(file string) {
	os.RemoveAll(file)
}

func validate(r *input) error {
	if !(getFormat(r.Format) == "mp3" || getFormat(r.Format) == "m4a") {
		return errors.Errorf("Unsuported format '%s'", r.Format)
	}
	if r.Data == "" {
		return errors.Errorf("No Audio")
	}
	return nil
}

func getFormat(f string) string {
	if f == "" {
		return "mp3"
	}
	return f
}

func live(data *Data) func(echo.Context) error {
	return func(c echo.Context) error {
		return c.JSONBlob(http.StatusOK, []byte(`{"service":"OK"}`))
	}
}
