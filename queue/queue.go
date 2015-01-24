// Package queue defines how a message queue and a message should be implemented
package queue

import "fmt"

type (
	// Message defines how a message should look like regardless of implementation
	Message interface {
		fmt.Stringer

		Original() interface{}
	}

	// MultiFetch implements an interface for message queues that support fetching multiple messages at once
	MultiFetch interface {
		MultipleFetch(int) ([]*Message, error)
	}

	// Queue defines how a message queue should look like regardless of implementation
	Queue interface {
		Add(*Message) error
		Fetch(int) ([]*Message, error)
		Confirm(*Message) error
		Delete(*Message) error
	}
)
