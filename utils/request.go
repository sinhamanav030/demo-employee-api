package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"githb.com/demo-employee-api/internal/config"
)

// type QueryParameters struct {
// 	IncludeArchieved bool   `json:"include_archieved"`
// 	Id               int    `json:"id"`
// 	Name             string `json:"name"`
// 	Page             int    `json:"page"`
// 	PerPage          int    `json:"perpage"`
// }

func ValidateRequest(r *http.Request, req interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(req); err != nil {
		return err
	}

	if validationResult := ValidateStruct(req); !validationResult.Success {
		log.Println(validationResult.FieldErrors)
		return validationResult.OriginalError
	}

	return nil
}

func ValidateParameters(conf *config.Config, r *http.Request) (map[string]string, error) {
	paramMap := make(map[string]string)
	if id := r.URL.Query().Get("id"); id != "" {
		paramMap["user_id"] = id
	}

	if fName := r.URL.Query().Get("first_name"); fName != "" {
		paramMap["first_name"] = fName
	}

	if lName := r.URL.Query().Get("last_name"); lName != "" {
		paramMap["first_name"] = lName
	}

	if archieved := r.URL.Query().Get("include_archieved"); archieved != "" {
		if archieved = strings.ToLower(archieved); archieved == "true" {
			paramMap["archieved"] = "true"
		} else {
			paramMap["archieved"] = "false"
		}

	} else {
		paramMap["archieved"] = "false"
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if _, err := strconv.Atoi(page); err == nil {
			paramMap["page"] = page
		} else {
			return nil, err
		}
	} else {
		paramMap["page"] = ""
	}

	if perPage := r.URL.Query().Get("per_page"); perPage != "" {
		if _, err := strconv.Atoi(perPage); err == nil {
			paramMap["per_page"] = perPage
		} else {
			return nil, err
		}
	} else {
		paramMap["per_page"] = ""
	}

	if sort_by := r.URL.Query().Get("sort_by"); sort_by != "" {
		paramMap["sort_by"] = sort_by

		if sort_order := r.URL.Query().Get("sort_order"); sort_order != "" {
			if sort_order == "asc" || sort_order == "desc" {
				paramMap["sort_order"] = sort_order
			} else {
				sort_order = "asc"
			}
		} else {
			paramMap["sort_order"] = "asc"
		}
	}

	return paramMap, nil

}
