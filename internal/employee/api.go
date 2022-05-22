package employee

import (
	"log"
	"net/http"
	"strconv"

	"githb.com/demo-employee-api/internal/config"
	"githb.com/demo-employee-api/internal/entity"
	"githb.com/demo-employee-api/internal/middleware"
	"githb.com/demo-employee-api/pkg/customErrors"
	"githb.com/demo-employee-api/pkg/token"
	"githb.com/demo-employee-api/utils"
	"github.com/gorilla/mux"
)

type resource struct {
	conf       *config.Config
	service    Service
	logger     *log.Logger
	tokenMaker token.Maker
}

type Status struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func RegisterHandlers(conf *config.Config, router *mux.Router, svc Service, logger *log.Logger, tokenMaker token.Maker) {
	res := resource{conf, svc, logger, tokenMaker}
	router.HandleFunc("/login", res.LoginEmployee).Methods("POST")
	router.HandleFunc("/employees", res.ListEmployee).Methods("GET")
	router.HandleFunc("/employees/params", res.ListEmployeeByParams).Methods("GET")
	router.HandleFunc("/employee", middleware.AuthorizeUser(res.tokenMaker, res.logger, res.CreateEmployee)).Methods("POST")
	router.HandleFunc("/employee/{id}", middleware.AuthorizeUser(res.tokenMaker, res.logger, res.UpdateEmployee)).Methods("PUT")
	router.HandleFunc("/employee/{id}", middleware.AuthorizeUser(res.tokenMaker, res.logger, res.DeleteEmployee)).Methods("DELETE")
}

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type LoginRes struct {
	Status      Status `json:"status"`
	AccessToken string `json:"token"`
}

func (res resource) LoginEmployee(w http.ResponseWriter, r *http.Request) {
	req := LoginReq{}
	resp := LoginRes{}
	err := utils.ValidateRequest(r, &req)
	if err != nil {
		res.logger.Println(err)
		errResp, code := customErrors.FindErrorType(customErrors.ErrorValidation)
		utils.JsonResponse(w, code, errResp)
		return
	}

	token, err := res.service.LoginEmployee(r.Context(), &req)
	if err != nil {
		res.logger.Println(err)
		errResp, code := customErrors.FindErrorType(err.Error())
		utils.JsonResponse(w, code, errResp)
		return
	}
	res.logger.Printf("%v Logged In", req.Email)
	resp.AccessToken = token
	resp.Status.Success = true
	utils.JsonResponse(w, http.StatusOK, resp)

}

func (res resource) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	resp := Status{}
	emp := entity.Employee{}
	err := utils.ValidateRequest(r, &emp)
	if err != nil {
		log.Print(err)
		errRes, code := customErrors.FindErrorType(customErrors.ErrorValidation)
		utils.JsonResponse(w, code, errRes)
		return
	}

	err = res.service.CreateEmployee(r.Context(), &emp)

	if err != nil {
		res.logger.Println(err)
		errRes, code := customErrors.FindErrorType(err.Error())
		utils.JsonResponse(w, code, errRes)
		return
	}
	res.logger.Printf("%v user created\n", emp.UserId)
	resp.Success = true
	utils.JsonResponse(w, http.StatusOK, resp)
}

type ListEmpRes struct {
	Status    Status            `json:"status"`
	Employees []entity.Employee `json:"employees"`
}

func (res resource) ListEmployee(w http.ResponseWriter, r *http.Request) {
	resp := ListEmpRes{}

	params, err := utils.ValidateParameters(res.conf, r)
	if err != nil {
		res.logger.Println(err)
		errRes, code := customErrors.FindErrorType(customErrors.ErrorValidation)
		utils.JsonResponse(w, code, errRes)
		return

	}
	// fmt.Println(params)

	emps, err := res.service.ListEmployee(r.Context(), params)
	if err != nil {
		res.logger.Println(err)
		errRes, code := customErrors.FindErrorType(err.Error())
		utils.JsonResponse(w, code, errRes)
		return
	}
	res.logger.Println("List Employee Query Exeuted")
	resp.Status.Success = true
	resp.Employees = emps
	utils.JsonResponse(w, http.StatusOK, resp)
}

func (res resource) ListEmployeeByParams(w http.ResponseWriter, r *http.Request) {

	resp := ListEmpRes{}

	params, err := utils.ValidateParameters(res.conf, r)
	if err != nil {
		res.logger.Println(err)
		errRes, code := customErrors.FindErrorType(customErrors.ErrorValidation)
		utils.JsonResponse(w, code, errRes)
		return
	}

	emps, err := res.service.ListEmployeeByParams(r.Context(), params)
	if err != nil {
		res.logger.Println(err)
		errRes, code := customErrors.FindErrorType(err.Error())
		utils.JsonResponse(w, code, errRes)
		return
	}
	res.logger.Println("List Employee Query with Params Exeuted")
	resp.Status.Success = true
	resp.Employees = emps
	utils.JsonResponse(w, http.StatusOK, resp)

}

type UpdateEmpReq struct {
	UserId    int    `json:"user_id,omitempty"`
	FirstName string `json:"first_name,omitempty" validate:"required"`
	LastName  string `json:"last_name,omitempty" validate:"required"`
	Role      string `json:"role" validate:"required,role"`
}

func (res resource) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	resp := Status{}
	id := mux.Vars(r)["id"]
	idv, err := strconv.Atoi(id)
	if err != nil {
		res.logger.Println(err)
		errResp, code := customErrors.FindErrorType(customErrors.ErrorInvalidRequest)
		utils.JsonResponse(w, code, errResp)
		return
	}
	var emp UpdateEmpReq
	err = utils.ValidateRequest(r, &emp)
	if err != nil {
		res.logger.Println(err)
		errRes, code := customErrors.FindErrorType(customErrors.ErrorValidation)
		utils.JsonResponse(w, code, errRes)
		return
	}

	err = res.service.UpdateEmployee(r.Context(), idv, &emp)
	if err != nil {
		res.logger.Println(err)
		errRes, code := customErrors.FindErrorType(err.Error())
		utils.JsonResponse(w, code, errRes)
		return
	}
	res.logger.Printf("id:%v Employee updated\n", idv)
	resp.Success = true
	utils.JsonResponse(w, http.StatusOK, resp)
}

func (res resource) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	resp := Status{}
	id := mux.Vars(r)["id"]
	err := res.service.DeleteEmployee(r.Context(), id)
	if err != nil {
		res.logger.Println(err)
		errRes, code := customErrors.FindErrorType(err.Error())
		utils.JsonResponse(w, code, errRes)
		return
	}
	res.logger.Printf("id:%v Employee Deleted\n", id)
	resp.Success = true
	utils.JsonResponse(w, http.StatusOK, resp)
}
