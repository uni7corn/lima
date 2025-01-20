package limayaml

import (
	"testing"

	"github.com/lima-vm/lima/pkg/ptr"
	"gotest.tools/v3/assert"
)

func TestMarshalEmpty(t *testing.T) {
	_, err := Marshal(&LimaYAML{}, false)
	assert.NilError(t, err)
}

func TestMarshalTilde(t *testing.T) {
	y := LimaYAML{
		Mounts: []Mount{
			{Location: "~", Writable: ptr.Of(false)},
			{Location: "/tmp/lima", Writable: ptr.Of(true)},
			{Location: "null"},
		},
	}
	b, err := Marshal(&y, true)
	assert.NilError(t, err)
	// yaml will load ~ (or null) as null
	// make sure that it is always quoted
	assert.Equal(t, string(b), `---
images: []
mounts:
- location: "~"
  writable: false
- location: /tmp/lima
  writable: true
- location: "null"
...
`)
}
