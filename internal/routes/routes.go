package routes

import (
	"SnippetsTESTBYGUIDE/internal/handlers"
	"net/http"
)

func Route() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/snippet", handlers.ShowSnippet)
	mux.HandleFunc("/snippet/creation", handlers.SnippetCreation)
	mux.HandleFunc("/snippet/creation/create", handlers.CreateSnippet)
	//mux.HandleFunc("/snippet/update", handlers.CreateSnippet)
	mux.HandleFunc("/snippet/delete", handlers.DeleteSnippet)

	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/loginpage", handlers.LoginPage)
	mux.HandleFunc("/logout", handlers.LogoutHandler)

	fileServer := http.FileServer(http.Dir("./web/static/"))
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}

//======ТУТ ДОЛЖНА БЫТЬ ИЗОЛЯЦИЯ ДЛЯ ПАПКИ СТАТИК, НО Я ЕЕ НЕ ПОНЯЛ, СЛЕДОВАТЕЛЬНО НЕ ДОБАВИЛ
