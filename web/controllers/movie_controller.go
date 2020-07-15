// file: web/controllers/movie_controller.go

package controllers

import (
	"errors"
	"fmt"
	"strconv"

	"swagDemo/datamodels"
	"swagDemo/services"

	"github.com/kataras/iris/v12"
)

// MovieController is our /movies controller.
type MovieController struct {
	// Our MovieService, it's an interface which
	// is binded from the main application.
	Service services.MovieService
}

// @Summary 获取movies列表
// @Description 获取movies列表
// @ID getMovies
// @tags 演示增删改查API
// @Accept  json
// @Produce  json
// @Success 200 {array} datamodels.Movie
// @Header 200 {string} Token "qwerty"
// @Router /movies [get]
func (c *MovieController) Get() (results []datamodels.Movie) {
	return c.Service.GetAll()
}

// @Summary 查询指定Id的movie
// @Description 查询指定Id的movie
// @ID getMovieById
// @tags 演示增删改查API
// @Accept  json
// @Produce  json
// @Param   id     path    int     true   "ID"
// @Success 200 {object} datamodels.Movie
// @Header 200 {string} Token "qwerty"
// @Router /movies/{id} [get]
func (c *MovieController) GetBy(id int64) (movie datamodels.Movie, found bool) {
	return c.Service.GetByID(id) // it will throw 404 if not found.
}

// @Summary 修改指定ID的movie
// @Description 修改指定ID的movie
// @ID modifyMovieById
// @tags 演示增删改查API
// @Accept  mpfd
// @Produce  json
// @Param   id     path    int     true   "ID"
// @Param   genre     formData    string     true   "Genre"
// @Param   poster     formData    file     true   "Poster"
// @Success 200 {object} datamodels.Movie
// @Header 200 {string} Token "qwerty"
// @Router /movies/{id} [put]
func (c *MovieController) PutBy(ctx iris.Context, id int64) (datamodels.Movie, error) {
	// get the request data for poster and genre
	file, info, err := ctx.FormFile("poster")
	if err != nil {
		return datamodels.Movie{}, errors.New("failed due form file 'poster' missing")
	}
	// we don't need the file so close it now.
	file.Close()

	// imagine that is the url of the uploaded file...
	poster := info.Filename
	genre := ctx.FormValue("genre")

	return c.Service.UpdatePosterAndGenreByID(id, poster, genre)
}

// @Summary 删除指定ID的movies
// @Description 删除指定ID的movies
// @ID deleteMovieById
// @tags 演示增删改查API
// @Accept  json
// @Produce  json
// @Param   id     path    int     true   "ID"
// @Success 200 {object} int
// @Header 200 {string} Token "qwerty"
// @Router /movies/{id} [delete]
func (c *MovieController) DeleteBy(id int64) interface{} {
	wasDel := c.Service.DeleteByID(id)
	if wasDel {
		// return the deleted movie's ID
		return iris.Map{"deleted": id}
	}
	// right here we can see that a method function can return any of those two types(map or int),
	// we don't have to specify the return type to a specific type.
	return iris.StatusBadRequest
}

// @Summary 添加一个Movie
// @Description 添加一个Movie
// @ID addMovie
// @tags 演示增删改查API
// @Accept  json
// @Produce  json
// @Param   name     formData    string     true   "Name"
// @Param   year     formData    int     	true   "Year"
// @Param   genre    formData    string     true   "Genre"
// @Param   poster   formData    string     true   "Poster"
// @Success 200 {object} datamodels.Movie
// @Header 200 {string} Token "qwerty"
// @Router /movies [post]
func (c *MovieController) Post(ctx iris.Context) (datamodels.Movie, error) {
	movie := new(datamodels.Movie)

	// 获取表单数据
	year, err := strconv.Atoi(ctx.FormValue("year"))
	if err != nil {
		year = 0
		// TODO: 返回参数错误
	}

	movie.Name = ctx.FormValue("name")
	movie.Year = year
	movie.Genre = ctx.FormValue("genre")
	movie.Poster = ctx.FormValue("poster")

	fmt.Println("add a movie", ctx.FormValue("name"), ctx.PostValue("name"))

	return c.Service.Add(movie)
}
