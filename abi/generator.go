//go:build ignore

package main

import (
	"fmt"
	"github.com/tonkeeper/tongo/abi/parser"
	"go/format"
	"os"
	"strings"
)

const HEADER = `package abi
// Code autogenerated. DO NOT EDIT. 

import (
%v
)

`

func main() {
	scheme, err := os.ReadFile("known.xml")
	if err != nil {
		panic(err)
	}
	abi, err := parser.ParseABI(scheme)
	if err != nil {
		panic(err)
	}

	gen, err := parser.NewGenerator(nil, abi)
	if err != nil {
		panic(err)
	}

	types := gen.CollectedTypes()
	msgDecoder := gen.GenerateMsgDecoder()

	getMethods, err := gen.GetMethods()
	if err != nil {
		panic(err)
	}
	invocationOrder, err := gen.RenderInvocationOrderList()
	if err != nil {
		panic(err)
	}

	for _, f := range [][]string{
		{types, "types.go", `"github.com/tonkeeper/tongo/tlb"`},
		{msgDecoder, "messages.go", `"fmt"`, `"github.com/tonkeeper/tongo/boc"`, `"github.com/tonkeeper/tongo/tlb"`},
		{getMethods, "get_methods.go", `"context"`, `"fmt"`, `"github.com/tonkeeper/tongo"`, `"github.com/tonkeeper/tongo/boc"`, `"github.com/tonkeeper/tongo/tlb"`},
		{invocationOrder, "ordering.go", `"context"`, `"github.com/tonkeeper/tongo"`},
	} {
		file, err := os.Create(f[1])
		if err != nil {
			panic(err)
		}
		_, err = file.WriteString(fmt.Sprintf(HEADER, strings.Join(f[2:], "\n")))
		if err != nil {
			panic(err)
		}
		code, err := format.Source([]byte(f[0]))
		if err != nil {
			panic(err)
		}
		_, err = file.Write(code)
		if err != nil {
			panic(err)
		}
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}
}
