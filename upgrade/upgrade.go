package upgrade

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sipt/shuttle"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const VersionEqual = "equal"
const VersionLess = "less"
const VersionGreater = "greater"

func CheckUpgrade(ver string) (latest, url, status string, err error) {
	r, err := GetLatestRelease()
	if err != nil {
		return
	}
	var tagVer, currentVer int
	currentVer, err = verStr2Int(ver)
	if err != nil {
		return
	}
	for i := len(r.TagName) - 1; i >= 0; i-- {
		if r.TagName[i] == 'v' || r.TagName[i] == 'V' {
			tagVer, err = verStr2Int(r.TagName[i:])
			latest = r.TagName[i:]
			break
		}
	}
	if err != nil {
		return
	}
	if currentVer == tagVer {
		status = VersionEqual
		return
	} else if currentVer > tagVer {
		status = VersionGreater
		return
	} else {
		status = VersionLess
	}
	goos := runtime.GOOS
	if goos == "darwin" {
		goos = "macos"
	}
	for _, v := range r.Assets {
		if strings.HasPrefix(v.Name, fmt.Sprintf("shuttle_%s_%s", goos, runtime.GOARCH)) {
			url = v.BrowserDownloadURL
			return
		}
	}
	err = errors.New("not support platform")
	return
}

func DownloadFile(name, downloadURL string) error {
	os.Remove(name)
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	proxyUrl, err := url.Parse("http://127.0.0.1:" + shuttle.HTTPProxyPort)
	myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	resp, err := myClient.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func verStr2Int(ver string) (int, error) {
	if ver[0] == 'v' || ver[0] == 'V' {
		ver = ver[1:]
	}
	vs := strings.Split(ver, ".")
	if len(vs) < 3 {
		return 0, errors.New(fmt.Sprintf("%s version string as : v0.0.1", ver))
	}
	verStr := ""
	if len(vs[0]) == 1 {
		verStr += "0"
	}
	verStr += vs[0]
	if len(vs[1]) == 1 {
		verStr += "0"
	}
	verStr += vs[1]
	if len(vs[2]) == 1 {
		verStr += "0"
	}
	verStr += vs[2]
	return strconv.Atoi(verStr)
}

func GetLatestRelease() (*Release, error) {
	resp, err := http.Get(LatestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	release := &Release{}
	err = json.Unmarshal(data, release)
	return release, err
}

const LatestURL = "https://api.github.com/repos/sipt/shuttle/releases/latest"

type Release struct {
	Assets []struct {
		BrowserDownloadURL string      `json:"browser_download_url"`
		ContentType        string      `json:"content_type"`
		CreatedAt          string      `json:"created_at"`
		DownloadCount      int         `json:"download_count"`
		ID                 int         `json:"id"`
		Label              interface{} `json:"label"`
		Name               string      `json:"name"`
		NodeID             string      `json:"node_id"`
		Size               int         `json:"size"`
		State              string      `json:"state"`
		UpdatedAt          string      `json:"updated_at"`
		Uploader           struct {
			AvatarURL         string `json:"avatar_url"`
			EventsURL         string `json:"events_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			GravatarID        string `json:"gravatar_id"`
			HTMLURL           string `json:"html_url"`
			ID                int    `json:"id"`
			Login             string `json:"login"`
			NodeID            string `json:"node_id"`
			OrganizationsURL  string `json:"organizations_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			ReposURL          string `json:"repos_url"`
			SiteAdmin         bool   `json:"site_admin"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			Type              string `json:"type"`
			URL               string `json:"url"`
		} `json:"uploader"`
		URL string `json:"url"`
	} `json:"assets"`
	AssetsURL string `json:"assets_url"`
	Author    struct {
		AvatarURL         string `json:"avatar_url"`
		EventsURL         string `json:"events_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		GravatarID        string `json:"gravatar_id"`
		HTMLURL           string `json:"html_url"`
		ID                int    `json:"id"`
		Login             string `json:"login"`
		NodeID            string `json:"node_id"`
		OrganizationsURL  string `json:"organizations_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		ReposURL          string `json:"repos_url"`
		SiteAdmin         bool   `json:"site_admin"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		URL               string `json:"url"`
	} `json:"author"`
	Body            string `json:"body"`
	CreatedAt       string `json:"created_at"`
	Draft           bool   `json:"draft"`
	HTMLURL         string `json:"html_url"`
	ID              int    `json:"id"`
	Name            string `json:"name"`
	NodeID          string `json:"node_id"`
	Prerelease      bool   `json:"prerelease"`
	PublishedAt     string `json:"published_at"`
	TagName         string `json:"tag_name"`
	TarballURL      string `json:"tarball_url"`
	TargetCommitish string `json:"target_commitish"`
	UploadURL       string `json:"upload_url"`
	URL             string `json:"url"`
	ZipballURL      string `json:"zipball_url"`
}
