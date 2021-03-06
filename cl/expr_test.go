/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/
package cl_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/goplus/gop/cl/cltest"
)

// -----------------------------------------------------------------------------

func TestNew(t *testing.T) {
	cltest.Expect(t, `
		a := new([2]int)
		println("a:", a)
		`,
		"a: &[0 0]\n",
	)
	cltest.Expect(t, `
		println(new())
		`,
		"",
		"missing argument to new\n",
	)
	cltest.Expect(t, `
		println(new(int, float64))
		`,
		"",
		"too many arguments to new(int)\n",
	)
}

func TestNew2(t *testing.T) {
	cltest.Expect(t, `
		a := new([2]int)
		a[0] = 2
		println("a:", a[0])
		`,
		"a: 2\n",
	)
	cltest.Expect(t, `
		a := new([2]float64)
		a[0] = 1.1
		println("a:", a[0])
		`,
		"a: 1.1\n",
	)
	cltest.Expect(t, `
		a := new([2]string)
		a[0] = "gop"
		println("a:", a[0])
		`,
		"a: gop\n",
	)
}

func TestBadIndex(t *testing.T) {
	cltest.Expect(t, `
		a := new(int)
		println(a[0])
		`,
		"",
		nil,
	)
	cltest.Expect(t, `
		a := new(int)
		a[0] = 2
		`,
		"",
		nil,
	)
}

// -------------------`----------------------------------------------------------

func TestAutoProperty(t *testing.T) {
	script := `
		import "io"

		func New() (*Bar, error) {
			return nil, io.EOF
		}

		bar, err := New()
		if err != nil {
			log.Println(err)
		}
	`
	gopcode := `
		import (
			"github.com/goplus/gop/ast/goptest"
		)

		script := %s

		doc := goptest.New(script)!
		println(doc.any.funcDecl.name)
	`
	cltest.Expect(t,
		fmt.Sprintf(gopcode, "`"+script+"`"),
		"[New main]\n",
	)
}

func TestAutoProperty2(t *testing.T) {
	cltest.Expect(t, `
		import "bytes"
		import "os"

		b := bytes.newBuffer(nil)
		b.writeString("Hello, ")
		b.writeString("Go+")
		println(b.string)
		`,
		"Hello, Go+\n",
	)
	cltest.Expect(t, `
		import "bytes"
		import "os"

		b := bytes.newBuffer(nil)
		b.writeString("Hello, ")
		b.writeString("Go+")
		println(b.string2)
		`,
		"",
		nil, // panic
	)
	cltest.Expect(t, `
		import "bytes"

		bytes.newBuf()
		`,
		"",
		nil, // panic
	)
}

func TestUnbound(t *testing.T) {
	cltest.Expect(t, `
		println("Hello " + "qiniu:", 123, 4.5, 7i)
		`,
		"Hello qiniu: 123 4.5 (0+7i)\n",
	)
}

func TestUnboundInt(t *testing.T) {
	cltest.Expect(t, `
	import "reflect"
	printf("%T",100)
	`,
		"int",
	)
	cltest.Expect(t, `
	import "reflect"
	printf("%T",-100)
	`,
		"int",
	)
}

func TestOverflowsInt(t *testing.T) {
	cltest.Expect(t, `
	println(9223372036854775807)
	`,
		"9223372036854775807\n",
	)
	cltest.Expect(t, `
	println(-9223372036854775808)
	`,
		"-9223372036854775808\n",
	)
	cltest.Expect(t, `
	println(9223372036854775808)
	`,
		"",
		nil,
	)
}

func TestOpLAndLOr(t *testing.T) {
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return true
}
func bar() bool {
	println("bar")
	return true
}

func fake() bool {
	println("fake")
	return false
}

if foo() || bar() {
}
println("---")
if foo() && bar() {
}
println("---")
if fake() && bar() {
}
	`, "foo\n---\nfoo\nbar\n---\nfake\n")
}

func TestOpLAndLOr2(t *testing.T) {
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return true
}
func bar() bool {
	println("bar")
	return true
}

func fake() bool {
	println("fake")
	return true
}

if foo() && bar() && fake() {
}
	`, "foo\nbar\nfake\n")
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return true
}
func bar() bool {
	println("bar")
	return false
}

func fake() bool {
	println("fake")
	return true
}

if foo() && bar() && fake() {
}
	`, "foo\nbar\n")
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return false
}
func bar() bool {
	println("bar")
	return true
}

func fake() bool {
	println("fake")
	return true
}

if foo() || bar() || fake() {
}
	`, "foo\nbar\n")
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return true
}
func bar() bool {
	println("bar")
	return true
}

func fake() bool {
	println("fake")
	return true
}

if foo() || bar() || fake() {
}
	`, "foo\n")
}

func TestOpLAndLOr3(t *testing.T) {
	cltest.Expect(t, `
func foo() int {
	println("foo")
	return 0
}
func bar() bool {
	println("bar")
	return true
}
if foo() || bar() {
}
	`, "", nil)
	cltest.Expect(t, `
func foo() int {
	println("foo")
	return 0
}
func bar() bool {
	println("bar")
	return true
}
if foo() && bar() {
}
	`, "", nil)
}

func TestOpLAndLOr4(t *testing.T) {
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return true
}
if true || foo() {
}
	`, "")
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return true
}
if false || foo() {
}
	`, "foo\n")
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return true
}
if true && foo() {
}
	`, "foo\n")
	cltest.Expect(t, `
func foo() bool {
	println("foo")
	return true
}
if false && foo() {
}
	`, "")
}

func TestPanic(t *testing.T) {
	cltest.Expect(t,
		`panic("Helo")`,
		"",
		"Helo", // panicMsg
	)
}

func TestTakeAddrMap(t *testing.T) {
	cltest.Expect(t, `
		m := {1:"hello",2:"ok"}
		println(m)
		println(&m)
		`,
		"map[1:hello 2:ok]\n&map[1:hello 2:ok]\n")
}

