package orchestrator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	q := queue{}
	q.Enqueue("5", "10")
	v, ok := q.TryDequeue()
	assert.Equal(t, ok, true)
	assert.Equal(t, v, "5")
	v, ok = q.TryDequeue()
	assert.Equal(t, ok, true)
	assert.Equal(t, v, "10")
	v, ok = q.TryDequeue()
	assert.Equal(t, ok, false)
}