package employee

import (
	"log"
	"net/http"
	"strconv"

	"githb.com/demo-employee-api/internal/config"
	"githb.com/demo-employee-api/utils"
	"github.com/gorilla/mux"
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
	Status    Status     `json:"status"`
	Employees []Employee `json:"employees"`
}

type ListEmpByIdRes struct {
	Status   Status   `json:"status"`
	Employee Employee `json:"employee"`
}

func RegisterHandlers(conf *config.Config, router *mux.Router, svc Service) {
	res := resource{conf, svc}
	router.HandleFunc("/employees", res.ListEmployee).Methods("GET")
	router.HandleFunc("/employees/{id}", res.ListEmployeeById).Methods("GET")
	router.HandleFunc("/employee", res.CreateEmployee).Methods("POST")
	router.HandleFunc("/employee/{id}", res.UpdateEmployee).Methods("PUT")
	router.HandleFunc("/employee/{id}", res.DeleteEmployee).Methods("DELETE")
	router.HandleFunc("/migrations", res.Migrations).Methods("GET")
}

func (res resource) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	resp := Status{}
	emp := Employee{}
	err := utils.SanitizeRequest(r, &emp)
	if err != nil {
		log.Print(err)
		resp.ErrorMessage = err.Error()
		utils.JsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	err = res.service.CreateEmployee(r.Context(), &emp)

	if err != nil {
		log.Println(err)
		resp.ErrorMessage = err.Error()
		utils.JsonResponse(w, http.StatusBadRequest, resp)
		return
	}
	resp.Success = true
	utils.JsonResponse(w, http.StatusOK, resp)
}

func (res resource) ListEmployee(w http.ResponseWriter, r *http.Request) {
	resp := ListEmpRes{}
	emps, err := res.service.ListEmployee(r.Context())
	if err != nil {
		log.Println(err)
		resp.Status.ErrorMessage = err.Error()
		utils.JsonResponse(w, http.StatusInternalServerError, resp)
		return
	}
	resp.Status.Success = true
	resp.Employees = emps
	utils.JsonResponse(w, http.StatusOK, resp)
}

func (res resource) ListEmployeeById(w http.ResponseWriter, r *http.Request) {
	// id := r.URL.Query().Get("id")
	id := mux.Vars(r)["id"]
	resp := ListEmpByIdRes{}
	if id == "" {
		log.Println("id is required ")
		resp.Status.ErrorMessage = "id is required"
		utils.JsonResponse(w, http.StatusBadRequest, resp)
		return
	}
	idv, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("failed to convert id: %s ", err)
		resp.Status.ErrorMessage = err.Error()
		utils.JsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	emp, err := res.service.ListEmployeeById(r.Context(), uint(idv))
	if err != nil {
		log.Println(err)
		resp.Status.ErrorMessage = err.Error()
		utils.JsonResponse(w, http.StatusInternalServerError, resp)
		return
	}
	resp.Status.Success = true
	resp.Employee = emp
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
	var emp Employee
	err := utils.SanitizeRequest(r, &emp)
	if err != nil {
		log.Println(err)
		resp.ErrorMessage = err.Error()
		utils.JsonResponse(w, http.StatusBadRequest, resp)
	}

	idv, err := strconv.Atoi(id)
	if err != nil {
		log.Panicln(err)
		resp.ErrorMessage = err.Error()
		utils.JsonResponse(w, http.StatusInternalServerError, resp)
		return
	}
	emp.UserId = idv
	err = res.service.UpdateEmployee(r.Context(), &emp)
	if err != nil {
		log.Println(err)
		resp.ErrorMessage = err.Error()
		utils.JsonResponse(w, http.StatusInternalServerError, resp)
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
		resp.ErrorMessage = err.Error()
		utils.JsonResponse(w, http.StatusInternalServerError, resp)
		return
	}

	resp.Success = true
	utils.JsonResponse(w, http.StatusOK, resp)
}

func (res resource) Migrations(w http.ResponseWriter, r *http.Request) {
	resp := Status{}
	err := res.service.Migrations(r.Context())
	if err != nil {
		resp.ErrorMessage = err.Error()
		utils.JsonResponse(w, http.StatusInternalServerError, resp)
		return
	}
	resp.Success = true
	utils.JsonResponse(w, http.StatusCreated, resp)
}
