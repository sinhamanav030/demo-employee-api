package employee

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"githb.com/demo-employee-api/internal/config"
	"githb.com/demo-employee-api/internal/entity"
	"githb.com/demo-employee-api/internal/middleware"
	"githb.com/demo-employee-api/pkg/customErrors"
	"githb.com/demo-employee-api/utils"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type resource struct {
	conf    *config.Config
	service Service
}

type Status struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type ListEmpRes struct {
	Status    Status            `json:"status"`
	Employees []entity.Employee `json:"employees"`
}

type UpdateEmpReq struct {
	UserId    int    `json:"user_id,omitempty"`
	FirstName string `json:"first_name,omitempty" validate:"required"`
	LastName  string `json:"last_name,omitempty" validate:"required"`
	Role      string `json:"role" validate:"required,role"`
}

func RegisterHandlers(conf *config.Config, router *mux.Router, svc Service) {
	res := resource{conf, svc}
	router.HandleFunc("/employees", res.ListEmployee).Methods("GET")
	router.HandleFunc("/employees/params", res.ListEmployeeByParams).Methods("GET")
	router.HandleFunc("/employee", middleware.AuthorizeUser(res.conf, res.CreateEmployee)).Methods("POST")
	router.HandleFunc("/employee/{id}", middleware.AuthorizeUser(res.conf, res.UpdateEmployee)).Methods("PUT")
	router.HandleFunc("/employee/{id}", middleware.AuthorizeUser(res.conf, res.DeleteEmployee)).Methods("DELETE")
	// router.HandleFunc("/migrations", middleware.AuthorizeUser(res.conf, res.Migrations)).Methods("GET")
}

func (res resource) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	resp := Status{}
	emp := entity.Employee{}
	err := utils.ValidateRequest(r, &emp)
	if err != nil {
		log.Print(err)
		resp.ErrorMessage = err.Error()
		utils.JsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	passByteSlice, err := bcrypt.GenerateFromPassword([]byte(emp.Password), bcrypt.MinCost)

	if err != nil {
		log.Println(err)
		resp.ErrorMessage = err.Error()
		utils.JsonResponse(w, http.StatusInternalServerError, resp)
		return
	}

	emp.Password = string(passByteSlice)

	err = res.service.CreateEmployee(r.Context(), &emp)

	if err != nil {
		log.Println(err)
		errRes, code := customErrors.FindErrorType(err.Error())
		utils.JsonResponse(w, code, errRes)
		return
	}
	resp.Success = true
	utils.JsonResponse(w, http.StatusOK, resp)
}

func (res resource) ListEmployee(w http.ResponseWriter, r *http.Request) {
	resp := ListEmpRes{}
	params, err := utils.ValidateParameters(res.conf, r)
	if err != nil {
		log.Println(err)
		errRes, code := customErrors.FindErrorType(err.Error())
		utils.JsonResponse(w, code, errRes)
		return

	}
	fmt.Println(params)
	var page, perPage int
	if page, err = strconv.Atoi(params["page"]); err != nil {
		page = res.conf.Pagination.DefaultPage
	}
	if perPage, err = strconv.Atoi(params["per_page"]); err != nil {
		perPage = res.conf.Pagination.PerPage
	}

	emps, err := res.service.ListEmployee(r.Context(), params["archieved"] == "true", page, perPage)
	if err != nil {
		log.Println(err)
		errRes, code := customErrors.FindErrorType(err.Error())
		utils.JsonResponse(w, code, errRes)
		return
	}
	resp.Status.Success = true
	resp.Employees = emps
	utils.JsonResponse(w, http.StatusOK, resp)
}

func (res resource) ListEmployeeByParams(w http.ResponseWriter, r *http.Request) {

	params, err := utils.ValidateParameters(res.conf, r)

	resp := ListEmpRes{}
	var page, perPage int
	if page, err = strconv.Atoi(params["page"]); err != nil {
		page = res.conf.Pagination.DefaultPage
	}
	if perPage, err = strconv.Atoi(params["per_page"]); err != nil {
		perPage = res.conf.Pagination.PerPage
	}

	emps, err := res.service.ListEmployeeByParams(r.Context(), params, page, perPage)
	if err != nil {
		log.Println(err)
		errRes, code := customErrors.FindErrorType(err.Error())
		utils.JsonResponse(w, code, errRes)
		return
	}
	resp.Status.Success = true
	resp.Employees = emps
	utils.JsonResponse(w, http.StatusOK, resp)

}

func (res resource) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	resp := Status{}
	id := mux.Vars(r)["id"]
	if id == "" {
		log.Println("id is required")
		resp.ErrorMessage = "id is required"
		utils.JsonResponse(w, http.StatusBadRequest, resp)
		return
	}
	var emp UpdateEmpReq
	err := utils.ValidateRequest(r, &emp)
	if err != nil {
		log.Println(err)
		resp.ErrorMessage = err.Error()
		utils.JsonResponse(w, http.StatusBadRequest, resp)
		return
	}
	idv, err := strconv.Atoi(id)
	if err != nil {
		log.Panicln(err)
		resp.ErrorMessage = err.Error()
		utils.JsonResponse(w, http.StatusInternalServerError, resp)
		return
	}
	employee := entity.Employee{
		UserId:    idv,
		FirstName: emp.FirstName,
		LastName:  emp.LastName,
		Role:      emp.Role,
	}
	err = res.service.UpdateEmployee(r.Context(), &employee)
	if err != nil {
		log.Println(err)
		errRes, code := customErrors.FindErrorType(err.Error())
		utils.JsonResponse(w, code, errRes)
		return
	}
	resp.Success = true
	utils.JsonResponse(w, http.StatusOK, resp)
}

func (res resource) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	resp := Status{}
	id := mux.Vars(r)["id"]
	if id == "" {
		log.Println("id is required")
		resp.ErrorMessage = "id is required"
		utils.JsonResponse(w, http.StatusBadRequest, resp)
		return
	}
	err := res.service.DeleteEmployee(r.Context(), id)
	if err != nil {
		log.Println(err)
		errRes, code := customErrors.FindErrorType(err.Error())
		utils.JsonResponse(w, code, errRes)
		return
	}

	resp.Success = true
	utils.JsonResponse(w, http.StatusOK, resp)
}

func (res resource) Migrations(w http.ResponseWriter, r *http.Request) {
	resp := Status{}
	err := res.service.Migrations(r.Context())
	if err != nil {
		errRes, code := customErrors.FindErrorType(err.Error())
		utils.JsonResponse(w, code, errRes)
		return
	}
	resp.Success = true
	utils.JsonResponse(w, http.StatusCreated, resp)
}
