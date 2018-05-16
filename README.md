# iif

iif file format exporter

### Behavior

There is the `DataLine` interface, to get the type of the struct passed as parameter to the `Export`:

```
type DataLine interface {
	GetType() Type
}
```

The `Export` function waits a list of `DataLine`s and the filename.


### Types

```
Accnt   Type = "ACCNT"
Invitem Type = "INVITEM"
Class   Type = "CLASS"
Cust    Type = "CUST"
Vend    Type = "VEND"
Trns    Type = "TRNS"
Spl     Type = "SPL"
```

### Benchmark

```
goos: darwin
goarch: amd64
pkg: github.com/PumpkinSeed/iif
BenchmarkExport-4   	   10000	    206726 ns/op
PASS
ok  	github.com/PumpkinSeed/iif	2.113s
```

### Usage

```
type accntData struct {
	Name      string `iif:"NAME"`
	AccntType string `iif:"ACCNTTYPE"`
	Desc      string `iif:"DESC"`
	Accnum    string `iif:"ACCNUM"`
	Extra     string `iif:"EXTRA"`
}

func (a accntData) GetType() Type {
	return iif.Accnt
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
	return iif.Vend
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
	return iif.Trns
}

type splData struct {
	SplID    string `iif:"SPLID"`
	TrnsType string `iif:"TRNSTYPE"`
	Date     string `iif:"DATE"`
	Accnt    string `iif:"ACCNT"`
	Name     string `iif:"NAME"`
	Class    string `iif:"CLASS"`
	Amount   string `iif:"AMOUNT"`
	Docnum   string `iif:"DOCNUM"`
	Memo     string `iif:"MEMP"`
	Clear    string `iif:"CLEAR"`
}

func (s splData) GetType() Type {
	return iif.Spl
}

data := []iif.DataLine{
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
	trnsData{
		TrnsID:   "",
		TrnsType: "BILLPMT",
		Date:     "7/16/1999",
		Accnt:    "Checking",
		Name:     "Joe",
		Amount:   "-35",
		Docnum:   "",
		Memo:     "Test Memo",
		Clear:    "N",
		ToPrint:  "Y",
	},
	splData{
		SplID:    "",
		TrnsType: "DEPOSIT",
		Date:     "7/1/1998",
		Accnt:    "Income",
		Name:     "Customer",
		Class:    "",
		Amount:   "-10000",
		Docnum:   "",
		Memo:     "",
		Clear:    "N",
	},
}

err := iif.Export(data, {FILENAME})
// handler err
```

**Output in the file**

```
!ACCNT	NAME	ACCNTTYPE	DESC	ACCNUM	EXTRA
ACCNT	Accounts	Payable	AP	2000	
!VEND	NAME	REFNUM	PRINTAS	ADDR1	ADDR2	ADDR3	ADDR4	ADDR5	VTYPE	CONT1	CONT2	PHONE1	PHONE2	FAXNUM	EMAIL	NOTE	TAXID	LIMIT	TERMS	NOTEPAD	SALUTATION	COMPANYNAME	FIRSTNAME	MIDINIT	LASTNAME
VEND	Vendor	1		Jon Vendor	555	Street St	"Anywhere, AZ 85730"	USA		Jon Vendor		5555555555											Jon		Vendor
!TRNS	TRNSID	TRNSTYPE	DATE	ACCNT	NAME	AMOUNT	DOCNUM	MEMO	CLEAR	TOPRINT
!SPL	SPLID	TRNSTYPE	DATE	ACCNT	NAME	CLASS	AMOUNT	DOCNUM	MEMP	CLEAR
!ENDTRNS
TRNS		BILLPMT	7/16/1998	Checking	Vendor	-35		Test Memo	N	Y
TRNS		BILLPMT	7/16/1999	Checking	Joe	-35		Test Memo	N	Y
SPL		DEPOSIT	7/1/1998	Income	Customer		-10000			N
ENDTRNS
```