package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gogocoding/micro-go-course/endpoint"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

const ContextReqUUid = "req_uuid"

func MakeHttpHandler(ctx context.Context, endpoints *endpoint.UserEndpoints) http.Handler {
	r := mux.NewRouter()

	kitLog := log.NewLogfmtLogger(os.Stderr)

	kitLog = log.With(kitLog, "ts", log.DefaultTimestampUTC)
	kitLog = log.With(kitLog, "caller", log.DefaultCaller)

	option := []kithttp.ServerOption{
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			UUID := uuid.NewV5(uuid.Must(uuid.NewV4(), nil), "req_uuid").String()
			kitLog.Log("给请求添加uuid=", UUID)
			ctx = context.WithValue(ctx, ContextReqUUid, UUID)
			return ctx
		}),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(kitLog)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/register").Handler(kithttp.NewServer(
		endpoints.RegisterEndpoint,
		decodeRegisterRequest,
		encodeJSONResponse,
		option...,
	))

	r.Methods("POST").Path("/login").Handler(kithttp.NewServer(
		endpoints.LoginEndpoint,
		decodeLoginRequest,
		encodeJSONResponse,
		option...,
	))
	return r
}

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}

func encodeJSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
