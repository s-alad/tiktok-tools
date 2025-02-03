package tikwm

import (
	"fmt"
	"net/http"

	"github.com/s-alad/tiktok-of-alexandria/internal/proxy"
	log "github.com/sirupsen/logrus"
	"resty.dev/v3"
)

const (
	BASE_URL = "https://tikwm.com/api/"
)

type TIKWM struct {
	id    string
	rest  *resty.Client
	proxy *proxy.Proxy
}

func (t *TIKWM) ID() string {
	return t.id
}

func (t *TIKWM) init() {
	t.rest = resty.New()
	t.proxy = proxy.Create("tikwm")
	/* t.rest.SetDebug(true) */
	t.rest.SetTransport(&http.Transport{
		Proxy:             nil,
		DisableKeepAlives: true,
	})
}

func (t *TIKWM) GetMedia(url string) (*TikwmRequestResponse, MediaType, error) {

	prox := t.proxy.GetProxy()
	log.WithFields(log.Fields{
		"id":    t.id,
		"proxy": prox,
	}).Info("tikwm.GetMedia")

	t.rest.SetProxy(prox)

	resp, err := t.rest.R().
		SetFormData(map[string]string{
			"url": url,
			/* "count":  fmt.Sprintf("%d", body.Count),
			"cursor": fmt.Sprintf("%d", body.Cursor),
			"web":    fmt.Sprintf("%d", body.Web),
			"hd":     fmt.Sprintf("%d", body.HD), */
		}).
		SetResult(&TikwmRequestResponse{}).
		Post(BASE_URL)
	if err != nil {
		return nil, MediaTypeNone, err
	}

	result := resp.Result().(*TikwmRequestResponse)
	if result.Msg != "success" {
		return nil, MediaTypeNone, fmt.Errorf("tikwm.GetMedia failed: %s on %s", result.Msg, url)
	}

	if result.Data.Images != nil {
		return result, MediaTypeSlides, nil
	} else {
		return result, MediaTypeVideo, nil
	}
}

func Create(id string) *TIKWM {
	t := &TIKWM{id: id}
	t.init()
	return t
}
