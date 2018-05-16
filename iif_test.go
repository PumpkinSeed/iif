package iif

import (
	"io/ioutil"
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
	err := Export(data, "example")

	assert.NoError(t, err)

	dat, err := ioutil.ReadFile("example.iif")
	assert.NoError(t, err)
	expected := "!ACCNT\tNAME\tACCNTTYPE\tDESC\tACCNUM\tEXTRA\nACCNT\tAccounts\tPayable\tAP\t2000\t\n!VEND\tNAME\tREFNUM\tPRINTAS\tADDR1\tADDR2\tADDR3\tADDR4\tADDR5\tVTYPE\tCONT1\tCONT2\tPHONE1\tPHONE2\tFAXNUM\tEMAIL\tNOTE\tTAXID\tLIMIT\tTERMS\tNOTEPAD\tSALUTATION\tCOMPANYNAME\tFIRSTNAME\tMIDINIT\tLASTNAME\nVEND\tVendor\t1\t\tJon Vendor\t555\tStreet St\t\"Anywhere, AZ 85730\"\tUSA\t\tJon Vendor\t\t5555555555\t\t\t\t\t\t\t\t\t\t\tJon\t\tVendor\n!TRNS\tTRNSID\tTRNSTYPE\tDATE\tACCNT\tNAME\tAMOUNT\tDOCNUM\tMEMO\tCLEAR\tTOPRINT\n!SPL\tSPLID\tTRNSTYPE\tDATE\tACCNT\tNAME\tCLASS\tAMOUNT\tDOCNUM\tMEMP\tCLEAR\n!ENDTRNS\nTRNS\t\tBILLPMT\t7/16/1998\tChecking\tVendor\t-35\t\tTest Memo\tN\tY\nTRNS\t\tBILLPMT\t7/16/1999\tChecking\tJoe\t-35\t\tTest Memo\tN\tY\nSPL\t\tDEPOSIT\t7/1/1998\tIncome\tCustomer\t\t-10000\t\t\tN\nENDTRNS"
	assert.Equal(t, expected, string(dat))
}

func TestSorting(t *testing.T) {
	var wrapper = []Wrapper{
		{
			Type:   Vend,
			Header: "test_header",
			Line:   "test_line",
		},
		{
			Type:   Accnt,
			Header: "test_header",
			Line:   "test_line",
		},
		{
			Type:   Trns,
			Header: "test_header",
			Line:   "test_line",
		},
		{
			Type:   Class,
			Header: "test_header",
			Line:   "test_line",
		},
	}

	wrapper = sorting(wrapper)
	assert.Equal(t, Accnt, wrapper[0].Type)
	assert.Equal(t, Class, wrapper[1].Type)
	assert.Equal(t, Vend, wrapper[2].Type)
	assert.Equal(t, Trns, wrapper[3].Type)
}

func TestGrouping(t *testing.T) {
	var wrapper = []Wrapper{
		{
			Type:   Vend,
			Header: "test_header",
			Line:   "test_line",
		},
		{
			Type:   Accnt,
			Header: "test_header",
			Line:   "test_line",
		},
		{
			Type:   Trns,
			Header: "test_header",
			Line:   "test_line",
		},
		{
			Type:   Trns,
			Header: "test_header_2",
			Line:   "test_line_2",
		},
		{
			Type:   Class,
			Header: "test_header",
			Line:   "test_line",
		},
	}

	wrapper = sorting(wrapper)
	gw := grouping(wrapper)
	assert.Equal(t, Accnt, gw[0].Type)
	assert.Equal(t, Trns, gw[3].Type)
	assert.Equal(t, "test_line", gw[3].Lines[0])
	assert.Equal(t, "test_line_2", gw[3].Lines[1])
}

func TestBuild(t *testing.T) {
	var wrapper = []Wrapper{
		{
			Type:   Vend,
			Header: "test_header",
			Line:   "test_line",
		},
		{
			Type:   Accnt,
			Header: "test_header",
			Line:   "test_line",
		},
		{
			Type:   Trns,
			Header: "test_header",
			Line:   "test_line",
		},
		{
			Type:   Trns,
			Header: "test_header_2",
			Line:   "test_line_2",
		},
		{
			Type:   Class,
			Header: "test_header",
			Line:   "test_line",
		},
	}

	wrapper = sorting(wrapper)
	gw := grouping(wrapper)
	result := build(gw)

	var expected = "test_header\ntest_line\ntest_header\ntest_line\ntest_header\ntest_line\ntest_header\n!ENDTRNS\ntest_line\ntest_line_2\nENDTRNS"

	assert.Equal(t, expected, string(result))
}

func TestGetFilename(t *testing.T) {
	type testCase struct {
		Value    string
		Expected string
	}

	cases := []testCase{
		{
			Value:    "test",
			Expected: "test.iif",
		},
		{
			Value:    "test.iif",
			Expected: "test.iif",
		},
		{
			Value:    "test.ii",
			Expected: "test.iif",
		},
	}

	for _, v := range cases {
		assert.Equal(t, v.Expected, getFilename(v.Value))
	}
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

	expected := "!VEND\tNAME\tREFNUM\tPRINTAS\tADDR1\tADDR2\tADDR3\tADDR4\tADDR5\tVTYPE\tCONT1\tCONT2\tPHONE1\tPHONE2\tFAXNUM\tEMAIL\tNOTE\tTAXID\tLIMIT\tTERMS\tNOTEPAD\tSALUTATION\tCOMPANYNAME\tFIRSTNAME\tMIDINIT\tLASTNAME"
	assert.Equal(t, expected, header)
}

func TestDataLineToString(t *testing.T) {
	vend := vendData{
		Name:      "Vendor",
		Refnum:    "1",
		PrintAs:   "",
		Addr1:     "Jon Vendor",
		Addr2:     "555",
		Addr3:     "Street St",
		Addr4:     `Anywhere, AZ 85730`,
		Addr5:     "USA",
		Cont1:     "Jon Vendor",
		Phone1:    "5555555555",
		FirstName: "Jon",
		LastName:  "Vendor",
	}

	result, err := dataLineToString(vend)

	assert.NoError(t, err)

	expected := "VEND\tVendor\t1\t\tJon Vendor\t555\tStreet St\tAnywhere, AZ 85730\tUSA\t\tJon Vendor\t\t5555555555\t\t\t\t\t\t\t\t\t\t\tJon\t\tVendor"
	assert.Equal(t, expected, result)
}

func TestGetEndTrns(t *testing.T) {
	endtrns := getEndTrns(true)
	assert.Equal(t, "!"+endTrns, endtrns)

	endtrns = getEndTrns(false)
	assert.Equal(t, endTrns, endtrns)
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
	return Spl
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
		err := Export(data, "example")
		if err != nil {
			panic(err)
		}
	}
}
