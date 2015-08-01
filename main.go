package main

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/googollee/go-socket.io"
	"github.com/pmylund/go-cache"
	"github.com/zenazn/goji"
	gojiweb "github.com/zenazn/goji/web"
)

var tmplt *template.Template
var storage *cache.Cache

func main() {
	rand.Seed(time.Now().UnixNano())
	loadTemplates()
	initSocketIO()
	storage = cache.New(time.Hour, 5*time.Minute)

	static := gojiweb.New()
	static.Get("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/static/", static)

	goji.Get("/", Index)
	goji.Get("/team", Team)
	goji.Get("/team/:team_id", Team)
	goji.Get("/solo", Solo)

	goji.Serve()
}

func tryToJoin(m *sync.Mutex, teamId string) (int, error) {
	m.Lock()
	defer m.Unlock()

	players := 0
	if val, ok := storage.Get(teamId); ok {
		if p, converted := val.(int); converted {
			players = p
		}
		if players >= 2 {
			return players, errors.New("team is full, try another one")
		}
	} else {
		return 0, errors.New("team is inactive, please create new one")
	}
	players++
	storage.Set(teamId, players, cache.DefaultExpiration)

	return players, nil
}

func initSocketIO() {
	sio, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	var m sync.Mutex
	sio.On("connection", func(so socketio.Socket) {
		teamId := so.Request().URL.Query().Get("team_id")
		if teamId == "" {
			so.Emit("error", "no team_id provided")
			return
		}

		playerNum, err := tryToJoin(&m, teamId)
		if err != nil {
			so.Emit("error", err.Error())
			return
		}

		log.Printf("team %s, player N %d", teamId, playerNum)
		so.Join("/" + teamId)

		so.Emit("connected", strconv.Itoa(playerNum))

		if playerNum == 2 {
			so.Emit("start", parseTemplate(tmplt, "game_0", nil))

			so.On("finish", func(message string) {
				log.Println("signalled message", message)
			})
		}

		so.On("disconnection", func() {
			so.Leave("/" + teamId)
			storage.Decrement(teamId, 1)
			log.Printf("player %d left team %s", playerNum, teamId)
		})
	})
	sio.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	// Sets up the handlers and listen on port 8080
	http.Handle("/socket.io/", sio)
}

func loadTemplates() {
	var templates []string

	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".html") {
			templates = append(templates, path)
		}
		return nil
	}

	err := filepath.Walk("./views/", fn)

	if err != nil {
		panic(err.Error())
	}

	tmplt = template.New("")
	tmplt.Delims("<<", ">>")
	tmplt = template.Must(tmplt.ParseFiles(templates...))
}

func parseTemplate(t *template.Template, name string, data interface{}) string {
	var doc bytes.Buffer
	err := t.ExecuteTemplate(&doc, name, data)
	if err != nil {
		log.Fatal(err.Error())
	}
	return doc.String()
}
