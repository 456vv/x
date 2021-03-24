// Package regexp provide Go+ "regexp" package, as "regexp" package in Go.
package regexp

import (
	io "io"
	reflect "reflect"
	regexp "regexp"

	gop "github.com/goplus/gop"
)

func execCompile(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := regexp.Compile(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execCompilePOSIX(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := regexp.CompilePOSIX(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execMatch(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := regexp.Match(args[0].(string), args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func toType0(v interface{}) io.RuneReader {
	if v == nil {
		return nil
	}
	return v.(io.RuneReader)
}

func execMatchReader(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := regexp.MatchReader(args[0].(string), toType0(args[1]))
	p.Ret(2, ret0, ret1)
}

func execMatchString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := regexp.MatchString(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execMustCompile(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := regexp.MustCompile(args[0].(string))
	p.Ret(1, ret0)
}

func execMustCompilePOSIX(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := regexp.MustCompilePOSIX(args[0].(string))
	p.Ret(1, ret0)
}

func execQuoteMeta(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := regexp.QuoteMeta(args[0].(string))
	p.Ret(1, ret0)
}

func execmRegexpString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*regexp.Regexp).String()
	p.Ret(1, ret0)
}

func execmRegexpCopy(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*regexp.Regexp).Copy()
	p.Ret(1, ret0)
}

func execmRegexpLongest(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*regexp.Regexp).Longest()
	p.PopN(1)
}

func execmRegexpNumSubexp(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*regexp.Regexp).NumSubexp()
	p.Ret(1, ret0)
}

func execmRegexpSubexpNames(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*regexp.Regexp).SubexpNames()
	p.Ret(1, ret0)
}

func execmRegexpSubexpIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*regexp.Regexp).SubexpIndex(args[1].(string))
	p.Ret(2, ret0)
}

func execmRegexpLiteralPrefix(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*regexp.Regexp).LiteralPrefix()
	p.Ret(1, ret0, ret1)
}

func execmRegexpMatchReader(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*regexp.Regexp).MatchReader(toType0(args[1]))
	p.Ret(2, ret0)
}

func execmRegexpMatchString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*regexp.Regexp).MatchString(args[1].(string))
	p.Ret(2, ret0)
}

func execmRegexpMatch(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*regexp.Regexp).Match(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmRegexpReplaceAllString(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*regexp.Regexp).ReplaceAllString(args[1].(string), args[2].(string))
	p.Ret(3, ret0)
}

func execmRegexpReplaceAllLiteralString(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*regexp.Regexp).ReplaceAllLiteralString(args[1].(string), args[2].(string))
	p.Ret(3, ret0)
}

func execmRegexpReplaceAllStringFunc(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*regexp.Regexp).ReplaceAllStringFunc(args[1].(string), args[2].(func(string) string))
	p.Ret(3, ret0)
}

func execmRegexpReplaceAll(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*regexp.Regexp).ReplaceAll(args[1].([]byte), args[2].([]byte))
	p.Ret(3, ret0)
}

func execmRegexpReplaceAllLiteral(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*regexp.Regexp).ReplaceAllLiteral(args[1].([]byte), args[2].([]byte))
	p.Ret(3, ret0)
}

func execmRegexpReplaceAllFunc(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*regexp.Regexp).ReplaceAllFunc(args[1].([]byte), args[2].(func([]byte) []byte))
	p.Ret(3, ret0)
}

func execmRegexpFind(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*regexp.Regexp).Find(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmRegexpFindIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*regexp.Regexp).FindIndex(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmRegexpFindString(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*regexp.Regexp).FindString(args[1].(string))
	p.Ret(2, ret0)
}

func execmRegexpFindStringIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*regexp.Regexp).FindStringIndex(args[1].(string))
	p.Ret(2, ret0)
}

func execmRegexpFindReaderIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*regexp.Regexp).FindReaderIndex(toType0(args[1]))
	p.Ret(2, ret0)
}

