package main

import (
	"database/sql"
	// "log"
	// "net"
	"net/http"
	// "os"
	// "os/exec"
	// "text/template"
	"time"

	// "golang.org/x/net/netutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	// "github.com/pkg/errors"
	// "github.com/takashabe/go-router"
	// "github.com/takashabe/go-session"
	// _ "github.com/takashabe/go-session/memory"
)

// DB table mapping
type UserModel struct {
	ID        int
	Name      string
	Email     string
	Salt      string
	Passhash  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DB table mapping
type Tweet struct {
	ID        int
	UserID    int
	UserName  string
	Content   string
	CreatedAt time.Time
}

// template content
type IndexContent struct {
	User      *UserModel
	Following int
	Followers int
	Tweets    []*Tweet
}

// template content
type LoginContent struct {
	Message string
}

// template content
type UserContent struct {
	Myself     *UserModel
	User       *UserModel
	Tweets     []*Tweet
	Followable bool
}

// template content
type FollowingContent struct {
	FollowingList []*Following
}

// template content
type FollowersContent struct {
	UserList []*UserModel
}

// for FollowingContent table mapping struct
type Following struct {
	UserId    int
	FollowId  int
	UserName  string
	CreatedAt time.Time
}

var db *sql.DB

func main() {
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
	// 	os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
	// 	os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
	// 	os.Getenv("DB_DATABASE"),
	// )

	// var err error
	// db, err = sql.Open("mysql", dsn)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	e := echo.New()
	// funcs := template.FuncMap{
	// 	"encode_json": func(v interface{}) string {
	// 		b, _ := json.Marshal(v)
	// 		return string(b)
	// 	},
	// }

	    // 全てのリクエストで差し込みたいミドルウェア（ログとか）はここ
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
	
		// ルーティング
		e.GET("/hello", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello World")
		})
	
		// サーバー起動
				e.Start(":8080")
	}