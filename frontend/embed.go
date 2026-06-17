package frontend

import "embed"

// FS holds the built React app (frontend/dist).
// The router uses this to serve the frontend as static files.
// Build the frontend first with: cd frontend && npm run build
//
//go:embed all:dist
var FS embed.FS
