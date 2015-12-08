package functions

import (
	"strconv"
	"strings"
	"text/template"
	"math/rand"
	"math"
)

func FunctionsMap() template.FuncMap {
	return template.FuncMap(fmap)
}

var fmap = map[string]interface{}{

	"trim":    strings.TrimSpace,
	"upper":   strings.ToUpper,
	"lower":   strings.ToLower,
	"title":   strings.Title,
	"substr":  func(start, length int, s string) string {return s[start:length]},
	"join":    func(sep string, ss []string) string { return strings.Join(ss, sep) },

	"repeat":  func(count int, str string) string { return strings.Repeat(str, count) },
	"trimall": func(c, s string) string { return strings.Trim(s, c) },
	"atoi":    func(a string) int { i, _ := strconv.Atoi(a); return i },
	"split":   func(sep, s string) []string { return strings.Split(s, sep)},

	"add":     func(a, b int) int { return a + b },
	"sub":     func(a, b int) int { return a - b },
	"div":     func(a, b int) int { return a / b },
	"mod":     func(a, b int) int { return a % b },
	"mul":     func(a, b int) int { return a * b },
	"max":     math.Max,
	"min":     math.Min,
	"rnd":     func(a, b int) int { return rand.Intn(b - a + 1) + a },
	"random":  func (s []string) (string) {return s[rand.Intn(len(s))]},
	"randoms":  func (s string) (string) {a := strings.Split(s, "|"); return a[rand.Intn(len(a))]},

}