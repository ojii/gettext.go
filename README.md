# gettext in golang

[![Build Status](https://travis-ci.org/ojii/gogettext.svg?branch=master)](https://travis-ci.org/ojii/gogettext)

## TODO

- [x] parse mofiles
- [x] compile plural forms
- [ ] non-utf8 mo files (possible wontfix)
- [x] gettext
- [x] ngettext
- [x] managing mo files / sane API


## Example


```go

translations := gogettext.NewTranslations("path/to/translations/", "messages", gogettext.DefaultResolver)

locale = translations.Locale("en")

fmt.Println(locale.Gettext("hello from gettext"))
```
