package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/dirker/familie/famsrv/media"
	"github.com/go-martini/martini"
)

type configuration struct {
	MediaRoot string
}

func loadConfiguration() (config configuration, err error) {
	f, err := os.Open("config.json")
	if err != nil {
		return
	}

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&config)
	return
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var templateFuncs = template.FuncMap{
	"ftime": templateFuncFtime,
	"mod":   templateFuncMod,
}

func templateFuncFtime(s time.Time, f string) string {
	return s.Format(f)
}

func templateFuncMod(a, b int) int {
	return a % b
}

type siteData struct {
	MediaItems []media.Item
}

func main() {
	config, err := loadConfiguration()
	check(err)

	mediaRoot := os.ExpandEnv(config.MediaRoot)
	if mediaRoot == "" {
		panic("config: MediaRoot not specified")
	}

	media.SetRoot(mediaRoot)

	assetBox := rice.MustFindBox("assets")
	templateBox := rice.MustFindBox("templates")

	tmplMain := template.New("main")
	tmplMain.Funcs(templateFuncs)
	template.Must(tmplMain.Parse(templateBox.MustString("main.tmpl")))

	m := martini.Classic()
	m.Get("/", func(w http.ResponseWriter) {
		var err error

		mediaItems, err := media.GetItems()
		check(err)

		/* only serve images until 12 month back */
		date := mediaItems[0].CreatedAt.AddDate(-1, 0, 0)
		n := sort.Search(len(mediaItems), func(i int) bool {
			return mediaItems[i].CreatedAt.Before(date)
		})

		sd := siteData{}
		sd.MediaItems = mediaItems[:n]

		tmplMain.Execute(w, sd)
	})
	m.Get("/assets/**", func(params martini.Params, w http.ResponseWriter, r *http.Request) {
		path := params["_1"]

		f, err := assetBox.Open(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		ext := filepath.Ext(path)
		contentType := mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/octet"
		}
		w.Header().Set("Content-Type", contentType)

		io.Copy(w, f)
	})
	m.Get("/media/original/**", func(params martini.Params, w http.ResponseWriter, r *http.Request) {
		item, err := media.NewItem(params["_1"])
		if err != nil {
			http.NotFound(w, r)
			return
		}

		f, err := item.Open()
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()

		w.Header().Set("Content-Type", "image/jpeg")
		io.Copy(w, f)
	})
	m.Get("/media/thumb/**", func(params martini.Params, w http.ResponseWriter, r *http.Request) {
		f, err := media.OpenThumb(params["_1"])
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, r)
			return
		}
		defer f.Close()

		w.Header().Set("Content-Type", "image/jpeg")
		io.Copy(w, f)
	})

	m.RunOnAddr("localhost:8080")
}
