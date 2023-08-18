package hassPassTokenGen

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash this function hashes password
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// CheckHash this check it password provided is qual to hashed password saved in the database
func CheckHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken this function generates user login / auth token / jwt token
/*func GenerateToken(data models.User) (string, error) {
	//	token is valid for 7 days
	date := time.Now().Add(time.Hour * 24 * 7)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":   data,
		"expiry": date.Unix(),
	})

	//	get path from root dir
	pwd, _ := os.Getwd()
	keyPath := pwd + "/jwtsecret.key"

	key, readErr := ioutil.ReadFile(keyPath)
	if readErr != nil {
		return "", readErr
	}

	tokenString, err := token.SignedString(key)

	return tokenString, err
}*/

//check api will renew the token when life is less than 3 days other wise, return null for token

/*func Check(context *gin.Context) {

	userRaw, ok := context.Get("user")
	if !ok {
		context.AbortWithStatus(401)
		return
	}

	user := userRaw.(models.User)

	tokenExpire := int64(context.MustGet("token_expire").(float64))
	now := time.Now().Unix()
	diff := tokenExpire - now

	fmt.Println(diff)

	if diff < 60*60*24*3 {
		//	renew token
		token, _ := GenerateToken(models.User{})
		context.SetCookie("token", token, 60*60*24*7, "/", "", false, true)
		context.JSON(200, gin.H{
			"token": token,
			"user":  user,
		})
		return
	}

	context.JSON(200, gin.H{
		"token": nil,
		"user":  user,
	})

}*/
