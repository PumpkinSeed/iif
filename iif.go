package iif

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
)

// Type defines the Type enum
type Type string

const (
	// Accnt type
	Accnt Type = "ACCNT"

	// Invitem type
	Invitem Type = "INVITEM"

	// Class type
	Class Type = "CLASS"

	// Cust type
	Cust Type = "CUST"

	// Vend type
	Vend Type = "VEND"

	// Trns type
	Trns Type = "TRNS"

	// Spl type
	Spl Type = "SPL"
)

const (
	tab     = "\t"
	newLine = "\n"
	endTrns = "ENDTRNS"
	tag     = "iif"
)

// DataLine provides an interface type to determine the structs type
type DataLine interface {
	GetType() Type
}

// Wrapper is the container for a line
type Wrapper struct {
	Type   Type
	Header string
	Line   string
}

// GroupedWrapper is a container for a specified type
type GroupedWrapper struct {
	Type   Type
	Header string
	Lines  []string
}

// Export is the main entrypoint
func Export(dataLines []DataLine, filename string) error {
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
	lines := build(gw)

	return writeFile(lines, filename)
}

// Types is the list of Types
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

// Location determine the location in the Type list
func (t Types) Location(t2 Type) int {
	for k, v := range t {
		if v == t2 {
			return k
		}
	}

	return -1
}

/*
	HELPERS
*/

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

func build(gw []GroupedWrapper) []byte {
	var result []string

	var trnsKey int
	for k, v := range gw {
		if v.Type == Trns {
			trnsKey = k
			break
		}
		result = append(result, v.Header)
		for _, line := range v.Lines {
			result = append(result, line)
		}
	}

	trns := buildTrns(gw[trnsKey:len(gw)])
	for _, line := range trns {
		result = append(result, line)
	}

	return []byte(strings.Join(result, newLine))
}

func buildTrns(gw []GroupedWrapper) []string {
	var result []string

	for _, v := range gw {
		result = append(result, v.Header)
	}
	result = append(result, getEndTrns(true))

	for _, v := range gw {
		for _, line := range v.Lines {
			result = append(result, line)
		}
	}
	result = append(result, getEndTrns(false))

	return result
}

func writeFile(data []byte, filename string) error {
	if _, err := os.Stat(getFilename(filename)); os.IsNotExist(err) {
		f, err := os.Create(getFilename(filename))
		if err != nil {
			return err
		}

		defer f.Close()

		_, err = f.Write(data)
		if err != nil {
			return err
		}
	}

	err := ioutil.WriteFile(getFilename(filename), data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getFilename(filename string) string {
	const sep = "."

	if strings.Contains(filename, sep) {
		parts := strings.Split(filename, sep)
		parts[len(parts)-1] = "iif"
		return strings.Join(parts, sep)
	}
	return filename + ".iif"
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

func getEndTrns(header bool) string {
	if header {
		return "!" + endTrns
	}
	return endTrns
}

func addType(line string, t Type, isHeader bool) string {
	if isHeader {
		t = "!" + t
	}
	return fmt.Sprintf("%s%s%s", t, tab, line)
}
