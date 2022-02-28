//nolint:typecheck
package service_test

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"

	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/service"
	"github.com/ONSdigital/dp-population-types-api/service/mock"
	serviceMock "github.com/ONSdigital/dp-population-types-api/service/mock"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx           = context.Background()
	testBuildTime = "BuildTime"
	testGitCommit = "GitCommit"
	testVersion   = "Version"
	errServer     = errors.New("HTTP Server error")
)

var (
	errHealthcheck = errors.New("healthCheck error")
)

var funcDoGetHealthcheckErr = func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	return nil, errHealthcheck
}

func TestRun(t *testing.T) {

	Convey("Having a set of mocked dependencies", t, func() {

		cfg, err := config.Get()
		So(err, ShouldBeNil)

		hcMock := &serviceMock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
		}

		serverWg := &sync.WaitGroup{}
		serverMock := &serviceMock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				serverWg.Done()
				return nil
			},
		}

		failingServerMock := &serviceMock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				serverWg.Done()
				return errServer
			},
		}

		cantabularClientMock := &serviceMock.CantabularClientMock{}

		funcDoGetHealthcheckOk := func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
			return hcMock, nil
		}

		funcDoGetHTTPServer := func(bindAddr string, router http.Handler) service.HTTPServer {
			return serverMock
		}

		funcDoGetFailingHTTPServer := func(bindAddr string, router http.Handler) service.HTTPServer {
			return failingServerMock
		}

		funcDoGetCantabularClientOk := func(ctx context.Context, cfg config.CantabularConfig) service.CantabularClient {
			return cantabularClientMock
		}

		Convey("Given that initialising healthcheck returns an error", func() {

			// setup (run before each `Convey` at this scope / indentation):
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc:       funcDoGetHTTPServer,
				DoGetHealthCheckFunc:      funcDoGetHealthcheckErr,
				DoGetCantabularClientFunc: funcDoGetCantabularClientOk,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails with the same error and the flag is not set", func() {
				So(err, ShouldResemble, errHealthcheck)
				So(svcList.HealthCheck, ShouldBeFalse)
			})

			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})

		Convey("When the service is initialised", func() {

			var actualConfig config.CantabularConfig
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc:  funcDoGetHTTPServer,
				DoGetHealthCheckFunc: funcDoGetHealthcheckOk,
				DoGetCantabularClientFunc: func(ctx context.Context, cfg config.CantabularConfig) service.CantabularClient {
					actualConfig = cfg
					return cantabularClientMock
				},
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			serverWg.Add(1)
			cfg.CantabularURL = "https://test.com/cantabularURL"
			cfg.CantabularExtURL = "https://test.com/cantabularExtURL"
			cfg.DefaultRequestTimeout = 42 * time.Hour

			svc, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			SoMsg("Then the cantabular client should be initialised using the expected configuration",
				actualConfig, ShouldResemble, config.CantabularConfig{
					CantabularURL:         "https://test.com/cantabularURL",
					CantabularExtURL:      "https://test.com/cantabularExtURL",
					DefaultRequestTimeout: 42 * time.Hour,
				})

			SoMsg("Then the service should have the created cantabular client",
				svc.CantabularClient, ShouldEqual, cantabularClientMock)
		})

		Convey("Given that all dependencies are successfully initialised", func() {

			// setup (run before each `Convey` at this scope / indentation):
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc:       funcDoGetHTTPServer,
				DoGetHealthCheckFunc:      funcDoGetHealthcheckOk,
				DoGetCantabularClientFunc: funcDoGetCantabularClientOk,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			serverWg.Add(1)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run succeeds and all the flags are set", func() {
				So(err, ShouldBeNil)
				So(svcList.HealthCheck, ShouldBeTrue)
			})

			Convey("The checkers are registered and the healthcheck and http server started", func() {
				So(len(hcMock.AddCheckCalls()), ShouldEqual, 1)
				So(len(initMock.DoGetHTTPServerCalls()), ShouldEqual, 1)
				So(initMock.DoGetHTTPServerCalls()[0].BindAddr, ShouldEqual, "localhost:12900")
				So(len(hcMock.StartCalls()), ShouldEqual, 1)
				//!!! a call needed to stop the server, maybe ?
				serverWg.Wait() // Wait for HTTP server go-routine to finish
				So(len(serverMock.ListenAndServeCalls()), ShouldEqual, 1)
			})

			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})

		Convey("Given that Checkers cannot be registered", func() {

			// setup (run before each `Convey` at this scope / indentation):
			errAddCheckFail := errors.New("Error(s) registering checkers for healthcheck")
			hcMockAddFail := &serviceMock.HealthCheckerMock{
				AddCheckFunc: func(name string, checker healthcheck.Checker) error { return errAddCheckFail },
				StartFunc:    func(ctx context.Context) {},
			}
			errCantabularClientCheckerFail := errors.New("Error registering checker for Cantabular client")
			cantabularClientMockCheckerFail := &serviceMock.CantabularClientMock{
				CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error { return errCantabularClientCheckerFail },
			}
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc: funcDoGetHTTPServer,
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMockAddFail, nil
				},
				DoGetCantabularClientFunc: func(ctx context.Context, cfg config.CantabularConfig) service.CantabularClient {
					return cantabularClientMockCheckerFail
				},
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			serverWg.Add(1)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails, but all checks try to register", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldResemble, fmt.Sprintf("unable to register checkers: %s", errAddCheckFail.Error()))
				So(svcList.HealthCheck, ShouldBeTrue)
				// ADD CODE: add code to confirm checkers exist
				So(len(hcMockAddFail.AddCheckCalls()), ShouldEqual, 1) // ADD CODE: change the '0' to the number of checkers you have registered
			})
			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})

		Convey("Given that all dependencies are successfully initialised but the http server fails", func() {

			// setup (run before each `Convey` at this scope / indentation):
			initMock := &serviceMock.InitialiserMock{
				DoGetHealthCheckFunc:      funcDoGetHealthcheckOk,
				DoGetHTTPServerFunc:       funcDoGetFailingHTTPServer,
				DoGetCantabularClientFunc: funcDoGetCantabularClientOk,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			serverWg.Add(1)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			Convey("Then the error is returned in the error channel", func() {
				sErr := <-svcErrors
				So(sErr.Error(), ShouldResemble, fmt.Sprintf("failure in http listen and serve: %s", errServer.Error()))
				So(len(failingServerMock.ListenAndServeCalls()), ShouldEqual, 1)
			})

			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})
	})
}

