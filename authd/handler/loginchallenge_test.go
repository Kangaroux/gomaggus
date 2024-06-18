package handler

import (
	"errors"
	"io"
	"testing"

	"github.com/kangaroux/gomaggus/authd"
	"github.com/kangaroux/gomaggus/authd/mock"
	"github.com/kangaroux/gomaggus/internal"
	"github.com/kangaroux/gomaggus/model"
	"github.com/mixcode/binarystruct"
	"github.com/stretchr/testify/assert"
)

func TestLoginChallenge(t *testing.T) {
	var client *authd.Client
	var accounts *mock.AccountService

	newHandler := func() *LoginChallenge {
		client = &authd.Client{
			Conn:  &mock.Conn{},
			State: authd.StateAuthChallenge,
		}
		accounts = &mock.AccountService{}
		return &LoginChallenge{
			Client:   client,
			Accounts: accounts,
		}
	}

	t.Run("invalid state", func(t *testing.T) {
		h := newHandler()
		client.State = authd.StateAuthenticated
		err, ok := h.Handle([]byte{}).(*ErrWrongState)

		assert.Error(t, err)
		assert.True(t, ok)
	})

	t.Run("malformed packet", func(t *testing.T) {
		assert.Equal(t, io.EOF, newHandler().Handle([]byte{}))
	})

	t.Run("account service error", func(t *testing.T) {
		h := newHandler()
		expectedErr := errors.New("fake")
		accounts.OnGet = func(_ *model.AccountGetParams) (*model.Account, error) {
			// Something unexpected happened with the DB
			return nil, expectedErr
		}
		packet := loginChallengeRequest{
			UsernameLength: 1,
			Username:       "a",
		}
		request := internal.MustMarshal(packet, binarystruct.LittleEndian)

		assert.Equal(t, expectedErr, h.Handle(request))
	})

	t.Run("unknown username fake response", func(t *testing.T) {
		packet := &loginChallengeRequest{
			UsernameLength: 4,
			Username:       "fake",
		}
		h := newHandler()
		accounts.OnGet = func(_ *model.AccountGetParams) (*model.Account, error) {
			// Account not found
			return nil, nil
		}
		request := internal.MustMarshal(packet, binarystruct.LittleEndian)

		assert.NoError(t, h.Handle(request))
		assert.Equal(t, authd.StateAuthProof, client.State)

		// The client's account/username will remain empty if it's faked
		assert.Nil(t, client.Account)
		assert.Empty(t, client.Username)
	})

	t.Run("success", func(t *testing.T) {
		packet := &loginChallengeRequest{
			UsernameLength: 3,
			Username:       "bob",
		}
		h := newHandler()
		mockAccount := &model.Account{}
		accounts.OnGet = func(params *model.AccountGetParams) (*model.Account, error) {
			assert.Equal(t, packet.Username, params.Username)
			return mockAccount, nil
		}
		request := internal.MustMarshal(packet, binarystruct.LittleEndian)

		assert.NoError(t, h.Handle(request))
		assert.Equal(t, authd.StateAuthProof, client.State)
		assert.Equal(t, mockAccount, client.Account)
		assert.Equal(t, packet.Username, client.Username)
	})
}
