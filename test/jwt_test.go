package test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go-api/helper"
	"go-api/model"
	"testing"
)

func TestJWT(t *testing.T) {
	uid, _ := uuid.NewUUID()
	token, err := helper.GenerateJWT(&model.UserResponse{
		Id:       uid,
		Username: "testjwt",
	})
	t.Log(token, err)
	assert.NoError(t, err)
	assert.NotNil(t, token)

	payload, err := helper.ValidateJWT(token)
	t.Log(payload, err)
	assert.NoError(t, err)
	assert.NotNil(t, payload)
}
