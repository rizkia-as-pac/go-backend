package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	mockdb "github.com/tech_school/simple_bank/db/mock"
	db "github.com/tech_school/simple_bank/db/sqlc"
	"github.com/tech_school/simple_bank/token"
	curr "github.com/tech_school/simple_bank/utils/currency"
	"go.uber.org/mock/gomock"
)

func TestTransferMoneyTechSchoolAPI(t *testing.T) {
	amount := int64(10)

	user1, _ := randomUser(t)
	user2, _ := randomUser(t)
	user3, _ := randomUser(t)

	account1 := randomAccount(user1.Username)
	account2 := randomAccount(user2.Username)
	account3 := randomAccount(user3.Username)

	account1.Currency = curr.USD
	account2.Currency = curr.USD
	account3.Currency = curr.JPY

	testCases := []struct {
		name          string
		reqBody       gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(store *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			reqBody: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        curr.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}

				store.EXPECT().
					TransferTxV2(gomock.Any(), gomock.Eq(arg)).
					Times(1)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, responseRecorder.Code)
			},
		},
		{
			name: "UnauthorizedUser",
			reqBody: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        curr.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user2.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					TransferTxV2(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			reqBody: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        curr.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					TransferTxV2(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
			},
		},
		{
			name: "FromAccountNotFound",
			reqBody: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        curr.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, db.ErrRecordNotFound)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(0)

				store.EXPECT().
					TransferTxV2(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, responseRecorder.Code)
			},
		},
		{
			name: "ToAccountNotFound",
			reqBody: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        curr.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, db.ErrRecordNotFound)

				store.EXPECT().
					TransferTxV2(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, responseRecorder.Code)
			},
		},
		{
			name: "FromAccountCurrencyMismatch",
			reqBody: gin.H{
				"from_account_id": account3.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        curr.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user3.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account3.ID)).
					Times(1).
					Return(account3, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					TransferTxV2(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
		{
			name: "ToAccountCurrencyMismatch",
			reqBody: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account3.ID,
				"amount":          amount,
				"currency":        curr.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account3.ID)).
					Times(1).
					Return(account3, nil)

				store.EXPECT().
					TransferTxV2(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, responseRecorder.Code)
			},
		},
		{
			name: "GetAccountError",
			reqBody: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        curr.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					TransferTxV2(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
			},
		},
		{
			name: "TransferTxError",
			reqBody: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        curr.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().TransferTxV2(gomock.Any(), gomock.Any()).Times(1).Return(db.TransferTxResult{}, sql.ErrTxDone)
			},
			checkResponse: func(store *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
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
			request, err := http.NewRequest(http.MethodPost, "/transfers", bytes.NewReader(marshalled))
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)

			// function dibawah mengirim request melalui server router dan merecord responsenya di recorder
			server.routers.ServeHTTP(responseRecorder, request)

			// CHECK
			testCase.checkResponse(t, responseRecorder)
		})

	}

}
