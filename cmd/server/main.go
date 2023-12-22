package server

import (
	"github.com/RuanScherer/journey-track-api/infrastructure/db"
)

func StartServer() {
	db.GetConnection()
}
