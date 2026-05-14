# API Path Configuration Guide

This guide explains how to correctly configure API paths in the Genset Monitoring project to avoid common issues like path duplication and CORS errors.

## 1. Why Path Duplication Happens

The error `POST http://localhost:5173/api/api/v1/login` occurs due to a mismatch between **Axios Base URL** and **Service Endpoints**.

*   **Axios Base URL**: `/api` (from `VITE_API_BASE_URL`)
*   **Service Endpoint**: `/api/v1/auth/login` (incorrectly prefixed)
*   **Result**: `/api` + `/api/v1/auth/login` = `/api/api/v1/auth/login`

Axios treats any path starting with `/` as relative to the `baseURL`. If you include the prefix in both places, it will be duplicated.

---

## 2. The Relationship Chain

To reach `http://localhost:8080/api/v1/auth/login`, the following chain must be correctly configured:

1.  **Axios Config (`axios.ts`)**:
    *   `baseURL: '/api'`
    *   Ensures all requests go through the local proxy.

2.  **API Constants (`api.constant.ts`)**:
    *   `LOGIN: '/v1/auth/login'`
    *   **Do not** include the `/api` prefix here.

3.  **Vite Proxy (`vite.config.ts`)**:
    *   Matches: `/api`
    *   Target: `http://localhost:8080`
    *   The browser sends `/api/v1/auth/login` to Vite. Vite sees `/api`, keeps it (since no rewrite is specified), and forwards the full path to the backend.

---

## 3. Correct Pattern for Scalability

We use a centralized constant system to manage endpoints. This ensures that if the API version changes (e.g., to `/v2`), you only need to update one file.

### Example: Adding a New Module

If you add an Alarms module:
1.  Add to `api.constant.ts`:
    ```typescript
    ALARMS: {
      LIST: '/v1/alarms',
      ACKNOWLEDGE: (id: string) => `/v1/alarms/${id}/ack`,
    }
    ```
2.  Use in service:
    ```typescript
    import { API_ENDPOINTS } from './api.constant'
    api.get(API_ENDPOINTS.ALARMS.LIST)
    ```

---

## 4. Common Mistakes to Avoid

*   **Hardcoding Localhost**: Never use `http://localhost:8080` inside your `.vue` or `.ts` files. Use relative paths and let the proxy/environment handle the target.
*   **Manual Path Joining**: Avoid string concatenation like `baseURL + '/path'`. Axios handles this automatically and correctly.
*   **Absolute URLs in Services**: If you pass a full URL like `api.get('http://google.com')`, Axios will ignore the `baseURL` entirely.
*   **Missing Proxy Rewrite Knowledge**: If your backend *doesn't* have an `/api` prefix but you want to use `/api` on frontend for proxying, you **must** use `rewrite: (path) => path.replace(/^\/api/, '')` in `vite.config.ts`. In this project, the backend **does** use `/api`, so no rewrite is needed.
