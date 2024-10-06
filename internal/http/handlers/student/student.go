package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ChandanGupta31/student-api/internal/storage"
	"github.com/ChandanGupta31/student-api/internal/types"
	"github.com/ChandanGupta31/student-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating a student variable
		var student types.Student

		// decoding the request data
		err := json.NewDecoder(r.Body).Decode(&student)
		// handling empty body error
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		// Handling General Error
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request validation
		if err := validator.New().Struct(student); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		id, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err!=nil{
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		// Providing Response
		response.WriteJson(w, http.StatusCreated, map[string]int64{"success": id})
	}
}

func GetStudentById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}
		student, err := storage.GetStudentById(intId)

		if err!= nil {
			response.WriteJson(w, http.StatusInternalServerError, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}
