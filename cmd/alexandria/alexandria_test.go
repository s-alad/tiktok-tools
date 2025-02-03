package alexandria

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/s-alad/tiktok-of-alexandria/internal/bucket"
	"github.com/s-alad/tiktok-of-alexandria/internal/tikwm"
	"github.com/stretchr/testify/assert"
	"resty.dev/v3"
)

func init() {
	err := godotenv.Load(filepath.Join("..", "..", ".env"))
	if err != nil {
		fmt.Println("Error loading .env file. Ensure the file exists and is in the correct directory.")
	}
}

func TestAlexandria(t *testing.T) {
	alexandria := &Library{
		id:     "000",
		tikwm:  tikwm.Create("000"),
		bucket: bucket.Create("000"),
		rest:   resty.New(),
	}

	media, mediaType, err := alexandria.tikwm.GetMedia("https://www.tiktok.com/@thecommissariat/video/7465037411429387526")
	assert.Equal(t, mediaType, tikwm.MediaTypeVideo)
	assert.Nil(t, err)
	assert.NotEmpty(t, media.Data.Play)
	assert.Contains(t, media.Data.Play, "https://")

	file, err := alexandria.Extract(media.Data.Play)
	assert.Nil(t, err)
	assert.NotNil(t, file)

	location, err := alexandria.bucket.Upload("7465037411429387526", file, "video/mp4")
	assert.Nil(t, err)
	assert.NotEmpty(t, location)

	fmt.Printf("file uploaded to: %s\n", location)
}