func TestTakeAddrMapIndexBad(t *testing.T) {
	cltest.Expect(t, `
		m := {1:"hello",2:"ok"}
		println(&m[1])
		`,
		"",
		"cannot take the address of m[1]\n")
}

func TestTakeAddrStringBad(t *testing.T) {
	cltest.Expect(t, `
		m := "hello"
		println(&m[1])
		`,
		"",
		"cannot take the address of m[1]\n")
}

func TestTypeCast(t *testing.T) {
	cltest.Call(t, `
	x := []byte("hello")
	x
	`).Equal([]byte("hello"))
}

func TestPkgTypeConv(t *testing.T) {
	cltest.Expect(t, `
	import "sort"
	ar := [1,5,3,2]
	sort.IntSlice(ar).Sort()
	println(ar)
	`, "[1 2 3 5]\n")
}

func TestRuneType(t *testing.T) {
	cltest.Expect(t, `
	a := 'a'
	printf("%T\n",a)
	`, "int32\n")
	cltest.Expect(t, `
	printf("%T\n",'a')
	`, "int32\n")
}

func TestAppendErr(t *testing.T) {
	cltest.Expect(t, `
		append()
		`,
		"",
		"append: argument count not enough\n",
	)
	cltest.Expect(t, `
		x := 1
		append(x, 2)
		`,
		"",
		"append: first argument not a slice\n",
	)
	cltest.Expect(t, `
		defer append([]int{1}, 2)
		`,
		"",
		"defer discards result of append([]int{1}, 2)\n",
	)
}

func TestLenErr(t *testing.T) {
	cltest.Expect(t, `
		len()
		`,
		"",
		"missing argument to len: len()\n",
	)
	cltest.Expect(t, `
		len("a", "b")
		`,
		"",
		`too many arguments to len: len("a", "b")`+"\n",
	)
}

func TestMake(t *testing.T) {
	cltest.Expect(t, `
		make()
		`,
		"",
		"missing argument to make: make()\n",
	)
	cltest.Expect(t, `
		a := make([]int, 0, 4)
		a = append(a, 1, 2, 3)
		println(a)
		`,
		"[1 2 3]\n",
	)
	cltest.Expect(t, `
		a := make([]int, 0, 4)
		a = append(a, [1, 2, 3]...)
		println(a)
		`,
		"[1 2 3]\n",
	)
	cltest.Expect(t, `
		n := 4
		a := make(map[string]interface{}, uint16(n))
		println(a)
		`,
		"map[]\n",
	)
	cltest.Expect(t, `
		import "reflect"

		a := make(chan *func(), uint16(4))
		println(reflect.TypeOf(a))
		`,
		"chan *func()\n",
	)
}

func TestOperator(t *testing.T) {
	cltest.Expect(t, `
		println("Hello", 123 * 4.5, 1 + 7i)
		`,
		"Hello 553.5 (1+7i)\n")
}

func TestVar(t *testing.T) {
	cltest.Expect(t, `
		x := 123.1
		println("Hello", x)
		`,
		"Hello 123.1\n")
}

func TestVarOp(t *testing.T) {
	cltest.Expect(t, `
		x := 123.1
		y := 1 + x
		println("Hello", y + 10)
		n, err := println("Hello", y + 10)
		println("ret:", n << 1, err)
		`,
		"Hello 134.1\nHello 134.1\nret: 24 <nil>\n",
	)
}

func TestGoPackage(t *testing.T) {
	cltest.Expect(t, `
		import "fmt"
		import gostrings "strings"

		x := gostrings.NewReplacer("?", "!").Replace("hello, world???")
		fmt.Println("x: " + x)
		`,
		"x: hello, world!!!\n",
	)
}

func TestSlice(t *testing.T) {
	cltest.Expect(t, `
		x := []float64{1, 2.3, 3.6}
		println("x:", x)
		`,
		"x: [1 2.3 3.6]\n",
	)
	cltest.Expect(t, `
		x := []float64{1, 2: 3.4, 5}
		println("x:", x)
		`,
		"x: [1 0 3.4 5]\n",
	)
	cltest.Expect(t, `
		x := []float64{1:1,3:3,4}
		println("x:", x)
		y := []float64{3:3,4,1:1}
		println("y:", x)
		`,
		"x: [0 1 0 3 4]\ny: [0 1 0 3 4]\n",
	)
}

func TestArray(t *testing.T) {
	cltest.Expect(t, `
		x := [4]float64{1, 2.3, 3.6}
		println("x:", x)
		y := [...]float64{1, 2.3, 3.6}
		println("y:", y)
		`,
		"x: [1 2.3 3.6 0]\ny: [1 2.3 3.6]\n",
	)
	cltest.Expect(t, `
		x := [5]float64{1:1,3:3,4}
		println("x:", x)
		y := [5]float64{3:3,4,1:1}
		println("y:", x)
		`,
		"x: [0 1 0 3 4]\ny: [0 1 0 3 4]\n",
	)
	cltest.Expect(t, `
		x := [...]float64{1, 3: 3.4, 5}
		x[1] = 217
		println("x:", x, "x[1]:", x[1])
		`,
		"x: [1 217 0 3.4 5] x[1]: 217\n",
	)
	cltest.Expect(t, `
		x := [...]float64{1, 2.3, 3, 4}
		x[2] = 3.1
		println("x[1:]:", x[1:])
		println(len(x))
	`,
		"x[1:]: [2.3 3.1 4]\n4\n")
}

func TestLoadVar(t *testing.T) {
	cltest.Expect(t, `
		var x1 int
		var x2 int = 10
		var x3 = 10
		println("x:",x1,x2,x3)
		`,
		"x: 0 10 10\n")
	cltest.Expect(t, `
		type Point struct {
			X int
			Y int
		}
		var x1 Point
		var x2 Point = Point{10,20}
		var x3 = Point{-10,-20}
		println("x:",x1,x2,x3)
		`,
		"x: {0 0} {10 20} {-10 -20}\n")
}

