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
		"/auth/register/local": {
			"post": {
				"description": "회원가입",
                "produces": [
                    "application/json"
                ],
                "summary": "회원가입",
                "parameters": [
                    {
                        "type": "string",
                        "description": "이메일",
                        "name": "email",
                        "required": true
                    },
					{
                        "type": "string",
                        "description": "유저명",
                        "name": "username",
                        "required": true						
					},
					{
                        "type": "string",
                        "description": "패스워드",
                        "name": "password",
                        "required": true	
					}
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "object",
                        }
                    },	
                    "400": {
                        "description": "VALIDATE_ERROR & WRONG_SCHEMA_BODY_DATA",
                        "schema": {
                            "type": "object",
                        }
                    },
                    "401": {
                        "description": "패스워드가 일치하지 않습니다",
                        "schema": {
                            "type": "object",
                        }
                    },
                    "403": {
                        "description": "계정을 찾을 수 없습니다.",
                        "schema": {
                            "type": "object",
                        }
                    },
                    "404": {
                        "description": "NOT FOUND",
                        "schema": {
                            "type": "object",
                        }
                    },
				}
			}
		}
	}
}`

type s struct{}

func (s *s) ReadDoc() string {
	return doc
}

func init() {
	swag.Register(swag.Name, &s{})
}
