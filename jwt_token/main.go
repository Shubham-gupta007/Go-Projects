package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

type User struct {
	// Username string `json:"username"`
	// Password string `json:"password"`
	Accesstoken string `json:"accesstoken"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

func CreateTokenEndpoint(w http.ResponseWriter, req *http.Request) {
	var user User
	fmt.Println("Inside Create token point")

	_ = json.NewDecoder(req.Body).Decode(&user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"accesstoken": user.Accesstoken,
		// "username": user.Username,
		// "password": user.Password,
	})
	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
}

func ProtectedEndpoint(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	fmt.Println("Inside Protected end point")
	token, _ := jwt.Parse(params["token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte("secret"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user User
		mapstructure.Decode(claims, &user)
		json.NewEncoder(w).Encode(user)
	} else {
		json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
	}
}

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	fmt.Println("Inside middle ware1")

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Inside middle ware2")
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				fmt.Println("Inside middle ware3")
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					fmt.Println("Inside middle ware4")
					return []byte("secret"), nil
				})
				if error != nil {
					json.NewEncoder(w).Encode(Exception{Message: error.Error()})
					return
				}
				if token.Valid {
					fmt.Println("Inside middle ware5")
					context.Set(req, "decoded", token.Claims)
					next(w, req)
					fmt.Println("Inside middle ware6")
				} else {
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
		}
	})
}

// func ValidateMiddleware(w http.ResponseWriter, req *http.Request) {
// 	// return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 	fmt.Println("INside middle ware")
// 	authorizationHeader := req.Header.Get("authorization")
// 	if authorizationHeader != "" {
// 		bearerToken := strings.Split(authorizationHeader, " ")
// 		if len(bearerToken) == 2 {
// 			token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
// 				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 					return nil, fmt.Errorf("There was an error")
// 				}
// 				return []byte("secret"), nil
// 			})
// 			if error != nil {
// 				json.NewEncoder(w).Encode(Exception{Message: error.Error()})
// 				return
// 			}
// 			if token.Valid {
// 				context.Set(req, "decoded", token.Claims)
// 				next(w, req)
// 			} else {
// 				json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
// 			}
// 		}
// 	} else {
// 		json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
// 		// }
// 	}

// }

func TestEndpoint(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Inside Test Endpoint1")
	decoded := context.Get(req, "decoded")
	var user User
	fmt.Println("user1:", user)

	mapstructure.Decode(decoded.(jwt.MapClaims), &user)

	fmt.Println("user2:", user)
	json.NewEncoder(w).Encode(user)
	// fmt.Println("json.NewEncoder(w).Encode(user):", json.NewEncoder(w).Encode(user))
	fmt.Println("Inside Test enpoint2")

}

func main() {
	router := mux.NewRouter()
	fmt.Println("Starting the application...")
	router.HandleFunc("/authenticate", CreateTokenEndpoint).Methods("POST")
	router.HandleFunc("/protected", ProtectedEndpoint).Methods("GET")
	router.HandleFunc("/test", ValidateMiddleware(TestEndpoint)).Methods("GET")
	log.Fatal(http.ListenAndServe(":12345", router))
}
