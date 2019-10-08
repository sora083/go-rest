package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// DB table mapping
type User struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Passhash  string    `json:"passhash,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
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
	User      *User
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
	Myself     *User
	User       *User
	Tweets     []*Tweet
	Followable bool
}

// template content
type FollowingContent struct {
	FollowingList []*Following
}

// template content
type FollowersContent struct {
	UserList []*User
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

	var err error
	db, err := sql.Open("mysql", "app:app@/app?parseTime=true&charset=utf8mb4")
	//db, err := sql.Open("mysql", "app:app@unix(/tmp/mysql.sock)/isucon?parseTime=true&charset=utf8mb4")
	if err != nil {
		log.Fatal(err)
	}

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

	e.GET("/users", func(c echo.Context) error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Commit()

		var user User
		if err := tx.QueryRow("SELECT * FROM user").Scan(&user.ID, &user.Name, &user.Passhash, &user.CreatedAt, &user.UpdatedAt); err != sql.ErrNoRows {
			//tx.Rollback()
			if err == nil {
				//return resError(c, "duplicated", 409)
				return err
			}
			return err
		}

		return c.JSON(200, echo.Map{
			"id":   user.ID,
			"name": user.Name,
		})
	})

	e.GET("/users/:id", func(c echo.Context) error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Commit()

		id := c.Param("id")
		var user User
		if err := tx.QueryRow("SELECT * FROM user WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Passhash, &user.CreatedAt, &user.UpdatedAt); err != sql.ErrNoRows {
			//tx.Rollback()
			if err == nil {
				//return resError(c, "duplicated", 409)
				return err
			}
			return err
		}

		return c.JSON(200, echo.Map{
			"id":   user.ID,
			"name": user.Name,
		})
	})

	e.POST("/users/:id/tweet", func(c echo.Context) error {
		var params struct {
			Content string `json:"content"`
		}
		c.Bind(&params)

		// require login
		// user, err := getCurrentUser(w, r)
		// if err != nil {
		// 	http.Redirect(w, r, "/login", 302)
		// 	return
		// }

		id := c.Param("id")

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		res, err := db.Exec("INSERT INTO tweet (user_id, content) VALUES (?,?)", id, params.Content)
		if err != nil {
			tx.Rollback()
			return resError(c, "", 0)
		}
		userID, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return resError(c, "", 0)
		}
		if err := tx.Commit(); err != nil {
			return err
		}

		return c.JSON(201, echo.Map{
			"id":      userID,
			"content": params.Content,
		})

		//http.Redirect(w, r, "/", 302)
	})

	// サーバー起動
	e.Start(":8080")
}

func resError(c echo.Context, e string, status int) error {
	if e == "" {
		e = "unknown"
	}
	if status < 100 {
		status = 500
	}
	return c.JSON(status, map[string]string{"error": e})
}
