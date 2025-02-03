package tikwm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVideo(t *testing.T) {
	tikwm := &TIKWM{id: "000"}

	video, _, err := tikwm.GetMedia("https://www.tiktok.com/@thecommissariat/video/7465037411429387526")
	if err != nil {
		t.Fatalf("tikwm.GetVideo returned an error: %v", err)
	}

	assert.NotNil(t, video)
	assert.Equal(t, "7465037411429387526", video.Data.ID)
}
