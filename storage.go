package firebasestorage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// Storage represent a Firebase Storage client
type Storage struct {
	Bucket       string
	Token        string
	RefreshToken string
	APIKey       string
}

func (s *Storage) resource(path string) string {
	return fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o?name=%s", s.Bucket, path)
}

func (s *Storage) auth(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Firebase %s", s.Token))
}

func (s *Storage) request(auth bool, verb string, loc string, data io.Reader) (map[string]interface{}, error) {
	req, err := http.NewRequest(verb, loc, data)
	if err != nil {
		return nil, err
	}

	if auth {
		s.auth(req)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	r := new(map[string]interface{})
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	return *r, nil
}

// Object will fetch the storage object metadata from Firebase
func (s *Storage) Object(path string) (map[string]interface{}, error) {
	return s.request(true, "GET", s.resource(path), nil)
}

// Download will download a file from Firebase
// `path` - object path in firebase
// `toFile` - local file to write to
func (s *Storage) Download(path, toFile string) error {
	obj, err := s.Object(path)
	if err != nil {
		return err
	}
	downloadToken := obj["downloadTokens"].(string)

	out, err := os.Create(toFile)
	if err != nil {
		return err
	}
	defer out.Close()

	reader, err := s.Read(path, downloadToken)
	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = io.Copy(out, reader)
	return err
}

// Read will read a file from Firebase storage, providing the bytes over an `io.ReadClose`
// `path` - object path in firebase
// `downloadToken` - a token retrieved from `Storage.Object`
func (s *Storage) Read(path, downloadToken string) (io.ReadCloser, error) {
	resp, err := http.Get(fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s",
		s.Bucket,
		url.QueryEscape(path),
		downloadToken,
	))
	return resp.Body, err
}

// Put will store a file in Firebase Storage
func (s *Storage) Put(file, path string) (map[string]interface{}, error) {
	data, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	res, err := s.request(true, "PUT", s.resource(path), data)
	if err != nil {
		return nil, err
	}
	return res, err
}

//Refresh will refresh your Firebase user access tokens. Typically these tokens are
//valid for an hour.
//
//You can call this when operations fail, when the time almost expires, or before
//every operation just to be safe.
func (s *Storage) Refresh() error {
	loc := fmt.Sprintf("https://securetoken.googleapis.com/v1/token?key=%s", s.APIKey)
	data := map[string]string{
		"grantType":    "refresh_token",
		"refreshToken": s.RefreshToken,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	res, err := s.request(false, "POST", loc, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	accessToken := res["access_token"].(string)
	if accessToken != "" {
		s.Token = accessToken
	}
	refreshToken := res["refresh_token"].(string)
	if refreshToken != "" {
		s.RefreshToken = refreshToken
	}

	return nil
}
