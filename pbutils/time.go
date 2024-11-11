package pbutils

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AsTime converts x to a time.Time.
func AsTime(x *timestamppb.Timestamp) (t time.Time) {
	if x == nil {
		return t
	}
	return x.AsTime()
}

// AsDuration converts x to a time.Duration,
// returning the closest duration value in the event of overflow.
func AsDuration(x *durationpb.Duration) (d time.Duration) {
	if x == nil {
		return d
	}
	return x.AsDuration()
}
