package model

import (
	"gorm.io/datatypes"
)

// RhcConnectionCreateRequest represents a request coming from the outside to create a Red Hat Connector connection.
type RhcConnectionCreateRequest struct {
	RhcId       string         `json:"rhc_id"`
	Extra       datatypes.JSON `json:"extra"`
	SourceIdRaw interface{}    `json:"source_id"`
	SourceId    int64
}

// RhcConnectionUpdateRequest represents a request coming from the outside to update a Red Hat Connector connection.
type RhcConnectionUpdateRequest struct {
	Extra datatypes.JSON `json:"extra"`
}

type RhcConnectionResponse struct {
	Id    *string        `json:"id"`
	RhcId *string        `json:"rhc_id"`
	Extra datatypes.JSON `json:"extra,omitempty"`
	AvailabilityStatus
	AvailabilityStatusError string   `json:"availability_status_error,omitempty"`
	SourceIds               []string `json:"source_ids,omitempty"`
}
