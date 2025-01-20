package cidata

import (
	"io"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

var defaultRemoveDefaults = false

func TestConfig(t *testing.T) {
	args := &TemplateArgs{
		Name:    "default",
		User:    "foo",
		UID:     501,
		Comment: "Foo",
		Home:    "/home/foo.linux",
		SSHPubKeys: []string{
			"ssh-rsa dummy foo@example.com",
		},
		MountType: "reverse-sshfs",
	}
	config, err := ExecuteTemplateCloudConfig(args)
	assert.NilError(t, err)
	t.Log(string(config))
	assert.Assert(t, !strings.Contains(string(config), "ca_certs:"))
}

func TestConfigCACerts(t *testing.T) {
	args := &TemplateArgs{
		Name:    "default",
		User:    "foo",
		UID:     501,
		Comment: "Foo",
		Home:    "/home/foo.linux",
		SSHPubKeys: []string{
			"ssh-rsa dummy foo@example.com",
		},
		MountType: "reverse-sshfs",
		CACerts: CACerts{
			RemoveDefaults: &defaultRemoveDefaults,
		},
	}
	config, err := ExecuteTemplateCloudConfig(args)
	assert.NilError(t, err)
	t.Log(string(config))
	assert.Assert(t, strings.Contains(string(config), "ca_certs:"))
}

func TestTemplate(t *testing.T) {
	args := &TemplateArgs{
		Name: "default",
		User: "foo",
		UID:  501,
		Home: "/home/foo.linux",
		SSHPubKeys: []string{
			"ssh-rsa dummy foo@example.com",
		},
		Mounts: []Mount{
			{MountPoint: "/Users/dummy"},
			{MountPoint: "/Users/dummy/lima"},
		},
		MountType: "reverse-sshfs",
		CACerts: CACerts{
			RemoveDefaults: &defaultRemoveDefaults,
			Trusted:        []Cert{},
		},
	}
	layout, err := ExecuteTemplateCIDataISO(args)
	assert.NilError(t, err)
	for _, f := range layout {
		t.Logf("=== %q ===", f.Path)
		b, err := io.ReadAll(f.Reader)
		assert.NilError(t, err)
		t.Log(string(b))
		if f.Path == "user-data" {
			// mounted later
			assert.Assert(t, !strings.Contains(string(b), "mounts:"))
			// ca_certs:
			assert.Assert(t, !strings.Contains(string(b), "trusted:"))
		}
	}
}

func TestTemplate9p(t *testing.T) {
	args := &TemplateArgs{
		Name: "default",
		User: "foo",
		UID:  501,
		Home: "/home/foo.linux",
		SSHPubKeys: []string{
			"ssh-rsa dummy foo@example.com",
		},
		Mounts: []Mount{
			{Tag: "mount0", MountPoint: "/Users/dummy", Type: "9p", Options: "ro,trans=virtio"},
			{Tag: "mount1", MountPoint: "/Users/dummy/lima", Type: "9p", Options: "rw,trans=virtio"},
		},
		MountType: "9p",
		CACerts: CACerts{
			RemoveDefaults: &defaultRemoveDefaults,
		},
	}
	layout, err := ExecuteTemplateCIDataISO(args)
	assert.NilError(t, err)
	for _, f := range layout {
		t.Logf("=== %q ===", f.Path)
		b, err := io.ReadAll(f.Reader)
		assert.NilError(t, err)
		t.Log(string(b))
		if f.Path == "user-data" {
			// mounted at boot
			assert.Assert(t, strings.Contains(string(b), "mounts:"))
		}
	}
}
