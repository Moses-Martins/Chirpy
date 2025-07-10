package main

import (
	"net/http"
	"encoding/json"
	"log"
	"time"
	"github.com/google/uuid"
	"github.com/Moses-Martins/Chirpy/internal/auth"
	"github.com/Moses-Martins/Chirpy/internal/database"
)


type AcceptsEmail struct {
	Password string `json:"password"`
	Email string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password string `json:"password"`
}


type UserDisplayed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}



func (cfg *apiConfig) Login(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := AcceptsEmail{}
	
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}


	UserDb, err := cfg.DB.GetUsers(req.Context())
	if err != nil {
    	http.Error(w, "Cannot Retrieve Users", http.StatusNotFound)
        return
	}



	dbToStruct := make([]User, 0, len(UserDb))
    for _, dbUser := range UserDb {
        userResp := User{
            ID:        dbUser.ID,
            CreatedAt: dbUser.CreatedAt,
            UpdatedAt: dbUser.UpdatedAt,
            Email:      dbUser.Email,
			Password: dbUser.HashedPassword,
        }
        dbToStruct = append(dbToStruct, userResp)
    }


	respBodyInitial, exist := findChirpByEmail(dbToStruct, params.Email)

	if exist != true {
		w.WriteHeader(404)
		w.Write([]byte("Incorrect email or password"))
		return
	}

	err = auth.CheckPasswordHash(params.Password, respBodyInitial.Password)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Incorrect email or password"))
		return
	}


	token, err := auth.MakeJWT(respBodyInitial.ID, cfg.JwtSecret, time.Duration(3600) * time.Second)
	if err != nil {
		log.Printf("Cannot generate token %s", err)
		w.WriteHeader(500)
		return
	}

	refreshtoken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error generating Refresh Token: %s", err)
		w.WriteHeader(500)
		return
	}

	_, err = cfg.DB.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token: refreshtoken,
		UserID: respBodyInitial.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})


	respBody := UserDisplayed{
		ID: respBodyInitial.ID,
		CreatedAt: respBodyInitial.CreatedAt,
		UpdatedAt: respBodyInitial.UpdatedAt,
		Email: respBodyInitial.Email,
		Token: token,
		RefreshToken: refreshtoken,
	}

	data, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)

}


func findChirpByEmail(users []User, email string) (User, bool) {
    for _, user := range users {
        if user.Email == email {
            return user, true // found
        }
    }
    return User{}, false // not found
}