func TestLoadVar2(t *testing.T) {
	cltest.Expect(t, `
		func main() {
			var x1 int
			var x2 int = 10
			var x3 = 10
			println("x:",x1,x2,x3)
		}`,
		"x: 0 10 10\n")
	cltest.Expect(t, `
		type Point struct {
			X int
			Y int
		}
		func main() {
			var x1 Point
			var x2 Point = Point{10,20}
			var x3 = Point{-10,-20}
			println("x:",x1,x2,x3)
		}`,
		"x: {0 0} {10 20} {-10 -20}\n")

	cltest.Expect(t, `
		func main() {
			type Point struct {
				X int
				Y int
			}
			var x1 Point
			var x2 Point = Point{10,20}
			var x3 = Point{-10,-20}
			println("x:",x1,x2,x3)
		}`,
		"x: {0 0} {10 20} {-10 -20}\n")
}

func TestLoadVar3(t *testing.T) {
	cltest.Expect(t, `
		var x1,x2,x3 = 10,20,x1+x2
		println("x:",x1,x2,x3)
		`,
		"x: 10 20 30\n")
	cltest.Expect(t, `
		var x1,x2,x3 = 10,20,x1+x2
		func main() {
			println("x:",x1,x2,x3)
		}
		`,
		"x: 10 20 30\n")
}

func TestMap(t *testing.T) {
	cltest.Expect(t, `
		x := map[string]float64{"Hello": 1, "xsw": 3.4}
		println("x:", x)
		`,
		"x: map[Hello:1 xsw:3.4]\n")
}

func TestMapLit(t *testing.T) {
	cltest.Expect(t, `
		x := {"Hello": 1, "xsw": 3.4}
		println("x:", x)
		`,
		"x: map[Hello:1 xsw:3.4]\n",
	)
	cltest.Expect(t, `
		x := {"Hello": 1, "xsw": "3.4"}
		println("x:", x)

		println("empty map:", {})
		`,
		"x: map[Hello:1 xsw:3.4]\nempty map: map[]\n",
	)
}

func TestMapIdx(t *testing.T) {
	cltest.Expect(t, `
		x := {"Hello": 1, "xsw": "3.4"}
		y := {1: "glang", 5: "Hi"}
		i := 1
		q := "Q"
		key := "xsw"
		x["xsw"], y[i] = 3.1415926, q
		println("x:", x, "y:", y)
		println("x[key]:", x[key], "y[1]:", y[1])
		`,
		"x: map[Hello:1 xsw:3.1415926] y: map[1:Q 5:Hi]\nx[key]: 3.1415926 y[1]: Q\n",
	)
}

func TestSliceLit(t *testing.T) {
	cltest.Expect(t, `
		x := [1, 3.4]
		println("x:", x)

		y := [1]
		println("y:", y)

		z := [1+2i, "xsw"]
		println("z:", z)

		println("empty slice:", [])
		`,
		"x: [1 3.4]\ny: [1]\nz: [(1+2i) xsw]\nempty slice: []\n")
}

func TestSliceIdx(t *testing.T) {
	cltest.Expect(t, `
		x := [1, 3.4, 17]
		n, m := 1, uint16(0)
		x[1] = 32.7
		x[m] = 36.86
		println("x:", x[2], x[m], x[n])
		`,
		"x: 17 36.86 32.7\n")
}

func TestListComprehension(t *testing.T) {
	cltest.Expect(t, `
		y := [i+x for i, x <- [1, 2, 3, 4]]
		println("y:", y)
		`,
		"y: [1 3 5 7]\n")
	cltest.Call(t, `
		y := [i+x for i, x <- {3: 1, 5: 2, 7: 3, 11: 4}]
		println("y:", y)
		`, -2).Equal(15)
	cltest.Call(t, `
		y := [i+x for i, x <- {3: 1, 5: 2, 7: 3, 11: 4}, x % 2 == 1]
		println("y:", y)
		`, -2).Equal(10)
}

func TestMapComprehension(t *testing.T) {
	cltest.Expect(t, `
		y := {x: i for i, x <- [3, 5, 7, 11, 13]}
		println("y:", y)
		`,
		"y: map[3:0 5:1 7:2 11:3 13:4]\n",
	)
	cltest.Expect(t, `
		y := {x: i for i, x <- [3, 5, 7, 11, 13], i % 2 == 1}
		println("y:", y)
		`,
		"y: map[5:1 11:3]\n",
	)
	cltest.Expect(t, `
		y := {v: k for k, v <- {"Hello": "xsw", "Hi": "glang"}}
		println("y:", y)
		`,
		"y: map[glang:Hi xsw:Hello]\n",
	)
	cltest.Expect(t, `
		println({x: i for i, x <- [3, 5, 7, 11, 13]})
		println({x: i for i, x <- [3, 5, 7, 11, 13]})
		`,
		"map[3:0 5:1 7:2 11:3 13:4]\nmap[3:0 5:1 7:2 11:3 13:4]\n",
	)
	cltest.Expect(t, `
		arr := [1, 2, 3, 4, 5, 6]
		x := [[a, b] for a <- arr, a < b for b <- arr, b > 2]
		println("x:", x)
		`,
		"x: [[1 3] [2 3] [1 4] [2 4] [3 4] [1 5] [2 5] [3 5] [4 5] [1 6] [2 6] [3 6] [4 6] [5 6]]\n")
}

func TestErrWrapExpr(t *testing.T) {
	cltest.Call(t, `
		x := println("Hello qiniu")!
		x
		`).Equal(12)
	cltest.Call(t, `
		import (
			"strconv"
		)
	
		func add(x, y string) (int, error) {
			return strconv.Atoi(x)? + strconv.Atoi(y)?, nil
		}
	
		x := add("100", "23")!
		x
		`).Equal(123)
}

