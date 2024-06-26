package format

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	flagGo = "GO"
	flagEt = "EASY"
)

const (
	unknown style = iota
	title
	lower
	upper
)

// ErrNamingFormat defines an error for unknown format
var ErrNamingFormat = errors.New("unsupported format")

type (
	styleFormat struct {
		before  string
		through string
		after   string
		goStyle style
		etStyle style
	}

	style int
)

func FileNamingFormat(format, content string) (string, error) {
	upperFormat := strings.ToUpper(format)
	indexGo := strings.Index(upperFormat, flagGo)
	indexEt := strings.Index(upperFormat, flagEt)
	if indexGo < 0 || indexEt < 0 || indexGo > indexEt {
		return "", ErrNamingFormat
	}
	var (
		before, through, after string
		flagGo, flagZero       string
		goStyle, zeroStyle     style
		err                    error
	)
	before = format[:indexGo]
	flagGo = format[indexGo : indexGo+2]
	through = format[indexGo+2 : indexEt]
	flagZero = format[indexEt : indexEt+4]
	after = format[indexEt+4:]

	goStyle, err = getStyle(flagGo)
	if err != nil {
		return "", err
	}

	zeroStyle, err = getStyle(flagZero)
	if err != nil {
		return "", err
	}
	var formatStyle styleFormat
	formatStyle.goStyle = goStyle
	formatStyle.etStyle = zeroStyle
	formatStyle.before = before
	formatStyle.through = through
	formatStyle.after = after
	return doFormat(formatStyle, content)
}

func doFormat(f styleFormat, content string) (string, error) {
	splits, err := split(content)
	if err != nil {
		return "", err
	}
	var join []string
	for index, split := range splits {
		if index == 0 {
			join = append(join, transferTo(split, f.goStyle))
			continue
		}
		join = append(join, transferTo(split, f.etStyle))
	}
	joined := strings.Join(join, f.through)
	return f.before + joined + f.after, nil
}

func transferTo(in string, style style) string {
	switch style {
	case upper:
		return strings.ToUpper(in)
	case lower:
		return strings.ToLower(in)
	case title:
		return strings.Title(in)
	default:
		return in
	}
}

func split(content string) ([]string, error) {
	var (
		list   []string
		reader = strings.NewReader(content)
		buffer = bytes.NewBuffer(nil)
	)
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				if buffer.Len() > 0 {
					list = append(list, buffer.String())
				}
				return list, nil
			}
			return nil, err
		}
		if r == '_' {
			if buffer.Len() > 0 {
				list = append(list, buffer.String())
			}
			buffer.Reset()
			continue
		}

		if r >= 'A' && r <= 'Z' {
			if buffer.Len() > 0 {
				list = append(list, buffer.String())
			}
			buffer.Reset()
		}
		buffer.WriteRune(r)
	}
}

func getStyle(flag string) (style, error) {
	compare := strings.ToLower(flag)
	switch flag {
	case strings.ToLower(compare):
		return lower, nil
	case strings.ToUpper(compare):
		return upper, nil
	case strings.Title(compare):
		return title, nil
	default:
		return unknown, fmt.Errorf("unexpected format: %s", flag)
	}
}
