// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/service"
	"net/http"
	"sync"
)

// Ensure, that InitialiserMock does implement service.Initialiser.
// If this is not the case, regenerate this file with moq.
var _ service.Initialiser = &InitialiserMock{}

// InitialiserMock is a mock implementation of service.Initialiser.
//
//	func TestSomethingThatUsesInitialiser(t *testing.T) {
//
//		// make and configure a mocked service.Initialiser
//		mockedInitialiser := &InitialiserMock{
//			GetCantabularClientFunc: func(cfg config.CantabularConfig) service.CantabularClient {
//				panic("mock out the GetCantabularClient method")
//			},
//			GetDatasetAPIClientFunc: func(cfg *config.Config) service.DatasetAPIClient {
//				panic("mock out the GetDatasetAPIClient method")
//			},
//			GetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer {
//				panic("mock out the GetHTTPServer method")
//			},
//			GetHTTPServerWithOtelFunc: func(bindAddr string, router http.Handler) service.HTTPServer {
//				panic("mock out the GetHTTPServerWithOtel method")
//			},
//			GetHealthCheckFunc: func(cfg *config.Config, time string, commit string, version string) (service.HealthChecker, error) {
//				panic("mock out the GetHealthCheck method")
//			},
//			GetMongoClientFunc: func(ctx context.Context, cfg *config.Config) (service.MongoClient, error) {
//				panic("mock out the GetMongoClient method")
//			},
//			GetResponderFunc: func() service.Responder {
//				panic("mock out the GetResponder method")
//			},
//		}
//
//		// use mockedInitialiser in code that requires service.Initialiser
//		// and then make assertions.
//
//	}
type InitialiserMock struct {
	// GetCantabularClientFunc mocks the GetCantabularClient method.
	GetCantabularClientFunc func(cfg config.CantabularConfig) service.CantabularClient

	// GetDatasetAPIClientFunc mocks the GetDatasetAPIClient method.
	GetDatasetAPIClientFunc func(cfg *config.Config) service.DatasetAPIClient

	// GetHTTPServerFunc mocks the GetHTTPServer method.
	GetHTTPServerFunc func(bindAddr string, router http.Handler) service.HTTPServer

	// GetHTTPServerWithOtelFunc mocks the GetHTTPServerWithOtel method.
	GetHTTPServerWithOtelFunc func(bindAddr string, router http.Handler) service.HTTPServer

	// GetHealthCheckFunc mocks the GetHealthCheck method.
	GetHealthCheckFunc func(cfg *config.Config, time string, commit string, version string) (service.HealthChecker, error)

	// GetMongoClientFunc mocks the GetMongoClient method.
	GetMongoClientFunc func(ctx context.Context, cfg *config.Config) (service.MongoClient, error)

	// GetResponderFunc mocks the GetResponder method.
	GetResponderFunc func() service.Responder

	// calls tracks calls to the methods.
	calls struct {
		// GetCantabularClient holds details about calls to the GetCantabularClient method.
		GetCantabularClient []struct {
			// Cfg is the cfg argument value.
			Cfg config.CantabularConfig
		}
		// GetDatasetAPIClient holds details about calls to the GetDatasetAPIClient method.
		GetDatasetAPIClient []struct {
			// Cfg is the cfg argument value.
			Cfg *config.Config
		}
		// GetHTTPServer holds details about calls to the GetHTTPServer method.
		GetHTTPServer []struct {
			// BindAddr is the bindAddr argument value.
			BindAddr string
			// Router is the router argument value.
			Router http.Handler
		}
		// GetHTTPServerWithOtel holds details about calls to the GetHTTPServerWithOtel method.
		GetHTTPServerWithOtel []struct {
			// BindAddr is the bindAddr argument value.
			BindAddr string
			// Router is the router argument value.
			Router http.Handler
		}
		// GetHealthCheck holds details about calls to the GetHealthCheck method.
		GetHealthCheck []struct {
			// Cfg is the cfg argument value.
			Cfg *config.Config
			// Time is the time argument value.
			Time string
			// Commit is the commit argument value.
			Commit string
			// Version is the version argument value.
			Version string
		}
		// GetMongoClient holds details about calls to the GetMongoClient method.
		GetMongoClient []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Cfg is the cfg argument value.
			Cfg *config.Config
		}
		// GetResponder holds details about calls to the GetResponder method.
		GetResponder []struct {
		}
	}
	lockGetCantabularClient   sync.RWMutex
	lockGetDatasetAPIClient   sync.RWMutex
	lockGetHTTPServer         sync.RWMutex
	lockGetHTTPServerWithOtel sync.RWMutex
	lockGetHealthCheck        sync.RWMutex
	lockGetMongoClient        sync.RWMutex
	lockGetResponder          sync.RWMutex
}

