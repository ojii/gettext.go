package gogettext

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ojii/gogettext/pluralforms"
	"log"
	"os"
	"strings"
)

const LE_MAGIC = 0x950412de
const BE_MAGIC = 0xde120495

type Header struct {
	Version          uint32
	NumStrings       uint32
	MasterIndex      uint32
	TranslationIndex uint32
}

func (header Header) GetMajorVersion() uint32 {
	return header.Version >> 16
}

func (header Header) GetMinorVersion() uint32 {
	return header.Version & 0xffff
}

type Catalog interface {
	Gettext(msgid string) string
	NGettext(msgid string, msgid_plural string, n uint32) string
}

type MoCatalog struct {
	Header      Header
	Language    string
	Messages    map[string][]string
	PluralForms pluralforms.Expression
	Info        map[string]string
	Charset     string
}

type NullCatalog struct{}

func (catalog NullCatalog) Gettext(msgid string) string {
	return msgid
}

func (catalog NullCatalog) NGettext(msgid string, msgid_plural string, n uint32) string {
	if n == 1 {
		return msgid
	} else {
		return msgid_plural
	}
}

func (catalog MoCatalog) Gettext(msgid string) string {
	msgstrs, ok := catalog.Messages[msgid]
	if !ok {
		return msgid
	}
	return msgstrs[0]
}

func (catalog MoCatalog) NGettext(msgid string, msgid_plural string, n uint32) string {
	msgstrs, ok := catalog.Messages[msgid]
	if !ok {
		if n == 1 {
			return msgid
		} else {
			return msgid_plural
		}
	} else {
		index := catalog.PluralForms.Eval(n)
		if index > len(msgstrs) {
			if n == 1 {
				return msgid
			} else {
				return msgid_plural
			}
		}
		return msgstrs[index]
	}
}

type len_offset struct {
	Len uint32
	Off uint32
}

func read_len_off(index uint32, file *os.File, order binary.ByteOrder) (len_offset, error) {
	lenoff := len_offset{}
	buf := make([]byte, 8)
	_, err := file.Seek(int64(index), os.SEEK_SET)
	if err != nil {
		return lenoff, err
	}
	_, err = file.Read(buf)
	if err != nil {
		return lenoff, err
	}
	buffer := bytes.NewBuffer(buf)
	err = binary.Read(buffer, order, &lenoff)
	if err != nil {
		return lenoff, err
	}
	return lenoff, nil
}

func read_message(file *os.File, lenoff len_offset) (string, error) {
	_, err := file.Seek(int64(lenoff.Off), os.SEEK_SET)
	if err != nil {
		return "", err
	}
	buf := make([]byte, lenoff.Len)
	_, err = file.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (catalog *MoCatalog) read_info(info string) error {
	lastk := ""
	for _, line := range strings.Split(info, "\n") {
		item := strings.TrimSpace(line)
		if len(item) == 0 {
			continue
		}
		var k string
		var v string
		if strings.Contains(item, ":") {
			tmp := strings.SplitN(item, ":", 2)
			k = strings.ToLower(strings.TrimSpace(tmp[0]))
			v = strings.TrimSpace(tmp[1])
			catalog.Info[k] = v
			lastk = k
		} else if len(lastk) != 0 {
			catalog.Info[lastk] += "\n" + item
		}
		if k == "content-type" {
			catalog.Charset = strings.Split(v, "charset=")[1]
		} else if k == "plural-forms" {
			p := strings.Split(v, ";")[1]
			s := strings.Split(p, "plural=")[1]
			expr, err := pluralforms.Compile(s)
			if err != nil {
				return err
			}
			catalog.PluralForms = expr
		}
	}
	return nil
}

func ParseMO(file *os.File) (Catalog, error) {
	var order binary.ByteOrder
	header := Header{}
	catalog := MoCatalog{
		Header:   header,
		Info:     make(map[string]string),
		Messages: make(map[string][]string),
	}
	magic := make([]byte, 4)
	_, err := file.Read(magic)
	if err != nil {
		return catalog, err
	}
	magic_number := binary.LittleEndian.Uint32(magic)
	switch magic_number {
	case LE_MAGIC:
		order = binary.LittleEndian
	case BE_MAGIC:
		order = binary.BigEndian
	default:
		return catalog, errors.New(fmt.Sprintf("Wrong magic %d", magic_number))
	}
	raw_headers := make([]byte, 32)
	_, err = file.Read(raw_headers)
	if err != nil {
		return catalog, err
	}
	buffer := bytes.NewBuffer(raw_headers)
	err = binary.Read(buffer, order, &header)
	if err != nil {
		return catalog, err
	}
	if (header.GetMajorVersion() != 0) && (header.GetMajorVersion() != 1) {
		log.Printf("major %d minor %d", header.GetMajorVersion(), header.GetMinorVersion())
		return catalog, errors.New(fmt.Sprintf("Unsupported version: %d.%d", header.GetMajorVersion(), header.GetMinorVersion()))
	}
	current_master_index := header.MasterIndex
	current_transl_index := header.TranslationIndex
	var index uint32 = 0
	for ; index < header.NumStrings; index++ {
		mlenoff, err := read_len_off(current_master_index, file, order)
		if err != nil {
			return catalog, err
		}
		tlenoff, err := read_len_off(current_transl_index, file, order)
		if err != nil {
			return catalog, err
		}
		msgid, err := read_message(file, mlenoff)
		if err != nil {
			return catalog, nil
		}
		msgstr, err := read_message(file, tlenoff)
		if err != nil {
			return catalog, err
		}
		if mlenoff.Len == 0 {
			err = catalog.read_info(msgstr)
			if err != nil {
				return catalog, err
			}
		}
		if strings.Contains(msgid, "\x00") {
			// Plural!
			msgidsingular := strings.Split(msgid, "\x00")[0]
			translations := strings.Split(msgstr, "\x00")
			catalog.Messages[msgidsingular] = translations
		} else {
			catalog.Messages[msgid] = []string{msgstr}
		}

		current_master_index += 8
		current_transl_index += 8
	}
	return catalog, nil
}