func TestRational(t *testing.T) {
	cltest.Call(t, `
		x := 3/4r + 5/7r
		x
	`).Equal(big.NewRat(41, 28))
	cltest.Call(t, `
		a := 3/4r
		x := a + 5/7r
		x
	`).Equal(big.NewRat(41, 28))
	y, _ := new(big.Float).SetString(
		"3.14159265358979323846264338327950288419716939937510582097494459")
	y.Mul(y, big.NewFloat(2))
	cltest.Call(t, `
		y := 3.14159265358979323846264338327950288419716939937510582097494459r
		y *= 2
		y
	`).Equal(y)
	cltest.Call(t, `
		a := 3/4r
		b := 5/7r
		if a > b {
			a = a + 1
		}
		a
	`).Equal(big.NewRat(7, 4))
	cltest.Call(t, `
		x := 1/3r + 1r*2r
		x
	`).Equal(big.NewRat(7, 3))
}

// TODO #575
func _TestIsNoExecCtx(t *testing.T) {
	cltest.Expect(t, `
	fns := make([]func() int, 3)
	for i, x <- [3, 15, 777] {
		var v = x
		var fn = func() int {
			return v
		}
		fns[i] = fn
	}
	println("values:", fns[0](), fns[1](), fns[2]())`, "values: 3 15 777\n")
}

func TestPkgMethod(t *testing.T) {
	cltest.Expect(t, `
	import "bytes"
	buf := bytes.NewBuffer([]byte("hello"))
	println(buf.String())
	`, "hello\n")
	cltest.Expect(t, `
	import "bytes"
	var buf bytes.Buffer
	buf.Write([]byte("hello"))
	println(buf.String())
	`, "hello\n")
	cltest.Expect(t, `
	import "reflect"
	v := reflect.ValueOf(100)
	println(v.Kind())
	`, "int\n")
	cltest.Expect(t, `
	import "reflect"
	v := reflect.ValueOf(100)
	p := &v
	println(p.Kind())
	`, "int\n")
}

func TestPkgMethodBadCall(t *testing.T) {
	cltest.Expect(t, `
	import "bytes"
	buf := bytes.NewBuffer([]byte("hello"))
	println((&buf).String())
	`, "", "calling method String with receiver &buf (type **bytes.Buffer) requires explicit dereference.")
	cltest.Expect(t, `
	import "reflect"
	v := reflect.ValueOf(100)
	p := &v
	println((&p).Kind())
	`, "", "calling method Kind with receiver &p (type **reflect.Value) requires explicit dereference.")
}

func TestComplex(t *testing.T) {
	cltest.Expect(t, `
	c := complex(1,2)
	printf("%v %T\n",c,c)
	`, "(1+2i) complex128\n")
	cltest.Expect(t, `
	c := complex(float64(1),2)
	printf("%v %T\n",c,c)
	`, "(1+2i) complex128\n")
	cltest.Expect(t, `
	c := complex(float32(1),2)
	printf("%v %T\n",c,c)
	`, "(1+2i) complex64\n")
	cltest.Expect(t, `
	func test() float64 { return 1 }
	c := complex(test(),2)
	printf("%v %T\n",c,c)
	`, "(1+2i) complex128\n")
	cltest.Expect(t, `
	func test() float32 { return 1 }
	c := complex(test(),2)
	printf("%v %T\n",c,c)
	`, "(1+2i) complex64\n")

	cltest.Expect(t, `
	c := real(1+2i)
	printf("%v %T\n",c,c)
	`, "1 float64\n")
	cltest.Expect(t, `
	c := real(complex128(1+2i))
	printf("%v %T\n",c,c)
	`, "1 float64\n")
	cltest.Expect(t, `
	c := real(complex64(1+2i))
	printf("%v %T\n",c,c)
	`, "1 float32\n")
	cltest.Expect(t, `
	c := real(complex(1,2))
	printf("%v %T\n",c,c)
	`, "1 float64\n")
	cltest.Expect(t, `
	c := real(complex(1,float32(2)))
	printf("%v %T\n",c,c)
	`, "1 float32\n")

	cltest.Expect(t, `
	c := imag(1+2i)
	printf("%v %T\n",c,c)
	`, "2 float64\n")
	cltest.Expect(t, `
	c := imag(complex128(1+2i))
	printf("%v %T\n",c,c)
	`, "2 float64\n")
	cltest.Expect(t, `
	c := imag(complex64(1+2i))
	printf("%v %T\n",c,c)
	`, "2 float32\n")
	cltest.Expect(t, `
	c := imag(complex(1,2))
	printf("%v %T\n",c,c)
	`, "2 float64\n")
	cltest.Expect(t, `
	c := imag(complex(float32(1),2))
	printf("%v %T\n",c,c)
	`, "2 float32\n")
}

func TestBadComplex(t *testing.T) {
	cltest.Expect(t, `
	complex(1)
	`, "", nil)
	cltest.Expect(t, `
	complex(1,2,3)
	`, "", nil)
	cltest.Expect(t, `
	complex(float32(1),float64(2))
	`, "", nil)
	cltest.Expect(t, `
	func test() int { return 100 }
	complex(test(),2)
	`, "", nil)
	cltest.Expect(t, `
	complex(1,int(2))
	`, "", nil)

	cltest.Expect(t, `
	real()
	`, "", nil)
	cltest.Expect(t, `
	real(1,2)
	`, "", nil)
	cltest.Expect(t, `
	real(int(1))
	`, "", nil)

	cltest.Expect(t, `
	imag()
	`, "", nil)
	cltest.Expect(t, `
	imag(1,2)
	`, "", nil)
	cltest.Expect(t, `
	imag(int(1))
	`, "", nil)
}

