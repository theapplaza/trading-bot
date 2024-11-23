package capital

import "errors"

var (
    ErrRecoverable   = errors.New("recoverable error")
    ErrNonRecoverable = errors.New("non-recoverable error")
)