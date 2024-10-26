package view

import (
	"fmt"

	"github.com/train-do/project-app-crud-golang-fernando/model"
)

func TitleView(title string) {
	fmt.Println("----", title, "----")
}
func MainView() {
	fmt.Println("1. Register")
	fmt.Println("2. Login")
	fmt.Println("0. Exit")
}
func DashboardView() {
	fmt.Println("1. Transfer")
	fmt.Println("2. Show Transactions")
	fmt.Println("3. Ganti Password")
	fmt.Println("4. Delete Account")
	fmt.Println("0. Logout")
}
func BanksView() {
	fmt.Println()
	for i, v := range model.Banks {
		fmt.Printf("%d. %s\n", i+1, v.Name)
	}
}
func FormInput(label string) string {
	var input string
	fmt.Print(label, " : ")
	fmt.Scanln(&input)
	return input
}
func PrintSucces(text string)  {
	fmt.Printf("\033[1;32m %s \033[0m\n", text)
}
func PrintWarning(text string)  {
	fmt.Printf("\033[33m %s \033[0m\n", text)
}
func PrintError(text string)  {
	fmt.Printf("\033[31m %s \033[0m\n", text)
}