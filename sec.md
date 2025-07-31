✅ 1. Authentication

Use secure login (username/email + password)

Passwords hashed (use bcrypt)

Session-based auth with secure cookies (e.g. gorilla/sessions or gin-contrib/sessions)

    CSRF protection (see below)

✅ 2. Authorization

Protect routes (e.g. admin-only, logged-in only)

Middleware to check user roles / permissions

    Avoid IDOR (don’t let users access resources they don’t own)

✅ 3. Session & Cookie Security

Set HttpOnly, Secure, and SameSite attributes

Encrypt cookie/session values (e.g. use securecookie or signed cookies)

    Rotate session IDs after login

✅ 4. CSRF Protection

Use CSRF tokens (manually or with gin-contrib/sessions + your own token system)

    HTMX has built-in CSRF support – send token in a header

<meta name="csrf-token" content="{{ .CSRF }}">
<script>
  htmx.defaults.headers['X-CSRF-Token'] = document.querySelector('meta[name="csrf-token"]').content
</script>

✅ 5. Input Validation & Sanitization

Validate all user input (use go-playground/validator or custom)

Prevent overposting by only binding allowed fields

    Escape or reject unsafe input if outputting HTML

✅ 6. Error Handling & Leaks

Never expose internal errors to users (e.g., GORM errors)

Sanitize error messages

    Return generic errors to UI, log the real cause

✅ 7. Secure Headers

Add headers like:

    Content-Security-Policy

    X-Frame-Options

    X-Content-Type-Options

    Strict-Transport-Security

    Can be added via middleware

✅ 8. Rate Limiting / Brute-force protection

Prevent login abuse with rate-limiting (e.g. IP-based middleware)

    Optional: CAPTCHA for sensitive actions

✅ 9. HTTPS

Always use HTTPS in production

    Redirect HTTP to HTTPS

✅ 10. Database Security

Use GORM parameterized queries (default behavior)

Avoid raw SQL unless fully controlled

Enforce foreign key constraints

    Validate before delete/update (e.g. check ownership)

✅ 11. File Uploads (if any)

Validate file type and size

Sanitize filenames

    Store uploaded files securely

✅ 12. Deployment/Secrets

Do not commit .env or secrets

Use a secure method to store and read secrets (env vars, Vault, etc.)

    Keep dependencies up-to-date


