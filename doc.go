package mdson

// To unmarshal a MDSon array into a slice, Unmarshal resets the slice length
// to zero and then appends each element to the slice.
// As a special case, to unmarshal an empty MDSon array into a slice,
// Unmarshal replaces the slice with a new empty slice.