func execmRegexpFindSubmatch(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*regexp.Regexp).FindSubmatch(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmRegexpExpand(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0 := args[0].(*regexp.Regexp).Expand(args[1].([]byte), args[2].([]byte), args[3].([]byte), args[4].([]int))
	p.Ret(5, ret0)
}

func execmRegexpExpandString(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	ret0 := args[0].(*regexp.Regexp).ExpandString(args[1].([]byte), args[2].(string), args[3].(string), args[4].([]int))
	p.Ret(5, ret0)
}

func execmRegexpFindSubmatchIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*regexp.Regexp).FindSubmatchIndex(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmRegexpFindStringSubmatch(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*regexp.Regexp).FindStringSubmatch(args[1].(string))
	p.Ret(2, ret0)
}

func execmRegexpFindStringSubmatchIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*regexp.Regexp).FindStringSubmatchIndex(args[1].(string))
	p.Ret(2, ret0)
}

func execmRegexpFindReaderSubmatchIndex(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*regexp.Regexp).FindReaderSubmatchIndex(toType0(args[1]))
	p.Ret(2, ret0)
}

func execmRegexpFindAll(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*regexp.Regexp).FindAll(args[1].([]byte), args[2].(int))
	p.Ret(3, ret0)
}

func execmRegexpFindAllIndex(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*regexp.Regexp).FindAllIndex(args[1].([]byte), args[2].(int))
	p.Ret(3, ret0)
}

func execmRegexpFindAllString(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*regexp.Regexp).FindAllString(args[1].(string), args[2].(int))
	p.Ret(3, ret0)
}

func execmRegexpFindAllStringIndex(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*regexp.Regexp).FindAllStringIndex(args[1].(string), args[2].(int))
	p.Ret(3, ret0)
}

func execmRegexpFindAllSubmatch(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*regexp.Regexp).FindAllSubmatch(args[1].([]byte), args[2].(int))
	p.Ret(3, ret0)
}

func execmRegexpFindAllSubmatchIndex(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*regexp.Regexp).FindAllSubmatchIndex(args[1].([]byte), args[2].(int))
	p.Ret(3, ret0)
}

func execmRegexpFindAllStringSubmatch(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*regexp.Regexp).FindAllStringSubmatch(args[1].(string), args[2].(int))
	p.Ret(3, ret0)
}

func execmRegexpFindAllStringSubmatchIndex(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*regexp.Regexp).FindAllStringSubmatchIndex(args[1].(string), args[2].(int))
	p.Ret(3, ret0)
}

