package internalgrpc

import (
	"bytes"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SashaMelva/calendar_service/internal/app"
	"github.com/SashaMelva/calendar_service/internal/config"
	"github.com/SashaMelva/calendar_service/internal/logger"
	memorystorage "github.com/SashaMelva/calendar_service/internal/storage/memory"
	sqlstorage "github.com/SashaMelva/calendar_service/internal/storage/sql"
	"go.uber.org/zap/zapcore"
)

func TestHendlerEvent(t *testing.T) {
	log := logger.NewLogger(&config.ConfigLogger{
		Level:       zapcore.InfoLevel,
		LogEncoding: "console",
	})
	connection := sqlstorage.New(&config.ConfigDB{
		NameDB:   "calendardb",
		Host:     "127.0.0.1",
		Port:     "5436",
		User:     "postgres",
		Password: "qwer",
	}, log)

	memstorage := memorystorage.New(connection.StorageDb)
	calendar := app.New(log, memstorage)
	server := NewGRPCServer(log, calendar)

	log.Info("Starting listening on port 8080")
	port := ":8080"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Info("Listening on %s", port)
	srv := server.NewGRPCServer()

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	testCase := []struct {
		name       string
		method     string
		path       string
		body       []byte
		want       string
		statusCode int
	}{
		{
			name:       "get one event for id",
			method:     http.MethodGet,
			path:       "/event?id=1",
			body:       []byte(""),
			want:       `{"id":1,"title":"qw","date_time_start":"2003-09-03T17:00:00+04:00","date_time_end":"2003-09-03T17:00:00+04:00","description":"qwqwq"}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "get one event for id",
			method:     http.MethodGet,
			path:       "/event?id=900",
			body:       []byte(""),
			want:       "",
			statusCode: http.StatusOK,
		},
		{
			name:       "create new event",
			method:     http.MethodPost,
			path:       "/event/",
			body:       []byte(`{"title":"testPost","date_time_start":"2003-09-03T17:00:00+04:00","date_time_end":"2003-09-03T17:00:00+04:00","description":"qwqwq"}`),
			want:       "",
			statusCode: http.StatusOK,
		},
		{
			name:       "fail create new nil event",
			method:     http.MethodPost,
			path:       "/event/",
			body:       []byte(""),
			want:       "row title empty;row date start empty;",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "fail create event empty strong data",
			method:     http.MethodPost,
			path:       "/event/",
			body:       []byte(`{"title":"","date_time_start":"","date_time_end":"17:00:00+04:00","description":"qwqwq"}`),
			want:       "row title empty;row date start empty;date end param empty date;",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "delete event",
			method:     http.MethodDelete,
			path:       "/event?id=12",
			body:       []byte(""),
			want:       "",
			statusCode: http.StatusOK,
		},
		{
			name:       "fail delete event",
			method:     http.MethodDelete,
			path:       "/event?id=110",
			body:       []byte(""),
			want:       "not found event by id == 110",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Edit event",
			method:     http.MethodPut,
			path:       "/event",
			body:       []byte(`{"id":2,"title":"updateEvent","date_time_start":"2003-09-03T17:00:00+04:00","date_time_end":"2003-09-03T17:00:00+04:00","description":"qwqwq"}`),
			want:       "",
			statusCode: http.StatusOK,
		},
		{
			name:       "edit event id not found",
			method:     http.MethodPut,
			path:       "/event",
			body:       []byte(`{"id":0,"title":"updateEvent","date_time_start":"2003-09-03T17:00:00+04:00","date_time_end":"2003-09-03T17:00:00+04:00","description":"qwqwq"}`),
			want:       "not found event by id",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "fail edit event",
			method:     http.MethodPut,
			path:       "/event",
			body:       []byte(`{"id":5,"title":"","date_time_start":"","date_time_end":"17:00:00+04:00","description":"qwqwq"}`),
			want:       "row title empty;row date start empty;date end param empty date;",
			statusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			bodyReader := bytes.NewReader(tc.body)
			request := httptest.NewRequest(tc.method, tc.path, bodyReader)
			responseRecorder := httptest.NewRecorder()

			serever.HendlerEvent(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}
		})
	}
}
