package drive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var DriveService *drive.Service

func GetDriveService() error {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		fmt.Printf("Unable to read credentials.json file. Err: %v\n", err)
		return err
	}

	// If you want to modifyt this scope, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveScope)

	if err != nil {
		return err
	}

	client := getClient(config)

	ctx := context.Background()
	service, err := drive.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		fmt.Printf("Cannot create the Google Drive service: %v\n", err)
		return err
	} else {
		log.Println("Initialized Google Drive service")
	}

	DriveService = service

	return err
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	fmt.Println("Paste Authrization code here :")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// func createFolder(service *drive.Service, name string, parentId string) (*drive.File, error) {
// 	d := &drive.File{
// 		Name:     name,
// 		MimeType: "application/vnd.google-apps.folder",
// 		Parents:  []string{parentId},
// 	}

// 	file, err := service.Files.Create(d).Do()

// 	if err != nil {
// 		log.Println("Could not create dir: " + err.Error())
// 		return nil, err
// 	}

// 	return file, nil
// }

func CreateFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	logrus.Info("=================> Processing upload file to Google Drive Server")

	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := service.Files.Create(f).Media(content).Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	logrus.Info("File ", file.Name, " was uploaded to Google Drive!")

	return file, nil
}

func RemoveFile(service *drive.Service, fileId string) (err error) {
	logrus.Info("=================> Removing file from Google Drive Server")

	err = service.Files.Delete(fileId).Do()
	if err != nil {
		return err
	}

	logrus.Info("Successfully deleleted file !")

	return err
}

func GetFile(service *drive.Service, fieldID string) (*drive.File, error) {
	file, err := service.Files.Get(fieldID).Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}

func ListFile(service *drive.Service, fID string) (*drive.FileList, error) {
	qFolderID := fmt.Sprintf("%s in parents", fID)

	files, err := service.Files.List().Q(qFolderID).Do()
	if err != nil {
		return nil, err
	}

	return files, nil
}

// GeneratePublicURL : generate public url for uploaded file
func GeneratePublicURL(service *drive.Service, fileID string) (*drive.File, error) {

	drivePermission := &drive.Permission{
		Role: "reader",
		Type: "anyone",
	}

	_, err := service.Permissions.Create(fileID, drivePermission).Do()
	if err != nil {
		return nil, err
	}

	file, err := service.Files.Get(fileID).Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}
