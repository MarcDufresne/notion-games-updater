# Security

## Overview

Game Tracker implements multiple layers of security to protect user data and ensure safe operation.

## Authentication & Authorization

### Firebase Authentication
- **Google Sign-In Only**: Currently supports Google OAuth 2.0 for authentication
- **JWT Token Verification**: All API requests (except health check) require valid Firebase ID tokens
- **Token Validation**: Firebase Admin SDK verifies tokens server-side on every request
- **User Isolation**: User ID extracted from verified token; cannot be spoofed

### Authorization Middleware
- Implemented in `internal/middleware/auth.go`
- Extracts and validates Bearer tokens from Authorization headers
- Injects verified user ID into request context
- Returns 401 Unauthorized for invalid/missing tokens

## Data Security

### Database Layer Protection
- **User ID Enforcement**: Database layer (`internal/database/client.go`) overwrites any user_id from client requests with the authenticated user's ID from context
- **Query Scoping**: All queries automatically filter by authenticated user's ID
- **No Cross-User Access**: Users cannot access or modify other users' data

### Firestore Security
- Application-level security enforced by Go middleware
- Consider adding Firestore Security Rules as defense-in-depth:
  ```
  rules_version = '2';
  service cloud.firestore {
    match /databases/{database}/documents {
      match /games/{gameId} {
        allow read, write: if request.auth != null && 
                           resource.data.user_id == request.auth.uid;
      }
    }
  }
  ```

## HTTPS/TLS

### Production Deployment
- **Reverse Proxy Required**: Deploy behind nginx/Caddy/Traefik for TLS termination
- **Environment Variables**: Never commit secrets or credentials to version control
- **Docker Secrets**: Use Docker secrets or environment variables for sensitive data

### Example Nginx Configuration
```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Secrets Management

### Required Secrets
1. **Firebase Service Account**: JSON key file with Firestore and Auth permissions
2. **IGDB API Credentials**: Client ID and Secret from IGDB
3. **Firebase Client Config**: API keys for frontend (can be public per Firebase docs)

### Best Practices
- Use environment variables or secrets management services (AWS Secrets Manager, HashiCorp Vault, etc.)
- Rotate credentials periodically
- Use separate Firebase projects for dev/staging/production
- Never commit `.env` files to version control

## Input Validation

### API Layer
- JSON schema validation via Go struct tags
- User input sanitized by Firebase Admin SDK
- Game data validated before database writes

### Frontend
- Firebase SDK handles authentication flow securely
- API tokens managed by Firebase Auth
- CORS policy enforced by browser

## Rate Limiting

### Recommendations
- Implement rate limiting at reverse proxy level (nginx limit_req module)
- Configure Firebase App Check for additional protection
- Monitor for unusual access patterns

### Example Rate Limit (nginx)
```nginx
limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;

location /api/ {
    limit_req zone=api burst=20;
}
```

## Background Worker Security

### IGDB Sync
- Runs with internal context (no external user input)
- Uses application credentials (not user credentials)
- Logs all operations for audit trail
- Handles API errors gracefully without exposing sensitive data

## Dependencies

### Go Modules
- All dependencies managed via `go.mod`
- Regular updates recommended: `go get -u && go mod tidy`
- Security scanning: `govulncheck ./...` (install: `go install golang.org/x/vuln/cmd/govulncheck@latest`)

### Node.js Dependencies
- Managed via `package.json` and `package-lock.json`
- Regular updates: `npm audit fix`
- Automated scanning via GitHub Dependabot (enable in repository settings)

## Security Scanning

### CodeQL Analysis
- No security vulnerabilities found in current codebase
- Continuous scanning recommended via GitHub Actions

### Recommendations
1. Enable GitHub Dependabot alerts
2. Enable GitHub Security Advisories
3. Run periodic security audits
4. Monitor Firebase security rules and logs

## Incident Response

### In Case of Security Issue
1. Rotate affected credentials immediately
2. Review Firebase Auth logs for unauthorized access
3. Check Firestore audit logs
4. Update affected dependencies
5. Deploy patched version
6. Notify affected users if data breach occurred

## Compliance

### Data Privacy
- User data stored in Firestore (Google Cloud Platform)
- Subject to GCP's privacy policies and certifications
- GDPR/CCPA: Implement data export/deletion features as needed
- User authentication via Google (review Google's privacy policy)

### Data Retention
- Implement data retention policies per your requirements
- Add user account deletion functionality
- Export user data on request

## Contact

For security concerns or to report vulnerabilities, please open an issue in the repository or contact the maintainers directly.
