package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

type Message struct {
	First   string `json:"first"`
	Second  string `json:"second"`
	Content string `json:"content"`
}

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
		teamIdFinish := teamId + "_finish"

		playerNum, err := tryToJoin(&m, teamId)
		if err != nil {
			so.Emit("error", err.Error())
			return
		}

		log.Printf("team %s, player N %d", teamId, playerNum)
		so.Join(teamId)

		so.Emit("connected", strconv.Itoa(playerNum))

		if playerNum == 2 {
			storage.Set(teamIdFinish, 0, cache.DefaultExpiration)
			nextGame(teamId, so)
		}

		so.On("event", func(message string) {
			println(message)
			so.BroadcastTo(teamId, "event", message)
		})

		so.On("finish", func(message string) {
			log.Println("signalled message", message)
			storage.Increment(teamIdFinish, 1)
			if v, ok := storage.Get(teamIdFinish); ok {
				if conv, converted := v.(int); converted {
					if conv >= 2 {
						nextGame(teamId, so)
					}
				} else {
					nextGame(teamId, so)
				}
			} else {
				nextGame(teamId, so)
			}
		})

		so.On("disconnection", func() {
			so.Leave(teamId)
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

func nextGame(teamId string, so socketio.Socket) {
	next := 0
	if val, ok := storage.Get(teamId + "_next"); ok {
		if conv, converted := val.(int); converted {
			next = conv
		}
	}

	storage.Set(teamId+"_next", next+1, cache.DefaultExpiration)

	var m Message
	switch next {
	case 0:
		m = game0()
	case 1:
		m = game1()
	default:
		so.Emit("finish")
		so.BroadcastTo(teamId, "finish")
		return
	}

	mess := encodeMessage(m)
	so.Emit("start", mess)
	so.BroadcastTo(teamId, "start", mess)
}

func encodeMessage(m Message) []byte {
	str, _ := json.Marshal(m)
	return str
}

type colorDef struct {
	Base   string
	First  string
	Second string
}

var colordefs []colorDef = []colorDef{
	{"#3498db", "#48a6ea", "#2890d2"},
	{"#f22613", "#ff4628", "#e30c07"},
}

func game1() Message {

	data := struct {
		BaseColor   string
		FirstColor  string
		SecondColor string
		Data        string
	}{}

	colorIndex := rand.Intn(len(colordefs))
	colordef := colordefs[colorIndex]
	firstMessage := "light"
	secondMessage := "dark"
	if rand.Intn(10) < 5 {
		firstMessage = "dark"
		secondMessage = "light"
		data.FirstColor = colordef.Second
		data.SecondColor = colordef.First
		data.BaseColor = colordef.Base
	} else {
		data.FirstColor = colordef.First
		data.SecondColor = colordef.Second
		data.BaseColor = colordef.Base
	}

	field := generateGame1Field(data.BaseColor, data.FirstColor, data.SecondColor)
	str, _ := json.Marshal(field)
	data.Data = string(str)

	content := parseTemplate(tmplt, "game_1", &data)
	return Message{
		First:   fmt.Sprintf("pick circles of %s shade", firstMessage),
		Second:  fmt.Sprintf("pick circles of %s shade", secondMessage),
		Content: content,
	}
}

func generateGame1Field(main, first, second string) [][]string {
	res := make([][]string, 6)
	for i := 0; i < 6; i++ {
		res[i] = make([]string, 7)
		for j := 0; j < 7; j++ {
			rnd := rand.Intn(100)
			if rnd < 20 {
				res[i][j] = first
			} else if rnd < 40 {
				res[i][j] = second
			} else {
				res[i][j] = main
			}
		}
	}

	return res
}

func game0() Message {

	data := struct {
		First  int
		Second int
		Data   string
	}{}

	nums := []int{6, 8, 9}
	shuffleInt(nums)
	data.First = nums[0]
	data.Second = nums[1]

	field := generateGame0Field(nums[2], nums[0], nums[1])
	str, _ := json.Marshal(field)
	data.Data = string(str)

	content := parseTemplate(tmplt, "game_0", &data)
	return Message{
		First:   "choose only numbers " + strconv.Itoa(data.First),
		Second:  "choose only numbers " + strconv.Itoa(data.Second),
		Content: content,
	}
}

func generateGame0Field(main, first, second int) [][]int {
	res := make([][]int, 8)
	for i := 0; i < 8; i++ {
		res[i] = make([]int, 25)
		for j := 0; j < 25; j++ {
			rnd := rand.Intn(100)
			if rnd < 10 {
				res[i][j] = first
			} else if rnd < 20 {
				res[i][j] = second
			} else {
				res[i][j] = main
			}
		}
	}

	return res
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
