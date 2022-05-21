package middleware

import (
	"fmt"
	"log"
	"net/http"

	"githb.com/demo-employee-api/internal/config"
	"githb.com/demo-employee-api/internal/entity"
	"githb.com/demo-employee-api/pkg/customErrors"
	"githb.com/demo-employee-api/utils"
)

func AuthorizeUser(conf *config.Config, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth-token")
		// fmt.Println(err)
		if err != nil {
			if err == http.ErrNoCookie {
				log.Println(err)
				errResp, code := customErrors.FindErrorType(customErrors.ErrorUnAuthorized)
				utils.JsonResponse(w, code, errResp)
				return
			}
			log.Println(err)
			errResp, code := customErrors.FindErrorType(err.Error())
			utils.JsonResponse(w, code, errResp)
			return

		}

		tokenStr := cookie.Value
		// fmt.Println(tokenStr)
		claims, err := utils.ExtractToken(tokenStr, conf.Auth.JwtKey)
		if err != nil {
			log.Println(err)
			errResp, code := customErrors.FindErrorType(err.Error())
			utils.JsonResponse(w, code, errResp)
			return
		}
		// fmt.Println(tokenStr, claims)

		if claims.Role != entity.RoleAdmin {
			log.Println(err)
			errResp, code := customErrors.FindErrorType(customErrors.ErrorUnAuthorized)
			utils.JsonResponse(w, code, errResp)
			return
		}
		fmt.Println(tokenStr)

		f(w, r)
	}
}
