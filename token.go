package gortf

import (
	"fmt"
	"strconv"
)

type tokenType int

const (
	tokenTypeNone tokenType = iota
	tokenTypeGroup
	tokenTypeGroupEnd
	tokenTypeText
	tokenTypeControlWord
)

type controlWordType int

const (
	controlWordTypeRtf controlWordType = iota + 1
	controlWordTypeAnsi
	controlWordTypeFontTable
	controlWordTypeFontNumber
	controlWordTypeFontSize
	controlWordTypeItalic
	controlWordTypeBold
	controlWordTypeUnderline
	controlWordTypeUnknown
)

func (c controlWordType) String() string {
	switch c {
	case controlWordTypeRtf:
		return "rtf"
	case controlWordTypeAnsi:
		return "ansi"
	case controlWordTypeFontTable:
		return "fonttbl"
	case controlWordTypeFontNumber:
		return "f"
	case controlWordTypeFontSize:
		return "fs"
	case controlWordTypeItalic:
		return "i"
	case controlWordTypeBold:
		return "b"
	case controlWordTypeUnderline:
		return "i"
	default:
		return "unknown"
	}
}

type token interface {
	tokenType() tokenType
}

type binaryToken struct {
	tokType tokenType
	value   []byte
}

func (b binaryToken) tokenType() tokenType {
	return b.tokType
}

func NewBinaryToken(value []byte) binaryToken {
	return binaryToken{
		tokType: tokenTypeText,
		value:   value,
	}
}

type groupToken struct {
	tokType  tokenType
	contents []token
}

func (g groupToken) tokenType() tokenType {
	return g.tokType
}

func (g groupToken) String() string {
	return "{Group}"
}

func newGroupToken() groupToken {
	return groupToken{
		tokType: tokenTypeGroup,
	}
}

func newGroupTokenWithContents(contents []token) groupToken {
	return groupToken{
		tokType:  tokenTypeGroup,
		contents: contents,
	}
}

type groupEndToken struct {
	tokType tokenType
}

func (g groupEndToken) tokenType() tokenType {
	return g.tokType
}

func (g groupEndToken) String() string {
	return "{GroupEnd}"
}

func newGroupEndToken() groupEndToken {
	return groupEndToken{
		tokType: tokenTypeGroupEnd,
	}
}

type controlWordToken struct {
	tokType         tokenType
	name            string
	controlWordType controlWordType
	parameter       int
}

func (c controlWordToken) tokenType() tokenType {
	return c.tokType
}

func (c controlWordToken) String() string {
	return fmt.Sprintf("{ControlWord %s %s %d}", c.name, c.controlWordType, c.parameter)
}

func newControlWordToken(input string) controlWordToken {
	// TODO: implement it better
	param, idx := extractControlParameterValueAndIndex(input)
	controlType := getControlTypeFromName(input, idx)

	return controlWordToken{
		tokType:         tokenTypeControlWord,
		name:            input,
		controlWordType: controlType,
		parameter:       param,
	}
}

// TODO: handle negative value
func extractControlParameterValueAndIndex(text string) (int, int) {
	startIdx := -1
	for idx := range text {
		if isNumber(text[idx]) || text[idx] == '-' {
			startIdx = idx
			break
		}
	}

	if startIdx > 0 {
		n, err := strconv.Atoi(text[startIdx:])
		if err != nil {
			panic(err)
		}

		return n, startIdx
	}

	return -1, startIdx
}

func getControlTypeFromName(name string, suffixIndex int) controlWordType {
	if suffixIndex < 0 {
		suffixIndex = len(name)
	}

	prefix := name[:suffixIndex]

	switch prefix {
	case `\rtf`:
		return controlWordTypeRtf
	case `\ansi`:
		return controlWordTypeAnsi
	case `\fonttbl`:
		return controlWordTypeFontTable
	case `\f`:
		return controlWordTypeFontNumber
	case `\fs`:
		return controlWordTypeFontSize
	case `\i`:
		return controlWordTypeItalic
	case `\b`:
		return controlWordTypeBold
	case `\u`:
		return controlWordTypeUnderline
	default:
		return controlWordTypeUnknown
	}
}

type textToken struct {
	tokType tokenType
	value   string
}

func (t textToken) tokenType() tokenType {
	return t.tokType
}

func (t textToken) String() string {
	return fmt.Sprintf("{Text %s}", t.value)
}

func newTextToken(value string) textToken {
	return textToken{
		tokType: tokenTypeText,
		value:   value,
	}
}