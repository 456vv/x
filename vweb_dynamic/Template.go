package vweb_dynamic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/456vv/vweb/v3"
)

// 并发安全说明：
//   - 同一 *Template 实例在 Parse / ParseFile / ParseText / SetPath 等写操作期间
//     不得被其他 goroutine 并发调用 Execute。
//   - 解析完成后，多个 goroutine 可安全并发调用 Execute（text/template 本身支持）。
type Template struct {
	rootPath  string // 文件目录
	pagePath  string // 文件名称
	entryName string

	fileName string
	t        *template.Template
}

func (T *Template) ParseText(name, content string) error {
	T.fileName = name
	r := strings.NewReader(content)
	return T.Parse(r)
}

func (T *Template) ParseFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()
	T.fileName = filepath.Base(path)
	return T.Parse(file)
}

func (T *Template) SetPath(root, page string) {
	T.rootPath = root
	T.pagePath = page
	T.fileName = filepath.Base(page)
}

func (T *Template) SetEntryName(name string) {
	T.entryName = name
}

func (T *Template) Parse(r io.Reader) error {
	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(r); err != nil {
		return err // 显式处理 ReadFrom 错误，避免静默失败
	}

	// 解析文件头和主体数据
	l := fileHeaderLine(buf)
	h := templateHeader(l)
	libs, err := h.OpenFile(T.rootPath, T.pagePath)
	if err != nil {
		return err
	}
	// 解析主体内容
	cs := T.format(h.DelimLeft, h.DelimRight, buf.String())

	// 模板文件内容载入Tmplate
	t := template.New(T.fileName)
	t.Delims(h.DelimLeft, h.DelimRight)
	t.Funcs(TemplateFunc)
	_, err = t.Parse(cs)
	if err != nil {
		return err
	}

	// 解析子内容
	T.t, err = T.loadTmpl(t, h.DelimLeft, h.DelimRight, libs)
	return err
}

func (T *Template) Execute(out io.Writer, in any) error {
	if T.t == nil {
		return errTemplateNotParse
	}
	// 执行模板
	if tdot, ok := in.(vweb.DotContexter); ok {
		tdot.WithContext(context.WithValue(tdot.Context(), "Template", &TemplateExtend{Template: T.t}))
	}
	return T.t.ExecuteTemplate(out, T.fileName, in)
}

// loadTmpl 模本载入
// 在父模板为 nil 时提前返回；对子模板解析错误保持原行为。
func (T *Template) loadTmpl(t *template.Template, delimLeft, delimRight string, f map[string]string) (*template.Template, error) {
	if t == nil {
		return nil, errors.New("the parent template is nil")
	}
	for k, v := range f {
		tmpl := t.New(k)
		v = T.format(delimLeft, delimRight, v)
		if _, err := tmpl.Parse(v); err != nil {
			return nil, err
		}
	}
	return t, nil
}

// format 语法整合
// 1. 修正换行符检测逻辑，确保优先使用 \r\n，其次 \n，最后 \r；2. 优化语法整合函数，实现
// 统一换行符处理逻辑，减少中间切片与重复 Trim，提升可读性与跨平台稳定性；保持原有“{{ ... }}”自动包裹行为。
func (T *Template) format(delimLeft, delimRight, c string) string {
	if delimLeft == "" {
		delimLeft = "{{"
	}
	if delimRight == "" {
		delimRight = "}}"
	}

	// 检测实际使用的换行符，优先 \r\n，其次 \n，最后 \r
	var lineSep string
	switch {
	case strings.Contains(c, "\r\n"):
		lineSep = "\r\n"
	case strings.Contains(c, "\n"):
		lineSep = "\n"
	case strings.Contains(c, "\r"):
		lineSep = "\r"
	default:
		// 无换行符，直接返回
		return c
	}

	var bytesBuffer strings.Builder
	clines := strings.Split(c, lineSep)
	clinesL := len(clines) - 1
	syntax := false

	for k, cline := range clines {
		clineTrim := strings.Trim(cline, " \t")
		leftHas := strings.HasSuffix(clineTrim, delimLeft)
		rightHas := strings.HasPrefix(clineTrim, delimRight)

		switch {
		case leftHas && rightHas:
			// 格式：\r\n    }}abcx{{
			clineTrim = strings.TrimPrefix(clineTrim, delimRight)
			clineTrim = strings.TrimSuffix(clineTrim, delimLeft)
			bytesBuffer.WriteString(clineTrim)
			syntax = true
			continue
		case leftHas:
			// 格式：abcx{{
			cline = strings.TrimRight(cline, " \t")
			cline = strings.TrimSuffix(cline, delimLeft)
			bytesBuffer.WriteString(cline)
			syntax = true
			continue
		case rightHas:
			// 格式：}}12345
			cline = strings.TrimLeft(cline, " \t")
			cline = strings.TrimPrefix(cline, delimRight)
			syntax = false
		}

		if syntax {
			if clineTrim == "" || strings.HasPrefix(clineTrim, "//") {
				// 空行或注释行跳过
				continue
			}
			// 格式：if eq 1 1 → 自动包裹分隔符
			bytesBuffer.WriteString(delimLeft)
			bytesBuffer.WriteString(cline)
			bytesBuffer.WriteString(delimRight)
		} else {
			// 文本行
			bytesBuffer.WriteString(cline)
			if k != clinesL {
				bytesBuffer.WriteString(lineSep)
			}
		}
	}

	if bytesBuffer.Len() != 0 {
		return bytesBuffer.String()
	}
	return c
}

