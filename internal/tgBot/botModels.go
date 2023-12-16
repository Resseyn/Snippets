package tgBot

var ShownSnippetMessages = make(map[int]int)
var CurrentShownLatestListID int

type JsonWithCommandAndData struct {
	Command string `json:"command"`
	ID      int    `json:"ID"`
}
