package use_case

import "errors"

var ErrInternalServerError = errors.New("internal server error")
var ErrNotFound = errors.New("not found")
var ErrFileSizeTooLarge = errors.New("file size too large")
var ErrOtpExpired = errors.New("otp expired")                    // should return 408 Expired
var ErrOtpInvalid = errors.New("otp invalid")                    // should return 401 Unauthorized
var ErrOtpRetryExceeded = errors.New("otp retry limit exceeded") // should return 429 Too many requests
var ErrUserLimitReached = errors.New("license user limit reached")
var ErrForbidden = errors.New("forbidden")
var ErrUnauthorized = errors.New("unauthorized")
var ErrExpired = errors.New("expired")
var ErrConflict = errors.New("conflict")
var ErrBadRequest = errors.New("bad request")