// GetCantabularClient calls GetCantabularClientFunc.
func (mock *InitialiserMock) GetCantabularClient(cfg config.CantabularConfig) service.CantabularClient {
	if mock.GetCantabularClientFunc == nil {
		panic("InitialiserMock.GetCantabularClientFunc: method is nil but Initialiser.GetCantabularClient was just called")
	}
	callInfo := struct {
		Cfg config.CantabularConfig
	}{
		Cfg: cfg,
	}
	mock.lockGetCantabularClient.Lock()
	mock.calls.GetCantabularClient = append(mock.calls.GetCantabularClient, callInfo)
	mock.lockGetCantabularClient.Unlock()
	return mock.GetCantabularClientFunc(cfg)
}

// GetCantabularClientCalls gets all the calls that were made to GetCantabularClient.
// Check the length with:
//
//	len(mockedInitialiser.GetCantabularClientCalls())
func (mock *InitialiserMock) GetCantabularClientCalls() []struct {
	Cfg config.CantabularConfig
} {
	var calls []struct {
		Cfg config.CantabularConfig
	}
	mock.lockGetCantabularClient.RLock()
	calls = mock.calls.GetCantabularClient
	mock.lockGetCantabularClient.RUnlock()
	return calls
}

// GetDatasetAPIClient calls GetDatasetAPIClientFunc.
func (mock *InitialiserMock) GetDatasetAPIClient(cfg *config.Config) service.DatasetAPIClient {
	if mock.GetDatasetAPIClientFunc == nil {
		panic("InitialiserMock.GetDatasetAPIClientFunc: method is nil but Initialiser.GetDatasetAPIClient was just called")
	}
	callInfo := struct {
		Cfg *config.Config
	}{
		Cfg: cfg,
	}
	mock.lockGetDatasetAPIClient.Lock()
	mock.calls.GetDatasetAPIClient = append(mock.calls.GetDatasetAPIClient, callInfo)
	mock.lockGetDatasetAPIClient.Unlock()
	return mock.GetDatasetAPIClientFunc(cfg)
}

// GetDatasetAPIClientCalls gets all the calls that were made to GetDatasetAPIClient.
// Check the length with:
//
//	len(mockedInitialiser.GetDatasetAPIClientCalls())
func (mock *InitialiserMock) GetDatasetAPIClientCalls() []struct {
	Cfg *config.Config
} {
	var calls []struct {
		Cfg *config.Config
	}
	mock.lockGetDatasetAPIClient.RLock()
	calls = mock.calls.GetDatasetAPIClient
	mock.lockGetDatasetAPIClient.RUnlock()
	return calls
}

// GetHTTPServer calls GetHTTPServerFunc.
func (mock *InitialiserMock) GetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	if mock.GetHTTPServerFunc == nil {
		panic("InitialiserMock.GetHTTPServerFunc: method is nil but Initialiser.GetHTTPServer was just called")
	}
	callInfo := struct {
		BindAddr string
		Router   http.Handler
	}{
		BindAddr: bindAddr,
		Router:   router,
	}
	mock.lockGetHTTPServer.Lock()
	mock.calls.GetHTTPServer = append(mock.calls.GetHTTPServer, callInfo)
	mock.lockGetHTTPServer.Unlock()
	return mock.GetHTTPServerFunc(bindAddr, router)
}

// GetHTTPServerCalls gets all the calls that were made to GetHTTPServer.
// Check the length with:
//
//	len(mockedInitialiser.GetHTTPServerCalls())
func (mock *InitialiserMock) GetHTTPServerCalls() []struct {
	BindAddr string
	Router   http.Handler
} {
	var calls []struct {
		BindAddr string
		Router   http.Handler
	}
	mock.lockGetHTTPServer.RLock()
	calls = mock.calls.GetHTTPServer
	mock.lockGetHTTPServer.RUnlock()
	return calls
}

// GetHTTPServerWithOtel calls GetHTTPServerWithOtelFunc.
func (mock *InitialiserMock) GetHTTPServerWithOtel(bindAddr string, router http.Handler) service.HTTPServer {
	if mock.GetHTTPServerWithOtelFunc == nil {
		panic("InitialiserMock.GetHTTPServerWithOtelFunc: method is nil but Initialiser.GetHTTPServerWithOtel was just called")
	}
	callInfo := struct {
		BindAddr string
		Router   http.Handler
	}{
		BindAddr: bindAddr,
		Router:   router,
	}
	mock.lockGetHTTPServerWithOtel.Lock()
	mock.calls.GetHTTPServerWithOtel = append(mock.calls.GetHTTPServerWithOtel, callInfo)
	mock.lockGetHTTPServerWithOtel.Unlock()
	return mock.GetHTTPServerWithOtelFunc(bindAddr, router)
}

