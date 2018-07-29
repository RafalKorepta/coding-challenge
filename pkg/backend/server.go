// Copyright [2018] [Rafał Korepta]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package backend

// TODO consider https://github.com/Stoakes/grpc-gateway-example
// TODO http2 testing can be done using this link
// TODO http://big-elephants.com/2017-09/this-programmer-tried-to-mock-an-http-slash-2-server-in-go-and-heres-what-happened/
// TODO Custom error handler for grpc-gateway https://mycodesmells.com/post/grpc-gateway-error-handler
// TODO How to Create a CSR and Key File for a SAN Certificate with Multiple Subject Alternate Names
// TODO https://support.citrix.com/article/CTX227983 AND https://gist.github.com/croxton/ebfb5f3ac143cd86542788f972434c96
// TODO Tip and tricks https://medium.com/@povilasve/go-advanced-tips-tricks-a872503ac859
// TODO

import (
	"context"

	"io"

	"fmt"

	"time"

	"mime"
	"net/http"
	"strings"

	"crypto/tls"

	"os"

	"crypto/x509"

	"net"

	pb "github.com/RafalKorepta/coding-challenge/pkg/api/email/v1alpha1"
	"github.com/RafalKorepta/coding-challenge/pkg/certs"
	"github.com/RafalKorepta/coding-challenge/pkg/log"
	"github.com/RafalKorepta/coding-challenge/pkg/services"
	"github.com/RafalKorepta/coding-challenge/pkg/ui/data/swagger"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/philips/go-bindata-assetfs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"github.com/veqryn/h2c"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	opts     *options
	listener net.Listener
	pb.EmailServiceServer
}

// NewServer constructor of Server
func NewServer(listener net.Listener, opts ...Option) Server {
	o := evaluateOptions(opts)
	return Server{
		opts:     o,
		listener: listener,
	}
}

// Serve will start gRPC and REST server on the same port with or without TLS
func (s *Server) Serve() error {
	closer, err := initializeGlobalTracer(zap.L(), zap.S())
	if err != nil {
		return err
	}
	defer closer.Close()

	srv, grpcServer, err := createHTTPServer(s.listener.Addr().String(),
		WithCertFile(s.opts.certFile),
		WithKeyFile(s.opts.keyFile),
		WithSecure(s.opts.secure))
	if err != nil {
		return err
	}

	grpc_prometheus.Register(grpcServer)

	defer s.listener.Close()
	if s.opts.secure {
		return srv.ServeTLS(s.listener, "", "") // The certificates are initialized already
	}

	return srv.Serve(s.listener)
}

// initializeGlobalTracer will set global tracer using jeager tracer
func initializeGlobalTracer(logger *zap.Logger, sugar *zap.SugaredLogger) (io.Closer, error) {
	zapWrapper := log.ZapWrapper{
		Logger: logger,
		Sugar:  sugar,
	}

	metricsFactory := prometheus.New()

	tracer, closer, err := config.Configuration{
		ServiceName: "portal-backend",
	}.NewTracer(
		config.Metrics(metricsFactory),
		config.Logger(zapWrapper),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to start tracer: %v", err)
	}
	opentracing.SetGlobalTracer(tracer)
	return closer, nil
}

func registerEmailService(serverOpts ...grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(serverOpts...)

	pb.RegisterEmailServiceServer(grpcServer, services.EmailService{})

	return grpcServer
}

func createGRPCOptions(addr string, secure bool, certFile string) ([]grpc.ServerOption, error) {
	var opts []grpc.ServerOption

	grpc_zap.ReplaceGrpcLogger(zap.L())

	optsCtx := []grpc_ctxtags.Option{
		grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
	}

	optZap := []grpc_zap.Option{
		// Add filed to logs that comes from gRPC middleware
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_ns", duration.Nanoseconds())
		}),
	}

	opts = append(opts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_ctxtags.StreamServerInterceptor(optsCtx...),
		grpc_opentracing.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
		grpc_zap.StreamServerInterceptor(zap.L(), optZap...),
		grpc_recovery.StreamServerInterceptor(),
	)))

	opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_ctxtags.UnaryServerInterceptor(optsCtx...),
		grpc_opentracing.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_zap.UnaryServerInterceptor(zap.L(), optZap...),
		grpc_recovery.UnaryServerInterceptor(),
	)))

	if secure {
		certPool, err := createPool(certFile)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.Creds(credentials.NewClientTLSFromCert(certPool, addr)))
	}
	return opts, nil
}