func TestResult(t *testing.T) {
	cltest.Expect(t, `
	import "fmt"
	type Writer struct {
	}
	func (w *Writer) Write(data string) (n int, err error) {
		return fmt.Println(data)
	}
	w := &Writer{}
	n, err := w.Write("hello")
	println(n,err)
	`, "hello\n6 <nil>\n")
	cltest.Expect(t, `
	import "fmt"
	type Writer struct {
	}
	func (w *Writer) Write(data string) (int, error) {
		fmt.Println(data)
		return len(data)+1,nil
	}
	w := &Writer{}
	n, err := w.Write("hello")
	println(n,err)
	`, "hello\n6 <nil>\n")
	cltest.Expect(t, `
	import "fmt"
	type Writer struct {
	}
	func myint(n int) int {
		return n
	}
	func myerr(err error) error {
		return err
	}
	func (w *Writer) Write(data string) (int, error) {
		n, err := fmt.Println(data)
		return myint(n),myerr(err)
	}
	w := &Writer{}
	n, err := w.Write("hello")
	println(n,err)
	`, "hello\n6 <nil>\n")
}

func TestBadResult(t *testing.T) {
	cltest.Expect(t, `
	import "fmt"
	type Writer struct {
	}
	func (w *Writer) Write(data string) (error) {
		err := fmt.Println(data)
		return err
	}
	w := &Writer{}
	n, err := w.Write("hello")
	println(n,err)
	`, "", nil)
	cltest.Expect(t, `
	import "fmt"
	type Writer struct {
	}
	func (w *Writer) Write(data string) (err error) {
		return fmt.Println(data)
	}
	w := &Writer{}
	n, err := w.Write("hello")
	println(n,err)
	`, "", nil)
}

func TestUnderscore(t *testing.T) {
	cltest.Expect(t, `
	import "fmt"
	var _ int
	var _ string
	_ = 100
	_ = "hello"
	_, a := 100,"world"
	b, _ := fmt.Println("Hello World")
	println(a,b)
	`, "Hello World\nworld 12\n")
}

func TestBadUnderscore(t *testing.T) {
	cltest.Expect(t, `
	println(_)
	`, "", nil)
	cltest.Expect(t, `
	_ := 100
	`, "", nil)
	cltest.Expect(t, `
	_,_ := 100,"hello"
	`, "", nil)
	cltest.Expect(t, `
	import "fmt"
	_, _ := fmt.Println("Hello World")
	`, "", nil)
}

func TestBadVar(t *testing.T) {
	cltest.Expect(t, `
	var a int
	var a string`, "", nil)
	cltest.Expect(t, `
	a = 10`, "", nil)
}

type testData struct {
	clause string
	want   string
	panic  bool
}

var testDeleteClauses = map[string]testData{
	"delete_int_key": {`
					m:={1:1,2:2}
					delete(m,1)
					println(m)
					delete(m,3)
					println(m)
					delete(m,2)
					println(m)
					`, "map[2:2]\nmap[2:2]\nmap[]\n", false},
	"delete_string_key": {`
					m:={"hello":1,"Go+":2}
					delete(m,"hello")
					println(m)
					delete(m,"hi")
					println(m)
					delete(m,"Go+")
					println(m)
					`, "map[Go+:2]\nmap[Go+:2]\nmap[]\n", false},
	"delete_var_string_key": {`
					m:={"hello":1,"Go+":2}
					delete(m,"hello")
					println(m)
					a:="hi"
					delete(m,a)
					println(m)
					arr:=["Go+"]
					delete(m,arr[0])
					println(m)
					`, "map[Go+:2]\nmap[Go+:2]\nmap[]\n", false},
	"delete_var_map_string_key": {`
					ma:=[{"hello":1,"Go+":2}]
					delete(ma[0],"hello")
					println(ma[0])
					a:="hi"
					delete(ma[0],a)
					println(ma[0])
					arr:=["Go+"]
					delete(ma[0],arr[0])
					println(ma[0])
					`, "map[Go+:2]\nmap[Go+:2]\nmap[]\n", false},
	"delete_no_key_panic": {`
					m:={"hello":1,"Go+":2}
					delete(m)
					`, "", true},
	"delete_multi_key_panic": {`
					m:={"hello":1,"Go+":2}
					delete(m,"hi","hi")
					`, "", true},
	"delete_not_map_panic": {`
					m:=[1,2,3]
					delete(m,1)
					`, "", true},
}

func TestDelete(t *testing.T) {
	testScripts(t, "TestDelete", testDeleteClauses)
}

// -----------------------------------------------------------------------------

var testCopyClauses = map[string]testData{
	"copy_int": {`
					a:=[1,2,3]
					b:=[4,5,6]
					n:=copy(b,a)
					println(n)
					println(b)
					`, "3\n[1 2 3]\n", false},
	"copy_string": {`
					a:=["hello"]
					b:=["hi"]
					n:=copy(b,a)
					println(n)
					println(b)
					`, "1\n[hello]\n", false},
	"copy_byte_string": {`
					a:=[byte(65),byte(66),byte(67)]
					println(string(a))
					n:=copy(a,"abc")
					println(n)
					println(a)
					println(string(a))
					`, "ABC\n3\n[97 98 99]\nabc\n", false},
	"copy_first_not_slice_panic": {`
					a:=1
					b:=[1,2,3]
					copy(a,b)
					println(a)
					`, "", true},
	"copy_second_not_slice_panic": {`
					a:=1
					b:=[1,2,3]
					copy(b,a)
					println(b)
					`, "", true},
	"copy_one_args_panic": {`
					a:=[1,2,3]
					copy(a)
					println(a)
					`, "", true},
	"copy_multi_args_panic": {`
					a:=[1,2,3]
					copy(a,a,a)
					println(a)
					`, "", true},
	"copy_string_panic": {`
					a:=[65,66,67]
					copy(a,"abc")
					println(a)
					`, "", true},
	"copy_different_type_panic": {`
					a:=[65,66,67]
					b:=[1.2,1.5,1.7]
					copy(b,a)
					copy(b,a)
					println(b)
					`, "", true},
	"copy_with_operation": {`
					a:=[65,66,67]
					b:=[1]
					println(copy(a,b)+copy(b,a)==2)
					`, "true\n", false},
}

