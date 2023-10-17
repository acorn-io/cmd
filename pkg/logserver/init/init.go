package init

import "github.com/acorn-io/cmd/pkg/logserver"

func init() {
	go logserver.StartServerWithDefaults()
}
