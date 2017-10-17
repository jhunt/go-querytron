package querytron_test

import (
	"net/url"
	"strings"
	"testing"

	qs "github.com/jhunt/go-querytron"
)

type Q map[string]string

func set(e Q) url.Values {
	q := make(url.Values)
	for name, val := range e {
		q.Set(name, val)
	}
	return q
}

type Shallow struct {
	Name    string `qs:"name"`
	Ignored string

	hidden string `qs:"secret"` /* noop */
}

type Deep struct {
	Family string `qs:"family"`
	Nested Shallow
}

type Values struct {
	Bool    bool    `qs:"some_bool"`
	Int     int     `qs:"some_int"`
	Int8    int8    `qs:"some_int_8"`
	Int16   int16   `qs:"some_int_16"`
	Int32   int32   `qs:"some_int_32"`
	Int64   int64   `qs:"some_int_64"`
	Uint    uint    `qs:"some_uint"`
	Uint8   uint8   `qs:"some_uint_8"`
	Uint16  uint16  `qs:"some_uint_16"`
	Uint32  uint32  `qs:"some_uint_32"`
	Uint64  uint64  `qs:"some_uint_64"`
	Float32 float32 `qs:"some_float_32"`
	Float64 float64 `qs:"some_float_64"`
}

func TestQuerytronShallow(t *testing.T) {
	is := func(got, expect, message string) {
		if got != expect {
			t.Errorf("%s failed - got '%s', expected '%s'\n", message, got, expect)
		}
	}

	q := set(Q{
		"name":    "overridden name",
		"secret":  "overridden secret",
		"ignored": "overridden ignored (BAD!)",
	})

	shallow := Shallow{
		Name:    "initial name",
		hidden:  "initial hidden",
		Ignored: "initial ignored",
	}
	is(shallow.Name, "initial name", "initial name is set before testing")
	is(shallow.Ignored, "initial ignored", "initial ignored is set before testing")
	is(shallow.hidden, "initial hidden", "initial hidden is set before testing")

	qs.Override(&shallow, q)
	is(shallow.Name, "overridden name", "Name is overridden from NAME env var")
	is(shallow.Ignored, "initial ignored", "initial ignored is still set")
	is(shallow.hidden, "initial hidden", "hidden fields cannot be overridden")
}

func TestQuerytronNested(t *testing.T) {
	is := func(got, expect, message string) {
		if got != expect {
			t.Errorf("%s failed - got '%s', expected '%s'\n", message, got, expect)
		}
	}

	q := set(Q{
		"family":  "overridden family",
		"name":    "overridden name",
		"secret":  "overridden secret",
		"ignored": "overridden ignored (BAD!)",
	})

	deep := Deep{
		Family: "initial family",
		Nested: Shallow{
			Name:    "initial name",
			hidden:  "initial hidden",
			Ignored: "initial ignored",
		},
	}
	is(deep.Family, "initial family", "initial family is set before testing")
	is(deep.Nested.Name, "initial name", "initial name is set before testing")
	is(deep.Nested.Ignored, "initial ignored", "initial ignored is set before testing")
	is(deep.Nested.hidden, "initial hidden", "initial hidden is set before testing")

	qs.Override(&deep, q)
	is(deep.Family, "overridden family", "Family is overridden from FAMILY env var")
	is(deep.Nested.Name, "overridden name", "Name is overridden from NAME env var")
	is(deep.Nested.Ignored, "initial ignored", "initial ignored is still set")
	is(deep.Nested.hidden, "initial hidden", "hidden fields cannot be overridden")
}

func TestQuerytronValues(t *testing.T) {
	ok := func(good bool, got, expect interface{}, message string) {
		if !good {
			t.Errorf("%s failed - got [%v], expected [%v]\n", message, got, expect)
		}
	}

	var (
		b   bool    = true
		u   uint    = 42
		u8  uint8   = 255
		u16 uint16  = 65535
		u32 uint32  = 4294967295
		u64 uint64  = 18446744073709551615
		i   int     = -42
		i8  int8    = 127
		i16 int16   = 32767
		i32 int32   = 2147483647
		i64 int64   = 9223372036854775807
		f32 float32 = 1.2345
		f64 float64 = 123456789.123456789123456789
	)

	q := set(Q{
		"some_bool":     "y",
		"some_uint":     "42",
		"some_uint_8":   "255",
		"some_uint_16":  "65535",
		"some_uint_32":  "4294967295",
		"some_uint_64":  "18446744073709551615",
		"some_int":      "-42",
		"some_int_8":    "127",
		"some_int_16":   "32767",
		"some_int_32":   "2147483647",
		"some_int_64":   "9223372036854775807",
		"some_float_32": "1.2345",
		"some_float_64": "123456789.123456789123456789",
	})

	values := Values{
		Bool:    false,
		Uint:    1,
		Uint8:   2,
		Uint16:  3,
		Uint32:  4,
		Uint64:  5,
		Int:     6,
		Int8:    7,
		Int16:   8,
		Int32:   9,
		Int64:   10,
		Float32: 11.12,
		Float64: 13.14,
	}
	ok(values.Bool == false, values.Bool, false, "initial Bool is set before testing")
	ok(values.Uint == 1, values.Uint, 1, "initial Uint is set before testing")
	ok(values.Uint8 == 2, values.Uint8, 2, "initial Uint8 is set before testing")
	ok(values.Uint16 == 3, values.Uint16, 3, "initial Uint16 is set before testing")
	ok(values.Uint32 == 4, values.Uint32, 4, "initial Uint32 is set before testing")
	ok(values.Uint64 == 5, values.Uint64, 5, "initial Uint64 is set before testing")
	ok(values.Int == 6, values.Int, 6, "initial Int is set before testing")
	ok(values.Int8 == 7, values.Int8, 7, "initial Int8 is set before testing")
	ok(values.Int16 == 8, values.Int16, 8, "initial Int16 is set before testing")
	ok(values.Int32 == 9, values.Int32, 9, "initial Int32 is set before testing")
	ok(values.Int64 == 10, values.Int64, 10, "initial Int64 is set before testing")
	ok(values.Float32 == 11.12, values.Float32, 11.12, "initial Float32 is set before testing")
	ok(values.Float64 == 13.14, values.Float64, 13.14, "initial Float64 is set before testing")

	qs.Override(&values, q)

	ok(values.Bool == b, values.Bool, b, "initial Bool is overridden from env")
	ok(values.Uint == u, values.Uint, u, "initial Uint is overridden from env")
	ok(values.Uint8 == u8, values.Uint8, u8, "initial Uint8 is overridden from env")
	ok(values.Uint16 == u16, values.Uint16, u16, "initial Uint16 is overridden from env")
	ok(values.Uint32 == u32, values.Uint32, u32, "initial Uint32 is overridden from env")
	ok(values.Uint64 == u64, values.Uint64, u64, "initial Uint64 is overridden from env")
	ok(values.Int == i, values.Int, i, "initial Int is overridden from env")
	ok(values.Int8 == i8, values.Int8, i8, "initial Int8 is overridden from env")
	ok(values.Int16 == i16, values.Int16, i16, "initial Int16 is overridden from env")
	ok(values.Int32 == i32, values.Int32, i32, "initial Int32 is overridden from env")
	ok(values.Int64 == i64, values.Int64, i64, "initial Int64 is overridden from env")
	ok(values.Float32 == f32, values.Float32, f32, "initial Float32 is overridden from env")
	ok(values.Float64 == f64, values.Float64, f64, "initial Float64 is overridden from env")
}

