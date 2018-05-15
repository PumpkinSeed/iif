package iif

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T) {
	data := []DataLine{
		accntData{
			Name:      "Accounts",
			AccntType: "Payable",
			Desc:      "AP",
			Accnum:    "2000",
		},
		vendData{
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
		trnsData{
			TrnsID:   "",
			TrnsType: "BILLPMT",
			Date:     "7/16/1998",
			Accnt:    "Checking",
			Name:     "Vendor",
			Amount:   "-35",
			Docnum:   "",
			Memo:     "Test Memo",
			Clear:    "N",
			ToPrint:  "Y",
		},
	}
	err := Export(data)

	assert.NoError(t, err)
}

func TestGetHeader(t *testing.T) {
	vend := vendData{
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

	header, err := getHeader(vend)

	assert.NoError(t, err)
	fmt.Println(header)
}

func TestDataLineToString(t *testing.T) {
	vend := vendData{
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

	result, err := dataLineToString(vend)

	assert.NoError(t, err)
	fmt.Println(result)
}

func TestOrderLocation(t *testing.T) {
	loc := orderOfTypes.Location(Accnt)
	assert.Equal(t, loc, 0)

	loc = orderOfTypes.Location(Trns)
	assert.Equal(t, loc, 5)
}

type accntData struct {
	Name      string `iif:"NAME"`
	AccntType string `iif:"ACCNTTYPE"`
	Desc      string `iif:"DESC"`
	Accnum    string `iif:"ACCNUM"`
	Extra     string `iif:"EXTRA"`
}

func (a accntData) GetType() Type {
	return Accnt
}

type vendData struct {
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

func (v vendData) GetType() Type {
	return Vend
}

type trnsData struct {
	TrnsID   string `iif:"TRNSID"`
	TrnsType string `iif:"TRNSTYPE"`
	Date     string `iif:"DATE"`
	Accnt    string `iif:"ACCNT"`
	Name     string `iif:"NAME"`
	Amount   string `iif:"AMOUNT"`
	Docnum   string `iif:"DOCNUM"`
	Memo     string `iif:"MEMO"`
	Clear    string `iif:"CLEAR"`
	ToPrint  string `iif:"TOPRINT"`
}

func (t trnsData) GetType() Type {
	return Trns
}

func BenchmarkExport(b *testing.B) {
	data := []DataLine{
		accntData{
			Name:      "Accounts",
			AccntType: "Payable",
			Desc:      "AP",
			Accnum:    "2000",
		},
		vendData{
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
		trnsData{
			TrnsID:   "",
			TrnsType: "BILLPMT",
			Date:     "7/16/1998",
			Accnt:    "Checking",
			Name:     "Vendor",
			Amount:   "-35",
			Docnum:   "",
			Memo:     "Test Memo",
			Clear:    "N",
			ToPrint:  "Y",
		},
	}

	for i := 0; i < b.N; i++ {
		err := Export(data)
		if err != nil {
			panic(err)
		}
	}
}
