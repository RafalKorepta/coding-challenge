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
	"net/http"

	"net/http/httptest"

	"mime"

	"io/ioutil"

	"io"

	"bytes"

	"context"

	"github.com/RafalKorepta/coding-challenge/pkg/api/email/v1alpha1"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/prometheus/util/promlint"
	"google.golang.org/grpc"
)

const emailURI = "/v1alpha1/email"

var _ = Describe("Server that register REST and gRPC endpoint on the same port", func() {
	var (
		srv             *http.Server
		err             error
		grpcServer      *grpc.Server
		response        *http.Response
		body            []byte
		newMockServer   *httptest.Server
		clientHTTP      *http.Client
		swaggerUIConent []byte
		requestedURI    string
		postBody        io.Reader
		marshaller      runtime.JSONPb
		clientRPC       email.EmailServiceClient
		responseGRPC    *email.EmailResponse
		conn            *grpc.ClientConn
	)

	JustBeforeEach(func() {
		srv, grpcServer, err = createHTTPServer(listener.Addr().String(),
			opts...)
	})

	JustBeforeHTTPAndGrpcContext := func() {
		if requestedURI != "" {
			clientHTTP = newMockServer.Client()
			if postBody != nil {
				response, err = clientHTTP.Post(newMockServer.URL+requestedURI, marshaller.ContentType(), postBody)
			} else {
				response, err = clientHTTP.Get(newMockServer.URL + requestedURI)
			}
			Expect(err).NotTo(HaveOccurred())

			body, err = ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())
		}
	}

	http1 := func() {
		Context("when swagger-ui URI is called on raw handler", func() {
			var (
				req *http.Request
				w   *httptest.ResponseRecorder
			)
			JustBeforeEach(func() {
				req = httptest.NewRequest("GET", "http://localhost:9091/swagger-ui", nil)
				w = httptest.NewRecorder()
				srv.Handler.ServeHTTP(w, req)

				response = w.Result()
				body, err = ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should response with moved permanently response", func() {
				Expect(response.StatusCode).To(Equal(http.StatusMovedPermanently))
				Expect(response.Header.Get("Content-Type")).To(Equal(mime.TypeByExtension(".html")))
				Expect(string(body)).To(ContainSubstring("Moved Permanently"))
			})
		})

		Context("when swagger-ui URI is called by client that handle redirect", func() {
			BeforeEach(func() {
				requestedURI = "/swagger-ui"
			})

			JustBeforeEach(func() {
				swaggerUIConent, err = ioutil.ReadFile("test_data/swagger-ui-main-page.txt")
				Expect(err).NotTo(HaveOccurred())
			})

			It("should response with main html page", func() {
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(response.Header.Get("Content-Type")).To(Equal(mime.TypeByExtension(".html")))
				Expect(body).To(Equal(swaggerUIConent))
			})
		})

		Context("when swagger.json URI is called", func() {
			BeforeEach(func() {
				requestedURI = "/swagger.json"
			})

			It("should return swagger-ui webpage", func() {
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(response.Header.Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
				Expect(string(body)).To(Equal(email.Swagger))
			})
		})

		Context("when /metrics URI is called", func() {
			var (
				problems []promlint.Problem
			)
			BeforeEach(func() {
				requestedURI = "/metrics"
			})

			JustBeforeEach(func() {
				problems, err = promlint.New(response.Body).Lint()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return valid prometheus response", func() {
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(response.Header.Get("Content-Type")).To(Equal("text/plain; version=0.0.4; charset=utf-8"))
				Expect(problems).To(BeEmpty())
			})
		})

		Context("when GET method on email URI is called", func() {
			BeforeEach(func() {
				requestedURI = emailURI
			})

			It("should return method not allowed", func() {
				Expect(response.StatusCode).To(Equal(http.StatusMethodNotAllowed))
				Expect(response.Header.Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
				Expect(string(body)).To(ContainSubstring(http.StatusText(http.StatusMethodNotAllowed)))
			})
		})

		Context("when POST method on email URI is called", func() {
			BeforeEach(func() {
				var marshaledProto []byte
				requestedURI = emailURI
				emailRequest := email.EmailRequest{
					Message: "Hello",
				}
				marshaller = runtime.JSONPb{}
				marshaledProto, err = marshaller.Marshal(emailRequest)
				Expect(err).NotTo(HaveOccurred())
				postBody = bytes.NewReader(marshaledProto)
			})

			It("should return correct response", func() {
				var r []byte
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(response.Header.Get("Content-Type")).To(Equal(marshaller.ContentType()))
				r, err = marshaller.Marshal(email.EmailResponse{
					Error: "Hello",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(body).To(Equal(r))
			})
		})
	}

	gRPC := func() {
		Context("when gRPC client call email service", func() {
			var (
				metadata runtime.ServerMetadata
			)
			BeforeEach(func() {
				var dialOpts []grpc.DialOption
				o := evaluateOptions(opts)
				dialOpts, err = createDialOpts("localhost", "test_data/server.pem", o.secure)
				Expect(err).NotTo(HaveOccurred())
				conn, err = grpc.Dial(newMockServer.Listener.Addr().String(), dialOpts...)
				clientRPC = email.NewEmailServiceClient(conn)
			})

			JustBeforeEach(func() {
				ctx := context.Background()
				responseGRPC, err = clientRPC.SendMail(ctx, &email.EmailRequest{
					Message: "Hello",
				}, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
			})

			AfterEach(func() {
				err = conn.Close()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return correct response", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(responseGRPC.Error).To(Equal("Hello"))
			})
		})
	}

	Describe("secure server", func() {
		BeforeEach(func() {
			opts = []Option{
				WithCertFile("test_data/server.pem"),
				WithKeyFile("test_data/server.key"),
				WithServerOverrideName("localhost"),
				WithSecure(true),
			}
		})

		JustBeforeEach(func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(srv).NotTo(BeNil())
			Expect(grpcServer).NotTo(BeNil())

			newMockServer = &httptest.Server{
				Listener: listener,
				TLS:      srv.TLSConfig,
				Config:   srv,
			}
			newMockServer.StartTLS()

			JustBeforeHTTPAndGrpcContext()
		})

		Describe("In combine test suite", func() {
			Describe("Http1.1", http1)

			Describe("gRPC", gRPC)
		})
	})

	Describe("Insecure server", func() {
		BeforeEach(func() {
			opts = []Option{
				WithSecure(false),
			}
		})

		JustBeforeEach(func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(srv).NotTo(BeNil())
			Expect(grpcServer).NotTo(BeNil())

			newMockServer = httptest.NewServer(srv.Handler)
			//newMockServer = &httptest.Server{
			//	Listener: listener,
			//	Config:   srv,
			//}
			//newMockServer.Start()

			JustBeforeHTTPAndGrpcContext()
		})

		//AfterEach(func() {
		//	newMockServer.CloseClientConnections()
		//	newMockServer.Close()
		//})

		Describe("In combine test suite", func() {
			Describe("Http1.1", http1)

			Describe("gRPC", gRPC)

		})
	})

	Describe("Not valid server key", func() {
		BeforeEach(func() {
			opts = []Option{
				WithCertFile("test_data/server.pem"),
				WithKeyFile("invalid.key"),
				WithSecure(true),
			}
		})

		It("should failed to initialize server", func() {
			Expect(srv).To(BeNil())
			Expect(grpcServer).To(BeNil())
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Not valid server certificate", func() {
		BeforeEach(func() {
			opts = []Option{
				WithCertFile("invalid.pem"),
				WithKeyFile("test_data/server.key"),
				WithSecure(true),
			}
		})

		It("should failed to initialize server", func() {
			Expect(srv).To(BeNil())
			Expect(grpcServer).To(BeNil())
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("pre-initialization of global tracer", func() {
		var (
			noopTracer opentracing.Tracer
			newTracer  opentracing.Tracer
			closer     io.Closer
		)
		BeforeEach(func() {
			noopTracer = opentracing.GlobalTracer()
			closer, err = initializeGlobalTracer(nil, nil)
			Expect(err).NotTo(HaveOccurred())

			newTracer = opentracing.GlobalTracer()
		})

		It("should global tracer be not the noop tracer", func() {
			Expect(newTracer).NotTo(Equal(noopTracer))
			Expect(closer).NotTo(BeNil())
		})
	})

	Describe("Server initialization", func() {
		var (
			newServer *Server
		)
		BeforeEach(func() {
			newServer = NewServer(nil)
		})
		It("Should created server", func() {
			Expect(newServer).NotTo(BeNil())
			Expect(newServer.listener).To(BeNil())
			Expect(newServer.opts).To(Equal(defaultOptions))
		})
	})
})

//func Test_initializeTracer(t *testing.T) {
//	// Arrange
//	noopTracer := opentracing.GlobalTracer()
//	t.Run("After initialize the global tracer is not NoopTracer", func(t *testing.T) {
//		// Act
//		closer, err := initializeGlobalTracer(nil, nil)
//
//		// Assert
//		tracer := opentracing.GlobalTracer()
//		assert.NoError(t, err, "Error should not occur")
//		assert.NotEqual(t, noopTracer, tracer, "Tracer must changed.")
//		assert.NotNil(t, closer, "Closer must exist.")
//	})
//}
//
//func Test_serveSwagger(t *testing.T) {
//	// Arrange
//	mux := http.NewServeMux()
//	serveSwagger(mux)
//	contentType := mime.TypeByExtension(".html")
//
//	t.Run("Mux has swagger-ui path registered", func(t *testing.T) {
//		// Arrange
//		req := httptest.NewRequest("GET", "http://localhost:9091/swagger-ui/", nil)
//		w := httptest.NewRecorder()
//
//		// Act
//		mux.ServeHTTP(w, req)
//
//		resp := w.Result()
//		body, _ := ioutil.ReadAll(resp.Body)
//
//		// Assert
//		assert.Equal(t, http.StatusOK, resp.StatusCode, "Response must be 200")
//		assert.Equal(t, contentType, resp.Header.Get("Content-Type"), "Content-Type must be text/html")
//		assert.NotEmpty(t, body, "Body must not be empty")
//
//	})
//	t.Run("Content moved permanently", func(t *testing.T) {
//		// Arrange
//		req := httptest.NewRequest("GET", "http://localhost:9091/swagger-ui", nil)
//		w := httptest.NewRecorder()
//
//		// Act
//		mux.ServeHTTP(w, req)
//
//		resp := w.Result()
//		body, _ := ioutil.ReadAll(resp.Body)
//
//		// Assert
//		assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode, "Response must be 200")
//		assert.Equal(t, contentType, resp.Header.Get("Content-Type"), "Content-Type must be text/html")
//		assert.NotEmpty(t, body, "Body must not be empty")
//		assert.Contains(t, string(body), "Moved Permanently", "Body must contains Moved Permanently")
//	})
//}

//func Test_createHTTPServer(t *testing.T) {
//	// Arrange
//	listener, err := net.Listen("tcp", "127.0.0.1:0")
//	if err != nil {
//		listener, err = net.Listen("tcp", "[::1]:0")
//		assert.NoError(t, err, "Listener must bind to address")
//	}
//	addr := listener.Addr().String()
//	t.Logf("Listener bind to %s", addr)
//
//	t.Run("Secure main handler", func(t *testing.T) {
//		// Arrange
//		//handler := createServerMainHandler(true,
//		//	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		//		grpc = true
//		//	}), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		//		http1 = true
//		//	}))
//
//		t.Run("Grpc request", func(t *testing.T) {
//			// Arrange
//			//req := httptest.NewRequest("GET", "http://localhost:9091/test", nil)
//			//w := httptest.NewRecorder()
//
//			// Act
//			//handler.ServeHTTP(w, req)
//
//			// Assert
//		})
//
//		t.Run("Http 1.1 request", func(t *testing.T) {
//			// Arrange
//			//req := httptest.NewRequest("GET", "http://localhost:9091/test", nil)
//			//w := httptest.NewRecorder()
//
//			// Act
//			//handler.ServeHTTP(w, req)
//
//			// Assert
//		})
//	})
//
//	t.Run("Insecure main handler", func(t *testing.T) {
//		// Arrange
//
//		srv, err := createHTTPServer(
//			WithCertFile(""),
//			WithKeyFile(""),
//			WithSecure(false),
//			WithListener(listener))
//		assert.Error(t, err, "HTTP server must be created without error")
//		assert.NotNil(t, srv, "HTTP server must be created")
//
//		t.Run("Grpc request", func(t *testing.T) {
//			// Arrange
//			//var ok bool
//			//req := httptest.NewRequest("GET", "http://localhost:9091/test", nil)
//			//w := httptest.NewRecorder()
//
//			// Act
//			//go func() { errc <- srv.Serve() }()
//			//select {
//			//case err := <-errc:
//			//	t.Logf("On try #%v: %v", try+1, err)
//			//
//			//case ln = <-lnc:
//			//	ok = true
//			//	t.Logf("Listening on %v", ln.Addr().String())
//			//	break
//			//}
//			//if !ok {
//			//	t.Fatalf("Failed to start up after %d tries", maxTries)
//			//}
//			//defer ln.Close()
//
//			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//				fmt.Fprintln(w, "Hello, client")
//			}))
//			defer ts.Close()
//
//			c, err := tls.DialWithDialer(&net.Dialer{
//				Timeout:  time.Minute,
//				Deadline: time.Now().Add(time.Minute),
//			}, "tcp", "localhost:9091", &tls.Config{
//				InsecureSkipVerify: true,
//				NextProtos:         []string{"h2", "http/1.1"},
//			})
//			if err != nil {
//				t.Fatal(err)
//			}
//			defer c.Close()
//
//			// Assert
//			assert.Equal(t, "h2", c.ConnectionState().NegotiatedProtocol, "Negotiated protocol must be h2")
//			assert.True(t, c.ConnectionState().NegotiatedProtocolIsMutual, "Negotiated protocol must be mutual")
//
//			//// Assert
//			//assert.True(t, grpc, "Http must be called")
//			//assert.False(t, http1, "Grpc must not be called")
//		})
//
//		t.Run("Http 1.1 request", func(t *testing.T) {
//			// Arrange
//			req := httptest.NewRequest("GET", "http://localhost:9091/test", nil)
//			w := httptest.NewRecorder()
//
//			// Act
//			srv.Handler.ServeHTTP(w, req)
//
//			// Assert
//
//		})
//	})
//}
