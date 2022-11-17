package utils

import (
	"bytes"
	"cryptoshare/conf"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func Create2fa(username string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      conf.AppHost,
		AccountName: username,
	})

	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, err
	}
	png.Encode(&buf, img)

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(dir, "storage", "qrcode", username+".png")
	if err := os.WriteFile(fullPath, buf.Bytes(), 0644); err != nil {
		return nil, err
	}

	return key, nil
}

func Validate2fa(otp, secret string) bool {
	trimmedToken := strings.TrimSpace(otp)

	return totp.Validate(trimmedToken, secret)
}