func TestQuerytronBools(t *testing.T) {
	is := func(got, expect bool, test, message string) {
		if got != expect {
			t.Errorf("%s ['%s' test] failed - got [%v], expected [%v]\n", message, test, got, expect)
		}
	}

	trues := strings.Split("Y y Yes YES yES true TrUe 1", " ")
	falses := strings.Split("N n No NO nO false fALse 0", " ")
	check := Values{}

	for _, yes := range trues {
		check.Bool = false
		q := set(Q{"some_bool": yes})
		is(check.Bool, false, yes, "structure bool is initially false, before override")
		qs.Override(&check, q)
		is(check.Bool, true, yes, "structure bool is overridden to be true")
	}

	for _, no := range falses {
		check.Bool = true
		q := set(Q{"some_bool": no})
		is(check.Bool, true, no, "structure bool is initially true, before override")
		qs.Override(&check, q)
		is(check.Bool, false, no, "structure bool is overridden to be false")
	}
}

type Strings struct {
	Name    string `qs:"name"`
	Version string `qs:"v"`
}

type Numbers struct {
	Number int `qs:"num"`
}

type Optional struct {
	Number *uint `qs:"opt"`
}

type Bools struct {
	Bool  *bool `qs:"bool"`
	YesNo *bool `qs:"bool:y:n"`
	TF    *bool `qs:"bool:t:f"`
	If    *bool `qs:"bool:yes"`
}

func TestQuerytronGenerate(t *testing.T) {
	is := func(got, expect string, message string) {
		if got != expect {
			t.Errorf("%s failed - got [%s], expected [%s]\n", message, got, expect)
		}
	}

	u := qs.Generate(Strings{Name: "generator"})
	is(u.Encode(), "name=generator", "string generator")

	u = qs.Generate(nil)
	is(u.Encode(), "", "nil generator")

	var nilStrings *Strings = nil
	u = qs.Generate(nilStrings)
	is(u.Encode(), "", "nil generator")

	u = qs.Generate(Strings{Version: "1.2.3"})
	is(u.Encode(), "v=1.2.3", "alternate name string generator")

	u = qs.Generate(Numbers{Number: 42})
	is(u.Encode(), "num=42", "number generator")

	u = qs.Generate(Optional{})
	is(u.Encode(), "", "number pointer generator")

	u = qs.Generate(Optional{Number: qs.Uint(42)})
	is(u.Encode(), "opt=42", "number pointer generator")

	u = qs.Generate(Bools{Bool: qs.True})
	is(u.Encode(), "bool=", "default bool generator")
	u = qs.Generate(Bools{Bool: qs.False})
	is(u.Encode(), "", "default bool generator")
	u = qs.Generate(Bools{Bool: nil})
	is(u.Encode(), "", "default bool generator")

	u = qs.Generate(Bools{YesNo: qs.True})
	is(u.Encode(), "bool=y", "y:n bool generator")
	u = qs.Generate(Bools{YesNo: qs.False})
	is(u.Encode(), "bool=n", "y:n bool generator")
	u = qs.Generate(Bools{YesNo: nil})
	is(u.Encode(), "", "y:n bool generator")

	u = qs.Generate(Bools{TF: qs.True})
	is(u.Encode(), "bool=t", "t:f bool generator")
	u = qs.Generate(Bools{TF: qs.False})
	is(u.Encode(), "bool=f", "t:f bool generator")
	u = qs.Generate(Bools{TF: nil})
	is(u.Encode(), "", "t:f bool generator")

	u = qs.Generate(Bools{If: qs.True})
	is(u.Encode(), "bool=yes", "2-arg bool generator")
	u = qs.Generate(Bools{If: qs.False})
	is(u.Encode(), "", "2-arg bool generator")
	u = qs.Generate(Bools{If: nil})
	is(u.Encode(), "", "2-arg bool generator")
}
