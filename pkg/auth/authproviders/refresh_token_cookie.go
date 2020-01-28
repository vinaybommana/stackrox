package authproviders

import (
	"net/url"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

const (
	refreshTokenCookieName = "RoxRefreshToken"
)

var (
	schemaEncoder = schema.NewEncoder()
	schemaDecoder = schema.NewDecoder()
)

type refreshTokenCookieData struct {
	ProviderType string `schema:"providerType,required"`
	ProviderID   string `schema:"providerId,required"`
	RefreshToken string `schema:"refreshToken,required"`
}

func (r *refreshTokenCookieData) Encode() (string, error) {
	vals := make(url.Values)
	if err := schemaEncoder.Encode(r, vals); err != nil {
		return "", err
	}
	return vals.Encode(), nil
}

func (r *refreshTokenCookieData) Decode(encoded string) error {
	vals, err := url.ParseQuery(encoded)
	if err != nil {
		return errors.Wrap(err, "parsing encoded cookie data")
	}
	return schemaDecoder.Decode(r, vals)
}
