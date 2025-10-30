package hack

import (
	"context"
	"fmt"
	"net/http"
	auth "vtb-apihack-2025/client-pilot/auth/hack"
)

// Authenticate implements client.Client.
func (c Client) Authenticate(ctx context.Context, clientId string, clientSecret string) error {
	resp, err := c.ac.CreateBankTokenAuthBankTokenPost(ctx, &auth.CreateBankTokenAuthBankTokenPostParams{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Auth error: exit code: %d\nbody:\n%s", resp.StatusCode, resp.Body)
	}
	aresp, err := auth.ParseCreateBankTokenAuthBankTokenPostResponse(resp)
	if err != nil {
		return err
	}
	c.authF = func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", "Bearer "+aresp.JSON200.AccessToken)
		return nil
	}
	return nil
}
