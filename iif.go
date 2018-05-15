package iif

import (
	"fmt"
	"reflect"
	"sort"
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
	Header string
	Line   string
}

type GroupedWrapper struct {
	Type   Type
	Header string
	Lines  []string
}

func Export(dataLines []DataLine) error {
	var wrapper []Wrapper

	for _, dataLine := range dataLines {
		header, err := getHeader(dataLine)
		if err != nil {
			return err
		}

		line, err := dataLineToString(dataLine)
		if err != nil {
			return err
		}

		wrapper = append(wrapper, Wrapper{
			Type:   dataLine.GetType(),
			Header: header,
			Line:   line,
		})
	}

	wrapper = sorting(wrapper)
	gw := grouping(wrapper)

	for _, v := range gw {
		fmt.Println(v)
	}

	return nil
}

func sorting(wrapper []Wrapper) []Wrapper {
	sorting := func(i, j int) bool {
		iLoc := orderOfTypes.Location(wrapper[i].Type)
		jLoc := orderOfTypes.Location(wrapper[j].Type)

		return iLoc < jLoc
	}
	sort.Slice(wrapper, sorting)
	return wrapper
}

func grouping(wrapper []Wrapper) []GroupedWrapper {
	var gw []GroupedWrapper

	if len(wrapper) < 1 {
		return nil
	}

	var temp = groupingTemp(wrapper[0])

	for k, v := range wrapper {
		if k != 0 {

			if v.Type != temp.Type {
				gw = append(gw, temp)
				temp = groupingTemp(v)
			} else {
				temp.Lines = append(temp.Lines, v.Line)
			}
		}
	}
	gw = append(gw, temp)

	return gw
}

func groupingTemp(wrapper Wrapper) GroupedWrapper {
	return GroupedWrapper{
		Type:   wrapper.Type,
		Header: wrapper.Header,
		Lines:  []string{wrapper.Line},
	}
}

func getHeader(dataLine DataLine) (string, error) {
	var header []string
	t := dataLine.GetType()
	rv := reflect.ValueOf(dataLine)
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
	var result []string
	t := dataLine.GetType()

	rv := reflect.ValueOf(dataLine)

	for i := 0; i < rv.NumField(); i++ {
		if rv.Field(i).Kind() != reflect.String {
			result = append(result, "")
		}
		result = append(result, rv.Field(i).String())
	}

	return addType(strings.Join(result, tab), t, false), nil
}

func addType(line string, t Type, isHeader bool) string {
	if isHeader {
		t = "!" + t
	}
	return fmt.Sprintf("%s%s%s", t, tab, line)
}

type Types map[int]Type

var orderOfTypes = Types{
	0: Accnt,
	1: Invitem,
	2: Class,
	3: Cust,
	4: Vend,
	5: Trns,
	6: Spl,
}

func (t Types) Location(t2 Type) int {
	for k, v := range t {
		if v == t2 {
			return k
		}
	}

	return -1
}
