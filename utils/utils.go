package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"github.com/train-do/project-app-crud-golang-fernando/model"
	"github.com/train-do/project-app-crud-golang-fernando/view"
)
func EncodeTransaction() {
	file, err := os.Create("Transaction.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	transactions := model.Transactions
	if err := encoder.Encode(transactions); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
}
func EncodeUsers() {
	file, err := os.Create("User.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	var users []model.User
	for _, v := range model.Users {
		users = append(users, *v)
	}
	if err := encoder.Encode(users); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
}
func DecodeUsers() {
	file, err := os.Open("User.json")
	if err != nil {
		// fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&model.Users); err != nil {
		// fmt.Println("Error decoding JSON:", err)
		return
	}
}
func DecodeTransactions() {
	file, err := os.Open("Transaction.json")
	if err != nil {
		// fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&model.Transactions); err != nil {
		// fmt.Println("Error decoding JSON:", err)
		return
	}
}
func GenerateNoRek(arr []interface{}) string {
	var noRek string
	for{
		noRek = strconv.Itoa(rand.Intn(99999999))
		err := ValidationUnique(arr, "NoRek", noRek)
		if err == nil {
			break
		}
	}
	return noRek
}
func GenerateBank() model.Bank {
	for{
		view.BanksView()
		bank := view.FormInput("Pilih Bank")
		v, _ := strconv.Atoi(bank)
		err := ValidationInput(bank, `^[123]$`, "Inputan hanya diisi angka 1,2 dan 3")
		if err == nil {
			return model.Banks[v-1]
		}
		view.PrintError(err.Error())
	}
}
func UserToInterface() []interface{} {
	var interfaceSlice []interface{}
	for _, bank := range model.Users {
	    interfaceSlice = append(interfaceSlice, *bank)
	}
	return interfaceSlice
}
func ValidationInput(input string, regex string, errMessage string) error {
	r := regexp.MustCompile(regex)
	if !r.MatchString(input) {
		return fmt.Errorf("%s", errMessage)
	}
	if input == "" || input == " "{
		return errors.New("Tidak boleh kosong")
	}
	return nil
}
func ValidationUnique(arr []interface{}, properti string, input string) error {
	for _, v := range arr {
		value, _ := ExtractValue(v, properti)
		if value == input {
			return fmt.Errorf("%s sudah ada", input)
		}
	}
	return nil
}
func ValidationIsMatch(arr []interface{}, properti string, input string) (string, error) {
	for _, v := range arr {
		value, _ := ExtractValue(v, properti)
		if value == input {
			return value, nil
		}
	}
	return "", fmt.Errorf("%s invalid", properti)
}
func HandleTransaction(iSend int, iReceive int, nominal int, fee int, namaPengirim string, namaPenerima string, rekPengirim string, rekPenerima string)  {
	for i, v := range model.Users {
		if i == iSend {
			(*v).Saldo -= (nominal + fee)
			break
		}
	}
	for i, v := range model.Users {
		if i == iReceive {
			(*v).Saldo += nominal
			break
		}
	}
	EncodeUsers()
	createdAt := time.Now().Format("2006-01-02T15:04:05.999999999Z07:00")
	model.Transactions = append(model.Transactions, model.Transaction{
		NamaPengirim: namaPengirim,
		Pengirim: rekPengirim,
		NamaPenerima: namaPenerima,
		Penerima: rekPenerima,
		Nominal: nominal,
		Type: "transfer",
		CreatedAt: createdAt,
	})
	EncodeTransaction()
}
func ExtractValue(obj interface{}, properti string) (string, error) {
    v := reflect.ValueOf(obj)
    field := v.FieldByName(properti)
    if !field.IsValid() {
        return "", fmt.Errorf("properti %s tidak ditemukan", properti)
    }
    return fmt.Sprintf("%v", field.Interface()), nil
}
func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}