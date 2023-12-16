package internal

import (
	myMySql2 "SnippetsTESTBYGUIDE/internal/database"
	"SnippetsTESTBYGUIDE/internal/loggers"
	"SnippetsTESTBYGUIDE/pkg/templates"
	"log"
)

type Application struct {
	Snippet myMySql2.SnippetModel
	User    myMySql2.UserModel
	Log     *log.Logger
}

var App *Application

func init() {
	loggers.InitLogger()
	err := myMySql2.InitMySqlORMdatabase()
	if err != nil {
		App.Log.Fatal(err)
	}
	templates.TemplateCache, err = templates.NewTemplateCache("/Users/romanovmaksim/GolandProjects/SnippetsTESTBYGUIDE/web/html/")
	if err != nil {
		App.Log.Fatal(err)
	}
	err = myMySql2.InitSnippetModel()
	if err != nil {
		App.Log.Fatal(err)
	}
	err = myMySql2.InitUserModel()
	if err != nil {
		App.Log.Fatal(err)
	}
	App = &Application{
		Snippet: myMySql2.Snippets,
		User:    myMySql2.Users,
		Log:     loggers.Logger,
	}
}
