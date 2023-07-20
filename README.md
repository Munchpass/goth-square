# Square `goth` Provider

An OAuth2 Code Flow Square `goth` provider. Inspired by the [Fitbit Goth Provider](https://github.com/markbates/goth/tree/master/providers/fitbit).

See https://github.com/markbates/goth for more information.

## Getting Started

```bash
go get github.com/munchpass/gothsquare
```

To use the provider:

```go
// Initialize the provider
// (replace the apiCtx values with your own values)
goth.UseProviders(square.New(apiCtx.SquareClientId, apiCtx.SquareSecret, apiCtx.SquareRedirectUrl))

// Create your HTTP Handlers.
r.GET("/square", square.OAuthStart)
r.GET("/square/callback", square.OAuthCallback)

// You're done!
// Just go to /square to start the OAuth flow.
```
