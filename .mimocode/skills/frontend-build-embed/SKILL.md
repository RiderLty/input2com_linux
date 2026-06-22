---
name: frontend-build-embed
description: Build Vite frontend and embed into Go binary
---

# Frontend Build and Embed

This skill covers the complete workflow for building the React/Vite frontend and embedding it into the Go binary.

## When to Use

- After any frontend code changes (React components, styles, API calls)
- When adding new frontend features or UI elements
- After modifying `server/src/` files

## Prerequisites

- Node.js and Yarn installed
- Go compiler available
- Frontend dependencies installed (`cd server && yarn install`)

## Workflow

### 1. Build Frontend

```bash
cd server && yarn build
```

This generates production files in `server/build/`.

### 2. Rebuild Go Binary

```bash
go build -ldflags="-s -w" -o input2com
```

The Go binary embeds the frontend via `//go:embed server/build` directive.

### 3. Verify

- Check that `server/build/` directory exists and contains files
- Verify binary size is reasonable (should include embedded frontend)

## Important Notes

- **Must rebuild Go binary after frontend changes** - Frontend is embedded at compile time
- **Build order matters** - Frontend must be built first, then Go binary
- **Development mode** - Use `cd server && yarn start` for frontend dev server (proxies `/api` to backend)
- **Cross-compilation** - Use `CGO_ENABLED=0 GOOS=linux GOARCH=amd64` for cross-compile

## Common Issues

1. **Frontend changes not showing** - Forgot to rebuild Go binary
2. **Build errors** - Check Node.js version and dependencies
3. **Large binary size** - Ensure proper optimization flags (`-ldflags="-s -w"`)

## Related Files

- `server/src/` - Frontend source code
- `server/build/` - Built frontend files
- `main.go` - Contains `//go:embed server/build` directive
- `config_server.go` - Serves embedded frontend