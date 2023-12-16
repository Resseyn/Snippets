package handlers

import (
	"SnippetsTESTBYGUIDE/internal/loggers"
	"SnippetsTESTBYGUIDE/pkg/models"
	"SnippetsTESTBYGUIDE/pkg/templates"
	"github.com/gorilla/sessions"
	"net/http"
	"strconv"
)

// Создаем глобальную переменную для хранения сессий.
var store = sessions.NewCookieStore([]byte("6467098865:AAHByMBybrT_pFOjySUOg960m6YiW7D7B4Y"))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		loggers.Logger.Println(err.Error())
		panic(err)
	}
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	user := models.User{
		UserID: id, // ID Telegram пользователя
		ChatID: int64(id),
	}
	session.Values["user"] = user
	err = session.Save(r, w)
	if err != nil {
		loggers.Logger.Println(err)
		panic(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/loginpage" {
		http.NotFound(w, r)
		return
	}
	err := templates.TemplateCache["login.page.tmpl.html"].Execute(w, nil)
	if err != nil {
		loggers.Logger.Fatal(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
}
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session
	session, _ := store.Get(r, "user-session")

	// Clear session data
	session.Values["user"] = nil // Clear user ID or other user data
	session.Options.MaxAge = -1  // Set session to expire immediately

	// Save the session
	err := session.Save(r, w)
	if err != nil {
		// Handle error if session cannot be saved
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to the login or home page
	http.Redirect(w, r, "/loginpage", http.StatusSeeOther)
}
