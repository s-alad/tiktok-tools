package bucket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucket(t *testing.T) {
	bucket := Create("000")
	assert.Equal(t, "000", bucket.ID())
}
