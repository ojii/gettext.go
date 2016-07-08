package gogettext

import (
	"fmt"
	"path"
	"testing"
)

func TestNullTranslations(t *testing.T) {
	translations := NewTranslations(".", "messages", DefaultResolver)
	en := translations.Locale("en")
	en_gettext := en.Gettext("mymsgid")
	assert_equal(t, en_gettext, "mymsgid")
	en_ngettext_0 := en.NGettext("mymsgid", "mymsgidp", 0)
	assert_equal(t, en_ngettext_0, "mymsgidp")
	en_ngettext_1 := en.NGettext("mymsgid", "mymsgidp", 1)
	assert_equal(t, en_ngettext_1, "mymsgid")
	en_ngettext_2 := en.NGettext("mymsgid", "mymsgidp", 2)
	assert_equal(t, en_ngettext_2, "mymsgidp")
	ja := translations.Locale("ja")
	ja_gettext := ja.Gettext("mymsgid")
	assert_equal(t, ja_gettext, "mymsgid")
	ja_ngettext_0 := ja.NGettext("mymsgid", "mymsgidp", 0)
	assert_equal(t, ja_ngettext_0, "mymsgidp")
	ja_ngettext_1 := ja.NGettext("mymsgid", "mymsgidp", 1)
	assert_equal(t, ja_ngettext_1, "mymsgid")
	ja_ngettext_2 := ja.NGettext("mymsgid", "mymsgidp", 2)
	assert_equal(t, ja_ngettext_2, "mymsgidp")
}

func my_resolver(root string, locale string, domain string) string {
	return path.Join(root, locale, fmt.Sprintf("%s.mo", domain))
}

func TestRealTranslations(t *testing.T) {
	translations := NewTranslations("testdata/", "messages", my_resolver)
	en := translations.Locale("en")
	assert_equal(t, en.Gettext("greeting"), "Hello")
	assert_equal(t,
		fmt.Sprintf(en.NGettext("order %d beer", "order %d beers", 0), 0),
		"0 beers please",
	)
	assert_equal(t,
		fmt.Sprintf(en.NGettext("order %d beer", "order %d beers", 1), 1),
		"1 beer please",
	)
	assert_equal(t,
		fmt.Sprintf(en.NGettext("order %d beer", "order %d beers", 2), 2),
		"2 beers please",
	)
	ja := translations.Locale("ja")
	assert_equal(t, ja.Gettext("greeting"), "こんいちは")
	assert_equal(t,
		fmt.Sprintf(ja.NGettext("order %d beer", "order %d beers", 0), 0),
		"ビールを0杯ください",
	)
	assert_equal(t,
		fmt.Sprintf(ja.NGettext("order %d beer", "order %d beers", 1), 1),
		"ビールを1杯ください",
	)
	assert_equal(t,
		fmt.Sprintf(ja.NGettext("order %d beer", "order %d beers", 2), 2),
		"ビールを2杯ください",
	)
	de := translations.Locale("de")
	assert_equal(t, de.Gettext("greeting"), "greeting")
	assert_equal(t,
		fmt.Sprintf(de.NGettext("order %d beer", "order %d beers", 0), 0),
		"order 0 beers",
	)
	assert_equal(t,
		fmt.Sprintf(de.NGettext("order %d beer", "order %d beers", 1), 1),
		"order 1 beer",
	)
	assert_equal(t,
		fmt.Sprintf(de.NGettext("order %d beer", "order %d beers", 2), 2),
		"order 2 beers",
	)
}
