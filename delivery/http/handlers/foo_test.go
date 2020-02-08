package handlers

import (
	"fmt"
	"github.com/hobord/golang-poc-rest/domain/entity"
	"github.com/hobord/golang-poc-rest/usecase/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/icrowley/fake"
	mock "github.com/stretchr/testify/mock"
)

func TestGetByID(t *testing.T) {
	type wantStruct struct {
		httpStatusCode int
		responseString string
	}

	fakeID := fake.Sentence()
	fakeTitle := fake.Sentence()
	var testCases = []struct {
		input *entity.Foo
		want  wantStruct
	}{
		{
			input: &entity.Foo{
				ID:    fakeID,
				Title: fakeTitle,
			},
			want: wantStruct{
				httpStatusCode: http.StatusOK,
				responseString: fmt.Sprintf(`{"id":"%s","title":"%s"}`, fakeID, fakeTitle),
			},
		},
		{
			input: nil,
			want: wantStruct{
				httpStatusCode: http.StatusNotFound,
				responseString: "no resource found\n",
			},
		},
	}

	for _, testCase := range testCases {

		// Create a mock uses case interactor and mock the results
		mockUsecase := &mocks.FooInteractorInterface{}
		mockUsecase.On("GetByID", mock.Anything, mock.Anything).Return(testCase.input, nil)

		// Create a test HTTPApp with moc usecase
		app := CreateFooRestHTTPModule(mockUsecase)

		var req *http.Request
		var err error
		// Create a test request
		if testCase.input == nil {
			req, err = http.NewRequest("GET", "/entity/NOTFOUND_SIMULATION", nil)
		} else {
			req, err = http.NewRequest("GET", "/entity/"+testCase.input.ID, nil)
		}
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()

		// Create a router and assign our handler func
		router := mux.NewRouter()
		router.HandleFunc("/entity/{id}", app.GetByID)

		// Make a request into the router
		router.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != testCase.want.httpStatusCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, testCase.want.httpStatusCode)
		}

		// Check the response body is what we expect.
		expected := testCase.want.responseString
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}
}

// TODO: implement all handler func tests