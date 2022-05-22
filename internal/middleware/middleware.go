package middleware

import (
	"log"
	"net/http"
	"strings"

	"githb.com/demo-employee-api/internal/entity"
	"githb.com/demo-employee-api/pkg/customErrors"
	"githb.com/demo-employee-api/pkg/token"
	"githb.com/demo-employee-api/utils"
)

func AuthorizeUser(tokenMaker token.Maker, logger *log.Logger, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if len(authHeader) == 0 {
			logger.Println("authorization not provided")
			errResp, code := customErrors.FindErrorType(customErrors.ErrorAuthFailed)
			utils.JsonResponse(w, code, errResp)
			return
		}

		fields := strings.Fields(authHeader)

		if len(fields) < 2 {
			logger.Println("authorization invalid format")
			errResp, code := customErrors.FindErrorType(customErrors.ErrorAuthFailed)
			utils.JsonResponse(w, code, errResp)
			return

		}
		authType := strings.ToLower(fields[0])

		if authType != "bearer" {
			logger.Println("auth type not supported")
			errResp, code := customErrors.FindErrorType(customErrors.ErrorAuthFailed)
			utils.JsonResponse(w, code, errResp)
			return
		}

		accessToken := fields[1]

		payload, err := tokenMaker.VerifyToken(accessToken)

		if err != nil {
			logger.Println("invalid token", err)
			errResp, code := customErrors.FindErrorType(err.Error())
			utils.JsonResponse(w, code, errResp)
			return
		}

		if payload.Role != entity.RoleAdmin {
			logger.Println(err)
			errResp, code := customErrors.FindErrorType(customErrors.ErrorUnAuthorized)
			utils.JsonResponse(w, code, errResp)
			return
		}

		f(w, r)
	}
}
