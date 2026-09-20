package io

import witgo "github.com/456vv/x/vweb_dynamic/internal/wazero/witgo"

type ErrorManager = witgo.ResourceManager[error]

func NewErrorManager() *ErrorManager {
	return witgo.NewResourceManager[error](nil)
}
