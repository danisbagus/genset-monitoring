# CORS Implementation & Deep Dive

This document explains the CORS (Cross-Origin Resource Sharing) implementation in the Genset Monitoring System and provides technical context on why these changes were necessary.

## 🚀 Implementation Summary

- **Package**: Used `github.com/gin-contrib/cors` for a production-ready implementation.
- **Middleware**: Created `internal/middleware/cors.go` to centralize CORS logic.
- **Configurable**: Allowed origins are now configurable via `CORS_ALLOWED_ORIGINS` in `.env`.
- **Global**: Applied globally in `cmd/api/main.go` to protect all endpoints.

---

## 🔍 Understanding CORS

### 1. Why the Preflight `OPTIONS` Request Happens
When a browser makes a "non-simple" cross-origin request, it first sends an `OPTIONS` request (the preflight) to the server.
- **Simple Requests**: GET/POST with basic headers (like `text/plain`).
- **Non-Simple Requests**: Any request using `Authorization` headers, `application/json` content type, or methods like `PUT`/`DELETE`.

The preflight checks if the server is willing to accept the actual request. If the server doesn't respond with the correct headers (e.g., `Access-Control-Allow-Origin`), the browser blocks the actual request.

### 2. Handling JWT Authorization Headers
Since we use JWT for authentication, the frontend sends the `Authorization: Bearer <token>` header.
- This automatically triggers a **Preflight** request.
- Our middleware explicitly allows the `Authorization` header in `AllowHeaders`.
- Without this, even if the origin is allowed, the browser would strip the header or block the request because it's considered "unsafe".

### 3. WebSocket Compatibility
WebSockets (`ws://`) don't use CORS in the same way as HTTP. However:
- The **initial handshake** is an HTTP request.
- The browser will include the `Origin` header in that handshake.
- If your CORS middleware is too restrictive, it might block the handshake if it doesn't recognize the origin.
- Our implementation uses `AllowCredentials: true`, which is essential for authenticated WebSocket handshakes if cookies or specific headers are used.

---

## ⚠️ Common CORS Mistakes

1. **Wildcard `*` with Credentials**: You cannot use `AllowOrigins: ["*"]` while also setting `AllowCredentials: true`. Browsers will block this for security reasons.
2. **Missing `OPTIONS` Handler**: Manual middleware often forgets to return a `204 No Content` or `200 OK` for `OPTIONS` requests, causing the preflight to fail.
3. **Ordering in Middleware**: CORS must be one of the **first** middlewares. If an earlier middleware (like Auth) returns a `401` before CORS headers are set, the browser will report a CORS error instead of the actual `401` error, making debugging difficult.
4. **Case Sensitivity**: Headers like `Authorization` are technically case-insensitive in HTTP, but some older proxies or custom logic might expect specific casing.

---

## 💡 Key Concepts

### Why Frontend alone cannot solve CORS?
CORS is a **browser-enforced security feature**. 
- It is NOT a server-side security mechanism (a hacker using `curl` doesn't care about CORS).
- It protects users from **Cross-Site Request Forgery (CSRF)** where a malicious site tries to make requests to another site where you are logged in.
- The frontend can only "request" permission; the **Backend** must explicitly grant it via response headers.

### Development vs. Production Setup

| Feature | Development | Production |
| :--- | :--- | :--- |
| **Allowed Origins** | `http://localhost:5173` | `https://dashboard.genset-monitoring.com` |
| **Credentials** | Often allowed for local testing | Only allowed if using Cookies/Sessions |
| **MaxAge** | Short (for frequent changes) | Long (to reduce preflight overhead) |
| **Security** | Flexible | Strict (no wildcards) |

---

## 🛠️ How to Add New Origins
If you deploy to a new domain, simply update your `.env` file:
```env
CORS_ALLOWED_ORIGINS=https://app.example.com,https://api.example.com
```
*Note: The implementation uses `v.GetStringSlice`, so it expects comma-separated values.*
