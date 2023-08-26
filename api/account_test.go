package api

import (
	"bytes"	
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	mockdb "github.com/tech_school/simple_bank/db/mock"
	db "github.com/tech_school/simple_bank/db/sqlc"
	"github.com/tech_school/simple_bank/token"
	"github.com/tech_school/simple_bank/utils/random"
	"go.uber.org/mock/gomock"
)

// tes untuk GetAccountApi handler
func TestGetAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(store *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		{
			name:      "UnAuthorized",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "Unauthorized_user", time.Minute)

			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				// check response status
				require.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
			},
		},
		{
			name:      "No Auth",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(0)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				// check response status
				require.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
			},
		},
		{
			name:      "Ok",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)

			},
			buildStubs: func(store *mockdb.MockStore) {
				// BUILD STUBS
				// gomock.Any() == karena ctx bisa bertipe any
				// argument selanjutnya harus berisi id yang sama dengan id yang ada pada randomAccount
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
				// function diatas jika ditranslatekan : i expected GetAccount of the store to be called with any context and this specific account id argument.
				// i expected this function to be called just one time.
				// i expected this function to return account and nil error
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				// check response status
				require.Equal(t, http.StatusOK, responseRecorder.Code)

				// check response body
				requireBodyMatchAccount(t, responseRecorder.Body, account)
			},
		},

		// TODO: add more cases
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)

			},
			buildStubs: func(store *mockdb.MockStore) {
				// BUILD STUBS
				// gomock.Any() == karena ctx bisa bertipe any
				// argument selanjutnya harus berisi id yang sama dengan id yang ada pada randomAccount
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				// check response status
				require.Equal(t, http.StatusNotFound, responseRecorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)

			},
			buildStubs: func(store *mockdb.MockStore) {
				// BUILD STUBS
				// gomock.Any() == karena ctx bisa bertipe any
				// argument selanjutnya harus berisi id yang sama dengan id yang ada pada randomAccount
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					// sql.errorConnDone is consider as internal error
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				// check response status
				require.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
			},
		},
		{
			name:      "InvalidId",
			accountID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)

			},
			buildStubs: func(store *mockdb.MockStore) {
				// BUILD STUBS
				// gomock.Any() == karena ctx bisa bertipe any
				// argument selanjutnya harus berisi id yang sama dengan id yang ada pada randomAccount
				store.EXPECT().
					GetAccount(gomock.All(), gomock.Any()).
					// karena dari req saja sudah salah maka getaccount tidak mungkin dipanggil
					Times(0)
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				// check response status
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

			// url path of the api we want to call
			url := fmt.Sprintf("/accounts/%d", testCase.accountID)
			// http.NewRequest(method, url,requestbody)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)

			// function dibawah mengirim request melalui server router dan merecord responsenya di recorder
			server.routers.ServeHTTP(responseRecorder, request)

			// CHECK
			testCase.checkResponse(t, responseRecorder)
		})

	}

}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       random.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  random.RandomMoney(),
		Currency: random.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, inputAccount db.Account) {
	// read data from response body
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var getAccount db.Account
	// unmarshal data to getAccount object
	err = json.Unmarshal(data, &getAccount)
	require.NoError(t, err)

	require.Equal(t, inputAccount, getAccount)
}
