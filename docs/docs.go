package docs

import (
	"github.com/swaggo/swag"
)

var doc = `{
	"swagger": "2.0",
	"info": {
        "description": "Lafu Server API 문서 정리 내용입니다.",
        "title": "Lafu Server API 문서",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0"
    },
    "host": "http://localhost:7000/api/v1.0",
	"basePath": "/api/v1.0",
	"paths": {

	}
}`

type s struct{}

func (s *s) ReadDoc() string {
	return doc
}

func init() {
	swag.Register(swag.Name, &s{})
}
