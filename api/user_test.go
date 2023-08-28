package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	mockdb "github.com/tech_school/simple_bank/db/mock"
	db "github.com/tech_school/simple_bank/db/sqlc"
	"github.com/tech_school/simple_bank/utils/pass"
	"github.com/tech_school/simple_bank/utils/random"
	"go.uber.org/mock/gomock"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

// parameter x pada function Matches menerima object apapun yang dikirimkan pada function CreateUser
// parameter x pada function Matches menerima db.CreateUserParams yang dikirikan pada function server.store.CreateUser di API handler
func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	argInCreateUser, ok := x.(db.CreateUserParams) // x dicasting
	if !ok {
		return false
	}

	err := pass.CheckPassword(e.password, argInCreateUser.HashedPassword)
	if err != nil {
		return false
	}

	// sebelum di cek field pada e.arg.HashedPassword yg tadinya kosong diisi dengan nilai dari  argInCreateUser.HashedPassword.
	// hal ini bisa kita lakukan karena sebelumnya kita sudah mengetahui bahwa e.password sama dengan argInCreateUser.HashedPassword
	e.arg.HashedPassword = argInCreateUser.HashedPassword

	return reflect.DeepEqual(e.arg, argInCreateUser)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

// create custom Equal function
func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		reqBody       gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(store *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			reqBody: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, responseRecorder.Code)
				requireBodyMatchUser(t, responseRecorder.Body, user)
			},
		},
		{
			name: "InternalError",
			reqBody: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, sql.ErrConnDone)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			reqBody: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrUniqueViolation)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, responseRecorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			reqBody: gin.H{
				"username":  "invalid_username#1",
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			reqBody: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     "invalid_email",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			reqBody: gin.H{
				"username":  user.Username,
				"password":  123,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
	}

	for i, _ := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// start test http server and send request
			server := newTestServer(t, store)

			// BUILD STUBS
			testCase.buildStubs(store)

			// httpnewrecorder berfungsi untuk merecord response of the api request
			responseRecorder := httptest.NewRecorder()

			// marshall data to json (like json_encode)
			marshalled, _ := json.Marshal(testCase.reqBody)

			// http.NewRequest(method, url,requestbody)
			request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(marshalled))
			require.NoError(t, err)

			// function dibawah mengirim request melalui server router dan merecord responsenya di recorder
			server.routers.ServeHTTP(responseRecorder, request)

			// CHECK
			testCase.checkResponse(t, responseRecorder)
		})

	}

}

func TestLoginUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		reqBody       gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(store *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			reqBody: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, responseRecorder.Code)
			},
		},
		{
			name: "UserNotFound",
			reqBody: gin.H{
				"username": "NonExistedUsername",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, responseRecorder.Code)
			},
		},
		{
			name: "IncorrectPassword",
			reqBody: gin.H{
				"username": user.Username,
				"password": "incorect_password",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
			},
		},
		{
			name: "InternalError",
			reqBody: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			reqBody: gin.H{
				"username": "invalid-username#1",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
	}

	for i, _ := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// start test http server and send request
			server := newTestServer(t, store)

			// BUILD STUBS
			testCase.buildStubs(store)

			// httpnewrecorder berfungsi untuk merecord response of the api request
			responseRecorder := httptest.NewRecorder()

			// marshall data to json (like json_encode)
			marshalled, _ := json.Marshal(testCase.reqBody)

			// http.NewRequest(method, url,requestbody)
			request, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(marshalled))
			require.NoError(t, err)

			// function dibawah mengirim request melalui server router dan merecord responsenya di recorder
			server.routers.ServeHTTP(responseRecorder, request)

			// CHECK
			testCase.checkResponse(t, responseRecorder)
		})

	}

}

func randomUser(t *testing.T) (user db.User, password string) {
	password = random.RandomString(6, "abcedfghijklmnopqrstuvwxyz")
	hashedPassword, err := pass.HashedPassword(password)
	require.NoError(t, err)

	person := random.RandomPerson()

	user = db.User{
		Username:       person.Username,
		HashedPassword: hashedPassword,
		FullName:       person.FullName,
		Email:          person.Email,
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, inputUser db.User) {
	// read data from response body
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	// unmarshal data to getAccount object
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.NoError(t, err)
	require.Equal(t, inputUser.Username, gotUser.Username)
	require.Equal(t, inputUser.FullName, gotUser.FullName)
	require.Equal(t, inputUser.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
