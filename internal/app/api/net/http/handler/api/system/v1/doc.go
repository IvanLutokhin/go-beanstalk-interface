//go:generate go run github.com/go-swagger/go-swagger/cmd/swagger generate spec -o $ROOT_DIR/api/system/v1/swagger.json --scan-models

// Package v1 Beanstalk API
//
// Documentation of Beanstalk API.
//
// Schemes: http, https
// BasePath: /api/system/v1
// Version: 1
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// SecurityDefinitions:
//    bearer:
//       type: apiKey
//       in: header
//       name: authorization
//    cookie:
//       type: apiKey
//       in: cookie
//       name: access_token
//
// swagger:meta
package v1
