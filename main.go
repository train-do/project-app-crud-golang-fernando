package main

import (
	"github.com/train-do/project-app-crud-golang-fernando/service"
	"github.com/train-do/project-app-crud-golang-fernando/utils"
	"github.com/train-do/project-app-crud-golang-fernando/view"
)

var IsLogin bool

func main() {
	App:
	for {
		view.TitleView("Centralization Bank App")
		view.MainView()
		input := view.FormInput("Pilih Menu")
		err := utils.ValidationInput(input, `^[012]$`, "Inputan hanya diisi angka 1, 2 dan 0")
		if err != nil {
			utils.ClearScreen()
			view.PrintError(err.Error())
			continue
		}
		switch input {
		case "1":
			utils.ClearScreen()
			service.Register()
		case "2":
			utils.ClearScreen()
			service.Login()
		case "0":
			break App
		}
	}
}
func init() {
	utils.ClearScreen()
	utils.DecodeUsers()
	utils.DecodeTransactions()
}