package hack

import (
	"context"
	"fmt"
	"net/http"
	"time"
	auth "vtb-apihack-2025/client-pilot/auth/hack"
	"vtb-apihack-2025/client-pilot/pe"
	"vtb-apihack-2025/internal/client"
	"vtb-apihack-2025/internal/storage/interfaces"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type Client struct {
	clientId string
	api      *ApiClient
}

func (c Client) Accounts(ctx context.Context) ([]client.Account, error) {
	resp, err := c.api.pec.GetAccounts(ctx, &pe.GetAccountsParams{
		Page:                   new(int32),
		XFapiAuthDate:          new(string),
		XFapiCustomerIpAddress: new(string),
		XFapiInteractionId:     uuid.UUID{},
		XCustomerUserAgent:     new(string),
	}, c.api.authF, func(ctx context.Context, req *http.Request) error {
		q := req.URL.Query()
		q.Add("client_id", c.clientId)
		req.URL.RawQuery = q.Encode()
		cons, err := c.api.CS.FirstValidFor(ctx, c.clientId, pe.ReadAccountsDetail)
		if err != nil {
			return err
		}
		req.Header.Add("x-consent-id", cons.ID)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmtHttpError("Error fetching user accounts: %w", resp)
	}
	respp, err := pe.ParseGetAccountsResponse(resp)
	if err != nil {
		return nil, err
	}
	accs := respp.JSON200.Data.Account
	if accs == nil {
		return nil, fmt.Errorf("probably unreachable: accounts list is null")
	}
	return lo.Map(*accs, func(item pe.Account, _ int) client.Account {
		return client.Account{
			ID: item.AccountId,
		}
	}), nil
}

// RequestConsents implements client.Client.
func (c Client) RequestConsents(ctx context.Context) error {
	perms := []pe.PermissionsType{
		pe.ReadAccountsBasic,
		pe.ReadAccountsDetail,
		pe.ReadBalances,
		pe.ReadParty,
		pe.ReadPaymentCards,
		pe.ReadProducts,
		pe.ReadStandingOrdersBasic,
		pe.ReadStandingOrdersDetail,
		pe.ReadStatementsBasic,
		pe.ReadStatementsDetail,
		pe.ReadTransactionsBasic,
		pe.ReadTransactionsCredits,
		pe.ReadTransactionsDebits,
		pe.ReadTransactionsDetail,
	}
	b := auth.ConsentRequestBody{
		ClientId:           c.clientId,
		Permissions:        lo.Map(perms, func(item pe.PermissionsType, _ int) string { return string(item) }),
		Reason:             &c.api.consentReason,
		RequestingBank:     &c.api.BankId,
		RequestingBankName: &c.api.bankName,
	}
	resp, err := c.api.ac.RequestConsentAccountConsentsRequestPost(
		ctx,
		&auth.RequestConsentAccountConsentsRequestPostParams{
			XRequestingBank: &c.api.BankId,
		},
		b,
		c.api.authF,
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmtHttpError("Auth2 error: %w", resp)
	}
	respp, err := auth.ParseRequestConsentAccountConsentsRequestPostResponse(resp)
	if err != nil {
		return err
	}
	return c.api.CS.InsertConsent(ctx, c.clientId, perms, interfaces.Consent{
		ID: respp.JSON200.ConsentId,
	})
}

// TransactionsPage implements client.Client.
func (c Client) TransactionsPage(ctx context.Context, accId string, page int32, fromBookingDateTime *time.Time, toBookingDateTime *time.Time) ([]pe.TransactionHistory, error) {
	resp, err := c.api.pec.GetAccountsaccountIdTransactions(ctx, accId, &pe.GetAccountsaccountIdTransactionsParams{
		Page:                &page,
		FromBookingDateTime: fromBookingDateTime,
		ToBookingDateTime:   toBookingDateTime,
		// XFapiAuthDate:          new(string),
		// XFapiCustomerIpAddress: new(string),
		XFapiInteractionId: uuid.New(),
		// XCustomerUserAgent:     new(string),
	}, c.api.authF)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmtHttpError("error fetching transactions history: %w", resp)
	}
	parsed, err := pe.ParseGetAccountsaccountIdTransactionsResponse(resp)
	if err != nil {
		return nil, err
	}
	return *parsed.JSON200.Data.Transaction, nil
}

// Client implements client.APIClient.
func (c *ApiClient) Client(clientId string) client.Client {
	return Client{
		clientId: clientId,
		api:      c,
	}
}
