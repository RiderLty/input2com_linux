---
name: add-api-endpoint
description: Add new REST API endpoint to config server
---

# Add API Endpoint

This skill covers adding new REST API endpoints to the config server.

## When to Use

- Adding new configuration options
- Creating new control features
- Adding status endpoints
- Implementing new functionality accessible via web UI

## Architecture Overview

- **Backend**: `config_server.go` - HTTP handlers
- **Frontend**: `server/src/api.js` - API client functions
- **UI**: `server/src/App.jsx` - React components

## Workflow

### 1. Add Backend Endpoint

In `config_server.go`, add new handler:

```go
http.HandleFunc("/api/your-endpoint", func(w http.ResponseWriter, r *http.Request) {
    // Parse query parameters
    param := r.URL.Query().Get("param")
    
    // Perform action
    result := doSomething(param)
    
    // Return JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
})
```

### 2. Add Frontend API Function

In `server/src/api.js`:

```javascript
export const yourEndpoint = async (param) => {
    const response = await fetch(`/api/your-endpoint?param=${param}`)
    return response.json()
}
```

### 3. Add UI Component

In `server/src/App.jsx`:

```jsx
import { yourEndpoint } from './api'

// In component
const handleYourAction = async () => {
    await yourEndpoint('value')
    // Update UI state
}
```

### 4. Rebuild and Test

```bash
# Build frontend
cd server && yarn build

# Rebuild Go binary
go build -ldflags="-s -w" -o input2com

# Test endpoint
curl "http://localhost:9264/api/your-endpoint?param=value"
```

## Important Notes

- **CORS**: Config server handles CORS automatically for same-origin requests
- **Error handling**: Always check for errors in both backend and frontend
- **Response format**: Use JSON for structured data, plain text for simple responses
- **URL encoding**: Use `r.URL.Query().Get()` for query parameters
- **Content-Type**: Set `application/json` for JSON responses

## Common Patterns

### GET Endpoint (Read)
```go
http.HandleFunc("/api/get/config", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(configData)
})
```

### POST Endpoint (Update)
```go
http.HandleFunc("/api/set/config", func(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    value := r.Form.Get("value")
    // Update config
    w.WriteHeader(http.StatusOK)
})
```

### Action Endpoint (Trigger)
```go
http.HandleFunc("/api/action/restart", func(w http.ResponseWriter, r *http.Request) {
    // Perform action
    go func() {
        time.Sleep(100 * time.Millisecond)
        os.Exit(0)
    }()
    w.WriteHeader(http.StatusOK)
})
```

## Related Files

- `config_server.go` - Backend HTTP handlers
- `server/src/api.js` - Frontend API functions
- `server/src/App.jsx` - React UI components
- `main.go` - Server startup and configuration

## Testing

1. Test with `curl` command
2. Check browser developer tools for network requests
3. Verify UI updates correctly
4. Test error conditions