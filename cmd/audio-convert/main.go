package main

import (
	"github.com/airenas/audio-convert/internal/pkg/encoder"
	"github.com/airenas/audio-convert/internal/pkg/file"
	"github.com/airenas/audio-convert/internal/pkg/service"
	"github.com/airenas/go-app/pkg/goapp"
	"github.com/pkg/errors"
	"github.com/labstack/gommon/color"
)

func main() {
	goapp.StartWithDefault()

	data := service.Data{}
	data.Port = goapp.Config.GetInt("port")

	var err error
	goapp.Log.Infof("Temp dir: %s", goapp.Config.GetString("tempDir"))
	data.Saver, err = file.NewSaver(goapp.Config.GetString("tempDir"))
	if err != nil {
		goapp.Log.Fatal(errors.Wrap(err, "can't init file saver"))
	}
	data.Coder, err = encoder.NewFFMpeg()
	if err != nil {
		goapp.Log.Fatal(errors.Wrap(err, "can't init ffmpeg wrapper"))
	}

	printBanner()

	err = service.StartWebServer(&data)
	if err != nil {
		goapp.Log.Fatal(errors.Wrap(err, "can't start the service"))
	}
}

var (
	version string
)

func printBanner() {
	banner := `
 _       _____ _    __     __       
| |     / /   | |  / /     \ \      
| | /| / / /| | | / /  _____\ \     
| |/ |/ / ___ | |/ /  /_____/ /     
|__/|__/_/  |_|___/        /_/      
								
       __  _______ _____       __   __  _____ __  ___ 
      /  |/  / __ \__  /     _/_/  /  |/  / // / /   |
     / /|_/ / /_/ //_ <    _/_/   / /|_/ / // /_/ /| |
    / /  / / ____/__/ /  _/_/    / /  / /__  __/ ___ |
   /_/  /_/_/   /____/  /_/     /_/  /_/  /_/ /_/  |_|  v: %s 

%s
________________________________________________________                                                 

`
	cl := color.New()
	cl.Printf(banner, cl.Red(version), cl.Green("https://github.com/airenas/audio-convert"))
}
