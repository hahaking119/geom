package tcase

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hahaking119/geom/encoding/wkb/internal/tcase/token"
)

var ErrMissingDesc = fmt.Errorf("missing desc field")

type Type uint8

const (
	TypeNone   Type = 0
	TypeEncode      = 1
	TypeDecode      = 2
	TypeBoth        = 3
)

func (ot Type) Is(t Type) bool { return ot&t == t }

type C struct {
	Desc        string
	BOM         binary.ByteOrder
	Skip        Type
	Expected    interface{}
	DecodeError string
	EncodeError string
	Bytes       []byte
}

func (c C) HasErrorFor(t Type) bool {
	switch t {
	case TypeEncode:
		return c.EncodeError != ""
	case TypeDecode:
		return c.DecodeError != ""
	}
	return false
}
func (c C) ErrorFor(t Type) string {
	switch t {
	case TypeEncode:
		return c.EncodeError
	case TypeDecode:
		return c.DecodeError
	}
	return ""
}
func (c C) DoesErrorMatch(t Type, e error) bool {
	if !c.HasErrorFor(t) {
		return e == nil
	}
	return e.Error() == c.ErrorFor(t)
}

func parse(r io.Reader, filename string) (cases []C, err error) {
	t := token.NewT(r)
	var cC *C
	for !t.AtEnd() {
		t.EatCommentsAndSpaces()
		if t.AtEnd() {
			break
		}
		label, err := t.ParseLabel()
		if err != nil {
			log.Printf("error trying to get label %#v", cC)
			return nil, err
		}
		switch strings.ToLower(string(label)) {
		case "desc":
			if cC != nil {
				cases = append(cases, *cC)
			}
			cC = new(C)
			cC.Desc = strings.TrimSpace(string(t.ParseTillEndIgnoreComments()))
		case "skip":
			if cC == nil {
				return cases, ErrMissingDesc
			}
			val := strings.TrimSpace(string(t.ParseTillEndIgnoreComments()))
			strings.ToLower(val)
			strings.TrimSpace(val)
			switch val {
			case "encode":
				cC.Skip = TypeEncode
			case "decode":
				cC.Skip = TypeDecode
			case "both":
				cC.Skip = TypeBoth
			}

			if err != nil {
				return cases, err
			}

		case "decode_error":
			if cC == nil {
				return cases, ErrMissingDesc
			}
			cC.DecodeError = strings.TrimSpace(string(t.ParseTillEndIgnoreComments()))
		case "encode_error":
			if cC == nil {
				return cases, ErrMissingDesc
			}
			cC.EncodeError = strings.TrimSpace(string(t.ParseTillEndIgnoreComments()))
		case "bytes":
			if cC == nil {
				return cases, ErrMissingDesc
			}
			bin, err := t.ParseBinaryField()
			if err != nil {
				return cases, err
			}
			cC.Bytes = bin
		case "bom":
			if cC == nil {
				return cases, ErrMissingDesc
			}
			bom := strings.ToLower(strings.TrimSpace(string(t.ParseTillEndIgnoreComments())))
			switch bom {
			case "little":
				cC.BOM = binary.LittleEndian
			case "big":
				cC.BOM = binary.BigEndian
			default:
				return cases, fmt.Errorf("invalid bom(%v), expect “little” or “big”", bom)
			}
		case "geometry":
			fallthrough
		case "expected":
			if cC == nil {
				return cases, ErrMissingDesc
			}
			geom, err := t.ParseExpectedField()
			if err != nil {
				return cases, err
			}
			cC.Expected = geom
		default:
			return cases, fmt.Errorf("unknown label:%v", string(label))
		}
	}
	cases = append(cases, *cC)
	return cases, nil
}

func Parse(r io.Reader, filename string) ([]C, error) {
	return parse(r, filename)
}
func ParseFile(filename string) ([]C, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parse(file, filename)
}

var isolatedFilenames = flag.String("tcase.Files", "", "List of comma separated file name to grab the test cases from; instead of all the files in the directory.")

func GetFiles(dir string) (filenames []string, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var ffiles []string
	if *isolatedFilenames != "" {
		ffiles = strings.Split(*isolatedFilenames, ",")
		for i := range ffiles {
			ffiles[i] = strings.TrimSpace(ffiles[i])
		}
	}
LOOP_FILES:
	for _, f := range files {
		fname := f.Name()
		fext := strings.ToLower(filepath.Ext(fname))
		if fext != ".tcase" {
			continue
		}
		if len(ffiles) != 0 {
			// need to filter out filenames.
			for i := range ffiles {
				if ffiles[i] == fname {
					goto ADD_FILE
				}
			}
			// We did not find a file matching this file so skip it.
			continue LOOP_FILES
		}
	ADD_FILE:
		filenames = append(filenames, filepath.Join(dir, fname))
	}
	return filenames, nil
}

func SprintBinary(bytes []byte, prefix string) (out string) {
	out = prefix + "//01 02 03 04 05 06 07 08"
	for i, b := range bytes {
		if i%8 == 0 {
			out += "\n" + prefix + "  "
		}
		out += fmt.Sprintf("%02x ", b)
	}
	return out
}
