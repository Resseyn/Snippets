package handlers

import (
	"SnippetsTESTBYGUIDE/internal/database"
	"SnippetsTESTBYGUIDE/internal/loggers"
	"SnippetsTESTBYGUIDE/pkg/models"
	"SnippetsTESTBYGUIDE/pkg/templates"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")
	user, ok := session.Values["user"].(models.User)
	if !ok {
		// Пользователь не авторизован, выполните соответствующие действия
		http.Redirect(w, r, "/loginpage", http.StatusSeeOther)
		return
	}

	//---------error by invalid PATH--------

	database.Snippets.DeleteExpired()
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	s, err := database.Snippets.Latest(0, 10, user.UserID)
	if err != nil {
		loggers.Logger.Println(err)
		return
	}

	data := templates.TemplateData{User: &user, Snippets: s}

	err = templates.TemplateCache["home.page.tmpl.html"].Execute(w, data)
	if err != nil {
		loggers.Logger.Fatal(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
}
