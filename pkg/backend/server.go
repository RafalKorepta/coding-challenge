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

import (
	"context"

	pb "github.com/RafalKorepta/coding-challenge/pkg/api/email/v1alpha1"
)

type EmailService struct {
	pb.EmailServiceServer
}

func (es EmailService) SendMail(ctx context.Context, req *pb.EmailRequest) (*pb.EmailResponse, error) {
	return &pb.EmailResponse{
		Error: req.Message,
	}, nil
}

func NewServer(portNumber int64) pb.EmailServiceServer {
	return EmailService{}
}

//func Server(portNumber int64, certDir, certFileName, keyFileName string) {
//	// Load KeyPair
//	pair, err := tls.LoadX509KeyPair(filepath.Join(certDir, certFileName), filepath.Join(certDir, keyFileName))
//	if err != nil {
//		log.Fatal("")
//	}
//
//	// Create grpc server
//
//	// Create Mux
//
//	// Register swagger json endpoint in mux
//
//	// Register gRPC gateway
//
//	// Register swagger ui into mux
//
//	// new https server
//
//	opts := []grpc.ServerOption{
//		grpc.Creds(credentials.NewClientTLSFromCert(demoCertPool, "localhost:10000"))}
//
//	grpcServer := grpc.NewServer(opts...)
//	pb.RegisterEmailServiceServer(grpcServer, newServer())
//	ctx := context.Background()
//
//	dcreds := credentials.NewTLS(&tls.Config{
//		ServerName: demoAddr,
//		RootCAs:    demoCertPool,
//	})
//	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
//		io.Copy(w, strings.NewReader(pb.Swagger))
//	})
//
//	gwmux := runtime.NewServeMux()
//	err := pb.RegisterEmailServiceHandlerFromEndpoint(ctx, gwmux, demoAddr, dopts)
//	if err != nil {
//		fmt.Printf("serve: %v\n", err)
//		return
//	}
//
//	mux.Handle("/", gwmux)
//	serveSwagger(mux)
//
//	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
//	if err != nil {
//		panic(err)
//	}
//
//	srv := &http.Server{
//		Addr:    demoAddr,
//		Handler: grpcHandlerFunc(grpcServer, mux),
//		TLSConfig: &tls.Config{
//			Certificates: []tls.Certificate{*demoKeyPair},
//			NextProtos:   []string{"h2"},
//		},
//	}
//
//	fmt.Printf("grpc on port: %d\n", portNumber)
//	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))
//
//	if err != nil {
//		log.Fatal("ListenAndServe: ", err)
//	}
//
//	return
//}
