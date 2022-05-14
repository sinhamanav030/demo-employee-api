package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Employee struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	Designation string `json:"designation"`
	DOB         string `json:"dob"`
}

const conn = `mongodb+srv://manav03:7bRpcqdQtiKaUvZb@gopherslab.uyaxi.mongodb.net/myFirstDatabase?retryWrites=true&w=majority`
const dbname = "company"
const colname = "employee"

var collection *mongo.Collection

func init() {

	clientOption := options.Client().ApplyURI(conn)
	cl, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatalln(err)
	}
	collection = cl.Database(dbname).Collection(colname)

}

func Rmain() {
	r := mux.NewRouter()
	r.Handle("/favicon.ico", http.NotFoundHandler())
	r.HandleFunc("/users", adduser).Methods("POST")
	r.HandleFunc("/", app)
	http.ListenAndServe(":8080", r)
}

func app(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func adduser(w http.ResponseWriter, r *http.Request) {
	E1 := Employee{r.FormValue("name"), r.FormValue("code"), r.FormValue("email"), r.FormValue("gender"), r.FormValue("designation"), r.FormValue("dob")}
	err := Create(E1)
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(E1)
}

func Create(emp Employee) error {
	in, err := collection.InsertOne(context.Background(), emp)
	if err != nil {
		return err
	}
	fmt.Println(in.InsertedID)
	return nil
}
