package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ChandanGupta31/student-api/internal/types"
	"github.com/ChandanGupta31/student-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating a student variable
		var student types.Student


		// decoding the request data
		err := json.NewDecoder(r.Body).Decode(&student)
		// handling empty body error
		if errors.Is(err, io.EOF){
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		// Handling General Error
		if err!=nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request validation
		if err:=validator.New().Struct(student); err!=nil{
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
		}


		// Providing Response
		response.WriteJson(w, http.StatusCreated, map[string]string {"success" : "OK"})
	}
}
