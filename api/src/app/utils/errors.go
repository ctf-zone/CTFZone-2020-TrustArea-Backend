package utils

import (
	"errors"
)

var ValidationError 					= errors.New("Validation errors!")
var UndefinedError 						= errors.New("Undefined error!")
var PageNotFoundError 					= errors.New("Page not found!")
var UserAlreadyExistsError 				= errors.New("User with this username already exists!")
var RefreshTokenNotExistsError 			= errors.New("Refresh token doesn't exist!")
var SessionTokenNotExistsError 			= errors.New("Session token doesn't exist!")
var SolutionIsNotCorrectError 			= errors.New("Solution is not correct!")
var TaskNotFoundError 					= errors.New("Task not found!")
var UserNotFoundError 					= errors.New("User not found!")