package main

import (
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/richjyoung/echopen"
	v310 "github.com/richjyoung/echopen/openapi/v3.1.0"
)

const Description = `
This is a sample server Petstore server.  You can find out more about
Swagger at [http://swagger.io](http://swagger.io) or on [irc.freenode.net, #swagger](http://swagger.io/irc/).  For
this sample, you can use the api key "special-key" to test the authorization filters.`

func main() {
	// Create a new echOpen wrapper
	api := echopen.New(
		"Swagger Petstore",
		"1.0.0",
		echopen.WithSpecDescription(Description),
		echopen.WithSpecLicense(&v310.License{Name: "Apache 2.0", URL: "http://www.apache.org/licenses/LICENSE-2.0.html"}),
		echopen.WithSpecTermsOfService("http://swagger.io/terms/"),
		echopen.WithSpecContact(&v310.Contact{Email: "apiteam@swagger.io"}),
		echopen.WithSpecServer(&v310.Server{URL: "http://petstore.swagger.io/v2"}),
		echopen.WithSpecExternalDocs(&v310.ExternalDocs{
			Description: "Find out more about Swagger",
			URL:         "http://swagger.io",
		}),
		echopen.WithSpecTag(&v310.Tag{
			Name:        "pet",
			Description: "Everything about your Pets",
			ExternalDocs: &v310.ExternalDocs{
				Description: "Find out more",
				URL:         "http://swagger.io",
			},
		}),
		echopen.WithSpecTag(&v310.Tag{
			Name:        "store",
			Description: "Access to Petstore orders",
		}),
		echopen.WithSpecTag(&v310.Tag{
			Name:        "user",
			Description: "Operations about user",
			ExternalDocs: &v310.ExternalDocs{
				Description: "Find out more about our store",
				URL:         "http://swagger.io",
			},
		}),
	)

	api.Spec.GetComponents().AddSecurityScheme("petstore_auth", &v310.SecurityScheme{
		Type: v310.OAuth2SecuritySchemeType,
		Flows: &v310.OAuthFlows{
			Implicit: &v310.OAuthFlow{
				AuthorizationURL: "http://petstore.swagger.io/oauth/dialog",
				Scopes: map[string]string{
					"write:pets": "modify pets in your account",
					"read:pets":  "read your pets",
				},
			},
		},
	})

	api.Spec.GetComponents().AddSecurityScheme("api_key", &v310.SecurityScheme{
		Type: v310.APIKeySecuritySchemeType,
		Name: "api_key",
		In:   "header",
	})

	api.Spec.GetComponents().AddRequestBody("Pet", &v310.RequestBody{
		Content: map[string]*v310.MediaTypeObject{
			echo.MIMEApplicationJSON: {
				Schema: api.ToSchemaRef(Pet{}),
			},
		},
		Required: true,
	})

	api.Spec.GetComponents().AddRequestBody("UserArray", &v310.RequestBody{
		Content: map[string]*v310.MediaTypeObject{
			echo.MIMEApplicationJSON: {
				Schema: api.ToSchemaRef([]User{}),
			},
		},
		Required: true,
	})

	petGroup := api.Group("/pet", echopen.WithGroupTags("pet"))

	petGroup.POST(
		"",
		noop,
		echopen.WithOperationID("addPet"),
		echopen.WithSummary("Add a new pet to the store"),
		echopen.WithSecurityRequirement("petstore_auth", []string{"write:pets", "read:pets"}),
		echopen.WithRequestBodyRef("Pet"),
		echopen.WithResponseDescription("405", "Invalid input"),
	)

	petGroup.PUT(
		"",
		noop,
		echopen.WithOperationID("updatePet"),
		echopen.WithSummary("Update an existing pet"),
		echopen.WithSecurityRequirement("petstore_auth", []string{"write:pets", "read:pets"}),
		echopen.WithRequestBodyRef("Pet"),
		echopen.WithResponseDescription("400", "Invalid ID supplied"),
		echopen.WithResponseDescription("404", "Pet not found"),
		echopen.WithResponseDescription("405", "Validation exception"),
	)

	petGroup.GET(
		"/findByStatus",
		noop,
		echopen.WithOperationID("findPetsByStatus"),
		echopen.WithSummary("Finds Pets by status"),
		echopen.WithDescription("Multiple status values can be provided with comma separated strings"),
		echopen.WithSecurityRequirement("petstore_auth", []string{"write:pets", "read:pets"}),
		echopen.WithQueryParameter(&echopen.QueryParameterConfig{
			Name:        "status",
			Description: "Status values that need to be considered for filter",
			Required:    true,
			Explode:     true,
			Schema: &v310.Schema{
				Type: "array",
				Items: &v310.Ref[v310.Schema]{
					Value: &v310.Schema{
						Type:    "string",
						Enum:    []string{"available", "pending", "sold"},
						Default: "available",
					},
				},
			},
		}),
		echopen.WithResponseStructConfig("200", &echopen.ResponseStructConfig{
			Description: "successful operation",
			Target:      []Pet{},
			JSON:        true,
			XML:         true,
		}),
		echopen.WithResponseDescription("400", "Invalid status value"),
	)

	petGroup.GET(
		"/findByTags",
		noop,
		echopen.WithOperationID("findPetsByTags"),
		echopen.WithSummary("Finds Pets by tags"),
		echopen.WithDescription("Muliple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing."),
		echopen.WithSecurityRequirement("petstore_auth", []string{"write:pets", "read:pets"}),
		echopen.WithQueryParameter(&echopen.QueryParameterConfig{
			Name:        "tags",
			Description: "Tags to filter by",
			Required:    true,
			Explode:     true,
			Schema:      api.TypeToSchema(reflect.TypeOf([]string{})),
		}),
		echopen.WithResponseStructConfig("200", &echopen.ResponseStructConfig{
			Description: "successful operation",
			Target:      []Pet{},
			JSON:        true,
			XML:         true,
		}),
		echopen.WithResponseDescription("400", "Invalid tag value"),
		echopen.WithDeprecated(),
	)

	petGroup.GET(
		"/{petId}",
		noop,
		echopen.WithOperationID("getPetById"),
		echopen.WithSummary("Find pet by ID"),
		echopen.WithDescription("Returns a single pet"),
		echopen.WithSecurityRequirement("api_key", []string{}),
		echopen.WithPathParameter(&echopen.PathParameterConfig{
			Name:        "petId",
			Description: "ID of pet to return",
			Schema:      api.TypeToSchema(reflect.TypeOf(int64(0))),
		}),
		echopen.WithResponseStruct("200", "successful operation", Pet{}),
		echopen.WithResponseDescription("400", "Invalid ID supplied"),
		echopen.WithResponseDescription("404", "Pet not found"),
		echopen.WithResponseDescription("default", "successful response"),
	)

	// Serve the generated schema
	api.ServeYAMLSpec("/openapi.yml")
	api.ServeUI("/", "/openapi.yml", "5.10.3")

	// Write the full generated spec
	api.WriteYAMLSpec("openapi_out.yml")

	// Start the server
	api.Start("localhost:3000")
}

func noop(c echo.Context) error {
	return nil
}
