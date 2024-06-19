package session

import (
	"context"
	"log"
	"time"

	"github.com/kangaroux/gomaggus/realmd"
)

type logoutResult uint32

const (
	logoutSuccess         logoutResult = 0
	logoutFailInCombat    logoutResult = 1
	logoutFailFrozenByGM  logoutResult = 2
	logoutFailNotOnGround logoutResult = 3
)

type logoutResponse struct {
	Result  logoutResult
	Instant bool
}

const (
	// The 20 second timer is set by the client
	logoutDelay = time.Second * 20
)

func LogoutHandler(client *realmd.Client) error {
	// TODO: lookup player in world, check rested state
	resp := logoutResponse{Result: logoutSuccess, Instant: true}
	if err := client.SendPacket(realmd.OpServerLogout, &resp); err != nil {
		return err
	}

	log.Println("sent logout response")

	if resp.Instant {
		return completeLogout(client)
	} else {
		go logoutAfterDelay(context.TODO(), client)
	}

	return nil
}

func LogoutCancelHandler(client *realmd.Client) error {
	client.CancelPendingLogout()

	// Always send an ACK even if there was no logout pending.
	if err := client.SendPacket(realmd.OpServerLogoutCancelACK, nil); err != nil {
		return err
	}

	log.Println("sent logout cancel ack response")
	return nil
}

// completeLogout notifies the client they should logout (or exit game) immediately.
func completeLogout(client *realmd.Client) error {
	// TODO: remove player from world
	if err := client.SendPacket(realmd.OpServerLogoutComplete, nil); err != nil {
		return err
	}

	log.Println("sent logout complete response")
	return nil
}

// logoutAfterDelay notifies the client to logout after a delay. During that delay, the logout is
// considered pending and can be cancelled by calling client.CancelPendingLogout.
func logoutAfterDelay(ctx context.Context, client *realmd.Client) {
	ctx, cancel := context.WithCancel(ctx)

	client.CancelPendingLogout = cancel
	client.LogoutPending = true

	defer func() {
		// The client is either logged out or the logout was cancelled
		client.LogoutPending = false

		// Cleanup context
		cancel()
	}()

	select {
	case <-ctx.Done():
		return
	case <-time.After(logoutDelay):
		// TODO: handle error
		// FIXME: concurrent writing is not safe
		completeLogout(client)
	}
}
