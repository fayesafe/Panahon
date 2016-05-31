package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	if Error != nil {
		t.Errorf("Logger Error not undefined on Start")
	}
	if Warning != nil {
		t.Errorf("Logger Warning not undefined on Start")
	}
	if Info != nil {
		t.Errorf("Logger Info not undefined on Start")
	}

	var buf bytes.Buffer
	t.Log("Init Logger ...")
	Init(&buf, &buf, &buf)

	if Error == nil {
		t.Errorf("Logger Error undefined after Init")
	}
	if Warning == nil {
		t.Errorf("Logger Warning undefined after Init")
	}
	if Info == nil {
		t.Errorf("Logger Info undefined after Init")
	}

	t.Log("Testing Logging Functions ...")
	Error.Println("This is an Error")
	ErrorLogStartsWith := strings.HasPrefix(buf.String(), "ERROR:")
	ErrorLogEndsWith := strings.HasSuffix(buf.String(), "This is an Error\n")
	if !ErrorLogStartsWith || !ErrorLogEndsWith {
		t.Errorf("Wrong Logging Output (Error)")
	}
	buf.Reset()

	Warning.Println("This is a Warning")
	WarningLogStartsWith := strings.HasPrefix(buf.String(), "WARNING:")
	WarningLogEndsWith := strings.HasSuffix(buf.String(), "This is a Warning\n")
	if !WarningLogStartsWith || !WarningLogEndsWith {
		t.Errorf("Wrong Logging Output (Warning)")
	}
	buf.Reset()

	Info.Println("This is an Info")
	InfoLogStartsWith := strings.HasPrefix(buf.String(), "INFO:")
	InfoLogEndsWith := strings.HasSuffix(buf.String(), "This is an Info\n")
	if !InfoLogStartsWith || !InfoLogEndsWith {
		t.Errorf("Wrong Logging Output (Info)")
	}
	buf.Reset()
}
