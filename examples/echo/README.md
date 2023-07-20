# Echo Examples

## How does it work?

Start the OAuth flow with the `OAuthStart` handler (i.e. at `/api/auth/square`).

Then, the redirect uri should be configured to redirect to whatever endpoint `OAuthCallback` is mapped to (i.e. `/api/auth/square/callback`).
