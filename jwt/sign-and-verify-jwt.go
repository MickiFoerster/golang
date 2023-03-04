package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

var hmacSampleSecret = []byte("AllYourBase")

func main() {

	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		Issuer:    "The holy Issuer",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("singed:string:  %v\n", ss)

	token, err = jwt.Parse(ss, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		for _, claim := range []string{"iss", "foo", "nbf"} {
			fmt.Println(claim, ":", claims[claim])
		}
	} else {
		fmt.Println(err)
	}

}

/*


	// Make request.  See func restrictedHandler for example request processor
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%v/restricted", serverPort), nil)
	fatal(err)
    // set BEARER token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	res, err := http.DefaultClient.Do(req)
	fatal(err)

	// Read the response body
	buf := new(bytes.Buffer)
	io.Copy(buf, res.Body)
	res.Body.Close()
	fmt.Println(buf.String())







func createToken(user string) (string, error) {
    signKey    *rsa.PrivateKey
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	t.Claims = &CustomClaimsExample{
		&jwt.StandardClaims{

			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		},
		"level1",
		CustomerInfo{user, "human"},
	}

	return t.SignedString(signKey)
}
*/
