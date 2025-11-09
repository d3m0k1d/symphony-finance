package hack

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	auth "vtb-apihack-2025/client-pilot/auth/hack"
)

// Authenticate implements client.Client.
func (c *ApiClient) Authenticate(ctx context.Context, clientId string, clientSecret string) error {
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
	c.AccessToken = aresp.JSON200.AccessToken
	c.authF = func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", "Bearer "+c.AccessToken)
		req.Header.Add("x-requesting-bank", c.bankId)
		return nil
	}
	c.authTimer = time.AfterFunc(
		time.Duration(aresp.JSON200.ExpiresIn)*time.Second,
		func() {
			if err := c.Authenticate(context.Background(), clientId, clientSecret); err != nil {
				log.Println(err)
			}
		})
	return nil
}
