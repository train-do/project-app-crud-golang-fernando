package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/train-do/project-app-crud-golang-fernando/model"
	"github.com/train-do/project-app-crud-golang-fernando/utils"
	"github.com/train-do/project-app-crud-golang-fernando/view"
)

func Register() {
	for{
		var err error
		view.TitleView("Register")
		var name, username, password string
		for{
			name = view.FormInput("Nama")
			err = utils.ValidationInput(name, `^[A-Za-z]*\s?[A-Za-z]*$|^[A-Za-z]+$|^\S$`, `Invalid input`)
			if err != nil {
				view.PrintError(err.Error())
				continue
			}
			break
		}
		for{
			username = view.FormInput("Username")
			err = utils.ValidationUnique(utils.UserToInterface(), "Username", username)
			if err != nil {
				view.PrintError(err.Error())
				continue
			}
			err = utils.ValidationInput(username, ``, ``)
			if err != nil {
				view.PrintError(err.Error())
				continue
			}
			break
		}
		for{
			password = view.FormInput("Password")
			err = utils.ValidationInput(password, ``, "")
			if err != nil {
				view.PrintError(err.Error())
				continue
			}
			break
		}
		noRek := utils.GenerateNoRek(utils.UserToInterface())
		var saldo int
		for{
			s := view.FormInput("Saldo")
			saldo, _ = strconv.Atoi(s)
			err = utils.ValidationInput(s, `^(0|[1-9][0-9]*)$`, "Inputan hanya menerima angka positif dan bulat")
			if err != nil {
				view.PrintError(err.Error())
				continue
			}
			break
		}
		bank := utils.GenerateBank()
		newUser := model.User{Name: name, Username: username, Password: password, NoRek: noRek, Saldo: saldo, Bank: bank}
		model.Users = append(model.Users, &newUser)
		utils.EncodeUsers()
		utils.ClearScreen()
		view.PrintSucces("Akun Berhasil dibuat")
		// fmt.Println(model.Users)
		return
	}
}
var session bool
func Login() {
	var username string
	for{
		view.TitleView("Login")
		username = view.FormInput("Username")
		_, err := utils.ValidationIsMatch(utils.UserToInterface(), "Username", username)
		password := view.FormInput("Password")
		_, err2 := utils.ValidationIsMatch(utils.UserToInterface(), "Password", password)
		if err != nil && err2 != nil {
			view.PrintError("Username dan Password Invalid")
			continue
		} else if err != nil {
			view.PrintError(err.Error())
			continue
		} else if err2 != nil {
			view.PrintError(err2.Error())
			continue
		}
		break
	}
	utils.ClearScreen()
	ctx, cancel := context.WithTimeout(context.Background(), 120 * time.Second)
	session = true
	fmt.Println(session)
	go func() {
		<-ctx.Done()
		session = false
	}()
	runningDashboard(ctx, cancel, username)
	defer cancel()
}
func runningDashboard(ctx context.Context, cancel context.CancelFunc, username string)  {
	for{
		if session {
			view.DashboardView()
			input := view.FormInput("Pilih Menu")
			err := utils.ValidationInput(input, `^[01234]$`, "Inputan hanya diisi angka 1 sampai 3 dan 0")
			if err != nil {
				utils.ClearScreen()
				view.PrintError(err.Error())
				continue
			}
			switch input {
			case "1":
				transfer(username, ctx)
			case "2":
				showTransaction(username, ctx)
			case "3":
				updatePassword(username, ctx)
			case "4":
				deleteAccount(username, ctx)
			case "0":
				utils.ClearScreen()
				cancel()
				utils.ClearScreen()
				fmt.Printf("Berhasil logout \033[31m%v\033[0m\n", ctx.Err())
				return
			}
		}else{
			fmt.Printf("Session habis, silahkan login kembali \033[31m%v\033[0m Dash\n", ctx.Err())
			return
		}
	}
}
func transfer(username string, ctx context.Context)  {
	var user, user2 model.User
	var nominal, iPengirim, iPenerima int
	for i, v := range model.Users {
		if v.Username == username {
			user = *v
			iPengirim = i
		}
	}
	for{
		penerima := view.FormInput("No. Rekening Penerima")
		if session {
			_, err := utils.ValidationIsMatch(utils.UserToInterface(), "NoRek", penerima)
			if err != nil {
				view.PrintError(err.Error())
				continue
			}
			for i, v := range model.Users {
				if v.NoRek == penerima {
					user2 = *v
					iPenerima = i
				}
			}
			break
		} else {
			fmt.Printf("Session habis, silahkan login kembali \033[31m%v\033[0m Transf\n", ctx.Err())
			return
		}
	}
	for{
		n := view.FormInput("Nominal Transfer")
		if session {
			err := utils.ValidationInput(n, `^([1-9][0-9]*)$`, "Inputan hanya menerima angka positif, bilangan bulat dan lebih dari 0")
			if err != nil {
				view.PrintError(err.Error())
				continue
			}
			nominal, _ = strconv.Atoi(n)
			if user.Saldo + user.Bank.FeeTransfer < nominal {
				view.PrintError("Saldo tidak cukup, saldo anda Rp. "+ fmt.Sprintf("%d", user2.Saldo))
				continue
			}
			break
		}else{
			fmt.Printf("Session habis, silahkan login kembali \033[31m%v\033[0m Transf\n", ctx.Err())
			return
		}
	}
	utils.HandleTransaction(iPengirim, iPenerima, nominal, user.Bank.FeeTransfer, user.Name, user2.Name, user.NoRek, user2.NoRek)
}
func showTransaction(username string, ctx context.Context)  {
	interfaceUser := make(map[string]interface{})
	var user model.User
	for _, v := range model.Users {
		if v.Username == username {
			user = *v
			interfaceUser["nameUser"] = v.Name
			interfaceUser["noRekUser"] = v.NoRek
			interfaceUser["saldoUser"] = v.Saldo
			interfaceUser["bank"] = v.Bank.Name
			break
		}
	}
	interfaceTemp := make(map[string]interface{})
	var arrInterface []interface{}
	for _, v := range model.Transactions {
		if v.Pengirim == user.NoRek {
			interfaceTemp = map[string]interface{}{"name" : v.NamaPenerima}
			interfaceTemp["noRek"] = v.Penerima
			interfaceTemp["nominal"] = v.Nominal
			interfaceTemp["typeTransc"] = v.Type +" Out"
			interfaceTemp["createdAt"] = v.CreatedAt
			arrInterface = append(arrInterface, interfaceTemp)
		}
		if v.Penerima == user.NoRek {
			interfaceTemp = map[string]interface{}{"name" : v.NamaPengirim}
			interfaceTemp["noRek"] = v.Pengirim
			interfaceTemp["nominal"] = v.Nominal
			interfaceTemp["typeTransc"] = v.Type +" In"
			interfaceTemp["createdAt"] = v.CreatedAt
			arrInterface = append(arrInterface, interfaceTemp)
		}
	}
	interfaceUser["historyTransaction"] = arrInterface
	dataJSON, err := json.MarshalIndent(interfaceUser, "", " ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	fmt.Println(string(dataJSON))
}
func updatePassword(username string, ctx context.Context)  {
	var newPassword string
	for{
		newPassword = view.FormInput("New Password")
		if session {
			err := utils.ValidationInput(newPassword, ``, "")
			if err != nil {
				view.PrintError(err.Error())
				continue
			}
			break
		}else{
			fmt.Printf("Session habis, silahkan login kembali \033[31m%v\033[0m UpdatePass\n", ctx.Err())
			return
		}
	}
	for{
		password := view.FormInput("Password")
		if session {
			_, err := utils.ValidationIsMatch(utils.UserToInterface(), "Password", password)
			if err != nil {
				view.PrintError(err.Error())
				continue
			}
			for _, v := range model.Users {
				if v.Username == username {
					(*v).Password = newPassword
					go utils.EncodeUsers()
					break
				}
			}
			break
		}else{
			fmt.Printf("Session habis, silahkan login kembali \033[31m%v\033[0m UpdatePass\n", ctx.Err())
			return
		}
	}
	view.PrintSucces("Berhasil Ganti Password")
}
func deleteAccount(username string, ctx context.Context)  {
	for i, v := range model.Users {
		if v.Username == username{
			if session {
				for{
					input := view.FormInput("Apakah anda yakin delete akun? (y/n)")
					utils.ClearScreen()
					if input == "y" {
						if v.Saldo > v.Bank.Fee {
							view.PrintWarning("Gagal delete, saldo melebihi batas minimum.\n Saldo anda Rp " + fmt.Sprintf("%d", v.Saldo))
							return
						}
					}else if input == "n"{
						return
					}else{
						view.PrintError(`Inputan hanya bisa "y" dan "n"`)
						continue
					}
					model.Users = append(model.Users[:i], model.Users[i+1:]...)
					go utils.EncodeUsers()
					break
				}
			}else{
				fmt.Printf("Session habis, silahkan login kembali \033[31m%v\033[0m Delete\n", ctx.Err())
				return
			}
		}
	}
	utils.ClearScreen()
	view.PrintSucces("Berhasil Delete Account")
}