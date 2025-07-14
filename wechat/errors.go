package wechat

import "fmt"

type MsgError[T any] struct {
	v *T
}

func (m *MsgError[T]) Error() string {
	return fmt.Sprintf("%#v", m.v)
}

func (m *MsgError[T]) Value() *T {
	return m.v
}