func TestCopy(t *testing.T) {
	testScripts(t, "TestCopy", testCopyClauses)
}

var testStructClauses = map[string]testData{
	"struct": {`
			println(struct {
				A int
				B string
			}{1, "Hello"})	
					`, "{1 Hello}\n", false},
	"struct_key_value": {`
			println(struct {
				A int
				B string
			}{A:1,B: "Hello"})	
					`, "{1 Hello}\n", false},
	"struct_ptr": {`
			println(&struct {
				A int
				B string
			}{1, "Hello"})
					`, "&{1 Hello}\n", false},
	"struct_key_value_ptr": {`
			println(&struct {
				A int  ` + "`json:\"a\"`" + `
				B string
			}{A: 1,B: "Hello"})
					`, "&{1 Hello}\n", false},
	"struct_key_value_ptr_unexport_field": {`
			println(&struct {
				a int  ` + "`json:\"a\"`" + `
				b string
			}{a: 1,b: "Hello"})
					`, "&{1 Hello}\n", false},
	"struct_key_value_unexport_field": {`
			println(struct {
				a int  ` + "`json:\"a\"`" + `
				b string
			}{a: 1,b: "Hello"})
					`, "{1 Hello}\n", false},
	"struct_unexport_field": {`
			println(struct {
				a int
				b string
			}{1, "Hello"})	
					`, "{1 Hello}\n", false},
	"struct_ptr_unexport_field": {`
			println(&struct {
				a int
				b string
			}{1, "Hello"})	
					`, "&{1 Hello}\n", false},
	"struct_store_field_panic": {`
				import "sync"

				mu := sync.WaitGroup{}
				
				mu.noCopy = struct{}{}
					`, "", true},
	"struct_array": {`
	type Point struct {
		X int
		Y int
	}
	ar := []Point{}
	ar = append(ar,Point{10,20})
	println(ar)
	`, "[{10 20}]\n", false},
	"struct_ptr_array": {`
	type Point struct {
		X int
		Y int
	}
	ar := []*Point{}
	ar = append(ar,&Point{10,20})
	println(ar[0])
	`, "&{10 20}\n", false},
}

func TestStruct2(t *testing.T) {
	testScripts(t, "TestStruct", testStructClauses)
}

// -----------------------------------------------------------------------------
var testMethodClauses = map[string]testData{
	"method set": {`
					type Person struct {
						Name string
						Age  int
					}
					func (p *Person) SetName(name string) {
						p.Name = name
					}

					p := &Person{
						Name: "bar",
						Age:  30,
					}

					p.SetName("foo")
					println(p.Name)
					`, "foo\n", false},
	"method get": {`
					type Person struct {
						Name string
						Age  int
					}
					func (p *Person) GetName() string {
						return p.Name
					}

					p := &Person{
						Name: "bar",
						Age:  30,
					}

					println(p.GetName())
					`, "bar\n", false},

	"struct set ptr": {`
	type Person struct {
		Name string
		Age  int
	}

	p := &Person{
		Name: "bar",
		Age:  30,
	}
	p.Name = "foo"

	println(p)
	`, "&{foo 30}\n", false},

	"struct set": {`
	type Person struct {
		Name string
		Age  int
	}

	p := Person{
		Name: "bar",
		Age:  30,
	}
	p.Name = "foo"

	println(p)
	`, "{foo 30}\n", false},

	"struct set ptr arg": {`
	type Person struct {
		Name string
		Age  int
	}
	func SetName(p *Person,name string) {
		p.Name = name
	}

	p := Person{
		Name: "bar",
		Age:  30,
	}
	SetName(&p,"foo")

	println(p)
	`, "{foo 30}\n", false},

	"method func no args": {`
					type Person struct {
						Name string ` + "`json:\"name\"`" + `
						Age  int
					}
					func (p *Person) PrintName() {
						println(p.Name)
					}

					p := &Person{
						Name: "bar",
						Age:  30,
					}

					p.PrintName()
					`, "bar\n", false},
	"method ptr struct no prt": {`
					type Person struct {
						Name string ` + "`json:\"name\"`" + `
						Age  int
					}
					func (p *Person) PrintName() {
						println(p.Name)
					}

					p := Person{
						Name: "bar",
						Age:  30,
					}

					p.PrintName()
					`, "bar\n", false},

	"method load field": {`
					type Person struct {
						Name string
						Age  int
					}
					func (p *Person) SetName(name string,age int) {
						p.Name = name
						p.Age = age
						println(name)
						println(p.Age)
					}

					p := Person{
						Name: "bar",
						Age:  30,
					}

					p.SetName("foo",31)
					`, "foo\n31\n", false},
	"method int type": {`
					
					type M int

					func (m M) Foo() {
						println("foo", m)
					}

					m := M(0)
					m.Foo()
					println(m)
					`, "foo 0\n0\n", false},
}

func TestMethodCases(t *testing.T) {
	testScripts(t, "TestMethod", testMethodClauses)
}

