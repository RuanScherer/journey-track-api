package main

import (
	"github.com/asaskevich/govalidator"
)

func main() {
	govalidator.SetFieldsRequiredByDefault(true)
}
