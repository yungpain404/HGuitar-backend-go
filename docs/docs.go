// Package docs cung cấp Swagger document cho HGuitar Backend API.
// Để sinh lại file này đầy đủ các path, hãy cài đặt swag CLI và chạy:
//   go install github.com/swaggo/swag/cmd/swag@latest
//   swag init -g cmd/server/main.go -o docs
package docs

import "github.com/swaggo/swag"

// @title HGuitar API
// @version 1.0.0
// @description Backend API for HGuitar
// @termsOfService http://localhost:5501/terms/
// @contact.name API Support
// @contact.url http://localhost:5501/support
// @contact.email support@hguitar.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:5501
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

type swaggerDoc struct{}

func (d *swaggerDoc) ReadDoc() string {
	return `{
  "swagger": "2.0",
  "info": {
    "title": "HGuitar API",
    "description": "Backend API for HGuitar",
    "version": "1.0.0",
    "termsOfService": "http://localhost:5501/terms/",
    "contact": {
      "name": "API Support",
      "url": "http://localhost:5501/support",
      "email": "support@hguitar.com"
    },
    "license": {
      "name": "MIT",
      "url": "https://opensource.org/licenses/MIT"
    }
  },
  "host": "localhost:5501",
  "basePath": "/api/v1",
  "schemes": ["http"],
  "paths": {},
  "securityDefinitions": {
    "BearerAuth": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header"
    }
  }
}
`
}

func init() {
	swag.Register(swag.Name, &swaggerDoc{})
}
