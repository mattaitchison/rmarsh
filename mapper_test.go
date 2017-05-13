package rmarsh_test

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"testing"

	"github.com/samcday/rmarsh"
)

func testMapperWriteValue(t *testing.T, exp string, v interface{}) {
	b := new(bytes.Buffer)
	gen := rmarsh.NewGenerator(b)

	if err := rmarsh.NewMapper().WriteValue(gen, v); err != nil {
		t.Fatal(err)
	}

	str := rbDecode(t, b.Bytes())
	if str != exp {
		t.Fatalf("Generated stream %s != %s\nRaw marshal:\n%s\n", str, exp, hex.Dump(b.Bytes()))
	}
}

func testMapperReadValue(t *testing.T, expr string, v interface{}) {
	b := rbEncode(t, expr)
	p := rmarsh.NewParser(bytes.NewReader(b))
	if err := rmarsh.NewMapper().ReadValue(p, v); err != nil {
		t.Fatal(err)
	}
}

func TestMapperWriteValueBool(t *testing.T) {
	testMapperWriteValue(t, `true`, true)
	v := true
	testMapperWriteValue(t, `true`, &v)
}

func BenchmarkMapperWriteTrue(b *testing.B) {
	gen := rmarsh.NewGenerator(ioutil.Discard)
	mapper := rmarsh.NewMapper()
	v := true

	for i := 0; i < b.N; i++ {
		gen.Reset(nil)

		if err := mapper.WriteValue(gen, &v); err != nil {
			b.Fatal(err)
		}
	}
}

func TestMapperReadValueBool(t *testing.T) {
	var v bool
	testMapperReadValue(t, "true", &v)

	if v != true {
		t.Errorf("%v != true", v)
	}

	var ptr *bool
	testMapperReadValue(t, "true", &ptr)

	if *ptr != true {
		t.Errorf("%v != true", ptr)
	}

	var silly *****bool
	testMapperReadValue(t, "true", &silly)

	if *****silly != true {
		t.Errorf("%v != true", silly)
	}
}

func BenchmarkMapperReadTrue(b *testing.B) {
	r := newCyclicReader(rbEncode(b, "true"))
	p := rmarsh.NewParser(r)
	mapper := rmarsh.NewMapper()

	var v bool

	for i := 0; i < b.N; i++ {
		v = false
		p.Reset()

		if err := mapper.ReadValue(p, &v); err != nil {
			b.Fatal(err)
		} else if v != true {
			b.Fatalf("%v != true", v)
		}
	}
}

func TestMapperReadValueInt(t *testing.T) {
	var n uint8
	testMapperReadValue(t, "254", &n)
	if n != 254 {
		t.Errorf("%v != 254", n)
	}

	var un uint16
	testMapperReadValue(t, "666", &un)
	if un != 666 {
		t.Errorf("%v != 666", un)
	}
}

func BenchmarkMapperReadUint(b *testing.B) {
	r := newCyclicReader(rbEncode(b, "0xDEAD"))
	p := rmarsh.NewParser(r)
	mapper := rmarsh.NewMapper()

	var n int32

	for i := 0; i < b.N; i++ {
		n = 0
		p.Reset()

		if err := mapper.ReadValue(p, &n); err != nil {
			b.Fatal(err)
		} else if n != 0xDEAD {
			b.Fatalf("%X != 0xDEAD", n)
		}
	}
}