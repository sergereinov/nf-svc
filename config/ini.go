package config

import (
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

type iniFile struct {
	*ini.File
}

func NewIniFile(path string) (*iniFile, error) {
	file, err := ini.Load(path)
	if err != nil {
		return &iniFile{File: ini.Empty()}, err
	}

	return &iniFile{File: file}, nil
}

// Some helpers make it look more like `flag.String` etc

func (f *iniFile) String(section, key, defaultVal string) string {
	return f.Section(section).Key(key).MustString(defaultVal)
}

func (f *iniFile) Int(section, key string, defaultVal int) int {
	return f.Section(section).Key(key).MustInt(defaultVal)
}

func (f *iniFile) Strings(section, key string, defaultVal []string, optDelim ...string) []string {
	var delim string
	if len(optDelim) > 0 {
		delim = optDelim[0]
	} else {
		delim = ","
	}

	s := f.Section(section)
	if s.HasKey(key) {
		return s.Key(key).Strings(delim)
	}

	s.Key(key).SetValue(strings.Join(defaultVal, delim+" "))
	return defaultVal
}

func (f *iniFile) Ints(section, key string, defaultVal []int, optDelim ...string) []int {
	var delim string
	if len(optDelim) > 0 {
		delim = optDelim[0]
	} else {
		delim = ","
	}

	s := f.Section(section)
	if s.HasKey(key) {
		return s.Key(key).Ints(delim)
	}

	var sb strings.Builder
	for i, v := range defaultVal {
		s := strconv.FormatInt(int64(v), 10)
		sb.WriteString(s)
		if i < len(defaultVal)-1 {
			sb.WriteString(delim + " ")
		}
	}
	s.Key(key).SetValue(sb.String())
	return defaultVal
}
