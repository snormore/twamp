package twamp_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/snormore/twamp"
	"github.com/stretchr/testify/require"
)

func TestLightSenderToReflector(t *testing.T) {
	_, addr := startListener(t)
	sender := newSender(t, addr)
	res, err := sender.SendProbeWithPadding(64)
	require.NoError(t, err)
	requireRTTPositive(t, res)
}

func TestSenderTimeout(t *testing.T) {
	sender, err := twamp.NewSender(":0", "127.0.0.1:65000", 100*time.Millisecond, log, nil)
	require.NoError(t, err)

	start := time.Now()
	_, err = sender.SendProbe()
	require.Error(t, err)
	require.WithinDuration(t, time.Now(), start.Add(100*time.Millisecond), 50*time.Millisecond)
}

func TestConcurrentProbes(t *testing.T) {
	_, addr := startListener(t)
	sender := newSender(t, addr)

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := sender.SendProbeWithPadding(32)
			require.NoError(t, err)
			requireRTTPositive(t, res)
		}()
	}
	wg.Wait()
}

func TestMultiplePaddingSizes(t *testing.T) {
	_, addr := startListener(t)
	sender := newSender(t, addr)
	for _, pad := range []int{0, 16, 128, 512} {
		t.Run(fmt.Sprintf("pad_%d", pad), func(t *testing.T) {
			res, err := sender.SendProbeWithPadding(pad)
			require.NoError(t, err)
			requireRTTPositive(t, res)
		})
	}
}

func startListener(t *testing.T) (twamp.Listener, string) {
	t.Helper()
	listener, err := twamp.NewListener("127.0.0.1:0", 128, log, nil)
	require.NoError(t, err)
	go func() {
		if err := listener.Run(context.Background()); err != nil {
			if !errors.Is(err, context.Canceled) {
				t.Logf("failed to run listener: %v", err)
				t.Fail()
			}
		}
	}()
	t.Cleanup(func() { _ = listener.Close() })
	return listener, listener.LocalAddr().String()
}

func newSender(t *testing.T, remote string) twamp.Sender {
	t.Helper()
	sender, err := twamp.NewSender("127.0.0.1:0", remote, time.Second, log, nil)
	require.NoError(t, err)
	return sender
}

func requireRTTPositive(t *testing.T, res *twamp.ProbeResult) {
	t.Helper()
	require.True(t, res.RTT > 0)
	require.Less(t, res.RTT, 100*time.Millisecond)
}
