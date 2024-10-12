package rest

import (
	"github.com/RuanScherer/journey-track-api/adapters/restadptr"
)

func StartAPI() {
	restadptr.StartServer()
}
