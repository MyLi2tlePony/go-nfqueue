//+build integration,linux

package nfqueue

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestLinuxNfqueue(t *testing.T) {
	// Set configuration options for nfqueue
	config := Config{
		NfQueue:      100,
		MaxPacketLen: 0xFFFF,
		MaxQueueLen:  0xFF,
	}
	// Open a socket to the netfilter log subsystem
	nfq, err := Open(&config)
	if err != nil {
		t.Fatalf("failed to open nfqueue socket: %v", err)
	}
	defer nfq.Close()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	fn := func(m Msg) int {
		id := m[AttrPacketID].(uint32)
		// Just print out the id and payload of the nfqueue packet
		fmt.Printf("[%d]\t%v\n", id, m[AttrPayload])
		nfq.SetVerdict(id, NfAccept)
		return 0
	}

	// Register your function to listen on nflog group 100
	err = nfq.Register(ctx, NfQnlCopyPacket, fn)
	if err != nil {
		t.Fatalf("failed to register hook function: %v", err)
	}

	// Block till the context expires
	<-ctx.Done()
}