// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-population-types-api/service"
	"sync"
)

// Ensure, that CantabularClientMock does implement service.CantabularClient.
// If this is not the case, regenerate this file with moq.
var _ service.CantabularClient = &CantabularClientMock{}

// CantabularClientMock is a mock implementation of service.CantabularClient.
//
// 	func TestSomethingThatUsesCantabularClient(t *testing.T) {
//
// 		// make and configure a mocked service.CantabularClient
// 		mockedCantabularClient := &CantabularClientMock{
// 			CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error {
// 				panic("mock out the Checker method")
// 			},
// 			GetAreasFunc: func(contextMoqParam context.Context, getAreasRequest cantabular.GetAreasRequest) (*cantabular.GetAreasResponse, error) {
// 				panic("mock out the GetAreas method")
// 			},
// 			GetGeographyDimensionsFunc: func(ctx context.Context, req cantabular.GetGeographyDimensionsRequest) (*cantabular.GetGeographyDimensionsResponse, error) {
// 				panic("mock out the GetGeographyDimensions method")
// 			},
// 			ListDatasetsFunc: func(ctx context.Context) ([]string, error) {
// 				panic("mock out the ListDatasets method")
// 			},
// 			StatusCodeFunc: func(err error) int {
// 				panic("mock out the StatusCode method")
// 			},
// 		}
//
// 		// use mockedCantabularClient in code that requires service.CantabularClient
// 		// and then make assertions.
//
// 	}
type CantabularClientMock struct {
	// CheckerFunc mocks the Checker method.
	CheckerFunc func(ctx context.Context, state *healthcheck.CheckState) error

	// GetAreasFunc mocks the GetAreas method.
	GetAreasFunc func(contextMoqParam context.Context, getAreasRequest cantabular.GetAreasRequest) (*cantabular.GetAreasResponse, error)

	// GetGeographyDimensionsFunc mocks the GetGeographyDimensions method.
	GetGeographyDimensionsFunc func(ctx context.Context, req cantabular.GetGeographyDimensionsRequest) (*cantabular.GetGeographyDimensionsResponse, error)

	// ListDatasetsFunc mocks the ListDatasets method.
	ListDatasetsFunc func(ctx context.Context) ([]string, error)

	// StatusCodeFunc mocks the StatusCode method.
	StatusCodeFunc func(err error) int

	// calls tracks calls to the methods.
	calls struct {
		// Checker holds details about calls to the Checker method.
		Checker []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// State is the state argument value.
			State *healthcheck.CheckState
		}
		// GetAreas holds details about calls to the GetAreas method.
		GetAreas []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// GetAreasRequest is the getAreasRequest argument value.
			GetAreasRequest cantabular.GetAreasRequest
		}
		// GetGeographyDimensions holds details about calls to the GetGeographyDimensions method.
		GetGeographyDimensions []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Req is the req argument value.
			Req cantabular.GetGeographyDimensionsRequest
		}
		// ListDatasets holds details about calls to the ListDatasets method.
		ListDatasets []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// StatusCode holds details about calls to the StatusCode method.
		StatusCode []struct {
			// Err is the err argument value.
			Err error
		}
	}
	lockChecker                sync.RWMutex
	lockGetAreas               sync.RWMutex
	lockGetGeographyDimensions sync.RWMutex
	lockListDatasets           sync.RWMutex
	lockStatusCode             sync.RWMutex
}

// Checker calls CheckerFunc.
func (mock *CantabularClientMock) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	if mock.CheckerFunc == nil {
		panic("CantabularClientMock.CheckerFunc: method is nil but CantabularClient.Checker was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		State *healthcheck.CheckState
	}{
		Ctx:   ctx,
		State: state,
	}
	mock.lockChecker.Lock()
	mock.calls.Checker = append(mock.calls.Checker, callInfo)
	mock.lockChecker.Unlock()
	return mock.CheckerFunc(ctx, state)
}

// CheckerCalls gets all the calls that were made to Checker.
// Check the length with:
//     len(mockedCantabularClient.CheckerCalls())
func (mock *CantabularClientMock) CheckerCalls() []struct {
	Ctx   context.Context
	State *healthcheck.CheckState
} {
	var calls []struct {
		Ctx   context.Context
		State *healthcheck.CheckState
	}
	mock.lockChecker.RLock()
	calls = mock.calls.Checker
	mock.lockChecker.RUnlock()
	return calls
}

