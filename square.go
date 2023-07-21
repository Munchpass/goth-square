// Package square implements the OAuth protocol for authenticating users through square.
// This package can be used as a reference implementation of an OAuth provider for Goth.
package gothsquare

import (
	"context"
	"fmt"
	"net/http"

	"github.com/markbates/goth"
	"golang.org/x/oauth2"
)

const (
	authURL  string = "https://connect.squareupsandbox.com/oauth2/authorize"
	tokenURL string = "https://connect.squareupsandbox.com/oauth2/token"
	// endpointProfile string = "https://connect.squareup.com/v2/merchants" // '-' for logged in user
)

// Relevant scopes. This is not a comprehensive list, but is fairly complete for Munch Insight's usecases.
// See https://developer.squareup.com/docs/oauth-api/square-permissions for a complete list.
const (
	// Catalog scopes
	ScopeItemsRead  = "ITEMS_READ"
	ScopeItemsWrite = "ITEMS_WRITE"

	// Order, Payments & Checkout scopes
	ScopeOrdersWrite   = "ORDERS_WRITE"
	ScopeOrdersRead    = "ORDERS_READ"
	ScopePaymentsRead  = "PAYMENTS_READ"
	ScopePaymentsWrite = "PAYMENTS_WRITE"

	// Customer scopes
	ScopeCustomerRead   = "CUSTOMERS_READ"
	ScopeCustomersWrite = "CUSTOMERS_WRITE"

	// Employee scopes
	ScopeEmployeeRead = "EMPLOYEES_READ"

	// Gift Card scopes
	ScopeGiftCardsRead  = "GIFTCARDS_READ"
	ScopeGiftCardsWrite = "GIFTCARDS_WRITE"

	// Inventory scopes
	ScopeInventoryRead  = "INVENTORY_READ"
	ScopeInventoryWrite = "INVENTORY_WRITE"

	// Invoice scopes
	ScopeInvoiceRead  = "INVOICES_READ"
	ScopeInvoiceWrite = "INVOICES_WRITE"

	// Labor (Timecard) scopes
	ScopeTimecardsRead  = "TIMECARDS_READ"
	ScopeTimecardsWrite = "TIMECARDS_WRITE"

	// Merchant & Locations scopes
	ScopeMerchantRead  = "MERCHANT_PROFILE_READ"
	ScopeMerchantWrite = "MERCHANT_PROFILE_WRITE"

	// Loyalty scopes
	ScopeLoyaltyRead  = "LOYALTY_READ"
	ScopeLoyaltyWrite = "LOYALTY_WRITE"

	// Payout scopes
	ScopePayoutsRead = "PAYOUTS_READ"

	// Online sites scopes
	ScopeOnlineStoreSitesRead = "ONLINE_STORE_SITE_READ"

	// Subscriptions
	ScopeSubscriptionsRead  = "SUBSCRIPTIONS_READ"
	ScopeSubscriptionsWrite = "SUBSCRIPTIONS_WRITE"

	// Vendors
	ScopeVendorRead  = "VENDOR_READ"
	ScopeVendorWrite = "VENDOR_WRITE"
)

// New creates a new square provider, and sets up important connection details.
// You should always call `square.New` to get a new Provider. Never try to create
// one manually.
func New(clientKey, secret, callbackURL string, scopes ...string) *Provider {
	p := &Provider{
		ClientKey:    clientKey,
		Secret:       secret,
		CallbackURL:  callbackURL,
		providerName: "square",
	}
	p.config = newConfig(p, scopes)
	return p
}

// Provider is the implementation of `goth.Provider` for accessing square.
type Provider struct {
	ClientKey    string
	Secret       string
	CallbackURL  string
	HTTPClient   *http.Client
	config       *oauth2.Config
	providerName string
}

// Name is the name used to retrieve this provider later.
func (p *Provider) Name() string {
	return p.providerName
}

// SetName is to update the name of the provider (needed in case of multiple providers of 1 type)
func (p *Provider) SetName(name string) {
	p.providerName = name
}

func (p *Provider) Client() *http.Client {
	return goth.HTTPClientWithFallBack(p.HTTPClient)
}

// Debug is a no-op for the square package.
func (p *Provider) Debug(debug bool) {}

// BeginAuth asks square for an authentication end-point.
func (p *Provider) BeginAuth(state string) (goth.Session, error) {
	url := p.config.AuthCodeURL(state)
	session := &Session{
		AuthURL: url,
	}
	return session, nil
}

// FetchUser will go to square and access basic information about the user.
func (p *Provider) FetchUser(session goth.Session) (goth.User, error) {
	s := session.(*Session)
	user := goth.User{
		AccessToken:  s.AccessToken,
		Provider:     p.Name(),
		RefreshToken: s.RefreshToken,
		ExpiresAt:    s.ExpiresAt,
		UserID:       s.UserID,
	}

	if user.AccessToken == "" {
		// data is not yet retrieved since accessToken is still empty
		return user, fmt.Errorf("%s cannot get user information without accessToken", p.providerName)
	}

	return user, nil
}

func newConfig(provider *Provider, scopes []string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     provider.ClientKey,
		ClientSecret: provider.Secret,
		RedirectURL:  provider.CallbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: []string{
			ScopeMerchantRead,
		},
	}

	defaultScopes := map[string]struct{}{
		ScopeMerchantRead: {},
	}

	for _, scope := range scopes {
		if _, exists := defaultScopes[scope]; !exists {
			c.Scopes = append(c.Scopes, scope)
		}
	}

	return c
}

// RefreshToken get new access token based on the refresh token
func (p *Provider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	token := &oauth2.Token{RefreshToken: refreshToken}
	ts := p.config.TokenSource(context.Background(), token)
	newToken, err := ts.Token()
	if err != nil {
		return nil, err
	}
	return newToken, err
}

// RefreshTokenAvailable refresh token is not provided by square
func (p *Provider) RefreshTokenAvailable() bool {
	return true
}
