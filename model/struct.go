package model

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	NoRek    string `json:"noRek"`
	Saldo    int    `json:"saldo"`
	Bank     Bank   `json:"bank"`
}
type Bank struct {
	Name        string `json:"name"`
	FeeTransfer int    `json:"feeTransfer"`
	Fee         int    `json:"fee"`
}
type Transaction struct {
	NamaPengirim string `json:"namaPengirim"`
	Pengirim     string `json:"noRekPengirim"`
	NamaPenerima string `json:"namaPenerima"`
	Penerima     string `json:"noRekPenerima"`
	Nominal      int    `json:"nominal"`
	Type         string `json:"type"`
	CreatedAt    string `json:"createdAt"`
}

var Banks = []Bank{
	{"BCA", 5000, 10000},
	{"Mandiri", 2000, 5000},
	{"CIMB", 3000, 15000},
}

var Transactions []Transaction
var Users []*User