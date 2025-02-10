package main

import (
	"diplom-chat-gost/internal/gui"
	"diplom-chat-gost/internal/model"
	"encoding/json"
	"errors"
	"fyne.io/fyne/v2/app"
	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
	"log"
	"os"
)

var config = model.Config{}

func main() {
	data, err := os.ReadFile("config/config.json")
	if err != nil {
		log.Println(errors.New("file open: " + err.Error()))
		return
	}

	if err = json.Unmarshal(data, &config); err != nil {
		log.Println(errors.New("unmarshal config: " + err.Error()))
		return
	}

	config.DB = reindexer.NewReindex(config.DBCredentials.Scheme + "://" + config.DBCredentials.Hostname + ":" +
		config.DBCredentials.Port + "/" + config.DBCredentials.Path)

	if err = config.DB.Status().Err; err != nil {
		log.Println(errors.New("reindexer connection: " + err.Error()))
		return
	}

	if err = config.DB.OpenNamespace("key", reindexer.DefaultNamespaceOptions(), model.KeyItem{}); err != nil {
		log.Println(errors.New("open namespace key: " + err.Error()))
		return
	}

	a := app.New()
	gui.SingInWindow(a, config)
}
