package gogettext

import (
	"testing"
	"os"
	"fmt"
)


func assert_equal(t *testing.T, expected string, got string){
	if expected != got {
		t.Logf("%s != %s", expected, got)
		t.Fail()
	}
}


func TestEnGettext(t *testing.T) {
	file, err := os.Open("testdata/en.mo")
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := ParseMO(file)
	if err != nil {
		t.Fatal(err)
	}
	assert_equal(t, catalog.Gettext("greeting"), "Hello")
}

func TestEnNGettext(t *testing.T) {
	file, err := os.Open("testdata/en.mo")
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := ParseMO(file)
	if err != nil {
		t.Fatal(err)
	}
	assert_equal(t,
		fmt.Sprintf(catalog.NGettext("order %d beer", "order %d beers", 0), 0),
		"0 beers please",
	)
	assert_equal(t,
		fmt.Sprintf(catalog.NGettext("order %d beer", "order %d beers", 1), 1),
		"1 beer please",
	)
	assert_equal(t,
		fmt.Sprintf(catalog.NGettext("order %d beer", "order %d beers", 2), 2),
		"2 beers please",
	)
}

func TestJaGettext(t *testing.T) {
	file, err := os.Open("testdata/ja.mo")
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := ParseMO(file)
	if err != nil {
		t.Fatal(err)
	}
	assert_equal(t, catalog.Gettext("greeting"), "こんいちは")
}

func TestJaNGettext(t *testing.T) {
	file, err := os.Open("testdata/ja.mo")
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := ParseMO(file)
	if err != nil {
		t.Fatal(err)
	}
	assert_equal(t,
		fmt.Sprintf(catalog.NGettext("order %d beer", "order %d beers", 0), 0),
		"ビールを0杯ください",
	)
	assert_equal(t,
		fmt.Sprintf(catalog.NGettext("order %d beer", "order %d beers", 1), 1),
		"ビールを1杯ください",
	)
	assert_equal(t,
		fmt.Sprintf(catalog.NGettext("order %d beer", "order %d beers", 2), 2),
		"ビールを2杯ください",
	)
}
