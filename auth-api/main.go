//package main
//
//import (
//	"github.com/mukhametkaly/Diploma/auth-api/models"
//	"golang.org/x/crypto/bcrypt"
//)
//
//func main() {
//	user := models.User{
//		UserName:   "sdfsdf",
//		Password:   "qwerty",
//		Salt:       "qwe",
//		Role:       "",
//		MerchantId: "",
//	}
//
//	password := "qwerty"
//
//	passString := user.Password + "@" + user.Salt
//
//	reqString := password + "@" + user.Salt
//	reqHash, err := bcrypt.GenerateFromPassword([]byte(reqString), 14)
//
//	err = bcrypt.CompareHashAndPassword(reqHash, []byte(passString))
//	if err != nil {
//		panic(err)
//	}
//
//	panic(err)
//}


import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func main() {

	//mySigningKey := []byte("AllYourBase")

	type MyCustomClaims struct {
		Foo string `json:"foo"`
		jwt.RegisteredClaims
	}

	// Create the claims
	claims := MyCustomClaims{
		"bar",
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}

	// Create claims while leaving out some of the optional fields
	//claims = MyCustomClaims{
	//	"bar",
	//	jwt.RegisteredClaims{
	//		// Also fixed dates can be used for the NumericDate
	//		ExpiresAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
	//		Issuer:    "test",
	//	},
	//}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Printf("%v %v\n", ss, err)

	tk, err := jwt.ParseWithClaims(ss, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})

	if claim, ok := tk.Claims.(*MyCustomClaims); ok && tk.Valid {
		fmt.Printf("%v %v", claim.Foo, claim.RegisteredClaims.Issuer)
	} else {
		fmt.Println(err)
	}

}
