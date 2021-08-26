package services

import (
	"context"
	"net/http"

	"google.golang.org/api/idtoken"
)

var httpClient = &http.Client{}

func VerifyIdToken(idToken string) (*idtoken.Payload, error) {
	client_id := "829385733688-1k8lri3tncr0bbdugqt4h09papv7vktj.apps.googleusercontent.com"
	payload, err := idtoken.Validate(context.Background(), idToken, client_id)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
