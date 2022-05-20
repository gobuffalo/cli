package plugdeps

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/gobuffalo/meta"
)

// bin, module version, and tag
var re = regexp.MustCompile(`.*(buffalo-[^/@]+)/?(v[0-9]+)?@?(.*)?`)

// Plugin represents a Go plugin for Buffalo applications
type Plugin struct {
	Binary   string         `toml:"binary" json:"binary"`
	GoGet    string         `toml:"go_get,omitempty" json:"go_get,omitempty"`
	Local    string         `toml:"local,omitempty" json:"local,omitempty"`
	Commands []Command      `toml:"command,omitempty" json:"commands,omitempty"`
	Tags     meta.BuildTags `toml:"tags,omitempty" json:"tags,omitempty"`
}

// String implementation of fmt.Stringer
func (p Plugin) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func (p Plugin) key() string {
	// p.Binary should be uniq and it can be used as the key
	return p.Binary
}

func NewPlugin(mod string, tags ...meta.BuildTags) Plugin {
	mod = strings.TrimSpace(mod)
	match := re.FindStringSubmatch(mod)
	bin := match[1]
	tag := match[3]
	plug := Plugin{
		Binary: bin,
		GoGet:  mod,
	}
	if len(tags) > 0 {
		plug.Tags = tags[0]
	}
	if _, err := os.Stat(mod); err == nil {
		plug.Local = mod
		plug.GoGet = ""
	}
	if plug.GoGet != "" && tag == "" {
		plug.GoGet = mod + "@latest"
	}
	return plug
}
