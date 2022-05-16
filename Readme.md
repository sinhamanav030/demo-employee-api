##### DEMO EMPLOYEE API

## TECH STACK 
    1. DB - POSTGRES
    2. GORILLA MUX

## FOLDER STRUCTURE:-

```
├── cmd
│   └── server
│       └── main.go
├── config.yaml
├── go.mod
├── go.sum
├── internal
│   ├── auth
│   │   ├── api.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── config
│   │   └── config.go
│   ├── employee
│   │   ├── api.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── entity
│   │   ├── claims.go
│   │   ├── employee.go
│   │   └── params.go
│   ├── healthcheck
│   │   └── healthcheck.go
│   └── middleware
│       └── middleware.go
├── local.yaml
├── pkg
│   ├── customErrors
│   │   └── customErrors.go
│   └── db
│       └── db.go
└── utils
    ├── request.go
    ├── response.go
    ├── validate.go
    └── verifyToken.go
```

## COMMON STRUCTS

    # STATUS
        type Status struct {
            Success      bool   `json:"success"`
            ErrorMessage string `json:"errorMessage,omitempty"`
        }

## ENDPOINTS
```
    1. (/employees) , METHOD -> GET AUTHORIZATION -> any
        GET ALL EMPLOYEES 
        REQUEST:
            QUERY PARAMETERS : 
                1. include_archieved : {true/false} DEFAULT: true OPTIONAL
                2. page : {integer} DEFAULT : 1 OPTIONAL
                3. per_page : {integer} DEFAULT : 5 OPTIONAL
        RESPONSE:
            STRUCT USED
                type Employee struct {
                    UserId       int    
                    FirstName    string 
                    LastName     string 
                    Email        string 
                    Password     string 
                    Role         string 
                    CreatedAt    string 
                    LastAccessAt string 
                    UpdatedAt    string 
                    Archieved    bool   
                }
            JSON RESPONSE : 
                {
                    status Status
                    employees []Employee
                }
    
    2. (/employees/params) , METHOD ->GET ,AUTHORIZATION -> any
        GET EMPLOYEES BASED ON QUERY PARAMS ,USED FOR (SEARCHING,SORTING)
        REQUEST:
            QUERY PARAMETERS : 
                1. include_archieved : {true/false} DEFAULT: true OPTIONAL
                2. page : {integer} DEFAULT : 1 OPTIONAL
                3. per_page : {integer} DEFAULT : 5 OPTIONAL
                4. first_name : full (or partial)matching string OPTIONAL
                5. last_name : full (or partial)matching string OPTIONAL
                6. sort_by : valid column name OPTIONAL
                7. sort_order : {asc/desc} DEFAULT(IN CASE sort_order is set) : asc
        RESPONSE:
            STRUCT USED
                type Employee struct {
                    UserId       int    
                    FirstName    string 
                    LastName     string 
                    Email        string 
                    Password     string 
                    Role         string 
                    CreatedAt    string 
                    LastAccessAt string 
                    UpdatedAt    string 
                    Archieved    bool   
                }
            JSON RESPONSE : 
                type ListEmpRes struct {
                    Status    Status            `json:"status"`
                    Employees []entity.Employee `json:"employees"`
                }
    
    3. (/employee), METHOD-> POST ,AUTHORIZATION -> admin
        TO CREATE NEW EMPLOYEE RECORD
        REQUEST:
            BODY:
                {
                    "first_name" :"...",
                    "last_name" :"...",
                    "email": "...",
                    "password":"...",
                    "role":"..."
                }
                ALL REQUIRED FIELDS 
        
        RESPONSE:
            STATUS STRUCT ACCORDING TO SUCCESS OR FAILURE

    4. (/employee/{id}), METHOD-> PUT ,AUTHORIZATION -> admin
        TO UPDATE EMPLOYEE DETAILS
        REQUEST:
            STRUCT:
                type UpdateEmpReq struct {
                    UserId    int    `json:"user_id,omitempty"`
                    FirstName string `json:"first_name,omitempty" validate:"required"`
                    LastName  string `json:"last_name,omitempty" validate:"required"`
                    Role      string `json:"role" validate:"required,role"`
                }

            BODY:
                {
                    "first_name" :"...",
                    "last_name" :"...",
                    "email": "...",
                    "role":"..."
                }
        
        RESPONSE:
            STATUS STRUCT ACCORDING TO SUCCESS OR FAILURE


    5. (/employee/{id}), METHOD-> DELETE, AUTHORIZATION -> admin
        TO DELETE EMPLOYEE DETAILS
        
        RESPONSE:
            STATUS STRUCT ACCORDING TO SUCCESS OR FAILURE

    6. (/auth/login) , METHOD -> POST, AUTHORIZATION -> any
        TO LOGG IN, SETS JWT TOKEN IN COOKIE 

        JWT CLAIMS STRUCT:
            type Claims struct {
                Email  string `json:"email"`
                UserId int    `json:"user_id"`
                Role   string `json:"role"`
                jwt.StandardClaims
            }


        REQUEST STRUCT:
            type LoginReq struct {
                Email    string `json:"email" validate:"required,email"`
                Password string `json:"password" validate:"required,password"`
            }

        REQUEST:
            {
                "email": "...",
                "password": "..."
            }
        RESPONSE STRUCT :
            type LoginRes struct {
                Status Status          `json:"status"`
                Emp    entity.Employee `json:"employee"`
            }

        RESPONSE:
            {
                "status": {
                    "success": {true|false}
                },
                "employee": {
                    "user_id": ...,
                    "email": "...",
                    "role": "..."
                }
            }

    7. (/auth/logout), METHOD ->POST ,AUTHORIZATION ->IF YOU'RE LOGGED IN
        TO LOG OUT

        RESPONSE:
             STATUS STRUCT ACCORDING TO SUCCESS OR FAILURE
```           

## LIMITATIONS:
    1. NOT IMPLEMENTED LOGGER FOR LOGGING USED STANDARD GO LOGGER
    2. FORGET PASSWORD AND PASSWORD UPDATION 
    3. REFRESHING JWT COOKIE ACCORDING TO LAST ACCESS TIME





