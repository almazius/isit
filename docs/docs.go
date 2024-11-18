// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Docs developer",
            "url": "https://t.me/sigy922"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/material": {
            "get": {
                "security": [
                    {
                        "AuthToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "File"
                ],
                "summary": "Ручка",
                "responses": {
                    "200": {
                        "description": "http.StatusOK"
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "AuthToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "File"
                ],
                "summary": "Добавление материала",
                "parameters": [
                    {
                        "description": "52",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Material"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "http.StatusOK"
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "AuthToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "File"
                ],
                "summary": "Обновелние материала",
                "parameters": [
                    {
                        "description": "52",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Product"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "http.StatusOK"
                    }
                }
            }
        },
        "/api/order": {
            "get": {
                "security": [
                    {
                        "AuthToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "File"
                ],
                "summary": "Ручка",
                "responses": {
                    "200": {
                        "description": "http.StatusOK"
                    }
                }
            }
        },
        "/api/order/": {
            "post": {
                "security": [
                    {
                        "AuthToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "File"
                ],
                "summary": "Добавление заявки",
                "parameters": [
                    {
                        "description": "52",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Order"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "http.StatusOK"
                    }
                }
            }
        },
        "/api/order/status": {
            "patch": {
                "security": [
                    {
                        "AuthToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "File"
                ],
                "summary": "Обновелние материала",
                "parameters": [
                    {
                        "description": "52",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UpdateOrder"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "http.StatusOK"
                    }
                }
            }
        },
        "/api/product": {
            "get": {
                "security": [
                    {
                        "AuthToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "File"
                ],
                "summary": "Ручка",
                "responses": {
                    "200": {
                        "description": "http.StatusOK"
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "AuthToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "File"
                ],
                "summary": "Добавление продукта",
                "parameters": [
                    {
                        "description": "52",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Product"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "http.StatusOK"
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Material": {
            "type": "object",
            "required": [
                "count",
                "name",
                "price",
                "reject_percent"
            ],
            "properties": {
                "address": {
                    "type": "string"
                },
                "count": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "reject_percent": {
                    "type": "number"
                }
            }
        },
        "models.MaterialSmallInfo": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "number"
                },
                "id": {
                    "type": "integer"
                }
            }
        },
        "models.Order": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "number"
                },
                "product_id": {
                    "type": "integer"
                }
            }
        },
        "models.Product": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "materials": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.MaterialSmallInfo"
                    }
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "reject_percent": {
                    "type": "number"
                }
            }
        },
        "models.UpdateOrder": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "AuthToken": {
            "description": "autorization token from auth_service",
            "type": "apiKey",
            "name": "sessionid",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.02",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Бэкенд Сервиса",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
