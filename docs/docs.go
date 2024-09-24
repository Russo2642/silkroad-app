// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/contact_form": {
            "post": {
                "description": "This method creates a new contact form and, if TourID is provided, sends tour details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "forms"
                ],
                "summary": "Create a contact form",
                "parameters": [
                    {
                        "description": "Contact form data",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/forms.ContactForm"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    }
                }
            }
        },
        "/countries": {
            "get": {
                "description": "This method returns a list of all available countries",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "countries"
                ],
                "summary": "Returns a list of countries",
                "responses": {
                    "200": {
                        "description": "List of countries",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    }
                }
            }
        },
        "/help_with_tour_form": {
            "post": {
                "description": "This method creates a new helpWithTour form",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "forms"
                ],
                "summary": "Create a help with tour form",
                "parameters": [
                    {
                        "description": "helpWithTour form data",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/forms.HelpWithTourForm"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    }
                }
            }
        },
        "/tour": {
            "get": {
                "description": "This method returns a list of tours. You can filter by price range, tour place, tour date, quantity, duration, and search by title. Pagination is also supported via limit and offset.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tour"
                ],
                "summary": "Returns a list of tours",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Filter by price range, example: '100-500'",
                        "name": "priceRange",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by tour place",
                        "name": "tour_place",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Filter by quantity of people",
                        "name": "quantity",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Filter by duration of the tour",
                        "name": "duration",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by date of the tour, format: YYYY-MM-DDT00:00:00+00:00",
                        "name": "tour_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Search tours by title",
                        "name": "search",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Limit the number of returned tours",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Offset for pagination",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of tours",
                        "schema": {
                            "$ref": "#/definitions/http.getAllToursResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    }
                }
            }
        },
        "/tour/{id}": {
            "get": {
                "description": "This method returns the details of a specific tour by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tour"
                ],
                "summary": "Returns tour by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Tour ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Tour details",
                        "schema": {
                            "$ref": "#/definitions/tour.Tour"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    }
                }
            }
        },
        "/tour_editor": {
            "post": {
                "description": "This method allows users to create a custom tour by submitting their details and tour preferences. The tour data will also be sent to Telegram.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tour"
                ],
                "summary": "Create a custom tour",
                "parameters": [
                    {
                        "description": "Custom tour data",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/tour.TourEditor"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ID and creation status of the tour",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/http.errorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "forms.ContactForm": {
            "type": "object",
            "required": [
                "description",
                "email",
                "name",
                "phone"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "tour_id": {
                    "type": "integer"
                }
            }
        },
        "forms.HelpWithTourForm": {
            "type": "object",
            "required": [
                "name",
                "phone",
                "place",
                "when_date"
            ],
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "place": {
                    "type": "string"
                },
                "when_date": {
                    "type": "string"
                }
            }
        },
        "http.errorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "http.getAllToursResponse": {
            "type": "object",
            "properties": {
                "tours": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/tour.Tour"
                    }
                }
            }
        },
        "tour.CalendarType": {
            "type": "object",
            "additionalProperties": {
                "type": "object",
                "additionalProperties": {
                    "$ref": "#/definitions/tour.TimeSlot"
                }
            }
        },
        "tour.TimeSlot": {
            "type": "object",
            "properties": {
                "from": {
                    "type": "string"
                },
                "to": {
                    "type": "string"
                }
            }
        },
        "tour.Tour": {
            "type": "object",
            "required": [
                "activity",
                "currency",
                "description_excursion",
                "description_route",
                "duration",
                "physical_rating",
                "price",
                "quantity",
                "season",
                "title",
                "tour_date",
                "tour_place",
                "tour_type"
            ],
            "properties": {
                "activity": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "calendar": {
                    "$ref": "#/definitions/tour.CalendarType"
                },
                "currency": {
                    "type": "string"
                },
                "description_excursion": {
                    "type": "string"
                },
                "description_route": {
                    "type": "string"
                },
                "duration": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "physical_rating": {
                    "type": "integer",
                    "maximum": 5,
                    "minimum": 1
                },
                "price": {
                    "type": "integer"
                },
                "quantity": {
                    "type": "integer"
                },
                "season": {
                    "type": "string"
                },
                "tariff": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "tour_date": {
                    "type": "string"
                },
                "tour_place": {
                    "type": "string"
                },
                "tour_type": {
                    "$ref": "#/definitions/tour.TourType"
                }
            }
        },
        "tour.TourEditor": {
            "type": "object",
            "required": [
                "activity",
                "email",
                "location",
                "name",
                "phone",
                "tour_date"
            ],
            "properties": {
                "activity": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "location": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "tour_date": {
                    "type": "string"
                }
            }
        },
        "tour.TourType": {
            "type": "string",
            "enum": [
                "Однодневный тур",
                "Многодневный тур",
                "Сити-тур",
                "Эксклюзивный тур",
                "Инфо-тур",
                "Авторский тур"
            ],
            "x-enum-varnames": [
                "OneDayTour",
                "MultiDayTour",
                "CityTour",
                "ExclusiveTour",
                "InfoTour",
                "AuthorsTour"
            ]
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "178.128.123.250:80",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "SilkRoad App API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
