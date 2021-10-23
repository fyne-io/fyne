//go:generate go run gen.go

// Package async provides unbounded channel and queue structures that are
// designed for caching unlimited number of a concrete type. For better
// performance, a given type should be less or euqal than 16 bytes.
//
// The difference of an unbounded channel or queue is that unbounde channels
// can utilize select and channel semantics, whereas queue cannot. A user of
// this package should balance this tradeoff. For instance, an unbounded
// channel can provide zero waiting cost when trying to receiving an object
// when the receiving select statement has a default case, and a queue can
// only receive the object with a time amount of time, but depending on the
// number of queue item producer, the receiving time may increase accordingly.
//
// Delicate dance: One must aware that an unbounded channel may lead to
// OOM when the consuming speed of the buffer is lower than the producing
// speed constantly. However, such a channel may be fairly used for event
// delivering if the consumer of the channel consumes the incoming
// forever, such as even processing.
//
// This package involves code generators, see gen.go for more details.
package async
