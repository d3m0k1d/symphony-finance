package impl

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"vtb-apihack-2025/internal/mail"
	"vtb-apihack-2025/internal/otp"
	"vtb-apihack-2025/internal/storage/interfaces"

	"golang.org/x/crypto/bcrypt"
)

var _ otp.OTPAuthenticator = otper{}

type otper struct {
	mailer mail.Mailer
	store  interfaces.OtpSessionStore
}

func NewOtper(
	mailer mail.Mailer,
	store interfaces.OtpSessionStore,
) otper {
	return otper{
		mailer: mailer,
		store:  store,
	}
}

// CompleteCodeAuth implements otp.OTPAuthenticator.
func (o otper) CompleteCodeAuth(ctx context.Context, session, code string) (user interfaces.UserIdentity, err error) {
	res, err := o.store.RetrieveCode(ctx, session)
	if err != nil {
		return
	}
	err = bcrypt.CompareHashAndPassword(res.Hash, []byte(code))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			err = otp.ErrOtpMismatch
			return
		}
		return
	}
	if err = o.store.DropCode(ctx, session); err != nil {
		return
	}
	return interfaces.UserIdentity{
		Email: res.Email,
	}, nil
}

// InitCodeAuth implements otp.OTPAuthenticator.
func (o otper) InitCodeAuth(ctx context.Context, email string, session string) error {
	var err error
	i, err := rand.Int(rand.Reader, big.NewInt(999999+1))
	if err != nil {
		return err
	}
	code := fmt.Sprintf("%.6d", i)
	err = o.mailer.SendCode(ctx, email, code)
	if err != nil {
		return err
	}
	codeHash, err := bcrypt.GenerateFromPassword([]byte(code), 0)
	if err != nil {
		return err
	}
	err = o.store.StoreCode(ctx, session, email, (codeHash))
	if err != nil {
		return err
	}

	return nil
}
