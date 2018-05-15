package iif

import (
	"fmt"
	"reflect"
	"strings"
)

type Type string

const (
	Accnt   Type = "ACCNT"
	Invitem Type = "INVITEM"
	Class   Type = "CLASS"
	Cust    Type = "CUST"
	Vend    Type = "VEND"
	Trns    Type = "TRNS"
	Spl     Type = "SPL"
)

const (
	tab = "\t"
	tag = "iif"
)

type DataLine interface {
	GetType() Type
}

type Wrapper struct {
	Type   Type
	Header []string
	Lines  []DataLine
}

func Export(dataLines []DataLine) error {
	for _, dataLine := range dataLines {
		rv := reflect.ValueOf(dataLine)
		//rt := rv.Type()

		if rv.Kind() != reflect.Struct {
			continue
		}

		fmt.Println(rv.Kind())

	}
	//fmt.Println(data)
	return nil
}

func getHeader(rv reflect.Value, t Type) (string, error) {
	var header []string

	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		value, ok := rt.Field(i).Tag.Lookup(tag)
		if ok {
			header = append(header, value)
		}
	}

	return addType(strings.Join(header, tab), t, true), nil
}

func dataLineToString(dataLine DataLine) (string, error) {
	var result string

	return result, nil
}

func addType(line string, t Type, isHeader bool) string {
	if isHeader {
		t = "!" + t
	}
	return fmt.Sprintf("%s%s%s", t, tab, line)
}

var orderOfTypes = map[int]Type{
	0: Accnt,
	1: Invitem,
	2: Class,
	3: Cust,
	4: Vend,
	5: Trns,
	6: Spl,
}
