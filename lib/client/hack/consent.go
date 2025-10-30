package hack

import (
	"context"
	"fmt"
	auth "vtb-apihack-2025/client-pilot/auth/hack"
	"vtb-apihack-2025/lib/client"
)
func sp(s string) *string {
	return &s
}

// RequestConsent implements client.Client.
func (c Client) RequestConsent(ctx context.Context, clientId string) (client.RequestedConsent, error) {

	b := auth.ConsentRequestBody{
		ClientId:           clientId,
		Permissions:        []string{"ReadAccountsDetail", "ReadBalances", "ReadTransactionsDetail"},
		Reason:             sp("WE NEEd YOUR dATA!!!!!"),
		RequestingBank:     &c.bankId,
		RequestingBankName: &c.bankName,
	}
	resp, err := c.ac.RequestConsentAccountConsentsRequestPost(
		ctx,
		&auth.RequestConsentAccountConsentsRequestPostParams{
			XRequestingBank: &c.bankId,
		},
		b,
		c.authF,
	)
	if err != nil {
		return client.RequestedConsent{}, err
	}
	if resp.StatusCode != 200 {
		return client.RequestedConsent{}, fmt.Errorf("Auth2 error: exit code: %d\nbody:\n%s", resp.StatusCode, resp.Body)
	}
	respp, err := auth.ParseRequestConsentAccountConsentsRequestPostResponse(resp)
	if err != nil {
		return client.RequestedConsent{}, err
	}
	return client.RequestedConsent{
		ExpiresAt: respp.JSON200.ExpiresAt.Time,
		ClientId:  respp.JSON200.ClientId,
	}, nil
}

