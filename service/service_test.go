package service_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/service"
	"github.com/ONSdigital/dp-population-types-api/service/mock"
)

var (
	ctx           = context.Background()
	testBuildTime = "12"
	testGitCommit = "GitCommit"
	testVersion   = "Version"
)

var errHealthcheck = errors.New("could not get healthcheck")

func TestInit(t *testing.T) {

	Convey("Having a set of mocked dependencies", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)

		cfgWithCantabularHealthcheckEnabled := *cfg
		cfgWithCantabularHealthcheckEnabled.CantabularHealthcheckEnabled = true

		initialiserMock := buildInitialiserMockWithNilDependencies()

		cantabularClientMock := &mock.CantabularClientMock{
			CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
		}

		datasetAPIClientMock := &mock.DatasetAPIClientMock{
			CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
		}

		initialiserMock.GetCantabularClientFunc = func(cfg config.CantabularConfig) service.CantabularClient { return cantabularClientMock }

		initialiserMock.GetDatasetAPIClientFunc = func(cfg *config.Config) service.DatasetAPIClient { return datasetAPIClientMock }

		serverMock := &mock.HTTPServerMock{
			ListenAndServeFunc: func() error { return nil },
			ShutdownFunc:       func(ctx context.Context) error { return nil },
		}
		initialiserMock.GetHTTPServerFunc = func(bindAddr string, router http.Handler) service.HTTPServer { return serverMock }

		svc := &service.Service{}

		Convey("Given that initialising healthcheck returns an error", func() {

			initialiserMock.GetHealthCheckFunc = func(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
				return nil, errHealthcheck
			}

			svc := service.New()
			err := svc.Init(ctx, &initialiserMock, cfg, testBuildTime, testGitCommit, testVersion)

			Convey("Then service Init fails with an error", func() {
				So(errors.Is(err, errHealthcheck), ShouldBeTrue)
			})
		})

		Convey("Given cantabular health check is enabled", func() {

			Convey("When the service is initialised", func() {

				hcMock := &mock.HealthCheckerMock{
					AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
				}

				initialiserMock.GetHealthCheckFunc = func(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
					return hcMock, nil
				}
				err := svc.Init(ctx, &initialiserMock, &cfgWithCantabularHealthcheckEnabled, testBuildTime, testGitCommit, testVersion)
				So(err, ShouldBeNil)

				Convey("Then the cantabular healthcheck should be added", func() {
					cantabularCall := findCantabularClientCheck(hcMock)
					So(cantabularCall.Checker, ShouldNotBeNil)

					checkState := healthcheck.CheckState{}
					err := cantabularCall.Checker(ctx, &checkState)
					So(err, ShouldBeNil)

					checkerCalls := cantabularClientMock.CheckerCalls()
					So(checkerCalls[0].State, ShouldPointTo, &checkState)
				})
			})
		})

		Convey("Given that the cantabular client health check fails to initialise", func() {

			hcMock := &mock.HealthCheckerMock{
				AddCheckFunc: func(name string, checker healthcheck.Checker) error {
					if name == "Cantabular client" {
						return errors.New("oops")
					}
					return nil
				},
			}

			initialiserMock.GetHealthCheckFunc = func(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
				return hcMock, nil
			}

			Convey("When the service is initialised", func() {

				err := svc.Init(ctx, &initialiserMock, &cfgWithCantabularHealthcheckEnabled, testBuildTime, testGitCommit, testVersion)
				Convey("Then the cantabular healthcheck error should be included in the returned errors", func() {

					So(strings.Contains(err.Error(), "cantabular client"), ShouldBeTrue)
				})
			})
		})

		Convey("Given that all dependencies are successfully initialised", func() {

			hcMock := &mock.HealthCheckerMock{
				AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			}

			initialiserMock.GetHealthCheckFunc = func(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
				return hcMock, nil
			}

			err := svc.Init(ctx, &initialiserMock, cfg, testBuildTime, testGitCommit, testVersion)

			Convey("Then service Init succeeds", func() {

				So(err, ShouldBeNil)
			})

			Convey("Then the cantabular healthcheck should not be added (as the flag is not set)", func() {

				cantabularCall := findCantabularClientCheck(hcMock)
				So(cantabularCall, ShouldBeNil)
			})
		})
	})
}

