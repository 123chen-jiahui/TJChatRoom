package tool

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func GenerateToken(account string) (string, error) {
	// 注意这里要选SigningMethodHS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Account": account,
		//"exp":     time.Now().Add(time.Hour.Truncate(2)),
		"exp": time.Now().Unix() + 3600*2, // 有效期2小时
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(MConfig.SecretStr))
	fmt.Println(tokenString)
	return tokenString, err
}

func ParseToken(tokenStr string) (account string) {
	account = ""
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(MConfig.SecretStr), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["Account"])
		account = fmt.Sprint(claims["Account"])
	} else {
		fmt.Println(err)
	}
	return
}
