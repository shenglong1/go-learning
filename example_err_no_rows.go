package golearning

import (
	"database/sql"
	"errors"
	"fmt"
	pkgerrors "github.com/pkg/errors"
)


// MockDaoQuery get primary id.
func MockDaoQueryID() (int, error) {

	// use pkg.errors.WithMessage
	return 0, pkgerrors.WithMessage(sql.ErrNoRows, "record no found")
}

// RunMockDaoQueryID 是MockDaoQueryID的调用方.
func RunMockDaoQueryID() {
	id, err := MockDaoQueryID()
	if errors.Is(err, sql.ErrNoRows) {
		// found no rows
		return
	}
	if err != nil {
		// error
		return
	}
	// found record, do sth.
	fmt.Printf("got id %d\n", id)
	return
}


// MyError 自定义的Wrapper
type MyError struct {
	err error
	msg string
}

func Wrap(err error, msg string) error {
	return &MyError{err: err, msg: msg}
}

func (e *MyError) Error() string {
	return e.msg + e.err.Error()
}

func (e *MyError) Unwrap() error {
	return e.err
}

// MockDaoQueryPID get primary id.
func MockDaoQueryPID() (int, error) {

	// use customized Wrapper
	return 0, Wrap(sql.ErrNoRows, "record no found")
}

// RunMockDaoQueryPID 是MockDaoQueryPID的调用方.
func RunMockDaoQueryPID() {
	id, err := MockDaoQueryPID()
	if errors.Is(err, sql.ErrNoRows) {
		// found no rows
		return
	}
	if err != nil {
		// error
		return
	}
	// found record, do sth.
	fmt.Printf("got id %d\n", id)
	return
}

