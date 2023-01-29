package models

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AccessToken struct {
	AccessToken string `json:"accessToken"`
}
