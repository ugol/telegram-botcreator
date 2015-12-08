package functions

import (
	"testing"
	"text/template"
	"bytes"
	"fmt"
)

func TestSubstr(t *testing.T) {
	tpl := `{{"fooo" | substr 0 3 }}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

func TestRandomInteger(t *testing.T) {
	tpl := `{{rnd 10 30}}`
	if err := runt(tpl, "30"); err != nil {
		t.Error(err)
	}
}

func TestRandomInteger1(t *testing.T) {
	tpl := `{{rnd 10 30}}`
	if err := runt(tpl, "19"); err != nil {
		t.Error(err)
	}
}

func TestSplit(t *testing.T) {
	tpl := `{{"a|b" | split "|"}}`
	if err := runt(tpl, "[a b]"); err != nil {
		t.Error(err)
	}
}

func TestChooseRandomString(t *testing.T) {
	tpl := `{{"a|b|c|d" | split "|" | random}}`
	if err := runt(tpl, "d"); err != nil {
		t.Error(err)
	}
}

func TestTitle(t *testing.T) {
	tpl := `{{"foo" | title}}`
	if err := runt(tpl, "Foo"); err != nil {
		t.Error(err)
	}
}

func TestMax(t *testing.T) {
	tpl := `{{max 1 4}}`
	if err := runt(tpl, "4"); err != nil {
		t.Error(err)
	}
}

func TestMin(t *testing.T) {
	tpl := `{{min 1 4}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}
}

func runt(tpl, expect string) error {
	return runtv(tpl, expect, "")
}

func runtv(tpl, expect string, vars interface{}) error {
	t := template.Must(template.New("test").Funcs(FunctionsMap()).Parse(tpl))
	var b bytes.Buffer
	err := t.Execute(&b, vars)
	if err != nil {
		return err
	}
	if expect != b.String() {
		return fmt.Errorf("Expected '%s', got '%s'", expect, b.String())
	}
	return nil
}
