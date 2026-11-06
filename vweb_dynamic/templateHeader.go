package vweb_dynamic

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TemplateHeader 标头-模本-处理动态页面文件
// 从文件头部 // 注释中解析配置。
type TemplateHeader struct {
	EntryName             string   // 入口函数名
	File                  []string // 额外加载的文件路径（相对或绝对）
	DelimLeft, DelimRight string   // 语法识别符（保留供扩展使用）
}

// OpenFile 打开文件内容
//
//	rootPath  string    	根目录
//	pagePath  string		页面路径
//	map[string][]byte   内容，map[文件名]文件内容
//	error               错误，如果文件不能打开读取
//
// 说明：无论 v 是否为绝对路径，最终都会以 rootPath 为根解析；
// 这与“沙箱内视图”的语义一致。若需真正访问宿主绝对路径，请显式配置 preopen。
func (T *TemplateHeader) OpenFile(rootPath, pagePath string) (map[string]string, error) {
	dirPath := filepath.Dir(pagePath)
	fileContent := make(map[string]string, len(T.File))

	for _, v := range T.File {
		if v == "" {
			continue
		}

		var filePath string
		if filepath.IsAbs(v) || (len(v) > 0 && (v[0] == '/' || v[0] == '\\')) {
			filePath = filepath.Clean(v)
		} else {
			filePath = filepath.Join(dirPath, v)
		}
		filePath = filepath.Join(rootPath, filePath)

		c, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("vweb_dynamic: Dynamically embedded template file read failed(%s)", err.Error())
		}
		fileContent[filepath.Base(filePath)] = string(c)
	}
	return fileContent, nil
}

// headerMap 将注释行解析为 key=value 映射（支持同 key 多值）
func headerMap(headerLine []string) map[string][]string {
	m := make(map[string][]string, len(headerLine))
	for _, line := range headerLine {
		i := strings.IndexByte(line, '=')
		if i <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:i])
		value := strings.TrimSpace(line[i+1:])
		if key == "" {
			continue
		}
		m[key] = append(m[key], value)
	}
	return m
}

// 解析模本,头
func templateHeader(headerLine []string) TemplateHeader {
	var h TemplateHeader
	for key, vals := range headerMap(headerLine) {
		switch key {
		case "entryName":
			if len(vals) > 0 {
				h.EntryName = vals[0]
			}
		case "file":
			for _, val := range vals {
				if val != "" {
					h.File = append(h.File, val)
				}
			}
		case "delimLeft":
			if len(vals) > 0 {
				h.DelimLeft = vals[0]
			}
		case "delimRight":
			if len(vals) > 0 {
				h.DelimRight = vals[0]
			}
		}
	}
	return h
}
