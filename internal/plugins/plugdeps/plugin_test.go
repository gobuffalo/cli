package plugdeps

import (
	"testing"

	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func TestNewPlugin_NoVersionNoTag(t *testing.T) {
	r := require.New(t)

	p := NewPlugin("example.com/user/buffalo-awesome")
	r.Equal("buffalo-awesome", p.key())
	r.Equal("buffalo-awesome", p.Binary)
	r.Equal("example.com/user/buffalo-awesome@latest", p.GoGet)

	r.Equal(meta.BuildTags(nil), p.Tags)
	r.Equal(`{"binary":"buffalo-awesome","go_get":"example.com/user/buffalo-awesome@latest"}`, p.String())
}

func TestNewPlugin_BuildTags(t *testing.T) {
	r := require.New(t)

	p := NewPlugin("example.com/user/buffalo-awesome", meta.BuildTags{"sqlite", "awesome"})
	r.Equal("buffalo-awesome", p.key())
	r.Equal("buffalo-awesome", p.Binary)
	r.Equal("example.com/user/buffalo-awesome@latest", p.GoGet)

	r.Equal(meta.BuildTags{"sqlite", "awesome"}, p.Tags)
	r.Equal(`{"binary":"buffalo-awesome","go_get":"example.com/user/buffalo-awesome@latest","tags":["sqlite","awesome"]}`, p.String())
}

func TestNewPlugin_NoVersionTag(t *testing.T) {
	r := require.New(t)

	p := NewPlugin("example.com/user/buffalo-awesome@v3.1")
	r.Equal("buffalo-awesome", p.key())
	r.Equal("buffalo-awesome", p.Binary)
	r.Equal("example.com/user/buffalo-awesome@v3.1", p.GoGet)
}

func TestNewPlugin_VersionNoTag(t *testing.T) {
	r := require.New(t)

	p := NewPlugin("example.com/user/buffalo-awesome/v3")
	r.Equal("buffalo-awesome", p.key())
	r.Equal("buffalo-awesome", p.Binary)
	r.Equal("example.com/user/buffalo-awesome/v3@latest", p.GoGet)
}

func TestNewPlugin_VersionTag(t *testing.T) {
	r := require.New(t)

	p := NewPlugin("example.com/user/buffalo-awesome/v3@v3.1")
	r.Equal("buffalo-awesome", p.key())
	r.Equal("buffalo-awesome", p.Binary)
	r.Equal("example.com/user/buffalo-awesome/v3@v3.1", p.GoGet)
}
