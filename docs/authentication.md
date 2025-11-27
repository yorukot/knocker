# Authentication Guide for Frontend Developers

## Overview

The Knocker authentication system supports:
- Email/password authentication
- OAuth (Google)
- Account linking (add OAuth to existing account)

**Token Strategy**: 
- Refresh tokens are stored in HTTP-only cookies (handled automatically by the browser)
- Access tokens are returned in API responses and should be stored in memory or sessionStorage
- Access tokens are used for authenticated API requests

## Authentication Flow

### Initial Authentication

```
User not authenticated
        ↓
[Register or Login]
        ↓
Receive refresh token cookie (automatic)
        ↓
Call /auth/refresh to get access token
        ↓
Store access token in memory/sessionStorage
        ↓
Use access token for API requests
```

### Token Refresh Flow

```
API request with access token
        ↓
    [401 Error?]
        ↓ Yes
Call /auth/refresh
        ↓
    [Success?]
        ↓ Yes
Get new access token + new refresh token cookie
        ↓
Retry failed API request
        ↓
Continue using new access token

    [401 on refresh?]
        ↓
Redirect to login
```

### When to Refresh Token

**Recommended Strategy**:

1. **On App Load/Page Refresh**
   - Always call `/auth/refresh` when app initializes
   - If successful: user is logged in
   - If fails: user needs to login

2. **On 401 Response**
   - When any API call returns 401
   - Try refreshing token once
   - If refresh succeeds: retry the original request
   - If refresh fails: redirect to login

3. **Proactive Refresh (Optional)**
   - Decode JWT to check `exp` (expiration)
   - Refresh 1-5 minutes before expiration
   - Prevents interruptions during user activity

**Example Implementation**:
```javascript
// Axios interceptor example
axios.interceptors.response.use(
  response => response,
  async error => {
    if (error.response?.status === 401 && !error.config._retry) {
      error.config._retry = true;
      try {
        await refreshToken(); // Call /auth/refresh
        return axios(error.config); // Retry original request
      } catch (refreshError) {
        // Redirect to login
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);
```

## API Endpoints

### 1. Register (Email/Password)

**POST** `/api/auth/register`

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "display_name": "John Doe"
}
```

**Validation**:
- Email: required, valid email format, max 255 chars
- Password: required, min 8 chars, max 255 chars
- Display name: required, min 3 chars, max 255 chars

**Response** (200 OK):
```json
{
  "success": true,
  "message": "User registered successfully"
}
```

**Side Effects**:
- Sets `refresh_token` HTTP-only cookie
- **Next Step**: Call `/auth/refresh` to get access token

---

### 2. Login (Email/Password)

**POST** `/api/auth/login`

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "message": "Login successful"
}
```

**Side Effects**:
- Sets `refresh_token` HTTP-only cookie
- **Next Step**: Call `/auth/refresh` to get access token

---

### 3. Refresh Token

**POST** `/api/auth/refresh`

**Request**: No body needed (uses `refresh_token` cookie automatically)

**Response** (201 Created):
```json
{
  "success": true,
  "message": "Access token refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Side Effects**:
- Issues new `refresh_token` cookie (refresh token rotation)
- Old refresh token is invalidated

**Error Responses**:
- 401: Refresh token invalid/expired/already used → Redirect to login
- 500: Server error → Retry or show error

**Usage**:
- Store the `access_token` from the response
- Use it in `Authorization: Bearer <access_token>` header for API requests
- Call this endpoint whenever you need a new access token

---

### 4. OAuth Login/Register (Google)

**GET** `/api/auth/oauth/google?next=/dashboard`

**Query Parameters**:
- `next` (optional): URL to redirect after successful authentication (default: `/`)

**Flow**:
1. Redirect user to this endpoint
2. User is redirected to Google for authentication
3. After Google auth, user is redirected to callback URL (handled by backend)
4. User is redirected back to your app at the `next` URL
5. `refresh_token` cookie is automatically set
6. **Your app should**: Call `/auth/refresh` to get access token

**Example Usage**:
```javascript
// Redirect to OAuth
window.location.href = '/api/auth/oauth/google?next=/dashboard';

// After redirect back to /dashboard:
// Call /auth/refresh to get access token
```

---

### 5. Link OAuth Account (for logged-in users)

**GET** `/api/auth/oauth/google?next=/settings`

**Requirements**:
- User must already be logged in (have valid refresh token cookie)
- Send request with credentials to include cookies

**Flow**:
Same as OAuth login, but links Google account to existing user instead of creating new user

**Use Case**: User registered with email/password and wants to add Google login

---

## Complete Authentication Flows

### Flow 1: New User Registration (Email)

```
1. POST /api/auth/register
   ↓
2. Receive refresh_token cookie
   ↓
3. POST /api/auth/refresh
   ↓
4. Store access_token from response
   ↓
