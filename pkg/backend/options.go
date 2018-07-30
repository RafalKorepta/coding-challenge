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

var (
	defaultOptions = &options{}
)

type options struct {
	certFile           string
	keyFile            string
	secure             bool
	serverOverrideName string
}

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

type Option func(*options)

// WithCertFile setup where the certificate should be found
func WithCertFile(c string) Option {
	return func(o *options) {
		o.certFile = c
	}
}

// WithCertFile setup where the key should be found
func WithKeyFile(k string) Option {
	return func(o *options) {
		o.keyFile = k
	}
}

// WithSecure if set to true, then it will start server with TLS encryption
func WithSecure(s bool) Option {
	return func(o *options) {
		o.secure = s
	}
}

// WithSecure if set to true, then it will start server with TLS encryption
func WithServerOverrideName(s string) Option {
	return func(o *options) {
		o.serverOverrideName = s
	}
}
