package core

import (
	"errors"
)

var (
	ErrInterleaveInvalidElementCount   = errors.New("each individual slice must have the same element count")
	ErrDeinterleaveInvalidElementCount = errors.New("number of elements in slice is not evenly divisible by slice count")
)

// InterleaveSlices takes one more slices as input and returns a new
// "interleaved" result. For example, given:
//
//	s0: [ 0, 3, 6, 9 ]
//	s1: [ 1, 4, 7, 10 ]
//	s2: [ 2, 5, 8, 11 ]
//
// InterleaveSlices(s1, s2, s3) will return:
//
//	[ 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11 ]
//
// InterleaveSlices will fail if any individual slice has a different length.
func InterleaveSlices[T any](s1 []T, slices ...[]T) ([]T, error) {

	elementsPerSlice := len(s1)
	totalSlices := 1 + len(slices)
	totalElements := elementsPerSlice * totalSlices

	res := make([]T, totalElements)
	for i := 0; i < elementsPerSlice; i++ {
		res[totalSlices*i] = s1[i]
	}
	for s := range slices {
		slice := slices[s]
		if len(slice) != elementsPerSlice {
			return nil, ErrInterleaveInvalidElementCount
		}
		for i := 0; i < elementsPerSlice; i++ {
			res[(totalSlices*i)+s+1] = slice[i]
		}
	}

	return res, nil
}

// DeinterleaveSlices is the inverse of InterleaveSlices. It divides a single
// slice into 'sliceCount' output slices. For example, given:
//
//	s: [ 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11 ]
//
// DeinterleaveSlices(s, 3) will return:
//
//	s0: [ 0, 3, 6, 9 ]
//	s1: [ 1, 4, 7, 10 ]
//	s2: [ 2, 5, 8, 11 ]
//
// DeinterleaveSlices will fail if 's' isn't evenly divisible by 'sliceCount'.
func DeinterleaveSlices[T any](s []T, sliceCount int) ([][]T, error) {

	remainder := len(s) % sliceCount
	if remainder != 0 {
		return nil, ErrDeinterleaveInvalidElementCount
	}
	elementsPerSlice := len(s) / sliceCount

	result := make([][]T, sliceCount)
	for i := 0; i < sliceCount; i++ {

		// Fill slice 'i' with strided data from 's'
		tmp := make([]T, elementsPerSlice)
		for j := 0; j < len(tmp); j++ {
			tmp[j] = s[sliceCount*j+i]
		}

		result[i] = tmp
	}

	return result, nil
}
