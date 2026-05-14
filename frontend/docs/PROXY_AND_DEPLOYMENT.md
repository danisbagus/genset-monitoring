# Frontend Proxy and Deployment Guide

This document explains the proxy configuration implemented for the Genset Monitoring frontend and how it handles CORS and production deployment.

## 1. How the Development Proxy Works

When you run `npm run dev`, Vite starts a development server (usually at `http://localhost:5173`). Without a proxy, if your frontend tries to call `http://localhost:8080/api/v1/login`, the browser will block the request due to **CORS (Cross-Origin Resource Sharing)** because the ports (5173 vs 8080) are different.

### The Proxy Flow:
1. **Frontend Request**: Your code calls `/api/v1/login`.
2. **Relative Path**: Since it's a relative path, the browser sends the request to the same origin: `http://localhost:5173/api/v1/login`.
3. **Vite Interception**: Vite sees the `/api` prefix and, based on `vite.config.ts`, intercepts the request.
4. **Forwarding**: Vite forwards the request to the `target` (`http://localhost:8080`).
5. **Response**: The backend responds to Vite, and Vite passes the response back to your browser.

**Result**: The browser thinks the request never left `localhost:5173`, so no CORS issues occur!

---

## 2. Why Only in Development?

The Vite proxy is part of the Vite **Development Server**. When you build your project for production (`npm run build`), Vite generates static HTML, JS, and CSS files. **There is no Vite server running in production** to handle proxying.

In production, you typically use a real web server like **Nginx** or **Apache** to serve your static files and handle the proxying to your backend.

---

## 3. Why Backend CORS is Still Required

Even with a proxy, you should keep CORS configured on your backend for several reasons:

1. **Direct Access**: Some clients might access the API directly (not through your frontend).
2. **Security**: CORS is a browser-side security feature. A properly configured backend ensures that only authorized domains can interact with your API even if someone tries to bypass your proxy.
3. **Multi-domain support**: If you ever decide to host your frontend on `dashboard.example.com` and your API on `api.example.com` without a shared reverse proxy, you will *need* backend CORS.

---

## 4. Production Deployment with Nginx

When deploying to production, you should configure Nginx to act as a Reverse Proxy, mimicking the behavior of the Vite proxy.

### Sample Nginx Configuration (`/etc/nginx/conf.d/default.conf`)

```nginx
server {
    listen 80;
    server_name monitoring.example.com;

    # 1. Serve Frontend Static Files
    location / {
        root /usr/share/nginx/html;
        index index.html;
        try_files $uri $uri/ /index.html;
    }

    # 2. Proxy API Requests
    location /api {
        proxy_pass http://backend-api:8080; # or http://localhost:8080
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # 3. Proxy WebSocket Requests
    location /ws {
        proxy_pass http://backend-api:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
    }
}
```

By using `/api` and `/ws` as relative paths in your frontend code (as we configured in `.env.production`), the same code works perfectly in both development (via Vite) and production (via Nginx).
