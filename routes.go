package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
)

type ValidAlbum struct {
	Name        string `form:"name"`
	Description string `form:"description" `
	Image       string `form:"image"`
	Years       string `form:"years" `
	Genre       string `form:"genre"`
}

type ValidUser struct {
	Pseudo string `form:"pseudo"`
	Email  string `form:"email" `
}

type ValidLogin struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type ValideFavorite struct {
	Album string `form:"album"`
}

type Cclaims struct {
	Userid uint   `json:"userid"`
	Grade  int    `json:"grade"`
	Email  string `json:email`
	jwt.StandardClaims
}

// Check dans la DB si les données correspondent et sert un JWT si tous se passe bien
func login(ctx *iris.Context) {
	json := ValidLogin{}
	err := ctx.ReadJSON(&json)
	if err != nil {
		ctx.JSON(200, iris.Map{"status": "Erreur", "info": "Requiert un json {username ,password"})
	} else {
		var Cuser User
		db.Model(&User{}).Where("Pseudo = ? ", json.Username).First(&Cuser)
		if json.Username == Cuser.Pseudo && passwordverif(Cuser.Password, json.Password) {
			claims := &Cclaims{
				Cuser.ID,
				Cuser.Grade,
				Cuser.Email,
				jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			signed, err := token.SignedString(jwtkey)
			if err != nil {
				fmt.Println(err)
			}
			ctx.JSON(200, iris.Map{"status": "Succes", "token": signed})
		} else {
			ctx.JSON(200, iris.Map{"status": "Erreur"})
		}
	}
}

// Route pour créer un album (nécessite le grade admin)
func newalbum(ctx *iris.Context) {
	rawtoken := ctx.Get("jwt").(*jwt.Token).Raw

	token, err := jwt.ParseWithClaims(rawtoken, &Cclaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtkey), nil
	})
	if err != nil {
		ctx.JSON(200, iris.Map{"status": "erreur", "info": "Probleme token"})
	}

	if claims, ok := token.Claims.(*Cclaims); ok && token.Valid {
		Grade := fmt.Sprintf("%d", claims.Grade)
		Intgrade, _ := strconv.Atoi(Grade)
		if Intgrade != 2 {
			ctx.JSON(200, iris.Map{"status": "erreur", "info": "Nécessite le rang administrateur"})
			return
		}
	}

	form := ValidAlbum{}
	// Si le bind ne retourne pas d'erreur on récupère les données
	err = ctx.ReadForm(&form)
	if err != nil {
		fmt.Println("Error when reading form: " + err.Error())
	} else {
		if form.Name != "" {
			Album := Album{}
			Album.Name = form.Name
			Album.Description = form.Description
			Album.Genre = form.Genre
			Album.Image = form.Image
			YearsInt, err := strconv.Atoi(form.Years)
			if err != nil {
				ctx.JSON(200, iris.Map{"status": "Le champ years nécéssite un entier"})
				return
			}
			Album.Years = YearsInt
			request := db.Create(&Album)
			if request.Error != nil {
				fmt.Println(request.Error)
			}
			ctx.JSON(200, iris.Map{"status": "Album enregistré"})

		} else {
			ctx.JSON(401, iris.Map{"status": "unauthorized"})
		}
	}

}

func reset(ctx *iris.Context) {
	randompass := generatepass()
	// hash le pass genere
	hashpass := encrypt(randompass)

	rawtoken := ctx.Get("jwt").(*jwt.Token).Raw

	token, err := jwt.ParseWithClaims(rawtoken, &Cclaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtkey), nil
	})
	if err != nil {
		ctx.JSON(409, iris.Map{"status": "erreur", "info": "Probleme token"})
	}

	if claims, ok := token.Claims.(*Cclaims); ok && token.Valid {
		Email := fmt.Sprintf("%s", claims.Email)

		var user User
		db.Where("email = ? ", Email).First(&user)
		user.Password = hashpass
		db.Save(&user)

		go resetpassword(Email, randompass)
	} else {
		ctx.JSON(409, iris.Map{"status": "erreur", "info": "Probleme token"})
	}
	ctx.JSON(200, iris.Map{"status": "OK", "info": "Mot de passe envoyé"})
}

