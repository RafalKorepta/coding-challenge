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
package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"mime"
	"os"
	"path/filepath"

	pb "github.com/RafalKorepta/coding-challenge/pkg/api/email/v1alpha1"
	"github.com/RafalKorepta/coding-challenge/pkg/backend"
	"github.com/RafalKorepta/coding-challenge/pkg/certs"
	"github.com/RafalKorepta/coding-challenge/pkg/ui/data/swagger"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/philips/go-bindata-assetfs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	portNumberFlag   = "port_number"
	certPathFlag     = "certsPath"
	certFileNameFlag = "certFileName"
	keyFileNameFlag  = "keyFileName"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launches the Email backend service",
	Run:   serve,
}

func init() {
	RootCmd.AddCommand(serveCmd)

	serveCmd.Flags().Int64P(portNumberFlag, "p", 9091,
		"the port on which the server will be listen on incoming requests")
	serveCmd.Flags().String(certPathFlag, "pkg/certs/local_certs", "the path where key and certificate are located")
	serveCmd.Flags().String(certFileNameFlag, "server.pem", "the path where key and certificate are located")
	serveCmd.Flags().String(keyFileNameFlag, "server.key", "the path where key and certificate are located")
	if err := viper.BindPFlags(serveCmd.Flags()); err != nil {
		logger.Error("Unable to bind flags")
	}
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

func serveSwagger(mux *http.ServeMux) {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		logger.Error("Unable to add extension type", zap.Error(err))
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

func serve(cmd *cobra.Command, args []string) {
	logger.Info("Initialize backend server", zap.Int64(portNumberFlag, viper.GetInt64(portNumberFlag)))

	addr := fmt.Sprintf("localhost:%d", viper.GetInt64(portNumberFlag))

	keyPair, err := tls.LoadX509KeyPair(
		filepath.Join(viper.GetString(certPathFlag), viper.GetString(certFileNameFlag)),
		filepath.Join(viper.GetString(certPathFlag), viper.GetString(keyFileNameFlag)))
	if err != nil {
		logger.Error("Unable to create x509 key pair certificate", zap.Error(err))
	}

	f, err := os.Open(filepath.Join(viper.GetString(certPathFlag), viper.GetString(certFileNameFlag)))
	if err != nil {
		logger.Error("Unable to open cert file", zap.Error(err))
	}
	certPool, err := certs.CreateX509Pool(f)
	if err != nil {
		logger.Error("Unable to create x509 cert pool", zap.Error(err))
	}

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewClientTLSFromCert(certPool, addr))}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterEmailServiceServer(grpcServer, backend.NewServer(32))
	ctx := context.Background()

	dcreds := credentials.NewTLS(&tls.Config{
		ServerName: addr, // Only connection from localhost will be accepted
		RootCAs:    certPool,
	})
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

	mux := http.NewServeMux()
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		var n int64
		n, err = io.Copy(w, strings.NewReader(pb.Swagger))
		if err != nil {
			logger.Error("Coping operation failed", zap.Int64("wrriten", n), zap.Error(err))
		}
	})

	gwmux := runtime.NewServeMux()
	err = pb.RegisterEmailServiceHandlerFromEndpoint(ctx, gwmux, addr, dopts)
	if err != nil {
		logger.Error("Unable to register gRPC gateway", zap.Error(err))
		return
	}

	mux.Handle("/", gwmux)
	serveSwagger(mux)

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.GetInt64(portNumberFlag)))
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: grpcHandlerFunc(grpcServer, mux),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{keyPair},
			NextProtos:   []string{"h2"},
		},
	}

	logger.Info("grpc server lunch", zap.Int64("port", viper.GetInt64(portNumberFlag)))
	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))

	if err != nil {
		logger.Fatal("ListenAndServe: ", zap.Error(err))
	}
}
