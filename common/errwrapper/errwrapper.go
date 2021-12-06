package errwrapper

import (
	"fmt"
	"sync"

	prisma_types "github.com/aporeto-se/prisma-sdk-go-v2/types"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

// ErrWrapper error wrapper utility
type ErrWrapper struct {
	errors []error
	sync   bool
	sync.Mutex
}

// NewErrWrapper returns new error wrapper
func NewErrWrapper() *ErrWrapper {
	return &ErrWrapper{}
}

// NewSyncErrWrapper returns new thread safe error wrapper
func NewSyncErrWrapper() *ErrWrapper {
	return &ErrWrapper{
		sync: true,
	}
}

// Add adds error and returns self
func (t *ErrWrapper) Add(err error) *ErrWrapper {

	if t.sync {
		t.Lock()
		defer t.Unlock()
	}

	if err != nil {
		t.errors = append(t.errors, err)
	}

	return t
}

// ErrorOrNil returns error or nil
func (t *ErrWrapper) ErrorOrNil() error {

	if t.sync {
		t.Lock()
		defer t.Unlock()
	}

	switch len(t.errors) {

	case 0:
		return nil

	case 1:
		return addConvenienceErr(t.errors[0])

	}

	var errors *multierror.Error

	for _, err := range t.errors {
		errors = multierror.Append(errors, addConvenienceErr(err))
	}

	return errors.ErrorOrNil()
}

func addConvenienceErr(err error) error {

	if err == nil {
		return nil
	}

	if e, ok := err.(*k8serrors.StatusError); ok {
		if k8serrors.IsUnauthorized(e) || k8serrors.IsForbidden(e) {
			zap.L().Debug("Adding Kubernetes convenience error")
			var errors *multierror.Error
			errors = multierror.Append(errors, err)
			errors = multierror.Append(errors, fmt.Errorf("A Kubernetes authorization policy is required"))
			return errors.ErrorOrNil()
		}
	}

	if e, ok := err.(*prisma_types.APIError); ok {
		if e.IsForbiddenOrUnauthorized() {
			zap.L().Debug("Adding Prisma convenience error")
			var errors *multierror.Error
			errors = multierror.Append(errors, err)
			errors = multierror.Append(errors, fmt.Errorf("A Prisma authorization policy is required"))
			return errors.ErrorOrNil()
		}
	}

	return err
}