func registerServerMux(addr string, dialOpts ...grpc.DialOption) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		var n int64
		n, err := io.Copy(w, strings.NewReader(pb.Swagger))
		if err != nil {
			zap.L().Error("Coping operation failed", zap.Int64("wrriten", n), zap.Error(err))
			http.Error(w, "swagger.json is currently unavailable", http.StatusInternalServerError)
		}
	})

	gwmux := runtime.NewServeMux()
	ctx := context.Background()
	err := pb.RegisterEmailServiceHandlerFromEndpoint(ctx, gwmux, addr, dialOpts)
	if err != nil {
		return nil, fmt.Errorf("unable to register gRPC gateway: %v", err)
	}

	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/", gwmux)
	serveSwagger(mux)

	return mux, nil
}

func createDialOpts(addr, certFile string, secure bool) ([]grpc.DialOption, error) {
	if secure {
		certPool, err := createPool(certFile)
		if err != nil {
			return nil, err
		}
		dcreds := credentials.NewTLS(&tls.Config{
			ServerName: addr, // Only connection from localhost will be accepted until certificate will have Subject Alternative Name init
			RootCAs:    certPool,
		})
		return []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}, nil
	}
	return []grpc.DialOption{grpc.WithInsecure()}, nil
}

func createPool(certFile string) (*x509.CertPool, error) {
	f, err := os.Open(certFile)
	if err != nil {
		zap.L().Error("Unable to open cert file", zap.Error(err))
	}
	certPool, err := certs.CreateX509Pool(f)
	if err != nil {
		return nil, fmt.Errorf("unable to create x509 cert pool: %v", err)
	}
	return certPool, nil
}

func createServerMainHandler(secure bool, grpcServer, mux http.Handler) http.Handler {
	if secure {
		return grpcHandlerFunc(grpcServer, mux)
	}
	// Wrap the Router
	return &h2c.HandlerH2C{
		Handler:  grpcHandlerFunc(grpcServer, mux),
		H2Server: &http2.Server{},
	}
}

func createHTTPServer(addr string, opts ...Option) (*http.Server, *grpc.Server, error) {
	o := evaluateOptions(opts)

	serverOpts, err := createGRPCOptions(addr, o.secure, o.certFile)
	if err != nil {
		return nil, nil, err
	}
	grpcServer := registerEmailService(serverOpts...)

	dialOpts, err := createDialOpts(addr, o.certFile, o.secure)
	if err != nil {
		return nil, nil, err
	}
	mux, err := registerServerMux(addr, dialOpts...)
	if err != nil {
		return nil, nil, err
	}

	rootHandler := createServerMainHandler(o.secure, grpcServer, mux)

	tlsCfg, err := createTLSConfig(o.secure, o.certFile, o.keyFile)
	if err != nil {
		return nil, nil, err
	}
	return &http.Server{
		Addr:      addr,
		Handler:   rootHandler,
		TLSConfig: tlsCfg,
	}, grpcServer, nil
}

func createTLSConfig(secure bool, certFile, keyFile string) (*tls.Config, error) {
	if secure {
		keyPair, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("unable to create x509 key pair certificate: %v", err)
		}

		return &tls.Config{
			Certificates: []tls.Certificate{keyPair},
			NextProtos:   []string{"h2"},
		}, nil
	}
	return nil, nil
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

// serveSwagger will register `/swagger-ui` endpoint into root mux.
// This will provide visual representation of gRPC contract
// The swagger-ui is auto generated by script located in `hack/build-ui.sh`
func serveSwagger(mux *http.ServeMux) {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		zap.L().Error("Unable to add extension type", zap.Error(err))
	}

	// Expose files in third_party/swagger-ui/ on <host>/swagger-ui
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}
