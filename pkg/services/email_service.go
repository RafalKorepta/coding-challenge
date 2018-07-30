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
package services

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
