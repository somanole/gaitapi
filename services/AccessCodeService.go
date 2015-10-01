// AccessCodeService
package services

import (
	"strings"
)

type AccessCode struct{
	AccessCode string
}

type ValidateAccessCodeResponse struct{
	AccessCode string
	IsValid bool
}

func ValidateAccessCode(accessCode AccessCode) ValidateAccessCodeResponse {
	var response ValidateAccessCodeResponse
	
	if strings.ToLower(accessCode.AccessCode) == "websummit2015" {
		response = ValidateAccessCodeResponse {accessCode.AccessCode, true}
	} else {
		response = ValidateAccessCodeResponse {accessCode.AccessCode, false}
	}
	
	return response
}

func GetAccessCode() AccessCode {
	response := AccessCode{"websummit2015"}
	
	return response
}
