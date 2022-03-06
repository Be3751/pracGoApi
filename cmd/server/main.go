package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

// ginのContextはリクエストの詳細や認証，シリアライズされたJSONを渡す
func getAlbums(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, albums) // ステータスコード200で構造体albumsをJSONにシリアライズしてレスポンス
}

func postAlbums(ctx *gin.Context) {
	var newAlbum album

	// リクエストボディの値をnewAlbumにバインド
	if err := ctx.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	ctx.IndentedJSON(http.StatusOK, newAlbum)
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums) // エンドポイントとハンドラーをマッピング（渡しているのはハンドラーであってハンドラー関数の結果ではない点に注意）
	router.POST("/albums", postAlbums)
	router.Run("localhost:8080")
}
