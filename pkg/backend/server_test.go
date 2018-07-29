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

	"github.com/RafalKorepta/coding-challenge/pkg/api/email/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
)

const swagger_ui_body = `<title>Swagger UI</title>
  <link rel="icon" type="image/png" href="images/favicon-32x32.png" sizes="32x32" />
  <link rel="icon" type="image/png" href="images/favicon-16x16.png" sizes="16x16" />
  <link href='css/typography.css' media='screen' rel='stylesheet' type='text/css'/>
  <link href='css/reset.css' media='screen' rel='stylesheet' type='text/css'/>
  <link href='css/screen.css' media='screen' rel='stylesheet' type='text/css'/>
  <link href='css/reset.css' media='print' rel='stylesheet' type='text/css'/>
  <link href='css/print.css' media='print' rel='stylesheet' type='text/css'/>
  <script src='lib/jquery-1.8.0.min.js' type='text/javascript'></script>
  <script src='lib/jquery.slideto.min.js' type='text/javascript'></script>
  <script src='lib/jquery.wiggle.min.js' type='text/javascript'></script>
  <script src='lib/jquery.ba-bbq.min.js' type='text/javascript'></script>
  <script src='lib/handlebars-2.0.0.js' type='text/javascript'></script>
  <script src='lib/underscore-min.js' type='text/javascript'></script>
  <script src='lib/backbone-min.js' type='text/javascript'></script>
  <script src='swagger-ui.js' type='text/javascript'></script>
  <script src='lib/highlight.7.3.pack.js' type='text/javascript'></script>
  <script src='lib/marked.js' type='text/javascript'></script>
  <script src='lib/swagger-oauth.js' type='text/javascript'></script>`

var _ = Describe("Server that register REST and gRPC endpoint on the same port", func() {
	var (
		srv        *http.Server
		err        error
		grpcServer *grpc.Server
		//listener net.Listener
		//request  *http.Request
		//writer   *httptest.ResponseRecorder
		response      *http.Response
		body          []byte
		newMockServer *httptest.Server
		client        *http.Client

		//opts     []Option
	)

	JustBeforeEach(func() {
		srv, grpcServer, err = createHTTPServer(listener.Addr().String(),
			opts...)

		Expect(err).NotTo(HaveOccurred())
		Expect(srv).NotTo(BeNil())
		Expect(grpcServer).NotTo(BeNil())
		//writer = httptest.NewRecorder()
		//srv.Handler.ServeHTTP(writer, request)
		//
		//response = writer.Result()
		//body, err = ioutil.ReadAll(response.Body)
		//Expect(err).NotTo(HaveOccurred())

		//httptest.NewServer(srv.Handler)
	})

	AfterEach(func() {

	})

	Describe("secure server", func() {
		BeforeEach(func() {
			opts = []Option{
				WithCertFile("test_data/server.pem"),
				WithKeyFile("test_data/server.key"),
				WithSecure(true),
			}
		})

		JustBeforeEach(func() {
			newMockServer = httptest.NewTLSServer(srv.Handler)
		})

		Describe("Http1.1", func() {
			Context("when swagger-ui URI is called", func() {
				JustBeforeEach(func() {
					client = newMockServer.Client()
					response, err = client.Get(newMockServer.URL + "/swagger-ui")
					body, err = ioutil.ReadAll(response.Body)
				})

				It("should response with moved permanently response", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
					Expect(response.Header.Get("Content-Type")).To(Equal(mime.TypeByExtension(".html")))
					Expect(string(body)).To(ContainSubstring(swagger_ui_body))
				})
			})

			Context("when swagger.json URI is correct", func() {
				JustBeforeEach(func() {
					client = newMockServer.Client()
					response, err = client.Get(newMockServer.URL + "/swagger.json")
					body, err = ioutil.ReadAll(response.Body)

				})

				It("should return swagger-ui webpage", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
					Expect(response.Header.Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
					Expect(string(body)).To(Equal(email.Swagger))
				})
			})

			It("Should return swagger.json", func() {

			})
		})

		Describe("gRPC", func() {

		})
	})

	Describe("Insecure server", func() {
		BeforeEach(func() {
			opts = []Option{
				WithSecure(false),
			}
		})

		JustBeforeEach(func() {
			httptest.NewServer(srv.Handler)
		})

		Describe("Rest", func() {

		})

		Describe("gRPC", func() {

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