5. Use access_token for API calls
```

### Flow 2: Existing User Login (Email)

```
1. POST /api/auth/login
   ↓
2. Receive refresh_token cookie
   ↓
3. POST /api/auth/refresh
   ↓
4. Store access_token from response
   ↓
5. Use access_token for API calls
```

### Flow 3: OAuth Login/Register

```
1. Redirect to GET /api/auth/oauth/google?next=/dashboard
   ↓
2. User authenticates with Google
   ↓
3. User redirected back to /dashboard with refresh_token cookie
   ↓
4. POST /api/auth/refresh
   ↓
5. Store access_token from response
   ↓
6. Use access_token for API calls
```

### Flow 4: App Initialization (Check Login Status)

```
On app load:
  ↓
1. POST /api/auth/refresh
   ↓
2. Success (200/201)?
   ↓ Yes              ↓ No
User is logged in    Show login page
   ↓
Store access_token
```

### Flow 5: Handling Expired Access Token

```
API call returns 401
   ↓
1. POST /api/auth/refresh
   ↓
2. Success?
   ↓ Yes              ↓ No
Get new access_token  Clear state & redirect to login
   ↓
Retry original API call
```

## Important Notes

### Cookies & CORS

- Refresh token cookies are HTTP-only (JavaScript cannot access them)
- Must include credentials in fetch/axios requests:
  ```javascript
  // Fetch API
  fetch('/api/auth/refresh', {
    credentials: 'include'
  })
  
  // Axios
  axios.defaults.withCredentials = true;
  ```

### Security Best Practices

1. **Store Access Tokens Securely**:
   - Memory (React state, Vue data) - most secure, lost on refresh
   - sessionStorage - survives page refresh, cleared on tab close
   - ❌ Never use localStorage (vulnerable to XSS)

2. **Token Refresh**:
   - Always call `/auth/refresh` on app initialization
   - Handle 401 errors by refreshing token once, then redirect to login
   - Refresh token rotation means old refresh tokens are invalidated

3. **Logout** (if implemented):
   - Clear stored access token
   - Clear application state
   - Redirect to login page
   - Note: Refresh token cookie will expire automatically

### Error Handling

**Common Error Scenarios**:

- **400 on /register or /login**: Invalid credentials or validation error
- **401 on /refresh**: Refresh token invalid/expired → redirect to login
- **500**: Server error → show error message, allow retry

### Token Reuse Detection

⚠️ **Security Feature**: If a refresh token is used twice, it indicates possible token theft. The backend logs this event. Always ensure you're not calling `/auth/refresh` concurrently from multiple requests.

## Example: React Hook

```javascript
import { useState, useEffect } from 'react';
import axios from 'axios';

// Configure axios
axios.defaults.baseURL = 'http://localhost:8080/api';
axios.defaults.withCredentials = true;

export function useAuth() {
  const [accessToken, setAccessToken] = useState(null);
  const [loading, setLoading] = useState(true);

  // Initialize: check if user is logged in
  useEffect(() => {
    refreshAccessToken();
  }, []);

  // Configure axios interceptor
  useEffect(() => {
    const interceptor = axios.interceptors.response.use(
      response => response,
      async error => {
        if (error.response?.status === 401 && !error.config._retry) {
          error.config._retry = true;
          const newToken = await refreshAccessToken();
          if (newToken) {
            error.config.headers.Authorization = `Bearer ${newToken}`;
            return axios(error.config);
          }
        }
        return Promise.reject(error);
      }
    );
    return () => axios.interceptors.response.eject(interceptor);
  }, []);

  async function refreshAccessToken() {
    try {
      const { data } = await axios.post('/auth/refresh');
      const token = data.data.access_token;
      setAccessToken(token);
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      return token;
    } catch (error) {
      setAccessToken(null);
      delete axios.defaults.headers.common['Authorization'];
      return null;
    } finally {
      setLoading(false);
    }
  }

  async function login(email, password) {
    await axios.post('/auth/login', { email, password });
    await refreshAccessToken();
  }

  async function register(email, password, display_name) {
    await axios.post('/auth/register', { email, password, display_name });
    await refreshAccessToken();
  }

  function logout() {
    setAccessToken(null);
    delete axios.defaults.headers.common['Authorization'];
    // Redirect to login or clear state
  }

  return {
    accessToken,
    loading,
    isAuthenticated: !!accessToken,
    login,
    register,
    logout,
  };
}
```

## Quick Start Checklist

- [ ] Configure axios/fetch to include credentials
- [ ] Call `/auth/refresh` on app initialization
- [ ] Store access token in memory or sessionStorage
- [ ] Add `Authorization: Bearer <token>` header to API requests
- [ ] Handle 401 errors by refreshing token, then retry
- [ ] Redirect to login if refresh fails
- [ ] For OAuth: redirect to `/api/auth/oauth/google?next=<url>`
