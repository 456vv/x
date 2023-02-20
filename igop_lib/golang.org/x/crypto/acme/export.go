// export by github.com/goplus/igop/cmd/qexp

//go:build igop_lib
// +build igop_lib

package acme

import (
	q "golang.org/x/crypto/acme"

	"go/constant"
	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "acme",
		Path: "golang.org/x/crypto/acme",
		Deps: map[string]string{
			"bytes":            "bytes",
			"context":          "context",
			"crypto":           "crypto",
			"crypto/ecdsa":     "ecdsa",
			"crypto/elliptic":  "elliptic",
			"crypto/hmac":      "hmac",
			"crypto/rand":      "rand",
			"crypto/rsa":       "rsa",
			"crypto/sha256":    "sha256",
			"crypto/sha512":    "sha512",
			"crypto/tls":       "tls",
			"crypto/x509":      "x509",
			"crypto/x509/pkix": "pkix",
			"encoding/asn1":    "asn1",
			"encoding/base64":  "base64",
			"encoding/hex":     "hex",
			"encoding/json":    "json",
			"encoding/pem":     "pem",
			"errors":           "errors",
			"fmt":              "fmt",
			"io":               "io",
			"math/big":         "big",
			"net/http":         "http",
			"runtime/debug":    "debug",
			"strconv":          "strconv",
			"strings":          "strings",
			"sync":             "sync",
			"time":             "time",
		},
		Interfaces: map[string]reflect.Type{
			"CertOption":  reflect.TypeOf((*q.CertOption)(nil)).Elem(),
			"OrderOption": reflect.TypeOf((*q.OrderOption)(nil)).Elem(),
		},
		NamedTypes: map[string]reflect.Type{
			"Account":                reflect.TypeOf((*q.Account)(nil)).Elem(),
			"Authorization":          reflect.TypeOf((*q.Authorization)(nil)).Elem(),
			"AuthorizationError":     reflect.TypeOf((*q.AuthorizationError)(nil)).Elem(),
			"AuthzID":                reflect.TypeOf((*q.AuthzID)(nil)).Elem(),
			"CRLReasonCode":          reflect.TypeOf((*q.CRLReasonCode)(nil)).Elem(),
			"Challenge":              reflect.TypeOf((*q.Challenge)(nil)).Elem(),
			"Client":                 reflect.TypeOf((*q.Client)(nil)).Elem(),
			"Directory":              reflect.TypeOf((*q.Directory)(nil)).Elem(),
			"Error":                  reflect.TypeOf((*q.Error)(nil)).Elem(),
			"ExternalAccountBinding": reflect.TypeOf((*q.ExternalAccountBinding)(nil)).Elem(),
			"KeyID":                  reflect.TypeOf((*q.KeyID)(nil)).Elem(),
			"Order":                  reflect.TypeOf((*q.Order)(nil)).Elem(),
			"OrderError":             reflect.TypeOf((*q.OrderError)(nil)).Elem(),
			"Subproblem":             reflect.TypeOf((*q.Subproblem)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"ErrAccountAlreadyExists": reflect.ValueOf(&q.ErrAccountAlreadyExists),
			"ErrNoAccount":            reflect.ValueOf(&q.ErrNoAccount),
			"ErrUnsupportedKey":       reflect.ValueOf(&q.ErrUnsupportedKey),
		},
		Funcs: map[string]reflect.Value{
			"AcceptTOS":          reflect.ValueOf(q.AcceptTOS),
			"DomainIDs":          reflect.ValueOf(q.DomainIDs),
			"IPIDs":              reflect.ValueOf(q.IPIDs),
			"JWKThumbprint":      reflect.ValueOf(q.JWKThumbprint),
			"RateLimit":          reflect.ValueOf(q.RateLimit),
			"WithKey":            reflect.ValueOf(q.WithKey),
			"WithOrderNotAfter":  reflect.ValueOf(q.WithOrderNotAfter),
			"WithOrderNotBefore": reflect.ValueOf(q.WithOrderNotBefore),
			"WithTemplate":       reflect.ValueOf(q.WithTemplate),
		},
		TypedConsts: map[string]igop.TypedConst{
			"CRLReasonAACompromise":         {reflect.TypeOf(q.CRLReasonAACompromise), constant.MakeInt64(int64(q.CRLReasonAACompromise))},
			"CRLReasonAffiliationChanged":   {reflect.TypeOf(q.CRLReasonAffiliationChanged), constant.MakeInt64(int64(q.CRLReasonAffiliationChanged))},
			"CRLReasonCACompromise":         {reflect.TypeOf(q.CRLReasonCACompromise), constant.MakeInt64(int64(q.CRLReasonCACompromise))},
			"CRLReasonCertificateHold":      {reflect.TypeOf(q.CRLReasonCertificateHold), constant.MakeInt64(int64(q.CRLReasonCertificateHold))},
			"CRLReasonCessationOfOperation": {reflect.TypeOf(q.CRLReasonCessationOfOperation), constant.MakeInt64(int64(q.CRLReasonCessationOfOperation))},
			"CRLReasonKeyCompromise":        {reflect.TypeOf(q.CRLReasonKeyCompromise), constant.MakeInt64(int64(q.CRLReasonKeyCompromise))},
			"CRLReasonPrivilegeWithdrawn":   {reflect.TypeOf(q.CRLReasonPrivilegeWithdrawn), constant.MakeInt64(int64(q.CRLReasonPrivilegeWithdrawn))},
			"CRLReasonRemoveFromCRL":        {reflect.TypeOf(q.CRLReasonRemoveFromCRL), constant.MakeInt64(int64(q.CRLReasonRemoveFromCRL))},
			"CRLReasonSuperseded":           {reflect.TypeOf(q.CRLReasonSuperseded), constant.MakeInt64(int64(q.CRLReasonSuperseded))},
			"CRLReasonUnspecified":          {reflect.TypeOf(q.CRLReasonUnspecified), constant.MakeInt64(int64(q.CRLReasonUnspecified))},
		},
		UntypedConsts: map[string]igop.UntypedConst{
			"ALPNProto":         {"untyped string", constant.MakeString(string(q.ALPNProto))},
			"LetsEncryptURL":    {"untyped string", constant.MakeString(string(q.LetsEncryptURL))},
			"StatusDeactivated": {"untyped string", constant.MakeString(string(q.StatusDeactivated))},
			"StatusExpired":     {"untyped string", constant.MakeString(string(q.StatusExpired))},
			"StatusInvalid":     {"untyped string", constant.MakeString(string(q.StatusInvalid))},
			"StatusPending":     {"untyped string", constant.MakeString(string(q.StatusPending))},
			"StatusProcessing":  {"untyped string", constant.MakeString(string(q.StatusProcessing))},
			"StatusReady":       {"untyped string", constant.MakeString(string(q.StatusReady))},
			"StatusRevoked":     {"untyped string", constant.MakeString(string(q.StatusRevoked))},
			"StatusUnknown":     {"untyped string", constant.MakeString(string(q.StatusUnknown))},
			"StatusValid":       {"untyped string", constant.MakeString(string(q.StatusValid))},
		},
	})
}
