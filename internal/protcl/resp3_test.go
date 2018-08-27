package protcl

import (
	"bufio"
	"strings"
	"testing"

	testifyAssert "github.com/stretchr/testify/assert"
)

func testResp3Parser(t *testing.T, input, expect string) {
	assert := testifyAssert.New(t)

	parser := NewResp3Parser(bufio.NewReader(strings.NewReader(input)))
	result, err := parser.Parse()
	assert.Nil(err)
	assert.Equal(expect, result.String())
}

func testResp3Error(t *testing.T, input, e string) {
	assert := testifyAssert.New(t)

	parser := NewResp3Parser(bufio.NewReader(strings.NewReader(input)))
	result, err := parser.Parse()
	assert.NotNil(err)
	assert.Nil(result)
	assert.Equal(e, err.Error())
}

func TestResp3Parser(t *testing.T) {
	// simple string
	testResp3Error(t, "+string", `EOF`)
	testResp3Parser(t, "+string\n", `"string"`)

	// blob string
	testResp3Parser(t, "$0\n\n", `""`)
	testResp3Error(t, "$1\n\n", `unexpected line end`)
	testResp3Error(t, "$1\naa\n", `unexpected line end`)
	testResp3Parser(t, "$10\n1234567890\n", `"1234567890"`)

	// simple error
	testResp3Parser(t, "-error\n", `(error) error`)
	testResp3Parser(t, "-Err error\n", `(error) Err error`)

	// blob error
	testResp3Parser(t, "!3\nerr\n", `(error) err`)
	testResp3Parser(t, "!17\nErr this is error\n", `(error) Err this is error`)

	// number
	testResp3Parser(t, ":-1\n", `(integer) -1`)
	testResp3Parser(t, ":0\n", `(integer) 0`)
	testResp3Parser(t, ":100\n", `(integer) 100`)

	// double
	testResp3Parser(t, ",-1\n", `(double) -1`)
	testResp3Parser(t, ",0\n", `(double) 0`)
	testResp3Parser(t, ",10\n", `(double) 10`)
	testResp3Parser(t, ",.1\n", `(double) 0.1`)
	testResp3Parser(t, ",1.23\n", `(double) 1.23`)
	testResp3Parser(t, ",1.\n", `(double) 1`)
	testResp3Error(t, ",invalid\n", `convert invalid to double fail, because of strconv.ParseFloat: parsing "invalid": invalid syntax`)

	// big number
	testResp3Parser(t, "(3492890328409238509324850943850943825024385\n", `(big number) 3492890328409238509324850943850943825024385`)
	testResp3Error(t, "(invalid string\n", `convert invalid string to Big Number fail`)

	// null
	testResp3Parser(t, "_\n", "(null)")

	// boolean
	testResp3Parser(t, "#t\n", `(boolean) true`)
	testResp3Parser(t, "#f\n", `(boolean) false`)
	testResp3Error(t, "#x\n", `unexpect string: t/f`)
	testResp3Error(t, "#\n", `unexpected line end`)

	// array
	testResp3Error(t, "*3\n:1\n:2\n", "EOF")
	testResp3Parser(t, "*3\n:1\n:2\n:3\n", "(array)\n\t(integer) 1\n\t(integer) 2\n\t(integer) 3")
	testResp3Parser(t, "*2\n*3\n:1\n$5\nhello\n:2\n#f\n", "(array)\n\t(array)\n\t\t(integer) 1\n\t\t\"hello\"\n\t\t(integer) 2\n\t(boolean) false")

	// set
	testResp3Error(t, "~3\n:1\n:2\n", "EOF")
	testResp3Parser(t, "~3\n:1\n:2\n:3\n", "(set)\n\t(integer) 1\n\t(integer) 2\n\t(integer) 3")
	testResp3Parser(t, "~2\n*3\n:1\n$5\nhello\n:2\n#f\n", "(set)\n\t(array)\n\t\t(integer) 1\n\t\t\"hello\"\n\t\t(integer) 2\n\t(boolean) false")
}