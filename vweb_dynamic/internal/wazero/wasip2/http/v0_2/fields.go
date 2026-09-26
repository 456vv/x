package v0_2

import (
	"context"
	"sort"
	"strings"

	manager_http "github.com/456vv/x/vweb_dynamic/internal/wazero/manager/http"
	witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"
)

// fieldsImpl 封装了 fields 资源的所有操作。
type fieldsImpl struct {
	hm *manager_http.HTTPManager
}

func newFieldsImpl(hm *manager_http.HTTPManager) *fieldsImpl {
	return &fieldsImpl{hm: hm}
}

// 对每个字符 ContainsRune 扫描分隔符串；256 表在 fields.set/append 热路径上更便宜，行为一致
var invalidFieldNameChar = func() (t [256]bool) {
	for _, c := range []byte("()<>@,;:\\\"/[]?={}") {
		t[c] = true
	}
	return t
}()

func validFieldName(name string) bool {
	if name == "" {
		return false
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		if c <= 32 || c >= 127 || invalidFieldNameChar[c] {
			return false
		}
	}
	return true
}

func validFieldValue(value []byte) bool {
	for i := 0; i < len(value); i++ {
		c := value[i]
		// CR/LF/NUL 会构成头注入，进入 net/http 或拼进请求行
		if c == 0 || c == '\r' || c == '\n' {
			return false
		}
	}
	return true
}

func (i *fieldsImpl) Constructor() Fields {
	return i.hm.Fields.Add(make(manager_http.Fields))
}

// FromListConstructor 实现了 [constructor]fields.from-list。
func (i *fieldsImpl) FromList(_ context.Context, entries []witgo.Tuple[FieldKey, FieldValue]) witgo.Result[Fields, HeaderError] {
	fields := make(manager_http.Fields)
	for _, entry := range entries {
		if !validFieldName(entry.F0) || !validFieldValue(entry.F1) {
			return witgo.Err[Fields, HeaderError](HeaderError{InvalidSyntax: &witgo.Unit{}})
		}
		if forbiddenFieldName(entry.F0) { // FromList 用 entry.F0
			return witgo.Err[Fields, HeaderError](HeaderError{Forbidden: &witgo.Unit{}})
		}
		key := strings.ToLower(entry.F0)
		fields[key] = append(fields[key], string(entry.F1))
	}
	return witgo.Ok[Fields, HeaderError](i.hm.Fields.Add(fields))
}

// Drop 是 fields 资源的析构函数。
func (i *fieldsImpl) Drop(_ context.Context, handle Fields) {
	i.hm.UnmarkFieldsImmutable(handle)
	i.hm.Fields.Remove(handle)
}

// Get 实现了 [method]fields.get。
func (i *fieldsImpl) Get(_ context.Context, this Fields, name FieldKey) []FieldValue {
	i.hm.LockFields()
	defer i.hm.UnlockFields()
	f, ok := i.hm.Fields.Get(this)
	if !ok || f == nil {
		return nil
	}
	values := f[strings.ToLower(name)]
	ret := make([]FieldValue, len(values))
	for j, v := range values {
		ret[j] = FieldValue(v)
	}
	return ret
}

// Get 实现了 [method]fields.has。
func (i *fieldsImpl) Has(_ context.Context, this Fields, name FieldKey) bool {
	i.hm.LockFields()
	defer i.hm.UnlockFields()
	f, ok := i.hm.Fields.Get(this)
	if !ok || f == nil {
		return false
	}
	_, exists := f[strings.ToLower(name)]
	return exists
}

