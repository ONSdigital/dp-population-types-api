package service_test

import (
	"context"
	"errors"
	"net/http"
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

		initialiserMock := buildInitialiserMockWithNilDependencies()

		hcMock := &mock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
			StopFunc:     func() {},
		}

		initialiserMock.GetHealthCheckFunc = func(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
			return hcMock, nil
		}

		serverMock := &mock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				return nil
			},
			ShutdownFunc: func(ctx context.Context) error {
				return nil
			},
		}
		initialiserMock.GetHTTPServerFunc = func(bindAddr string, router http.Handler) service.HTTPServer {
			return serverMock
		}

		svc := &service.Service{}

		Convey("Given that initialising healthcheck returns an error", func() {
			initialiserMock.GetHealthCheckFunc = func(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
				return nil, errHealthcheck
			}
			// setup (run before each `Convey` at this scope / indentation):
			svc := service.New()
			err := svc.Init(ctx, &initialiserMock, cfg, testBuildTime, testGitCommit, testVersion)

			Convey("Then service Init fails with an error", func() {
				So(errors.Is(err, errHealthcheck), ShouldBeTrue)
			})

			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})

		Convey("Given that all dependencies are successfully initialised", func() {

			// setup (run before each `Convey` at this scope / indentation):
			err := svc.Init(ctx, &initialiserMock, cfg, testBuildTime, testGitCommit, testVersion)

			Convey("Then service Init succeeds", func() {
				So(err, ShouldBeNil)
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
		initialiserMock := buildInitialiserMockWithNilDependencies()

		// healthcheck Stop does not depend on any other service being closed/stopped
		hcMock := &mock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
			StopFunc:     func() { hcStopped = true },
		}

		initialiserMock.GetHealthCheckFunc = func(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
			return hcMock, nil
		}

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
			failingserverMock := &mock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(ctx context.Context) error {
					return errors.New("failed to stop http server")
				},
			}

			initialiserMock.GetHTTPServerFunc = func(bindAddr string, router http.Handler) service.HTTPServer {
				return failingserverMock
			}

			svcErrors := make(chan error, 1)
			svc := service.New()
			err := svc.Init(ctx, &initialiserMock, cfg, testBuildTime, testGitCommit, testVersion)
			So(err, ShouldBeNil)

			svc.Start(context.Background(), svcErrors)

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
