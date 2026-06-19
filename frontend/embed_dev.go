//go:build !production

package frontend

import "embed"

// FS is an empty filesystem used in development and tests.
// The router handles a missing index.html gracefully — it returns a
// "frontend not built" message instead of crashing.
// For a real frontend, build with: cd frontend && npm run build && go build -tags production ./main
var FS embed.FS
