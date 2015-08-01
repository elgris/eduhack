package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/pmylund/go-cache"
	"github.com/zenazn/goji/web"
)

type ErrTemplateData struct {
	Error string
}

type GameTemplateData struct {
	TeamID string
}

type TeamData struct {
}

func Index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, parseTemplate(tmplt, "index", nil))
}

func Team(c web.C, w http.ResponseWriter, r *http.Request) {
	teamId, ok := c.URLParams["team_id"]
	if !ok {
		newTeam(w, r)
		return
	}

	if _, ok := storage.Get(teamId); !ok {
		data := ErrTemplateData{
			Error: fmt.Sprintf("Team with id %s is not active. Create new one", teamId),
		}
		io.WriteString(w, parseTemplate(tmplt, "index", &data))
		return
	}

	data := GameTemplateData{
		TeamID: teamId,
	}
	io.WriteString(w, parseTemplate(tmplt, "game", &data))
}

func getTeam(teamId string) (team TeamData, ok bool) {
	return
}

func newTeam(w http.ResponseWriter, r *http.Request) {
	key := generateKey()
	storage.Set(key, "true", cache.DefaultExpiration)
	http.Redirect(w, r, "/team/"+key, 302)
}

func generateKey() string {
	rnd := rand.Int()
	h := sha1.New()
	h.Write([]byte(strconv.Itoa(rnd)))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func Solo(c web.C, w http.ResponseWriter, r *http.Request) {
	data := ErrTemplateData{
		Error: "solo game will be started",
	}
	io.WriteString(w, parseTemplate(tmplt, "index", &data))
}
