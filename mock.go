package main

import (
	"errors"

	"github.com/golang/protobuf/proto"
)

// publisherMock can be used to mock AWS SQS publisher for testing.
type publisherMock struct {
	shouldFail   bool
	messageCount int
}

// Counts calls to send message methods and returns a new message id. If you pass "error" as queue name it will returns with
// an error and doesn't count this call.
func (mock *publisherMock) Send(message proto.Message) error {

	if mock.shouldFail {
		return errors.New("Unable to send message.")
	}

	mock.messageCount += 2
	return nil
}

// newPublisherMock returns a new mock for a AWS SQS publisher.
func newPublisherMock() *publisherMock {
	return &publisherMock{shouldFail: false, messageCount: 0}
}
