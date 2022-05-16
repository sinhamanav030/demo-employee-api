package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"githb.com/demo-employee-api/internal/config"
	"githb.com/demo-employee-api/internal/entity"
	"githb.com/demo-employee-api/pkg/customErrors"
	"githb.com/demo-employee-api/utils"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type resource struct {
	conf    *config.Config
	service Service
}

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type LoginRes struct {
	Status Status          `json:"status"`
	Emp    entity.Employee `json:"employee"`
}

type Status struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func RegisterHandlers(conf *config.Config, router *mux.Router, svc Service) {
	res := resource{conf, svc}
	router.HandleFunc("/auth/login", res.Login).Methods("POST")
	router.HandleFunc("/auth/logout", res.Logout).Methods("GET")
}

func (res resource) Login(w http.ResponseWriter, r *http.Request) {
	req := LoginReq{}
	resp := LoginRes{}
	err := utils.SanitizeRequest(r, &req)
	if err != nil {
		log.Println(err)
		errResp, code := customErrors.ErrorDisplayMode(err.Error())
		utils.JsonResponse(w, code, errResp)
		return
	}

	emp, err := res.service.LoginEmployee(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Println(err)
		errResp, code := customErrors.ErrorDisplayMode(err.Error())
		utils.JsonResponse(w, code, errResp)
		return
	}

	claims := &entity.Claims{
		Email:  emp.Email,
		UserId: emp.UserId,
		Role:   emp.Role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println(res.conf.Auth.JwtKey)
	tokenstring, err := token.SignedString([]byte(res.conf.Auth.JwtKey))
	if err != nil {
		log.Println(err)
		errResp, code := customErrors.ErrorDisplayMode(err.Error())
		utils.JsonResponse(w, code, errResp)
		return
	}

	cookie := &http.Cookie{
		Name:    "auth-token",
		Value:   tokenstring,
		Expires: time.Now().Add(time.Minute * 30),
		MaxAge:  int(time.Second * 60 * 30),
	}

	http.SetCookie(w, cookie)

	resp.Status.Success = true
	resp.Emp = emp
	utils.JsonResponse(w, http.StatusOK, resp)

}

func (res resource) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth-token")
	if err != nil {
		log.Println(err)
		errResp, code := customErrors.ErrorDisplayMode(customErrors.ErrorUnAuthorized)
		utils.JsonResponse(w, code, errResp)
		return
	}
	tokenStr := cookie.Value
	_, err = utils.ExtractToken(tokenStr, res.conf.Auth.JwtKey)

	if err != nil {
		log.Println(err)
		erResp, code := customErrors.ErrorDisplayMode(err.Error())
		utils.JsonResponse(w, code, erResp)
		return
	}

	cookie = &http.Cookie{
		Name:   "auth-token",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	status := Status{}
	status.Success = true
	utils.JsonResponse(w, http.StatusAccepted, status)
}
