package uuid

import uuid "github.com/satori/go.uuid"

//New - new uuid string
func New() string {
	newUUID := uuid.NewV4()
	return newUUID.String()
}
