package handlers

import (
	"SnippetsTESTBYGUIDE/internal/database"
	"SnippetsTESTBYGUIDE/internal/loggers"
	"SnippetsTESTBYGUIDE/pkg/models"
	"SnippetsTESTBYGUIDE/pkg/templates"
	"errors"
	"net/http"
	"strconv"
)

func ShowSnippet(w http.ResponseWriter, r *http.Request) {
	//---------error by invalid METHOD--------

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Метод запрещен!", 405)
		return
	}
	//----------error by invalid ID-----------

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		loggers.Logger.Println(err)
		http.NotFound(w, r)
		return
	}

	//==========BODY===========
	foundSnippet, err := database.Snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecordedSnippet) {
			w.Write([]byte("Запись не найдена("))
		} else {
			loggers.Logger.Println(err)
		}
		return
	}
	//files := []string{
	//	"./web/html/show.page.tmpl.html",
	//	"./web/html/base.layout.tmpl.html",
	//	"./web/html/basement.partial.tmpl.html",
	//}
	//
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	loggers.Logger.Println(err)
	//	return
	//}
	data := templates.TemplateData{Snippet: foundSnippet}
	err = templates.TemplateCache["show.page.tmpl.html"].Execute(w, data)
	if err != nil {
		loggers.Logger.Println(err)
		return
	}

}
func SnippetCreation(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")
	_, ok := session.Values["user"].(models.User)
	if !ok {
		// Пользователь не авторизован, выполните соответствующие действия
		http.Redirect(w, r, "/loginpage", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		w.Write([]byte("POST-Метод запрещен!"))
		return
	}
	err := templates.TemplateCache["creation.page.tmpl.html"].Execute(w, nil)
	if err != nil {
		loggers.Logger.Println(err)
		return
	}
}
func CreateSnippet(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")
	user, ok := session.Values["user"].(models.User)
	if !ok {
		// Пользователь не авторизован, выполните соответствующие действия
		http.Redirect(w, r, "/loginpage", http.StatusSeeOther)
		return
	}
	//---------error by invalid METHOD--------

	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		w.Write([]byte("GET-Метод запрещен!"))
		return
	}
	err := r.ParseForm()
	if err != nil {
		loggers.Logger.Println(err)
		return
	}
	form := r.Form
	database.Snippets.Insert(user.UserID, form.Get("title"),
		form.Get("content"),
		form.Get("expires"))

	http.Redirect(w, r, "http://127.0.0.1/", http.StatusSeeOther)
}

func DeleteSnippet(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodDelete {
	//	w.WriteHeader(405)
	//	w.Write([]byte("GET-Метод запрещен!"))
	//	return
	//} //Я ХУЙ ЗНАЕТ КАК УКАЗАТ В ШТМЛ МЕТОД ТАК ЧТО ПРОСТО ПО ГЕТ ОН УДАЛИТЬ ЧТО УГОДНО НАХУЙ
	idSt := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idSt)

	if err != nil {
		loggers.Logger.Println(err)
		return
	}
	database.Snippets.Delete(id)
	http.Redirect(w, r, "http://127.0.0.1/", http.StatusSeeOther)

}

//func UpdateSnippet(w http.ResponseWriter, r *http.Request)
//
//
