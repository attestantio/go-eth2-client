package v1

import "fmt"

const MetadataKeyFinalized = "finalized"
const MetadataKeyExecutionOptimistic = "execution_optimistic"

// Response is a response from the beacon API which may contain metadata.
type Response[T any] struct {
	Data     T
	Metadata map[string]interface{}
}

func (r Response[T]) Finalized() (bool, error) {
	val, ok := r.Metadata[MetadataKeyFinalized]
	if !ok {
		return false, fmt.Errorf("metadata key %s not found", MetadataKeyFinalized)
	}
	return val.(bool), nil
}

func (r Response[T]) ExecutionOptimistic() (bool, error) {
	val, ok := r.Metadata[MetadataKeyExecutionOptimistic]
	if !ok {
		return false, fmt.Errorf("metadata key %s not found", MetadataKeyExecutionOptimistic)
	}
	return val.(bool), nil
}