func findCantabularClientCheck(hcMock *mock.HealthCheckerMock) *struct {
	Name    string
	Checker healthcheck.Checker
} {
	addCheckCalls := hcMock.AddCheckCalls()
	for _, call := range addCheckCalls {
		if call.Name == "Cantabular client" {
			return &call
		}
	}
	return nil
}

func TestClose(t *testing.T) {

	Convey("Having a correctly initialised service", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)

		hcStopped := false
		initialiserMock := buildInitialiserMockWithNilDependencies()

		// healthcheck Stop does not depend on any other service being closed/stopped
		hcMock := &mock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
			StopFunc:     func() { hcStopped = true },
		}

		datasetAPIClientMock := &mock.DatasetAPIClientMock{
			CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
		}

		initialiserMock.GetHealthCheckFunc = func(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
			return hcMock, nil
		}

		initialiserMock.GetDatasetAPIClientFunc = func(cfg *config.Config) service.DatasetAPIClient { return datasetAPIClientMock }

		// server Shutdown will fail if healthcheck is not stopped
		serverMock := &mock.HTTPServerMock{
			ListenAndServeFunc: func() error { return nil },
			ShutdownFunc: func(ctx context.Context) error {
				if !hcStopped {
					return errors.New("server stopped before healthcheck")
				}
				return nil
			},
		}

		initialiserMock.GetHTTPServerFunc = func(bindAddr string, router http.Handler) service.HTTPServer {
			return serverMock
		}

		Convey("Closing the service results in all the dependencies being closed in the expected order", func() {
			svcErrors := make(chan error, 1)
			svc := service.New()
			err := svc.Init(ctx, &initialiserMock, cfg, testBuildTime, testGitCommit, testVersion)
			So(err, ShouldBeNil)

			svc.Start(context.Background(), svcErrors)

			err = svc.Close(context.Background())
			So(err, ShouldBeNil)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(serverMock.ShutdownCalls()), ShouldEqual, 1)
		})

		Convey("If services fail to stop, the Close operation tries to close all dependencies and returns an error", func() {
			failingServiceMock := &mock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(ctx context.Context) error {
					return errors.New("failed to stop http server")
				},
			}

			initialiserMock.GetHTTPServerFunc = func(bindAddr string, router http.Handler) service.HTTPServer {
				return failingServiceMock
			}

			svcErrors := make(chan error, 1)
			svc := service.New()
			err := svc.Init(ctx, &initialiserMock, cfg, testBuildTime, testGitCommit, testVersion)
			So(err, ShouldBeNil)

			svc.Start(context.Background(), svcErrors)

			err = svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(failingServiceMock.ShutdownCalls()), ShouldEqual, 1)
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

			svc := service.Service{
				Config:      cfg,
				Server:      timeoutServerMock,
				HealthCheck: hcMock,
			}

			err = svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(errors.Is(err, context.DeadlineExceeded), ShouldBeTrue)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(timeoutServerMock.ShutdownCalls()), ShouldEqual, 1)
		})
	})
}

func buildInitialiserMockWithNilDependencies() mock.InitialiserMock {
	return mock.InitialiserMock{
		GetCantabularClientFunc: func(cfg config.CantabularConfig) service.CantabularClient {
			return nil
		},
		GetResponderFunc: func() service.Responder {
			return nil
		},
		GetHealthCheckFunc: func(cfg *config.Config, time string, commit string, version string) (service.HealthChecker, error) {
			return nil, nil
		},
		GetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer {
			return nil
		},
	}
}
