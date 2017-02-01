package main

import (
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kataras/iris"
)

var (
	db        *gorm.DB
	jwtmiddle *jwtmiddleware.Middleware
	// Config pour mise en prod
	jwtkey     []byte = []byte("uneecle") // La clé pour signer les jwt
	sqldbname  string = "apialbum"        // Nom de la db qui accueille les tables
	smtpserver string = "smtp.free.fr"    // le serv smtp
	smtpport   int    = 587               // port smtp
	mailfrom   string = "lesite@domain.com"
)

func main() {
	db, _ = gorm.Open("mysql", "root:@/"+sqldbname+"?charset=utf8&parseTime=True&loc=Local")

	// Migre les tables sur la db
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Album{})
	db.AutoMigrate(&Favorite{})
	defer db.Close()

	// Vérifie l'authenticité de la clé passé par le header ("Authorization")
	jwtmiddle := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtkey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	//créer un user
	iris.Post("/signup", signup)

	//Sert un token
	iris.Post("/login", login)

	//reset le pass
	iris.Post("/reset", jwtmiddle.Serve, reset)

	//Récupère les albums ,le param pagination est égale à l'offset , ex : /album/50 retournera
	//10 albums à partir de l'id 50
	iris.Get("/album/:pagination", listalbums)

	//Récupère les favoris de l'utilisateur
	iris.Get("/favorite/:pagination/*userid", jwtmiddle.Serve, listfavorite)

	//ajoute ou retire des favoris (nécessite d'etre log)
	iris.Post("/favorite/:idalbum", jwtmiddle.Serve, favorite)

	//créer un album (nécessite le grade admin)
	iris.Post("/album", jwtmiddle.Serve, newalbum)

	iris.Listen("localhost:2020")
}