func TestEmbeddedField(t *testing.T) {
	cltest.Expect(t, `
	type Base struct {
		Info string
	}
	type Point struct {
		X int
		Y int
	}
	type My struct {
		Base
		Point
	}
	m := &My{Base:Base{"hello"},Point{10,20}}
	println(m.Info,m.X,m.Y)
	m.Info = "world"
	m.Point = Point{-10,-20}
	println(m.Info,m.X,m.Y)
	m.Base = Base{"goplus"}
	m.Point.X = 100
	m.Y = 200
	println(m.Info,m.X,m.Y)
	`, "hello 10 20\nworld -10 -20\ngoplus 100 200\n")
	cltest.Expect(t, `
	type Base struct {
		Info string
	}
	type Point struct {
		X int
		Y int
	}
	type My struct {
		Base
		*Point
	}
	m := &My{Base:Base{"hello"},&Point{10,20}}
	println(m.Info,m.X,m.Y)
	m.Info = "world"
	m.Point = &Point{-10,-20}
	println(m.Info,m.X,m.Y)
	m.Base = Base{"goplus"}
	m.Point.X = 100
	m.Y = 200
	println(m.Info,m.X,m.Y)
	`, "hello 10 20\nworld -10 -20\ngoplus 100 200\n")
	cltest.Expect(t, `
	import "bytes"
	type Buf struct {
		*bytes.Buffer
	}
	buf := &Buf{bytes.NewBufferString("hello")}
	println(buf)
	`, "&{hello}\n")
	cltest.Expect(t, `
	import "reflect"
	type Value struct {
		reflect.Value
	}
	v := Value{reflect.ValueOf(100)}
	println(v.Value)
	`, "100\n")
}

func TestEmbeddedMethod(t *testing.T) {
	cltest.Expect(t, `
	import "bytes"
	type Buf struct {
		*bytes.Buffer
	}
	buf := &Buf{&bytes.Buffer{}}
	buf.Write([]byte("hello"))
	println(buf.String())
	`, "hello\n")
	cltest.Expect(t, `
	import "bytes"
	type Buf struct {
		*bytes.Buffer
	}
	buf := Buf{&bytes.Buffer{}}
	buf.Write([]byte("hello"))
	println(buf.String())
	`, "hello\n")
	cltest.Expect(t, `
	import "bytes"
	type Buf struct {
		*bytes.Buffer
		size int
	}
	buf := Buf{&bytes.Buffer{},1}
	buf.Write([]byte("hello"))
	println(buf)
	`, "{hello 1}\n")
	cltest.Expect(t, `
	import "reflect"
	type Value struct {
		reflect.Value
	}
	v := Value{reflect.ValueOf(100)}
	println(v.Kind())
	`, "int\n")
	cltest.Expect(t, `
	import "reflect"
	type Value struct {
		reflect.Value
	}
	v := &Value{reflect.ValueOf(100)}
	println(v.Kind())
	`, "int\n")
	cltest.Expect(t, `
	type Point struct {
		X int
		Y int
	}
	func (p Point) Test() {
		println(p.X,p.Y)
	}
	type My struct {
		Point
	}
	m := &My{Point{10,20}}
	m.Test()
	`, "10 20\n")
	cltest.Expect(t, `
	type Point struct {
		X int
		Y int
	}
	func (p Point) Test() {
		println(p.X,p.Y)
	}
	type My struct {
		Point
	}
	m := My{Point{10,20}}
	m.Test()
	`, "10 20\n")
	cltest.Expect(t, `
	type Point struct {
		X int
		Y int
	}
	func (p *Point) Test() {
		println(p.X,p.Y)
	}
	type My struct {
		Point
	}
	m := &My{Point{10,20}}
	m.Test()
	`, "10 20\n")
	cltest.Expect(t, `
	type Point struct {
		X int
		Y int
	}
	func (p *Point) Test() {
		println(p.X,p.Y)
	}
	type My struct {
		*Point
	}
	m := &My{&Point{10,20}}
	m.Test()
	`, "10 20\n")
	cltest.Expect(t, `
	type Point struct {
		X int
		Y int
	}
	func (p *Point) Test() {
		println(p.X,p.Y)
	}
	type Base struct {
		*Point
	}
	type My struct {
		Base
	}
	m := &My{}
	m.Point = &Point{10,20}
	m.Test()
	`, "10 20\n")
	cltest.Expect(t, `
	type Point struct {
		X int
		Y int
	}
	func (p *Point) Test() {
		println(p.X,p.Y)
	}
	type Base struct {
		*Point
	}
	type My struct {
		*Base
		size int
	}
	m := &My{&Base{&Point{10,20}},1}
	m.Test()
	`, "10 20\n")
}

// -----------------------------------------------------------------------------

var testStarExprClauses = map[string]testData{
	"star expr": {`
				func A(a *int, c *struct {
					b *int
					m map[string]*int
					s []*int
				}) {
					*a = 5
					*c.b = 3
					*c.m["foo"] = 7
					*c.s[0] = 9
				}

				a1 := 6
				a2 := 6
				a3 := 6
				c := struct {
					b *int
					m map[string]*int
					s []*int
				}{
					b: &a1,
					m: map[string]*int{
						"foo": &a2,
					},
					s: []*int{&a3},
				}
				A(&a1, &c)
				*c.m["foo"] = 8
				*c.s[0] = 10
				*c.s[0+0] = 10
				println(a1, *c.b, *c.m["foo"], *c.s[0], *c.s[0+0])

					`, "3 3 8 10 10\n", false},
	"star expr exec": {`
					func A(a *int, c *struct {
						b *int
						m map[string]*int
						s []*int
					}) {
						*a = 5
						*c.b = 3
						*c.m["foo"] = 7
						*c.s[0] = 9
					}
	
					func main() {
						a1 := 6
						a2 := 6
						a3 := 6
						c := struct {
							b *int
							m map[string]*int
							s []*int
						}{
							b: &a1,
							m: map[string]*int{
								"foo": &a2,
							},
							s: []*int{&a3},
						}
						A(&a1, &c)
						*c.m["foo"] = 8
						*c.s[0] = 10
						*c.s[0+0] = 10
						println(a1, *c.b, *c.m["foo"], *c.s[0], *c.s[0+0])
					}
						`, "3 3 8 10 10\n", false},
	"star expr lhs slice index func": {`
					func A(a *int, c *struct {
						b *int
						m map[string]*int
						s []*int
					}) {
						*a = 5
						*c.b = 3
						*c.m["foo"] = 7
						*c.s[0] = 9
					}
					
					func Index() int {
						return 0
					}
					
					a1 := 6
					a2 := 6
					a3 := 6
					c := struct {
						b *int
						m map[string]*int
						s []*int
					}{
						b: &a1,
						m: map[string]*int{
							"foo": &a2,
						},
						s: []*int{&a3},
					}
					A(&a1, &c)
					*c.m["foo"] = 8
					*c.s[0] = 10
					*c.s[Index()] = 11
					println(a1, *c.b, *c.m["foo"], *c.s[0])
	
						`, "3 3 8 11\n", false},
	"start expr ptr conv": {`
					a := 10
					println(*(*int)(&a))
					`, "10\n", false},
}

