package main

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Cmder struct {
	prefix string
	cmds   []*entry
}

type entry struct {
	name    string
	desc    string
	usage   string
	handler func([]string) string
}

func New(prefix string) *Cmder {
	return &Cmder{prefix: prefix}
}

func (c *Cmder) Exec(s string) string {
	if !strings.HasPrefix(s, c.prefix) {
		return ""
	}
	s = strings.TrimPrefix(s, c.prefix)
	s = strings.TrimSpace(s)

	tokens, err := ParseTokens(s)
	if err != nil {
		return err.Error()
	}
	if len(tokens) == 0 {
		return c.usage()
	}

	cmd := tokens[0]
	args := tokens[1:]
	for _, e := range c.cmds {
		if e.name == cmd {
			return e.handler(args)
		}
	}
	return c.usage()
}

func (c *Cmder) usage() string {
	buf := new(bytes.Buffer)
	fmt.Fprintln(buf, "Commands:")
	for _, cmd := range c.cmds {
		fmt.Fprintf(buf, "\n%v%v: %v\n", c.prefix, cmd.name, cmd.desc)
		if cmd.usage != "" {
			fmt.Fprintln(buf, "usage:")
			fmt.Fprintln(buf, cmd.usage)
		}
	}
	return buf.String()
}

// A func must be like one of:
//
//	func(opts *options, args []string) string
//	func(opts *options) string
//	func(args []string) string
//	func() string
func (c *Cmder) Register(fn any, name, desc string) {
	val := reflect.ValueOf(fn)

	if name == "" {
		name = runtime.FuncForPC(val.Pointer()).Name()
		if dot := strings.LastIndexByte(name, '.'); dot >= 0 {
			name = name[dot+1:]
		}
		if minus := strings.IndexByte(name, '-'); minus > 0 {
			name = name[:minus]
		}
	}

	for _, e := range c.cmds {
		if e.name == name {
			panic(fmt.Errorf("cmder: duplicated register cmd: %v", name))
		}
	}

	typ := val.Type()
	if typ.NumOut() != 1 {
		panic(fmt.Errorf("cmder: register cmd %v: must returns a string as result to output", name))
	}
	if typ.Out(0).Kind() != reflect.String {
		panic(fmt.Errorf("cmder: register cmd %v: must returns a string as result to output", name))
	}

	if typ.NumIn() > 2 {
		panic(fmt.Errorf("cmder: register cmd %v: param in count cannot greater than 2", name))
	}

	if typ.NumIn() == 0 {
		c.cmds = append(c.cmds, &entry{
			name: name,
			desc: desc,
			handler: func([]string) string {
				return val.Call(nil)[0].String()
			},
		})
		return
	}

	if typ.NumIn() == 1 {
		if typ.In(0) == reflect.TypeOf([]string(nil)) {
			c.cmds = append(c.cmds, &entry{
				name:  name,
				desc:  desc,
				usage: fmt.Sprintf("  %v%v [args ...]", c.prefix, name),
				handler: func(args []string) string {
					return val.Call([]reflect.Value{reflect.ValueOf(args)})[0].String()
				},
			})
			return
		}

		handler, usage := wrapOptions(name, val)
		c.cmds = append(c.cmds, &entry{
			name:    name,
			desc:    desc,
			usage:   fmt.Sprintf("  %v%v options\noptions:\n%v", c.prefix, name, usage),
			handler: handler,
		})
		return
	}

	if typ.NumIn() == 2 {
		handler, usage := wrap2(name, val)
		c.cmds = append(c.cmds, &entry{
			name:    name,
			desc:    desc,
			usage:   fmt.Sprintf("  %v%v options [args ...]\noptions:\n%v", c.prefix, name, usage),
			handler: handler,
		})
		return
	}
}

func wrapOptions(cmd string, fn reflect.Value) (func([]string) string, string) {
	parse, usage := genParseOptions(cmd, fn.Type().In(0))
	return func(args []string) string {
		opts, _, err := parse(args)
		if err != nil {
			return err.Error()
		}
		return fn.Call([]reflect.Value{opts})[0].String()
	}, usage
}

func wrap2(cmd string, fn reflect.Value) (func([]string) string, string) {
	parse, usage := genParseOptions(cmd, fn.Type().In(0))
	return func(args []string) string {
		opts, args, err := parse(args)
		if err != nil {
			return err.Error()
		}
		return fn.Call([]reflect.Value{opts, reflect.ValueOf(args)})[0].String()
	}, usage
}

type (
	parseOptionFunc  func(reflect.Value, []string) ([]string, error)
	parseOptionsFunc func([]string) (reflect.Value, []string, error)
)

