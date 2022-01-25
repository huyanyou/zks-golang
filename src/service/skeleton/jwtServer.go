package skeleton

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//指定加密密钥
var jwtSecret = []byte("Henu.kingo154")

// token用户实体
type Claims struct {
	Jsessionid string `json:"jsessionid"`
	jwt.StandardClaims
}

func GenerateToken(jsessionid string) (string, error) {
	//	设置token有效时间
	nowTime := time.Now()
	expireTime := nowTime.Add(10000 * time.Hour)

	claims := Claims{
		Jsessionid: jsessionid,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: expireTime.Unix(),
			//	指定token发行人
			Issuer: "zks",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//	该方法内部生成签名字符串,用于获取完整、已签受的token
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

//	根据传入的token值获取到Claims对象信息,（进而获取其中的用户名和密码）
func ParseToken(token string) (*Claims, error) {
	// 用于解析健全的声明
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