func TestClose(t *testing.T) {

	Convey("Having a correctly initialised service", t, func() {

		cfg, err := config.Get()
		So(err, ShouldBeNil)

		hcStopped := false

		// healthcheck Stop does not depend on any other service being closed/stopped
		hcMock := &serviceMock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
			StopFunc:     func() { hcStopped = true },
		}

		// server Shutdown will fail if healthcheck is not stopped
		serverMock := &mock.HTTPServerMock{
			ListenAndServeFunc: func() error { return nil },
			ShutdownFunc: func(ctx context.Context) error {
				if !hcStopped {
					return errors.New("Server stopped before healthcheck")
				}
				return nil
			},
		}

		Convey("Closing the service results in all the dependencies being closed in the expected order", func() {

			initMock := &mock.InitialiserMock{
				DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer { return serverMock },
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMock, nil
				},
				DoGetCantabularClientFunc: func(ctx context.Context, cfg config.CantabularConfig) service.CantabularClient {
					return nil
				},
			}

			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			svc, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			err = svc.Close(context.Background())
			So(err, ShouldBeNil)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(serverMock.ShutdownCalls()), ShouldEqual, 1)
		})

		Convey("If services fail to stop, the Close operation tries to close all dependencies and returns an error", func() {

			failingserverMock := &mock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(ctx context.Context) error {
					return errors.New("Failed to stop http server")
				},
			}

			initMock := &mock.InitialiserMock{
				DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer { return failingserverMock },
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMock, nil
				},
				DoGetCantabularClientFunc: func(ctx context.Context, cfg config.CantabularConfig) service.CantabularClient {
					return nil
				},
			}

			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			svc, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			err = svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(failingserverMock.ShutdownCalls()), ShouldEqual, 1)
		})

		Convey("If service times out while shutting down, the Close operation fails with the expected error", func() {
			cfg.GracefulShutdownTimeout = 1 * time.Millisecond
			timeoutServerMock := &mock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(ctx context.Context) error {
					time.Sleep(2 * time.Millisecond)
					return nil
				},
			}

			svcList := service.NewServiceList(nil)
			svcList.HealthCheck = true
			svc := service.Service{
				Config:      cfg,
				ServiceList: svcList,
				Server:      timeoutServerMock,
				HealthCheck: hcMock,
			}

			err = svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldResemble, "context deadline exceeded")
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(timeoutServerMock.ShutdownCalls()), ShouldEqual, 1)
		})
	})
}
