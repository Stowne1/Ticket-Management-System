//go:build production

package frontend

import "embed"

// FS holds the built React app embedded at compile time.
// Only active when built with -tags production (see dockerfile).
//
//go:embed all:dist
var FS embed.FS