// Set 实现了 [method]fields.set。
func (i *fieldsImpl) Set(_ context.Context, this Fields, name FieldKey, value []FieldValue) witgo.Result[witgo.Unit, HeaderError] {
	i.hm.LockFields()
	defer i.hm.UnlockFields()
	if !validFieldName(name) {
		return witgo.Err[witgo.Unit, HeaderError](HeaderError{InvalidSyntax: &witgo.Unit{}})
	}
	if forbiddenFieldName(name) {
		return witgo.Err[witgo.Unit, HeaderError](HeaderError{Forbidden: &witgo.Unit{}})
	}
	for _, v := range value {
		if !validFieldValue(v) {
			return witgo.Err[witgo.Unit, HeaderError](HeaderError{InvalidSyntax: &witgo.Unit{}})
		}
	}
	f, ok := i.hm.Fields.Get(this)
	if !ok || f == nil || i.hm.IsFieldsImmutable(this) {
		return witgo.Err[witgo.Unit, HeaderError](HeaderError{Immutable: &witgo.Unit{}})
	}
	values := make([]string, len(value))
	for j, v := range value {
		values[j] = string(v)
	}
	f[strings.ToLower(name)] = values
	return witgo.Ok[witgo.Unit, HeaderError](witgo.Unit{})
}

// Delete 实现了 [method]fields.delete。
func (i *fieldsImpl) Delete(_ context.Context, this Fields, name FieldKey) witgo.Result[witgo.Unit, HeaderError] {
	i.hm.LockFields()
	defer i.hm.UnlockFields()
	f, ok := i.hm.Fields.Get(this)
	if !ok || f == nil || i.hm.IsFieldsImmutable(this) {
		return witgo.Err[witgo.Unit, HeaderError](HeaderError{Immutable: &witgo.Unit{}})
	}
	delete(f, strings.ToLower(name))
	return witgo.Ok[witgo.Unit, HeaderError](witgo.Unit{})
}

// Append 实现了 [method]fields.append。
func (i *fieldsImpl) Append(_ context.Context, this Fields, name FieldKey, value FieldValue) witgo.Result[witgo.Unit, HeaderError] {
	i.hm.LockFields()
	defer i.hm.UnlockFields()
	if !validFieldName(name) || !validFieldValue(value) {
		return witgo.Err[witgo.Unit, HeaderError](HeaderError{InvalidSyntax: &witgo.Unit{}})
	}
	if forbiddenFieldName(name) {
		return witgo.Err[witgo.Unit, HeaderError](HeaderError{Forbidden: &witgo.Unit{}})
	}
	f, ok := i.hm.Fields.Get(this)
	if !ok || f == nil || i.hm.IsFieldsImmutable(this) {
		return witgo.Err[witgo.Unit, HeaderError](HeaderError{Immutable: &witgo.Unit{}})
	}
	key := strings.ToLower(name)
	f[key] = append(f[key], string(value))
	return witgo.Ok[witgo.Unit, HeaderError](witgo.Unit{})
}

// Entries 实现了 [method]fields.entries。
func (i *fieldsImpl) Entries(_ context.Context, this Fields) []witgo.Tuple[FieldKey, FieldValue] {
	i.hm.LockFields()
	defer i.hm.UnlockFields()
	f, ok := i.hm.Fields.Get(this)
	if !ok || f == nil {
		return nil
	}
	// 为了确保稳定的返回顺序，我们对 key 进行排序。
	keys := make([]FieldKey, 0, len(f))
	for k := range f {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var entries []witgo.Tuple[FieldKey, FieldValue]
	for _, k := range keys {
		for _, v := range f[k] {
			entries = append(entries, witgo.Tuple[FieldKey, FieldValue]{F0: k, F1: FieldValue(v)})
		}
	}
	return entries
}

// Clone 实现了 [method]fields.clone。
func (i *fieldsImpl) Clone(_ context.Context, this Fields) Fields {
	i.hm.LockFields()
	defer i.hm.UnlockFields()
	f, ok := i.hm.Fields.Get(this)
	if !ok || f == nil {
		// 如果源句柄无效，创建一个空的 fields 并返回。
		return i.Constructor()
	}

	// 创建一个新的 map 并深拷贝所有键值对。
	newFields := make(manager_http.Fields, len(f))
	for k, v := range f {
		newValues := make([]string, len(v))
		copy(newValues, v)
		newFields[k] = newValues
	}
	// clone 得到的是可变副本，不继承 immutable
	return i.hm.Fields.Add(newFields)
}

func forbiddenFieldName(name string) bool {
	switch strings.ToLower(name) {
	case "connection", "keep-alive", "host", "te", "trailer", "transfer-encoding", "upgrade", "proxy-connection":
		return true
	default:
		return false
	}
}