func execmRegexpSplit(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*regexp.Regexp).Split(args[1].(string), args[2].(int))
	p.Ret(3, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("regexp")

func init() {
	I.RegisterFuncs(
		I.Func("Compile", regexp.Compile, execCompile),
		I.Func("CompilePOSIX", regexp.CompilePOSIX, execCompilePOSIX),
		I.Func("Match", regexp.Match, execMatch),
		I.Func("MatchReader", regexp.MatchReader, execMatchReader),
		I.Func("MatchString", regexp.MatchString, execMatchString),
		I.Func("MustCompile", regexp.MustCompile, execMustCompile),
		I.Func("MustCompilePOSIX", regexp.MustCompilePOSIX, execMustCompilePOSIX),
		I.Func("QuoteMeta", regexp.QuoteMeta, execQuoteMeta),
		I.Func("(*Regexp).String", (*regexp.Regexp).String, execmRegexpString),
		I.Func("(*Regexp).Copy", (*regexp.Regexp).Copy, execmRegexpCopy),
		I.Func("(*Regexp).Longest", (*regexp.Regexp).Longest, execmRegexpLongest),
		I.Func("(*Regexp).NumSubexp", (*regexp.Regexp).NumSubexp, execmRegexpNumSubexp),
		I.Func("(*Regexp).SubexpNames", (*regexp.Regexp).SubexpNames, execmRegexpSubexpNames),
		I.Func("(*Regexp).SubexpIndex", (*regexp.Regexp).SubexpIndex, execmRegexpSubexpIndex),
		I.Func("(*Regexp).LiteralPrefix", (*regexp.Regexp).LiteralPrefix, execmRegexpLiteralPrefix),
		I.Func("(*Regexp).MatchReader", (*regexp.Regexp).MatchReader, execmRegexpMatchReader),
		I.Func("(*Regexp).MatchString", (*regexp.Regexp).MatchString, execmRegexpMatchString),
		I.Func("(*Regexp).Match", (*regexp.Regexp).Match, execmRegexpMatch),
		I.Func("(*Regexp).ReplaceAllString", (*regexp.Regexp).ReplaceAllString, execmRegexpReplaceAllString),
		I.Func("(*Regexp).ReplaceAllLiteralString", (*regexp.Regexp).ReplaceAllLiteralString, execmRegexpReplaceAllLiteralString),
		I.Func("(*Regexp).ReplaceAllStringFunc", (*regexp.Regexp).ReplaceAllStringFunc, execmRegexpReplaceAllStringFunc),
		I.Func("(*Regexp).ReplaceAll", (*regexp.Regexp).ReplaceAll, execmRegexpReplaceAll),
		I.Func("(*Regexp).ReplaceAllLiteral", (*regexp.Regexp).ReplaceAllLiteral, execmRegexpReplaceAllLiteral),
		I.Func("(*Regexp).ReplaceAllFunc", (*regexp.Regexp).ReplaceAllFunc, execmRegexpReplaceAllFunc),
		I.Func("(*Regexp).Find", (*regexp.Regexp).Find, execmRegexpFind),
		I.Func("(*Regexp).FindIndex", (*regexp.Regexp).FindIndex, execmRegexpFindIndex),
		I.Func("(*Regexp).FindString", (*regexp.Regexp).FindString, execmRegexpFindString),
		I.Func("(*Regexp).FindStringIndex", (*regexp.Regexp).FindStringIndex, execmRegexpFindStringIndex),
		I.Func("(*Regexp).FindReaderIndex", (*regexp.Regexp).FindReaderIndex, execmRegexpFindReaderIndex),
		I.Func("(*Regexp).FindSubmatch", (*regexp.Regexp).FindSubmatch, execmRegexpFindSubmatch),
		I.Func("(*Regexp).Expand", (*regexp.Regexp).Expand, execmRegexpExpand),
		I.Func("(*Regexp).ExpandString", (*regexp.Regexp).ExpandString, execmRegexpExpandString),
		I.Func("(*Regexp).FindSubmatchIndex", (*regexp.Regexp).FindSubmatchIndex, execmRegexpFindSubmatchIndex),
		I.Func("(*Regexp).FindStringSubmatch", (*regexp.Regexp).FindStringSubmatch, execmRegexpFindStringSubmatch),
		I.Func("(*Regexp).FindStringSubmatchIndex", (*regexp.Regexp).FindStringSubmatchIndex, execmRegexpFindStringSubmatchIndex),
		I.Func("(*Regexp).FindReaderSubmatchIndex", (*regexp.Regexp).FindReaderSubmatchIndex, execmRegexpFindReaderSubmatchIndex),
		I.Func("(*Regexp).FindAll", (*regexp.Regexp).FindAll, execmRegexpFindAll),
		I.Func("(*Regexp).FindAllIndex", (*regexp.Regexp).FindAllIndex, execmRegexpFindAllIndex),
		I.Func("(*Regexp).FindAllString", (*regexp.Regexp).FindAllString, execmRegexpFindAllString),
		I.Func("(*Regexp).FindAllStringIndex", (*regexp.Regexp).FindAllStringIndex, execmRegexpFindAllStringIndex),
		I.Func("(*Regexp).FindAllSubmatch", (*regexp.Regexp).FindAllSubmatch, execmRegexpFindAllSubmatch),
		I.Func("(*Regexp).FindAllSubmatchIndex", (*regexp.Regexp).FindAllSubmatchIndex, execmRegexpFindAllSubmatchIndex),
		I.Func("(*Regexp).FindAllStringSubmatch", (*regexp.Regexp).FindAllStringSubmatch, execmRegexpFindAllStringSubmatch),
		I.Func("(*Regexp).FindAllStringSubmatchIndex", (*regexp.Regexp).FindAllStringSubmatchIndex, execmRegexpFindAllStringSubmatchIndex),
		I.Func("(*Regexp).Split", (*regexp.Regexp).Split, execmRegexpSplit),
	)
	I.RegisterTypes(
		I.Type("Regexp", reflect.TypeOf((*regexp.Regexp)(nil)).Elem()),
	)
}
