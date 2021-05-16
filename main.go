package main

import (
	"database/sql"
	"fmt"
	"goblog/bootstrap"
	"goblog/pkg/database"
	"goblog/pkg/logger"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

var router *mux.Router
var db *sql.DB

type Article struct {
	Title, Body string
	ID          int64
}

func (a Article) Delete() (rowsAffected int64, err error) {
	rs, err := db.Exec("DELETE FROM articles where id = " + strconv.FormatInt(a.ID, 10))
	if err != nil {
		return 0, err
	}

	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}
	return 0, nil
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// 1.设置标头
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 2.继续处理请求
		next.ServeHTTP(rw, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Printf(r.URL.Path)
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(rw, r)
	})
}

func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)
	article, err := getArticleByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		rowsAffected, err := article.Delete()
		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		} else {
			if rowsAffected > 0 {
				indexURL, _ := router.Get("article.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}
}

func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

func main() {

	database.Initialize()
	db = database.DB

	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()

	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
