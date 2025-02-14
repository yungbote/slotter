package clients

import (
  "bytes"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "github.com/yungbote/slotter/backend/services/auth/internal/models"
)

type DatabaseClient interface {
  CreateUser(u *models.User) error
  GetUserByEmail(email string) (*models.User, error)
}

type dbClient struct {
  baseURL     string
  httpClient  *http.Client
}

func NewDatabaseClient(baseURL string) DatabaseClient {
  return &dbClient{
    baseURL: baseURL,
    httpClient: &http.Client{},
  }
}

func (d *dbClient) CreateUser(u *models.User) error {
  payload, err := json.Marshal(u)
  if err != nil {
    return err
  }
  url := fmt.Sprintf("%s/v1/user", d.baseURL)
  resp, err := d.httpClient.Post(url, "application/json", bytes.NewReader(payload))
  if err != nil {
    return err
  }
  defer resp.Body.Close()
  if resp.StatusCode != http.StatusCreated {
    body, _ := ioutil.ReadAll(resp.Body)
    return fmt.Errorf("CreateUser failed: %s", string(body))
  }
  return nil
}

func (d *dbClient) GetUserByEmail(email string) (*models.User, error) {
  url := fmt.Sprintf("%s/v1/user/email/%s", d.baseURL, email)
  resp, err := d.httpClient.Get(url)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  if resp.StatusCode == http.StatusNotFound {
    return nil, nil
  }
  if resp.StatusCode != http.StatusOK {
    body, _ := ioutil.ReadAll(resp.Body)
    return nil, fmt.Errorf("GetUserByEmail failed with status %d: %s", resp.StatusCode, string(body))
  }
  var user models.User
  if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
    return nil, err
  }
  return &user, nil
}
