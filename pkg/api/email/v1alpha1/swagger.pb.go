package email

const (
	swagger = `{
  "swagger": "2.0",
  "info": {
    "title": "email.proto",
    "version": "version not set"
  },
  "schemes": [
    "http",
    "https"
  ],
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "paths": {
    "/v1alpha1/email": {
      "post": {
        "summary": "SendMail",
        "operationId": "SendMail",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/v1alpha1EmailResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/v1alpha1EmailRequest"
            }
          }
        ],
        "tags": [
          "EmailService"
        ]
      }
    }
  },
  "definitions": {
    "v1alpha1EmailRequest": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string"
        }
      }
    },
    "v1alpha1EmailResponse": {
      "type": "object",
      "properties": {
        "error": {
          "type": "string"
        }
      }
    }
  }
}
`
)
