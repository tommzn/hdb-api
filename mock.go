package main

import (
	"context"
	"sync"
	"time"

	core "github.com/tommzn/hdb-core"
	events "github.com/tommzn/hdb-events-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type datasourceMock struct {
	messages       []proto.Message
	offset         int
	delay          time.Duration
	dataSourceChan chan proto.Message
}

func (mock *datasourceMock) Run(ctx context.Context, waitGroup *sync.WaitGroup) error {

	ticker := time.NewTicker(mock.delay)
	go func() {
		for {
			select {
			case <-ticker.C:
				mock.dataSourceChan <- mock.nextMessage()
			case <-ctx.Done():
				waitGroup.Done()
				return
			}
		}
	}()
	return nil
}

func (mock *datasourceMock) Observe(filter *[]core.DataSource) <-chan proto.Message {
	return mock.dataSourceChan
}

func (mock *datasourceMock) nextMessage() proto.Message {
	message := mock.messages[mock.offset]
	if mock.offset < len(mock.messages)-1 {
		mock.offset++
	} else {
		mock.offset = 0
	}
	indoorClimate, _ := message.(*events.IndoorClimate)
	indoorClimate.Timestamp = timestamppb.New(time.Now())
	return indoorClimate

}
