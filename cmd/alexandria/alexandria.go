package alexandria

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"

	tiktokdata "github.com/s-alad/tiktok-of-alexandria/data"
	"github.com/s-alad/tiktok-of-alexandria/internal/bucket"
	"github.com/s-alad/tiktok-of-alexandria/internal/database"
	"github.com/s-alad/tiktok-of-alexandria/internal/tikwm"
	log "github.com/sirupsen/logrus"
	"resty.dev/v3"
)

type Library struct {
	id     string
	tikwm  *tikwm.TIKWM
	bucket *bucket.Bucket
	rest   *resty.Client
	mdb    *database.MDB
}

func (a *Library) ID() string {
	return a.id
}

func Run(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.WithField("file", file).Errorf("failed to open file: %v", err)
		return
	}
	defer f.Close()

	var data tiktokdata.TikTokData
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		log.WithField("file", file).Errorf("failed to decode file: %v", err)
		return
	}

	alexandria := &Library{
		id:     "000",
		tikwm:  tikwm.Create("000"),
		bucket: bucket.Create("000"),
		rest:   resty.New(),
		mdb:    database.Create("000"),
	}
	fmt.Println(alexandria)

	log.Info("processing media started")
	alexandria.ProcessMedia(data.Activity.FavoriteVideos.FavoriteVideoList)
}

func (a *Library) ProcessMedia(favorites []tiktokdata.FavoriteVideo) {
	for _, m := range favorites {
		log.WithField("media", m).Info("processing media")

		/* exists, _, _ := a.bucket.Exists(ExtractID(m.Link))
		if exists {
			log.WithField("media", m).Warn("media already exists")
			continue
		} */

		time.Sleep(125 * time.Millisecond)

		media, mediaType, err := a.tikwm.GetMedia(m.Link)
		if err != nil || mediaType == tikwm.MediaTypeNone {
			log.WithField("media", m).Errorf("failed to get media: %v", err)
			a.mdb.Save(a.id, ExtractID(m.Link), "-1", m.Date, database.Failure, mediaType)
			continue
		} else {
			log.WithFields(log.Fields{"id": media.Data.ID, "type": mediaType}).Info("tikwm media retrieved")
		}

		if mediaType == tikwm.MediaTypeSlides {
			log.WithFields(log.Fields{"id": media.Data.ID, "type": mediaType}).Info("CONSUME SLIDES")
			location, _ := a.ConsumeSlides(media)
			a.mdb.Save(a.id, media.Data.ID, location, m.Date, database.Success, mediaType)
		}

		if mediaType == tikwm.MediaTypeVideo {
			log.WithFields(log.Fields{"id": media.Data.ID, "type": mediaType}).Info("CONSUME VIDEO")
			location, _ := a.ConsumeVideo(media)
			a.mdb.Save(a.id, media.Data.ID, location, m.Date, database.Success, mediaType)
		}
	}
}

func (a *Library) ConsumeSlides(media *tikwm.TikwmRequestResponse) (string, error) {
	base := media.Data.ID + "/"

	if media.Data.Music != "" {
		data, err := a.Extract(media.Data.Music)
		if err != nil {
			log.WithFields(log.Fields{
				"id":    media.Data.ID,
				"type":  tikwm.MediaTypeSlides,
				"music": media.Data.Music,
			}).Errorf("failed to extract music: %v", err)
			return "", err
		}

		location, err := a.bucket.Upload(base+"music", data, "audio/mpeg")
		if err != nil {
			log.WithFields(log.Fields{
				"id":    media.Data.ID,
				"type":  tikwm.MediaTypeSlides,
				"music": media.Data.Music,
			}).Errorf("failed to upload music: %v", err)
			return "", err
		}

		log.WithFields(log.Fields{
			"id":       media.Data.ID,
			"type":     tikwm.MediaTypeSlides,
			"location": location,
		}).Info("location result success")
	}

	for i, slide := range media.Data.Images {
		data, err := a.Extract(slide)
		if err != nil {
			log.WithFields(log.Fields{
				"id":    media.Data.ID,
				"type":  tikwm.MediaTypeSlides,
				"slide": slide,
			}).Errorf("failed to extract slide: %v", err)
			return "", err
		}

		location, err := a.bucket.Upload(base+fmt.Sprintf("slide_%d", i), data, "image/jpeg")
		if err != nil {
			log.WithFields(log.Fields{
				"id":    media.Data.ID,
				"type":  tikwm.MediaTypeSlides,
				"slide": slide,
			}).Errorf("failed to upload slide: %v", err)
			return "", err
		}

		log.WithFields(log.Fields{
			"id":       media.Data.ID,
			"type":     tikwm.MediaTypeSlides,
			"location": location,
		}).Info("location result success")
	}

	return a.bucket.Path(base), nil
}

func (a *Library) ConsumeVideo(media *tikwm.TikwmRequestResponse) (location string, err error) {
	data, err := a.Extract(media.Data.Play)
	if err != nil {
		log.WithFields(log.Fields{
			"id":    media.Data.ID,
			"type":  tikwm.MediaTypeVideo,
			"video": media.Data.Play,
		}).Errorf("failed to extract video: %v", err)
		return "", err
	}

	location, err = a.bucket.Upload(media.Data.ID, data, "video/mp4")
	if err != nil {
		log.WithFields(log.Fields{
			"id":    media.Data.ID,
			"type":  tikwm.MediaTypeVideo,
			"video": media.Data.Play,
		}).Errorf("failed to upload video: %v", err)
		return "", err
	}

	log.WithFields(log.Fields{
		"id":       media.Data.ID,
		"type":     tikwm.MediaTypeVideo,
		"location": location,
	}).Info("location result success")

	return location, nil
}

func (a *Library) Extract(url string) (io.ReadCloser, error) {
	resp, err := a.rest.R().SetDoNotParseResponse(true).Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: %s", resp.Status())
	}

	return resp.RawResponse.Body, nil
}

func ExtractID(url string) string {
	pattern := `^https:\/\/www\.tiktokv\.com\/share\/video\/(\d+)\/?$`
	re, _ := regexp.Compile(pattern)
	matches := re.FindStringSubmatch(url)
	return matches[1]
}
