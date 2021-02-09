package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ersmith/mailgun-coding-challenge/models"
	"github.com/ersmith/mailgun-coding-challenge/test"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func TestGetDomainNewDomain(t *testing.T) {
	app := getApp()
	domainName := test.RandomDomainName(20)
	req, err := http.NewRequest("GET", fmt.Sprintf("/domains/%s", domainName), nil)
	test.CheckError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/domains/{domain}", app.getDomainHandler).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v",
			status, http.StatusOK)
	}

	expected := fmt.Sprintf(`{"id":0,"domain_name":"%v","delivered":0,"bounced":0,"catch_all":"unknown"}`, domainName)
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetDomainWithExistingDomainDelivered(t *testing.T) {
	dbPool := test.CreateTestPgxPool(t)
	app := getApp()
	domainName := test.RandomDomainName(20)
	domain := &models.Domain{
		DomainName: domainName,
	}
	domain.IncrementDelivered(dbPool)

	domain, err := models.GetDomain(dbPool, zap.NewNop().Sugar(), domainName)
	test.CheckError(t, err)

	req, err := http.NewRequest("GET", fmt.Sprintf("/domains/%s", domainName), nil)
	test.CheckError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/domains/{domain}", app.getDomainHandler).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v",
			status, http.StatusOK)
	}

	expected := fmt.Sprintf(`{"id":%d,"domain_name":"%v","delivered":1,"bounced":0,"catch_all":"unknown"}`, domain.Id, domainName)
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetDomainWithExistingDomainBounced(t *testing.T) {
	dbPool := test.CreateTestPgxPool(t)
	app := getApp()
	domainName := test.RandomDomainName(20)
	domain := &models.Domain{
		DomainName: domainName,
	}
	domain.IncrementBounced(dbPool)

	domain, err := models.GetDomain(dbPool, zap.NewNop().Sugar(), domainName)
	test.CheckError(t, err)

	req, err := http.NewRequest("GET", fmt.Sprintf("/domains/%s", domainName), nil)
	test.CheckError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/domains/{domain}", app.getDomainHandler).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v",
			status, http.StatusOK)
	}

	expected := fmt.Sprintf(`{"id":%d,"domain_name":"%v","delivered":0,"bounced":1,"catch_all":"%s"}`,
		domain.Id,
		domainName,
		models.IsNotCatchAllStatus)
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPutEventBounced(t *testing.T) {
	dbPool := test.CreateTestPgxPool(t)
	app := getApp()
	domainName := test.RandomDomainName(20)

	req, err := http.NewRequest("PUT", fmt.Sprintf("/events/%s/bounced", domainName), nil)
	test.CheckError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/events/{domain}/bounced", app.putEventBouncedHandler).Methods("PUT")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v",
			status, http.StatusOK)
	}

	domain, err := models.GetDomain(dbPool, zap.NewNop().Sugar(), domainName)
	test.CheckError(t, err)

	assert.Equal(t, 1, domain.Bounced)
}

func TestPutEventDelivered(t *testing.T) {
	dbPool := test.CreateTestPgxPool(t)
	app := getApp()
	domainName := test.RandomDomainName(20)

	req, err := http.NewRequest("PUT", fmt.Sprintf("/events/%s/delivered", domainName), nil)
	test.CheckError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/events/{domain}/delivered", app.putEventDeliveredHandler).Methods("PUT")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v",
			status, http.StatusOK)
	}

	domain, err := models.GetDomain(dbPool, zap.NewNop().Sugar(), domainName)
	test.CheckError(t, err)

	assert.Equal(t, 1, domain.Delivered)
}

func getApp() *App {
	app := App{}
	app.Initialize(&test.DatabaseConfig, zap.NewNop().Sugar())
	return &app
}
