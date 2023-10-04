package example

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestList(t *testing.T) {

	l := NewList()
	l.PushBack(1)
	l.PushBack(2)
	assert.Equal(t, 1, l.Front())
	assert.Equal(t, 2, l.Back())

	l.PopFront()
	l.PopFront()
}