func TestStarExpr(t *testing.T) {
	testScripts(t, "TestStarExpr", testStarExprClauses)
}

// -----------------------------------------------------------------------------
var testRefTypeClauses = map[string]testData{
	"ref type": {`
	func foo() []int {
		return make([]int, 10)
	}
	
	func foo1() map[int]int {
		return make(map[int]int, 10)
	}
	
	func foo2() chan int {
		return make(chan int, 10)
	}
	a := foo()
	if a != nil {
		println("foo")
	}
	
	a1 := foo1()
	if a1 != nil {
		println("foo1")
	}
	a2 := foo2()
	if a2 != nil {
		println("foo2")
	}
						`, "foo\nfoo1\nfoo2\n", false},
	"ref type 2": {`
	func foo() []int {
		return nil
	}
	
	func foo1() map[int]int {
		return make(map[int]int, 10)
	}
	
	func foo2() chan int {
		return make(chan int, 10)
	}

	func foo3() *int {
		return nil
	}
	
	println(foo() == nil)
	println(nil == foo())
	println(foo() != nil)
	println(nil != foo())
	
	println(foo1() == nil)
	println(nil == foo1())
	println(foo1() != nil)
	println(nil != foo1())
	
	println(foo2() == nil)
	println(nil == foo2())
	println(foo2() != nil)
	println(nil != foo2())
	
	println(foo3() == nil)
	println(nil == foo3())
	println(foo3() != nil)
	println(nil != foo3())
						`, "true\ntrue\nfalse\nfalse\nfalse\nfalse\ntrue\ntrue\nfalse\nfalse\ntrue\ntrue\ntrue\ntrue\nfalse\nfalse\n", false},
}

func TestRefType(t *testing.T) {
	testScripts(t, "TestRefType", testRefTypeClauses)
}

func TestMatchType(t *testing.T) {
	cltest.Expect(t, `
		println(nil == nil,nil != nil)
	`, "true false\n")
	cltest.Expect(t, `
		var i interface{}
		println(i == nil,i != nil)
	`, "true false\n")
	cltest.Expect(t, `
		var i *int
		println(i == nil,i != nil)
	`, "true false\n")
	cltest.Expect(t, `
		var i chan int
		println(i == nil,i != nil)
	`, "true false\n")
	cltest.Expect(t, `
		var i func()
		println(i == nil,i != nil)
	`, "true false\n")
	cltest.Expect(t, `
		var i []int
		println(i == nil,i != nil)
	`, "true false\n")
	cltest.Expect(t, `
		var i map[int]string
		println(i == nil,i != nil)
	`, "true false\n")
	cltest.Expect(t, `
		var v int
		println(v == nil)
	`, "", "invalid operator: v == nil (mismatched types int and nil)")
	cltest.Expect(t, `
		var v int
		println(nil == v)
	`, "", "invalid operator: nil == v (mismatched types nil and int)")
	cltest.Expect(t, `
		var a int
		var b uint8
		println(a == b)
	`, "", "invalid operator: a == b (mismatched types int and uint8)")
	cltest.Expect(t, `
		var a int
		switch a {
			case 0:
			case uint8(1):
		}
	`, "", "invalid case uint8(1) in switch on a (mismatched types int and uint8)")
}

func testScripts(t *testing.T, testName string, scripts map[string]testData) {
	for name, script := range scripts {
		t.Log("Run " + testName + "---" + name)
		var panicMsg []interface{}
		if script.panic {
			panicMsg = append(panicMsg, nil)
		}
		cltest.Expect(t, script.clause, script.want, panicMsg...)
	}
}

// -----------------------------------------------------------------------------

func TestTwoValueExpr(t *testing.T) {
	clause := `m:={2:3,1:2}
			if v,ok:=m[m[1]];ok{
				println(1,v,ok)
			}
			if v,ok:=m[m[3]];!ok{
				println(3,v,ok)
			}`
	cltest.Expect(t, clause, "1 3 true\n3 0 false\n")
}

func TestUnsafe(t *testing.T) {
	cltest.Expect(t, `
	import (
		"unsafe"
	)
	type SliceHeader struct {
		Data uintptr
		Len  int
		Cap  int
	}
	type StringHeader struct {
		Data uintptr
		Len  int
	}
	a := "hello"
	b := []byte("world")
	v := (*StringHeader)(unsafe.Pointer(&a))
	v2 := (*SliceHeader)(unsafe.Pointer(&b))
	v3 := (*StringHeader)(unsafe.Pointer(&b))
	println(*(*string)(unsafe.Pointer(v)))
	println(string(*(*[]byte)(unsafe.Pointer(v2))))
	println(*(*string)(unsafe.Pointer(v2)))
	println(*(*string)(unsafe.Pointer(v3)))
	`, "hello\nworld\nworld\nworld\n")
	cltest.Expect(t, `
	import "unsafe"
	type Point struct {
		X int
		Y int
	}
	pt := Point{10, 20}
	println(unsafe.Sizeof(pt),unsafe.Alignof(pt))
	println(unsafe.Offsetof(pt.X), unsafe.Offsetof(pt.Y))
	`, "16 8\n0 8\n")
	cltest.Expect(t, `
	import "unsafe"
	ar := [unsafe.Sizeof(true)]int{}
	println(len(ar))
	`, "1\n")
}
