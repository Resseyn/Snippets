package main

import (
	. "SnippetsTESTBYGUIDE/internal"
	"SnippetsTESTBYGUIDE/internal/routes"
	"SnippetsTESTBYGUIDE/internal/tgBot"
	"SnippetsTESTBYGUIDE/pkg/models"
	"encoding/gob"
	"fmt"
	"net/http"
)

func main() {
	gob.Register(models.User{})
	fmt.Println("start - Start\ncreate - Create a new snippet\nshowLatest - Fully show 5 latest snippets")
	mux := routes.Route()
	go func() {
		err := http.ListenAndServe(":80", mux)
		if err != nil {
			App.Log.Fatal(err)
		}
	}()
	n := 1
	p := &n
	*p = 2
	tgBot.BotStart()
}
