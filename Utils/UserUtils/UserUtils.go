package UserUtils

import (
	"github.com/Dontpingforgank/AuthenticationService/Models"
	"github.com/Dontpingforgank/AuthenticationService/Utils/PasswordUtils"
)

func VerifyUserData(userModel *Models.UserRegisterModel) (bool, error) {
	verified, passErr := PasswordUtils.VerifyPassword(userModel.Password)
	if passErr != nil {
		return false, passErr
	}

	if verified &&
		len(userModel.Email) > 0 &&
		len(userModel.Name) > 0 &&
		len(userModel.Country) > 0 &&
		len(userModel.City) > 0 {
		return true, nil
	}

	return false, nil
}
