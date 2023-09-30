package types

import "github.com/google/uuid"

type SubscriberInfo struct {
	ID     *uuid.UUID
	UserID int
	Token  string
}
