package vweb_dynamic

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 标头-模本-处理动态页面文件
type TemplateHeader struct {
	File                  []string // 文件路径, map[文件名或别名]文件路径
	DelimLeft, DelimRight string   // 语法识别符
}

// 打开文件内容
//
//	rootPath  string    	根目录
//	pagePath  string		页面路径
//	map[string][]byte   内容，map[文件名]文件内容
//	error               错误，如果文件不能打开读取
func (T *TemplateHeader) OpenFile(rootPath, pagePath string) (map[string]string, error) {
	var (
		dirPath     = filepath.Dir(pagePath)
		filePath    string
		fileBase    string
		fileContent = make(map[string]string)
		c           []byte
		err         error
	)
	for _, v := range T.File {
		if v[0] == '/' || v[0] == '\\' {
			filePath = filepath.Clean(v)
		} else {
			filePath = filepath.Join(dirPath, v)
		}
		filePath = filepath.Join(rootPath, filePath)
		c, err = os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("vweb_dynamic: Dynamically embedded template file read failed(%s)", err.Error())
		}
		fileBase = filepath.Base(filePath)
		fileContent[fileBase] = string(c)
	}
	return fileContent, nil
}

func headerMap(headerLine []string) map[string][]string {
	m := map[string][]string{}
	for _, line := range headerLine {
		// 跳过空key
		i := strings.IndexByte(line, '=')
		if i <= 0 {
			continue
		}

		key := strings.Trim(line[:i], "\t ")
		i++ // 跳过=符号
		value := strings.Trim(line[i:], "\t ")
		m[key] = append(m[key], value)
	}
	return m
}

// 解析模本,头
func templateHeader(headerLine []string) TemplateHeader {
	var h TemplateHeader
	for key, vals := range headerMap(headerLine) {
		switch key {
		case "file":
			for _, val := range vals {
				if val == "" {
					continue
				}
				h.File = append(h.File, val)
			}
		case "delimLeft":
			h.DelimLeft = vals[0]
		case "delimRight":
			h.DelimRight = vals[0]
		}
	}
	return h
}
