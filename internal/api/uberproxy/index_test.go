package uberproxy

import(
	"net/http"
	"testing"
	"vtb-apihack-2025/internal/client/hack"
	"vtb-apihack-2025/internal/mail"
)

func Test_server_SetHandlers(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		jwtSecret string
		apis      map[string]*hack.ApiClient
		mailer    mail.Mailer
		origin    string
		// Named input parameters for target function.
		mux *http.ServeMux
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer(tt.jwtSecret, tt.apis, tt.mailer, tt.origin)
			s.SetHandlers(tt.mux)
		})
	}
}
