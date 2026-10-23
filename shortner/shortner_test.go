package shortner

import (
	"testing"
)

func TestCreatShortLinkFunc(t *testing.T) {
	longUrl, err := RandomString(30)
	if err != nil {
		t.Fatal("radomString error %w", err)
	}
	_, err = CreateShortUrl(longUrl, 0)
	if err != nil {
		t.Fatal("CreateShortUrl error %w", err)

	}

}