type part struct {
	input  []reflect.Value
	output []reflect.Value
}

func (T *part) Args(i int) any {
	if i == -1 {
		var ret []any
		for _, in := range T.input {
			ret = append(ret, in.Interface())
		}
		return ret
	}
	if len(T.input) > i {
		v := T.input[i]
		return v.Interface()
	}
	return nil
}

func (T *part) Result(out ...any) {
	for _, arg := range out {
		T.output = append(T.output, reflect.ValueOf(arg))
	}
}

// TemplateExtend 这是个额外扩展，由于模板不能实现函数创建，只能在这里构造一个支持创建函数。
// 在创建的函数内部，需要使用 Args 方法读取参数，使用 Result 方法返回结果。
// 仅限用于template模板，自定义模板不支持
type TemplateExtend struct {
	*template.Template
}

// NewFunc 构建一个新的函数，仅限在template中使用
//
//	func([]reflect.Value) []reflect.Value)	新的函数
func (T *TemplateExtend) NewFunc(name string) (f func([]reflect.Value) []reflect.Value, err error) {
	if T.Template == nil {
		return nil, errTemplateNotParse
	}
	if T.Template.Lookup(name) == nil {
		return nil, fmt.Errorf("Template content not found {{define \"%s\"}} ... {{end}} , Calling this method does not support", name)
	}
	return func(in []reflect.Value) []reflect.Value {
		p := &part{input: in}
		err := T.Template.ExecuteTemplate(io.Discard, name, p)
		if err != nil {
			panic(err)
		}
		return p.output
	}, nil
}

// Call 执行函数
//
//	f func([]reflect.Value) []reflect.Value	由NewFunc创建的函数
//	args ...interface{}						可变参数
//	[]interface{}							返回结果
func (T *TemplateExtend) Call(f func([]reflect.Value) []reflect.Value, args ...any) []any {
	var (
		inv []reflect.Value
		ret []any
		ec  vweb.ExecCall
	)
	for _, arg := range args {
		inv = append(inv, reflect.ValueOf(arg))
	}
	if err := ec.Func(f, inv); err != nil {
		panic(err)
	}
	// NewFunc 执行后返回是[]reflect.Value
	for _, result := range ec.Exec() {
		// 已100%确认变量的类型为reflect.Value
		// 1，查看func (T *serverHandlerDynamicTemplateExtend) NewFunc(name string) (f func([]reflect.Value) []reflect.Value, err error)
		// 2，查看func (T *execFunc) exec() (ret []interface{})
		for _, rv := range result.([]reflect.Value) {
			ret = append(ret, rv.Interface())
		}
	}
	return ret
}

// ExecuteTemplate 调用模板
//
//	out io.Writer	模板内容返回内容写入到out
//	name string		模板名称
//	in interface{}	模板点，带外模板的内容
//	error			执行语法错误
func (T *TemplateExtend) ExecuteTemplate(out io.Writer, name string, in any) error {
	if T.Template == nil {
		return errTemplateNotParse
	}
	if T.Template.Lookup(name) == nil {
		return fmt.Errorf("Template content not found {{define \"%s\"}} ... {{end}} , Calling this method does not support", name)
	}
	// 执行模板
	if tdot, ok := in.(vweb.DotContexter); ok {
		tdot.WithContext(context.WithValue(tdot.Context(), "Template", T))
	}
	return T.Template.ExecuteTemplate(out, name, in)
}
