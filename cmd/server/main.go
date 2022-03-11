package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// タグによりJSONにシリアライズ
// タグがないと，フィールド名がそのまま用いられる
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

var Db *sql.DB // 構造体sql.DBの宣言

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=gwp dbname=gwp password=ppp sslmode=disable") // DB接続
	if err != nil {
		panic(err)
	}

	for _, album := range albums {
		album.CreateAlbum()
	}
}

// ginのContextはリクエストの詳細や認証，シリアライズされたJSONを渡す
func getAlbums(ctx *gin.Context) {
	albums, err := GetAlbums()
	if err == nil {
		ctx.IndentedJSON(http.StatusOK, albums) // ステータスコード200で構造体albumsをJSONにシリアライズしてレスポンス
	}
}

func GetAlbums() ([]album, error) {
	var albums []album
	rows, err := Db.Query("select * from albums")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		album := album{}
		err = rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
		if err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}
	rows.Close()

	return albums, nil
}

func postAlbums(ctx *gin.Context) {
	var newAlbum album

	// リクエストボディの値をnewAlbumにバインド
	if err := ctx.BindJSON(&newAlbum); err != nil {
		return
	}
	newAlbum.CreateAlbum()

	albums = append(albums, newAlbum)
	ctx.IndentedJSON(http.StatusOK, newAlbum)
}

func (album *album) CreateAlbum() (err error) {
	statement := "insert into albums (title, artist, price) values ($1, $2, $3) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(album.Title, album.Artist, album.Price).Scan(&album.ID)
	return
}

func getAlbumByID(ctx *gin.Context) {
	id := ctx.Param("id") // URLパラメータの値を取得

	album, err := GetAlbumByID(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"}) // 該当データがない場合，ステータスコード404でメッセージを返す
		return
	}
	ctx.IndentedJSON(http.StatusOK, album)
}

func GetAlbumByID(id string) (album album, err error) {
	err = Db.QueryRow("select id, title, artist, price from albums where id = $1", id).Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
	return
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums) // エンドポイントとハンドラーをマッピング（渡しているのはハンドラーであってハンドラー関数の結果ではない点に注意）
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}
