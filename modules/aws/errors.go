package aws

import (
	"fmt"
)

// NotFoundError is returned when an expected object is not found
type NotFoundError struct {
	objectType string
	objectID   string
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("Object of type %s with id %s not found", err.objectType, err.objectID)
}

func NewNotFoundError(objectType string, objectID string) NotFoundError {
	return NotFoundError{objectType, objectID}
}

// AsgCapacityNotMetError is returned when the ASG capacity is not yet at the desired capacity.
type AsgCapacityNotMetError struct {
	asgName         string
	desiredCapacity int64
	currentCapacity int64
}

func (err AsgCapacityNotMetError) Error() string {
	return fmt.Sprintf(
		"ASG %s not yet at desired capacity %d (current %d)",
		err.asgName,
		err.desiredCapacity,
		err.currentCapacity,
	)
}

func NewAsgCapacityNotMetError(asgName string, desiredCapacity int64, currentCapacity int64) AsgCapacityNotMetError {
	return AsgCapacityNotMetError{asgName, desiredCapacity, currentCapacity}
}
