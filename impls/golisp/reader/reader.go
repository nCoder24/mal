package reader

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/nCoder24/mal/impls/golisp/types"
)

type Reader struct {
	tokens []string
	pos    int
}

func NewReader(tokens []string) *Reader {
	return &Reader{tokens: tokens}
}

func (r *Reader) peak() string {
	return r.tokens[r.pos]
}

func (r *Reader) next() string {
	current := r.peak()
	r.pos++

	return current
}

func (r *Reader) isEOF() bool {
	return r.pos >= len(r.tokens)
}

func ReadStr(str string) (types.MalValue, error) {
	return readForm(NewReader(tokenize(str)))
}

var tokenRegexp = regexp.MustCompile(
	"[\\s,]*(~@|[\\[\\]{}()'`~^@]|\"" +
		"(?:\\\\.|[^\\\\\"])*\"?|;.*|[^\\s\\[\\]{}('\"`,;)]*)",
)

func tokenize(str string) []string {
	tokens := make([]string, 0)

	for _, match := range tokenRegexp.FindAllStringSubmatch(str, -1) {
		tokens = append(tokens, match[1:]...)
	}

	return tokens
}

func readForm(reader *Reader) (types.MalValue, error) {
	if reader.isEOF() {
		return nil, io.EOF
	}

	switch reader.peak() {
	case "(":
		return readList(reader)
	case "[":
		return readVector(reader)
	case "{":
		return readMap(reader)
	case "'":
		return readReaderMacro(reader, "quote")
	case "`":
		return readReaderMacro(reader, "quasiquote")
	case "~":
		return readReaderMacro(reader, "unquote")
	case "@":
		return readReaderMacro(reader, "deref")
	case "~@":
		return readReaderMacro(reader, "splice-unquote")

	default:
		return readAtom(reader)
	}
}

func readReaderMacro(reader *Reader, symbol string) (types.MalValue, error) {
	reader.next()
	form, err := readForm(reader)
	if err != nil {
		return nil, err
	}

	return types.List{symbol, form}, nil
}

func readForms(reader *Reader, end string) ([]types.MalValue, error) {
	forms := make([]types.MalValue, 0)

	for reader.next(); !reader.isEOF(); reader.next() {
		if reader.peak() == end {
			return forms, nil
		}

		form, err := readForm(reader)
		if err != nil {
			return nil, err
		}

		forms = append(forms, form)
	}

	return nil, fmt.Errorf("expected '%s', got EOF", end)
}

func readList(reader *Reader) (types.List, error) {
	return readForms(reader, ")")
}

func readVector(reader *Reader) (types.Vector, error) {
	return readForms(reader, "]")
}

func readMap(reader *Reader) (types.Map, error) {
	return readForms(reader, "}")
}

var (
	strRegexp = regexp.MustCompile(`^"`)
	numRegexp = regexp.MustCompile(`^-?\d`)
)

func readAtom(reader *Reader) (types.MalValue, error) {
	token := reader.peak()

	switch {
	case token == "nil":
		return types.Nil, nil

	case token == "true":
		return types.True, nil

	case token == "false":
		return types.False, nil

	case numRegexp.MatchString(token):
		i, err := strconv.ParseFloat(token, 64)
		if err == nil {
			return types.Number(i), nil
		}

		return nil, fmt.Errorf("cannot process number '%v': %w", token, err)

	case strRegexp.MatchString(token):
		if !closed(token) {
			return nil, fmt.Errorf("expected '\"', got EOF")
		}

		return types.String(token), nil

	case strings.HasPrefix(token, ":"):
		return types.Keyword(token), nil
	}

	return types.Symbol(token), nil
}

func closed(strToken string) (closed bool) {
	for i := 1; i < len(strToken); i++ {
		char := string(strToken[i])
		if char == "\"" {
			closed = true
			break
		}

		if char == "\\" {
			i++
		}
	}

	return
}
