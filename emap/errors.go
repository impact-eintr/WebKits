package emap

import "fmt"

type IllegalParameterError struct {
	msg string
}

func newIllegalParameterError(errMsg string) IllegalParameterError {
	return IllegalParameterError{
		msg: fmt.Sprintf("[EMAP]: illgal parameter: %s", errMsg),
	}
}

func (i IllegalParameterError) Error() string {
	return i.msg
}

type PairRedistributorError struct {
	msg string
}

func newPairRedistributorError(errMsg string) PairRedistributorError {
	return PairRedistributorError{
		msg: fmt.Sprintf("[EMAP]: failing pair redistribution: %s", errMsg),
	}
}

func (p PairRedistributorError) Error() string {
	return p.msg
}

type IllegalPairTypeError struct {
	msg string
}

func newIllegalPairTypeError(pair Pair) IllegalPairTypeError {
	return IllegalPairTypeError{
		msg: fmt.Sprintf("[EMAP]: illegal pair type: %T", pair),
	}
}

func (i IllegalPairTypeError) Error() string {
	return i.msg
}