// GetHTTPServerWithOtelCalls gets all the calls that were made to GetHTTPServerWithOtel.
// Check the length with:
//
//	len(mockedInitialiser.GetHTTPServerWithOtelCalls())
func (mock *InitialiserMock) GetHTTPServerWithOtelCalls() []struct {
	BindAddr string
	Router   http.Handler
} {
	var calls []struct {
		BindAddr string
		Router   http.Handler
	}
	mock.lockGetHTTPServerWithOtel.RLock()
	calls = mock.calls.GetHTTPServerWithOtel
	mock.lockGetHTTPServerWithOtel.RUnlock()
	return calls
}

// GetHealthCheck calls GetHealthCheckFunc.
func (mock *InitialiserMock) GetHealthCheck(cfg *config.Config, time string, commit string, version string) (service.HealthChecker, error) {
	if mock.GetHealthCheckFunc == nil {
		panic("InitialiserMock.GetHealthCheckFunc: method is nil but Initialiser.GetHealthCheck was just called")
	}
	callInfo := struct {
		Cfg     *config.Config
		Time    string
		Commit  string
		Version string
	}{
		Cfg:     cfg,
		Time:    time,
		Commit:  commit,
		Version: version,
	}
	mock.lockGetHealthCheck.Lock()
	mock.calls.GetHealthCheck = append(mock.calls.GetHealthCheck, callInfo)
	mock.lockGetHealthCheck.Unlock()
	return mock.GetHealthCheckFunc(cfg, time, commit, version)
}

// GetHealthCheckCalls gets all the calls that were made to GetHealthCheck.
// Check the length with:
//
//	len(mockedInitialiser.GetHealthCheckCalls())
func (mock *InitialiserMock) GetHealthCheckCalls() []struct {
	Cfg     *config.Config
	Time    string
	Commit  string
	Version string
} {
	var calls []struct {
		Cfg     *config.Config
		Time    string
		Commit  string
		Version string
	}
	mock.lockGetHealthCheck.RLock()
	calls = mock.calls.GetHealthCheck
	mock.lockGetHealthCheck.RUnlock()
	return calls
}

// GetMongoClient calls GetMongoClientFunc.
func (mock *InitialiserMock) GetMongoClient(ctx context.Context, cfg *config.Config) (service.MongoClient, error) {
	if mock.GetMongoClientFunc == nil {
		panic("InitialiserMock.GetMongoClientFunc: method is nil but Initialiser.GetMongoClient was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Cfg *config.Config
	}{
		Ctx: ctx,
		Cfg: cfg,
	}
	mock.lockGetMongoClient.Lock()
	mock.calls.GetMongoClient = append(mock.calls.GetMongoClient, callInfo)
	mock.lockGetMongoClient.Unlock()
	return mock.GetMongoClientFunc(ctx, cfg)
}

// GetMongoClientCalls gets all the calls that were made to GetMongoClient.
// Check the length with:
//
//	len(mockedInitialiser.GetMongoClientCalls())
func (mock *InitialiserMock) GetMongoClientCalls() []struct {
	Ctx context.Context
	Cfg *config.Config
} {
	var calls []struct {
		Ctx context.Context
		Cfg *config.Config
	}
	mock.lockGetMongoClient.RLock()
	calls = mock.calls.GetMongoClient
	mock.lockGetMongoClient.RUnlock()
	return calls
}

// GetResponder calls GetResponderFunc.
func (mock *InitialiserMock) GetResponder() service.Responder {
	if mock.GetResponderFunc == nil {
		panic("InitialiserMock.GetResponderFunc: method is nil but Initialiser.GetResponder was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetResponder.Lock()
	mock.calls.GetResponder = append(mock.calls.GetResponder, callInfo)
	mock.lockGetResponder.Unlock()
	return mock.GetResponderFunc()
}

// GetResponderCalls gets all the calls that were made to GetResponder.
// Check the length with:
//
//	len(mockedInitialiser.GetResponderCalls())
func (mock *InitialiserMock) GetResponderCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetResponder.RLock()
	calls = mock.calls.GetResponder
	mock.lockGetResponder.RUnlock()
	return calls
}