func genParseOptions(cmd string, typ reflect.Type) (parseOptionsFunc, string) {
	if typ.Kind() != reflect.Pointer {
		panic(fmt.Errorf("cmd %v option arg must be a pointer", cmd))
	}
	typ = typ.Elem()
	if typ.Kind() != reflect.Struct {
		panic(fmt.Errorf("cmd %v option arg must be a pointer of struct", cmd))
	}

	var usages []string
	var setDefalts []func(reflect.Value)
	parse := make(map[string]parseOptionFunc)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		opt := field.Tag.Get("option")
		opt = strings.TrimLeft(opt, "-")
		if opt == "" {
			continue
		}
		dft := field.Tag.Get("default")
		usage := field.Tag.Get("usage")

		if dft != "" {
			usages = append(usages, fmt.Sprintf("  -%v %v %v (default: %v)", opt, field.Type, usage, dft))
		} else {
			usages = append(usages, fmt.Sprintf("  -%v %v %v", opt, field.Type, usage))
		}

		if field.Type.String() == "time.Duration" {
			parse[opt] = genParseDuration(i, opt)
		} else if field.Type.String() == "time.Time" {
			parse[opt] = genParseTime(i, opt)
		} else {
			switch field.Type.Kind() {
			default:
				panic(fmt.Errorf("cmd %v unsuported option %v type: %v", cmd, opt, field.Type))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				parse[opt] = genParseInt(i, opt)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				parse[opt] = genParseUint(i, opt)
			case reflect.String:
				parse[opt] = genParseString(i, opt)
			case reflect.Float32, reflect.Float64:
				parse[opt] = genParseFloat(i, opt)
			case reflect.Bool:
				parse[opt] = genParseBool(i)
			}
		}

		if dft != "" {
			_, err := parse[opt](reflect.New(typ), []string{dft})
			if err != nil {
				panic(fmt.Errorf("cmd %v option %v parse default value %q: %v", cmd, opt, dft, err))
			}
			setDefalts = append(setDefalts, func(v reflect.Value) {
				parse[opt](v, []string{dft})
			})
		}
	}

	return genParseOptionsFunc(typ, setDefalts, parse), strings.Join(usages, "\n")
}

func genParseDuration(i int, opt string) func(reflect.Value, []string) ([]string, error) {
	return func(val reflect.Value, args []string) ([]string, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("%v: miss param", opt)
		}
		dur, err := time.ParseDuration(args[0])
		if err == nil {
			val.Elem().Field(i).Set(reflect.ValueOf(dur))
		}
		return args[1:], err
	}
}

func genParseTime(i int, opt string) func(reflect.Value, []string) ([]string, error) {
	return func(val reflect.Value, args []string) ([]string, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("%v: miss param", opt)
		}
		layouts := []string{
			"20060102150405",
			"20060102T150405",
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05",
			"2006/01/02 15:04:05",
			"2006/01/02T15:04:05",
			time.RFC3339,
			"20060102",
			"2006-01-02",
			"2006/01/02",
		}
		for _, layout := range layouts {
			t, err := time.Parse(layout, args[0])
			if err == nil {
				val.Elem().Field(i).Set(reflect.ValueOf(t))
				return args[1:], nil
			}
		}
		v, err := strconv.ParseInt(args[0], 10, 64)
		if err == nil {
			val.Elem().Field(i).Set(reflect.ValueOf(time.Unix(v, 0)))
			return args[1:], err
		}
		return args, fmt.Errorf("%v: time layout invalid", opt)
	}
}

func genParseInt(i int, opt string) func(reflect.Value, []string) ([]string, error) {
	return func(val reflect.Value, args []string) ([]string, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("%v: miss param", opt)
		}
		v, err := strconv.ParseInt(args[0], 10, 64)
		if err == nil {
			val.Elem().Field(i).SetInt(v)
		}
		return args[1:], err
	}
}

func genParseUint(i int, opt string) func(reflect.Value, []string) ([]string, error) {
	return func(val reflect.Value, args []string) ([]string, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("%v: miss param", opt)
		}
		v, err := strconv.ParseUint(args[0], 10, 64)
		if err == nil {
			val.Elem().Field(i).SetUint(v)
		}
		return args[1:], err
	}
}

func genParseString(i int, opt string) func(reflect.Value, []string) ([]string, error) {
	return func(val reflect.Value, args []string) ([]string, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("%v: miss param", opt)
		}
		val.Elem().Field(i).SetString(args[0])
		return args[1:], nil
	}
}

func genParseFloat(i int, opt string) func(reflect.Value, []string) ([]string, error) {
	return func(val reflect.Value, args []string) ([]string, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("%v: miss param", opt)
		}
		v, err := strconv.ParseFloat(args[0], 64)
		if err == nil {
			val.Elem().Field(i).SetFloat(v)
		}
		return args[1:], err
	}
}

func genParseBool(i int) func(reflect.Value, []string) ([]string, error) {
	return func(val reflect.Value, args []string) ([]string, error) {
		if len(args) == 0 {
			val.Elem().Field(i).SetBool(true)
			return nil, nil
		}
		if args[0] == "true" {
			val.Elem().Field(i).SetBool(true)
			return args[1:], nil
		}
		if args[0] == "false" {
			val.Elem().Field(i).SetBool(false)
			return args[1:], nil
		}
		val.Elem().Field(i).SetBool(true)
		return args, nil
	}
}

func genParseOptionsFunc(typ reflect.Type, setDefalts []func(reflect.Value),
	parse map[string]parseOptionFunc) parseOptionsFunc {

	return func(args []string) (reflect.Value, []string, error) {
		var last []string
		val := reflect.New(typ)

		for _, d := range setDefalts {
			d(val)
		}

		for len(args) > 0 {
			arg := args[0]
			if !strings.HasPrefix(arg, "-") {
				last = append(last, arg)
				args = args[1:]
				continue
			}

			arg = strings.TrimLeft(arg, "-")
			if eq := strings.IndexByte(arg, '='); eq > 0 {
				opt := arg[:eq]
				if p := parse[opt]; p != nil {
					_, err := p(val, []string{arg[eq+1:]})
					if err != nil {
						return val, nil, err
					}
				} else {
					last = append(last, args[0])
				}
				args = args[1:]
				continue
			}

			if p := parse[arg]; p != nil {
				var err error
				args, err = p(val, args[1:])
				if err != nil {
					return reflect.Value{}, nil, err
				}
			} else {
				last = append(last, args[0])
				args = args[1:]
			}
			continue
		}
		return val, last, nil
	}
}
