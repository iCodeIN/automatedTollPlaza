package pkg

import "automatedTollPlaze/pkg/toll"

// ServiceHandler is a handler with multiple service handlers
// in this case toll handler service
type ServiceHandler struct {
	TollHandler toll.Service
}