// GetAreas calls GetAreasFunc.
func (mock *CantabularClientMock) GetAreas(contextMoqParam context.Context, getAreasRequest cantabular.GetAreasRequest) (*cantabular.GetAreasResponse, error) {
	if mock.GetAreasFunc == nil {
		panic("CantabularClientMock.GetAreasFunc: method is nil but CantabularClient.GetAreas was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		GetAreasRequest cantabular.GetAreasRequest
	}{
		ContextMoqParam: contextMoqParam,
		GetAreasRequest: getAreasRequest,
	}
	mock.lockGetAreas.Lock()
	mock.calls.GetAreas = append(mock.calls.GetAreas, callInfo)
	mock.lockGetAreas.Unlock()
	return mock.GetAreasFunc(contextMoqParam, getAreasRequest)
}

// GetAreasCalls gets all the calls that were made to GetAreas.
// Check the length with:
//     len(mockedCantabularClient.GetAreasCalls())
func (mock *CantabularClientMock) GetAreasCalls() []struct {
	ContextMoqParam context.Context
	GetAreasRequest cantabular.GetAreasRequest
} {
	var calls []struct {
		ContextMoqParam context.Context
		GetAreasRequest cantabular.GetAreasRequest
	}
	mock.lockGetAreas.RLock()
	calls = mock.calls.GetAreas
	mock.lockGetAreas.RUnlock()
	return calls
}

// GetGeographyDimensions calls GetGeographyDimensionsFunc.
func (mock *CantabularClientMock) GetGeographyDimensions(ctx context.Context, req cantabular.GetGeographyDimensionsRequest) (*cantabular.GetGeographyDimensionsResponse, error) {
	if mock.GetGeographyDimensionsFunc == nil {
		panic("CantabularClientMock.GetGeographyDimensionsFunc: method is nil but CantabularClient.GetGeographyDimensions was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Req cantabular.GetGeographyDimensionsRequest
	}{
		Ctx: ctx,
		Req: req,
	}
	mock.lockGetGeographyDimensions.Lock()
	mock.calls.GetGeographyDimensions = append(mock.calls.GetGeographyDimensions, callInfo)
	mock.lockGetGeographyDimensions.Unlock()
	return mock.GetGeographyDimensionsFunc(ctx, req)
}

// GetGeographyDimensionsCalls gets all the calls that were made to GetGeographyDimensions.
// Check the length with:
//     len(mockedCantabularClient.GetGeographyDimensionsCalls())
func (mock *CantabularClientMock) GetGeographyDimensionsCalls() []struct {
	Ctx context.Context
	Req cantabular.GetGeographyDimensionsRequest
} {
	var calls []struct {
		Ctx context.Context
		Req cantabular.GetGeographyDimensionsRequest
	}
	mock.lockGetGeographyDimensions.RLock()
	calls = mock.calls.GetGeographyDimensions
	mock.lockGetGeographyDimensions.RUnlock()
	return calls
}

// ListDatasets calls ListDatasetsFunc.
func (mock *CantabularClientMock) ListDatasets(ctx context.Context) ([]string, error) {
	if mock.ListDatasetsFunc == nil {
		panic("CantabularClientMock.ListDatasetsFunc: method is nil but CantabularClient.ListDatasets was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockListDatasets.Lock()
	mock.calls.ListDatasets = append(mock.calls.ListDatasets, callInfo)
	mock.lockListDatasets.Unlock()
	return mock.ListDatasetsFunc(ctx)
}

// ListDatasetsCalls gets all the calls that were made to ListDatasets.
// Check the length with:
//     len(mockedCantabularClient.ListDatasetsCalls())
func (mock *CantabularClientMock) ListDatasetsCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockListDatasets.RLock()
	calls = mock.calls.ListDatasets
	mock.lockListDatasets.RUnlock()
	return calls
}

// StatusCode calls StatusCodeFunc.
func (mock *CantabularClientMock) StatusCode(err error) int {
	if mock.StatusCodeFunc == nil {
		panic("CantabularClientMock.StatusCodeFunc: method is nil but CantabularClient.StatusCode was just called")
	}
	callInfo := struct {
		Err error
	}{
		Err: err,
	}
	mock.lockStatusCode.Lock()
	mock.calls.StatusCode = append(mock.calls.StatusCode, callInfo)
	mock.lockStatusCode.Unlock()
	return mock.StatusCodeFunc(err)
}

// StatusCodeCalls gets all the calls that were made to StatusCode.
// Check the length with:
//     len(mockedCantabularClient.StatusCodeCalls())
func (mock *CantabularClientMock) StatusCodeCalls() []struct {
	Err error
} {
	var calls []struct {
		Err error
	}
	mock.lockStatusCode.RLock()
	calls = mock.calls.StatusCode
	mock.lockStatusCode.RUnlock()
	return calls
}