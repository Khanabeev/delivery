package main

import (
	"delivery/cmd"
	httpin "delivery/internal/adapters/in/http"
	"delivery/internal/generated/servers"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	oam "github.com/oapi-codegen/echo-middleware"
	"github.com/robfig/cron/v3"

	"net/http"

	"os"
)

func main() {
	config := getConfigs()

	compositionRoot := cmd.NewCompositionRoot(
		config,
	)
	defer compositionRoot.CloseAll()

	startWebServer(compositionRoot, config.HttpPort)
	startCron(compositionRoot)
}

func getConfigs() cmd.Config {
	config := cmd.Config{
		HttpPort:                  goDotEnvVariable("HTTP_PORT"),
		DbHost:                    goDotEnvVariable("DB_HOST"),
		DbPort:                    goDotEnvVariable("DB_PORT"),
		DbUser:                    goDotEnvVariable("DB_USER"),
		DbPassword:                goDotEnvVariable("DB_PASSWORD"),
		DbName:                    goDotEnvVariable("DB_NAME"),
		DbSslMode:                 goDotEnvVariable("DB_SSLMODE"),
		GeoServiceGrpcHost:        goDotEnvVariable("GEO_SERVICE_GRPC_HOST"),
		KafkaHost:                 goDotEnvVariable("KAFKA_HOST"),
		KafkaConsumerGroup:        goDotEnvVariable("KAFKA_CONSUMER_GROUP"),
		KafkaBasketConfirmedTopic: goDotEnvVariable("KAFKA_BASKET_CONFIRMED_TOPIC"),
		KafkaOrderChangedTopic:    goDotEnvVariable("KAFKA_ORDER_CHANGED_TOPIC"),
	}
	return config
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func startWebServer(compositionRoot *cmd.CompositionRoot, port string) {
	handlers, err := httpin.NewServer(
		compositionRoot.NewCreateOrderHandler(),
	)

	if err != nil {
		log.Fatalf("Ошибка инициализации HTTP Server: %v", err)
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
	}))

	spec, err := servers.GetSwagger()
	if err != nil {
		log.Fatalf("Error reading OpenAPI spec: %v", err)
	}
	e.Use(oam.OapiRequestValidator(spec))
	e.Pre(middleware.RemoveTrailingSlash())
	registerSwaggerOpenApi(e)
	registerSwaggerUi(e)
	servers.RegisterHandlers(e, handlers)
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", port)))
}

func registerSwaggerOpenApi(e *echo.Echo) {
	e.GET("/openapi.json", func(c echo.Context) error {
		swagger, err := servers.GetSwagger()
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to load swagger: "+err.Error())
		}

		data, err := swagger.MarshalJSON()
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to marshal swagger: "+err.Error())
		}

		return c.Blob(http.StatusOK, "application/json", data)
	})
}

func registerSwaggerUi(e *echo.Echo) {
	e.GET("/docs", func(c echo.Context) error {
		html := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		  <meta charset="UTF-8">
		  <title>Swagger UI</title>
		  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
		</head>
		<body>
		  <div id="swagger-ui"></div>
		  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
		  <script>
			window.onload = () => {
			  SwaggerUIBundle({
				url: "/openapi.json",
				dom_id: "#swagger-ui",
			  });
			};
		  </script>
		</body>
		</html>`
		return c.HTML(http.StatusOK, html)
	})
}

func startCron(compositionRoot *cmd.CompositionRoot) {
	c := cron.New()
	_, err := c.AddJob("@every 1s", compositionRoot.NewAssignOrdersJob())
	if err != nil {
		log.Fatalf("ошибка при добавлении задачи: %v", err)
	}
	_, err = c.AddJob("@every 1s", compositionRoot.NewMoveCouriersJob())
	if err != nil {
		log.Fatalf("ошибка при добавлении задачи: %v", err)
	}
	c.Start()
}
