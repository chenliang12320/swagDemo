package main

import (
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"swagDemo/datasource"
	"swagDemo/repositories"
	"swagDemo/services"
	"swagDemo/web/controllers"

	_ "swagDemo/docs"
)

// @title Iris 和 Swagger 演示
// @version 1.0
// @description demo 演示.

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	app := iris.New()

	url := swagger.URL("http://localhost:8080/swagger/doc.json") //The url pointing to API definition
	app.Get("/swagger/{any:path}", swagger.WrapHandler(swaggerFiles.Handler, url))

	// Serve our controllers.
	mvc.New(app.Party("/hello")).Handle(new(controllers.HelloController))
	// You can also split the code you write to configure an mvc.Application
	// using the `mvc.Configure` method, as shown below.
	mvc.Configure(app.Party("/api/v1/movies"), movies)

	app.Run(iris.Addr(":8080"))
}

// note the mvc.Application, it's not iris.Application.
func movies(app *mvc.Application) {
	// Add the basic authentication(admin:password) middleware
	// for the /movies based requests.
	// app.Router.Use(middleware.BasicAuth)

	// Create our movie repository with some (memory) data from the datasource.
	repo := repositories.NewMovieRepository(datasource.Movies)
	// Create our movie service, we will bind it to the movie app's dependencies.
	movieService := services.NewMovieService(repo)
	app.Register(movieService)

	// serve our movies controller.
	// Note that you can serve more than one controller
	// you can also create child mvc apps using the `movies.Party(relativePath)` or `movies.Clone(app.Party(...))`
	// if you want.
	app.Handle(new(controllers.MovieController))
}
