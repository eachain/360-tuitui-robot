package main

import (
	"reflect"
	"testing"
	"time"
)

type defaultOptions struct {
	A int           `option:"a" default:"123"`
	B string        `option:"b" default:"xyz"`
	C bool          `option:"c" default:"true"`
	D time.Duration `option:"d" default:"10s"`
	E time.Time     `option:"e" default:"1704038400"`
}

func TestParseOptionsDefault(t *testing.T) {
	parse, _ := genParseOptions("", reflect.TypeOf((*defaultOptions)(nil)))
	args, _ := ParseTokens("")
	val, _, _ := parse(args)
	opts := val.Interface().(*defaultOptions)

	if opts.A != 123 {
		t.Fatalf("options default a except 123, return %v", opts.A)
	}
	if opts.B != "xyz" {
		t.Fatalf("options default b except xyz, return %v", opts.B)
	}
	if opts.C != true {
		t.Fatalf("options default c except true, return false")
	}
	if opts.D != 10*time.Second {
		t.Fatalf("options default d except 10s, return %v", opts.D)
	}
	if opts.E.Unix() != 1704038400 {
		t.Fatalf("options default e except 1704038400, return %v", opts.E.Unix())
	}
}

func TestParseOptions(t *testing.T) {
	parse, _ := genParseOptions("", reflect.TypeOf((*defaultOptions)(nil)))
	args, _ := ParseTokens("-e 2024-06-01 -c=false -b zxc -a 798 -d 3s")
	val, _, _ := parse(args)
	opts := val.Interface().(*defaultOptions)

	if opts.A != 798 {
		t.Fatalf("options default a except 798, return %v", opts.A)
	}
	if opts.B != "zxc" {
		t.Fatalf("options default b except zxc, return %v", opts.B)
	}
	if opts.C != false {
		t.Fatalf("options default c except false, return true")
	}
	if opts.D != 3*time.Second {
		t.Fatalf("options default d except 3s, return %v", opts.D)
	}
	if opts.E.Unix() != 1717200000 {
		t.Fatalf("options default e except 1717200000, return %v", opts.E.Unix())
	}
}
