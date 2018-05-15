package iif

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T) {
	data := []DataLine{
		AccntData{
			Name:      "Accounts",
			AccntType: "Payable",
			Desc:      "AP",
			Accnum:    "2000",
		},
		VendData{
			Name:      "Vendor",
			Refnum:    "1",
			PrintAs:   "",
			Addr1:     "Jon Vendor",
			Addr2:     "555",
			Addr3:     "Street St",
			Addr4:     `"Anywhere, AZ 85730"`,
			Addr5:     "USA",
			Cont1:     "Jon Vendor",
			Phone1:    "5555555555",
			FirstName: "Jon",
			LastName:  "Vendor",
		},
	}
	err := Export(data)

	assert.NoError(t, err)
}

func TestGetHeader(t *testing.T) {
	vend := VendData{
		Name:      "Vendor",
		Refnum:    "1",
		PrintAs:   "",
		Addr1:     "Jon Vendor",
		Addr2:     "555",
		Addr3:     "Street St",
		Addr4:     `"Anywhere, AZ 85730"`,
		Addr5:     "USA",
		Cont1:     "Jon Vendor",
		Phone1:    "5555555555",
		FirstName: "Jon",
		LastName:  "Vendor",
	}

	header, err := getHeader(reflect.ValueOf(vend), Vend)

	assert.NoError(t, err)
	fmt.Println(header)
}

type AccntData struct {
	Name      string `iif:"NAME"`
	AccntType string `iif:"ACCNTTYPE"`
	Desc      string `iif:"DESC"`
	Accnum    string `iif:"ACCNUM"`
	Extra     string `iif:"EXTRA"`
}

func (a AccntData) GetType() Type {
	return Accnt
}

type VendData struct {
	Name        string `iif:"NAME"`
	Refnum      string `iif:"REFNUM"`
	PrintAs     string `iif:"PRINTAS"`
	Addr1       string `iif:"ADDR1"`
	Addr2       string `iif:"ADDR2"`
	Addr3       string `iif:"ADDR3"`
	Addr4       string `iif:"ADDR4"`
	Addr5       string `iif:"ADDR5"`
	Vtype       string `iif:"VTYPE"`
	Cont1       string `iif:"CONT1"`
	Cont2       string `iif:"CONT2"`
	Phone1      string `iif:"PHONE1"`
	Phone2      string `iif:"PHONE2"`
	Faxnum      string `iif:"FAXNUM"`
	Email       string `iif:"EMAIL"`
	Note        string `iif:"NOTE"`
	TaxID       string `iif:"TAXID"`
	Limit       string `iif:"LIMIT"`
	Terms       string `iif:"TERMS"`
	Notepad     string `iif:"NOTEPAD"`
	Salutation  string `iif:"SALUTATION"`
	CompanyName string `iif:"COMPANYNAME"`
	FirstName   string `iif:"FIRSTNAME"`
	Midinit     string `iif:"MIDINIT"`
	LastName    string `iif:"LASTNAME"`
}

func (v VendData) GetType() Type {
	return Vend
}
