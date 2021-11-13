package opts

import (
	"fmt"
	"testing"
)

var strToPoint = "hello world"

func TestOptionExtraction(t *testing.T) {
	values := []interface{}{
		WithInvalidOption(true),
		WithBool(true),
		WithItem("hello"),
		WithItem("world"),
		WithUsername("userbob"),
		WithPhoneNum(8675309),
		WithPtrString(&strToPoint),
		WithList([]string{"hello", "world"}),
	}
	opts := testoptions{}
	err := Extract(&opts, values...)

	if err != nil {
		t.Fatalf("%s", err)
	}

	opts = testoptions{}
	err = MustExtract(&opts, values...)
	if err == nil {
		t.Fatalf("MustExtract should have failed, but err is nil")
	}

	eString := "invalid option WithInvalidOption"
	if fmt.Sprintf("%s", err) != eString {
		t.Fatalf("MustExtract should have failed, with an error of '%s' but failed with '%s' instead", eString, err)
	}
}

type WithBool bool
type WithItem string
type WithUsername string
type WithPhoneNum int
type WithPtrString *string
type WithList []string
type WithInvalidOption bool

type testoptions struct {
	Items     []string `optname:"WithItem"`
	PhoneNum  int      `optname:"WithPhoneNum"`
	Username  string   `optname:"WithUsername"`
	PtrString *string  `optname:"WithPtrString"`
	List      []string `optname:"WithList"`
	Boolean   bool     `optname:"WithBool"`
}
