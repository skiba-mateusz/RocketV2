package commandeer

import (
	"flag"
	"strconv"
)

type Flags struct {
	set *flag.FlagSet
}

func (f *Flags) Init(name string) {
	if f.set == nil {
		f.set = flag.NewFlagSet(name, flag.ContinueOnError)
	}
}

func (f *Flags) SetString(name string, value string, usage string) {
	f.set.String(name, value, usage)
}

func (f *Flags) SetInt(name string, value int, usage string) {
	f.set.Int(name, value, usage)
}

func (f *Flags) SetBool(name string, value bool, usage string) {
	f.set.Bool(name, value, usage)
}

func (f *Flags) GetString(name string) string {
	if f.set == nil {
		return ""
	}

	flag := f.set.Lookup(name)
	if flag == nil {
		return ""
	}
	
	return flag.Value.String()
}

func (f *Flags) GetInt(name string) int {
	if f.set == nil {
		return 0
	}

	flag := f.set.Lookup(name)
	if flag == nil {
		return 0
	}

	val, err := strconv.Atoi(flag.Value.String())
	if err != nil {
		return 0
	}

	return val
}

func (f *Flags) GetBool(name string) bool {
	if f.set == nil {
		return false
	}

	flag := f.set.Lookup(name)
	if flag == nil {
		return false
	}

	val, err := strconv.ParseBool(flag.Value.String())
	if err != nil {
		return false
	}

	return val
}

func (f *Flags) isEmpty() bool {
	empty := true

	f.set.VisitAll(func(f *flag.Flag) {
		empty = false
	})

	return empty
}