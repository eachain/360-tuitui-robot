package chain

import (
	"fmt"
	"testing"

	"github.com/eachain/360-tuitui-robot/webhook"
)

func TestCallbacksPointer(t *testing.T) {
	c1 := webhook.Callback{
		OnReceiveSingleMessage: func(webhook.SingleMessageEvent) {},
	}
	c2 := webhook.Callback{
		OnReceiveGroupMessage: func(webhook.GroupMessageEvent) {},
	}

	cb := Callbacks(c1, c2)

	p1 := fmt.Sprintf("%p", c1.OnReceiveSingleMessage)
	p := fmt.Sprintf("%p", cb.OnReceiveSingleMessage)
	if p1 != p {
		t.Fatal("OnReceiveSingleMessage pointer not equal")
	}

	p2 := fmt.Sprintf("%p", c2.OnReceiveGroupMessage)
	p = fmt.Sprintf("%p", cb.OnReceiveGroupMessage)
	if p2 != p {
		t.Fatal("OnReceiveGroupMessage pointer not equal")
	}

	if cb.OnCreateGroup != nil {
		t.Fatal("OnCreateGroup is not nil")
	}
}

func TestCallbacksOrder(t *testing.T) {
	order := make(chan int, 2)
	c1 := webhook.Callback{
		OnReceiveSingleMessage: func(webhook.SingleMessageEvent) {
			order <- 1
		},
	}
	c2 := webhook.Callback{
		OnReceiveSingleMessage: func(webhook.SingleMessageEvent) {
			order <- 2
		},
	}

	cb := Callbacks(c1, c2)
	cb.OnReceiveSingleMessage(webhook.SingleMessageEvent{})

	o1 := <-order
	if o1 != 1 {
		t.Fatalf("the first called is: %v", o1)
	}

	o2 := <-order
	if o2 != 2 {
		t.Fatalf("the second called is: %v", o2)
	}
}

func TestFuncs(t *testing.T) {
	order := make(chan int, 2)
	OnReceiveSingleMessage1 := func(webhook.SingleMessageEvent) {
		order <- 1
	}
	OnReceiveSingleMessage2 := func(webhook.SingleMessageEvent) {
		order <- 2
	}

	on := Funcs(OnReceiveSingleMessage1, OnReceiveSingleMessage2)
	on(webhook.SingleMessageEvent{})

	o1 := <-order
	if o1 != 1 {
		t.Fatalf("the first called is: %v", o1)
	}

	o2 := <-order
	if o2 != 2 {
		t.Fatalf("the second called is: %v", o2)
	}
}
