package lib

import (
	"github.com/hyperioxx/goscript/pkg/core/virtualmachine"
)

type System struct {
	virtualmachine.Module
}

func NewSystemModule() virtualmachine.Object {
	return &System{}
}