// Récupère un json ex: {"pseudo":"LEPSEUDO","email":"LEMAIL"} et enregistre un user
func signup(ctx *iris.Context) {
	form := ValidUser{}
	// Si le bind ne retourne pas d'erreur on récupère les données
	err := ctx.ReadJSON(&form)
	if err != nil {
		ctx.JSON(200, iris.Map{"status": "Erreur", "info": "Requiert un json {pseudo ,email}"})
	} else {
		if form.Pseudo != "" {
			aUser := User{}
			aUser.Pseudo = form.Pseudo
			aUser.Email = form.Email
			// genere un string password
			randompass := generatepass()
			// hash le pass genere
			hashpass := encrypt(randompass)
			aUser.Password = hashpass

			var tuser User
			var count int
			db.Where("pseudo = ? AND email = ?", form.Pseudo, form.Email).First(&tuser).Count(&count)
			if count == 0 {
				db.Create(&aUser)
			} else {
				ctx.JSON(200, iris.Map{"status": "Pseudo ou Email déja existant"})
				return
			}

			go sendpassword(form.Email, randompass)

			ctx.JSON(200, iris.Map{"status": "Nouvel Utilisateur enregistré"})
		} else {
			ctx.JSON(401, iris.Map{"status": "unauthorized"})
		}
	}

}

//Récupère les albums
//le param pagination est égale à l'offset de la req mysql , ex : /album/50 retournera 10 albums à partir de l'id 50
func listalbums(ctx *iris.Context) {
	pagination := ctx.Param("pagination")
	//vérifie que le param est un entier
	pagin, err := strconv.Atoi(pagination)
	if err != nil {
		ctx.JSON(401, iris.Map{"status": "Nécessite un entier après /album/"})
		return
	}
	var album []Album
	db.Offset(pagin).Limit(10).Find(&album)
	ctx.JSON(200, album)
}

func listfavorite(ctx *iris.Context) {
	pagination := ctx.Param("pagination")
	targetornot := false
	var intuserid int
	var err error
	if ctx.Param("userid") != "/" {
		puserid := strings.Replace(ctx.Param("userid"), "/", "", -1)
		intuserid, err = strconv.Atoi(puserid)
		if err != nil {
			ctx.JSON(401, iris.Map{"status": "Nécessite un entier après /favorite/:pagination/"})
			return
		} else {
			targetornot = true
		}
	}

	//vérifie que le param est un entier
	pagin, err := strconv.Atoi(pagination)
	if err != nil {
		ctx.JSON(401, iris.Map{"status": "Nécessite un entier après /favorite/"})
		return
	}

	rawtoken := ctx.Get("jwt").(*jwt.Token).Raw

	token, err := jwt.ParseWithClaims(rawtoken, &Cclaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtkey), nil
	})
	if err != nil {
		ctx.JSON(200, iris.Map{"status": "erreur", "info": "Probleme token"})
	}

	if claims, ok := token.Claims.(*Cclaims); ok && token.Valid {
		Userid := fmt.Sprintf("%d", claims.Userid)
		Uid, _ := strconv.Atoi(Userid)

		var favorite []Favorite
		if targetornot {
			db.Where("user_id = ?", uint(intuserid)).Offset(pagin).Limit(10).Find(&favorite)
			ctx.JSON(200, favorite)
			return
		} else {
			db.Where("user_id = ?", uint(Uid)).Offset(pagin).Limit(10).Find(&favorite)
			ctx.JSON(200, favorite)
			return
		}

	} else {
		fmt.Println(err)
	}
}

func favorite(ctx *iris.Context) {
	idalbum := ctx.Param("idalbum")
	idalbumint, errint := strconv.Atoi(idalbum)
	if errint != nil {
		ctx.JSON(200, iris.Map{"status": "Erreur", "info": "Requiert un entier après /favorite/"})
		return
	}

	rawtoken := ctx.Get("jwt").(*jwt.Token).Raw

	token, err := jwt.ParseWithClaims(rawtoken, &Cclaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtkey), nil
	})
	if err != nil {
		ctx.JSON(200, iris.Map{"status": "erreur", "info": "Probleme token"})
	}

	if claims, ok := token.Claims.(*Cclaims); ok && token.Valid {
		Userid := fmt.Sprintf("%d", claims.Userid)
		Uid, _ := strconv.Atoi(Userid)

		var alb Album
		var counta int
		db.Where("id= ?", idalbumint).First(&alb).Count(&counta)
		if counta == 0 {
			ctx.JSON(200, iris.Map{"status": "error", "info": "Album inexistant"})
			return
		}

		var onefav Favorite
		var count int
		db.Where("user_id = ? AND album = ?", Uid, idalbumint).First(&onefav).Count(&count)

		if count == 0 {
			fav := Favorite{UserId: uint(Uid), Album: uint(idalbumint)}
			request := db.Create(&fav)
			if request.Error != nil {
				fmt.Println(request.Error)
			}
			ctx.JSON(200, iris.Map{"status": "OK", "info": "Favoris ajouté"})
			return
		} else {
			fav := Favorite{UserId: uint(Uid), Album: uint(idalbumint)}
			request := db.Delete(&fav)
			if request.Error != nil {
				fmt.Println(request.Error)
			}
			ctx.JSON(200, iris.Map{"status": "OK", "info": "Favoris retiré"})
			return
		}

	} else {
		fmt.Println(err)
	}
	ctx.JSON(200, token)

}
