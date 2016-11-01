package googlespreadsheet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	sheets "google.golang.org/api/sheets/v4"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	readRange = "Mep 11/16!A3:H173"
	//readRange     = "Mep 11/16!A78:H173"
	spreadsheetID = "1NBh34AQfSrIEzzRa8xL5Mlth9XlIfBsyK8e0R0X_UcQ"
)

// UserStory is the internal representation of UserStory mapped on spreadsheet backlog
type UserStory struct {
	Priority    int32
	Number      int32
	Description string
	ManDay      float32
	Bot         []string
}

// Print allow to print an user story in stdout
func (us UserStory) Print() {
	formated := fmt.Sprintf("Priority : %d, Number : %d, ManDay : %f, Bot : %s \n %s", us.Priority, us.Number, us.ManDay, us.Bot, us.Description)
	fmt.Println(formated)
}

// DescriptionToTitle allow to cut description to create a title
func (us UserStory) DescriptionToTitle() string {
	split := strings.Split(us.Description, "\n")[0]
	var result string
	if len(split) > 120 {
		result = fmt.Sprintf("#%d : %s", us.Number, split[:120])
	} else {
		result = fmt.Sprintf("#%d : %s", us.Number, split)
	}
	return result
}

// ReadSheet read the google backlog spreadsheet and convert it into internet UserStory struct
func ReadSheet() []UserStory {
	var result = []UserStory{}
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-quickstart.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}

	if len(resp.Values) > 0 {
		for _, row := range resp.Values {

			if len(row) > 4 {
				priority, _ := strconv.ParseInt(row[1].(string), 10, 32)
				number, _ := strconv.ParseInt(row[2].(string), 10, 32)
				taskDescription := row[3].(string)
				manDayStr := strings.Replace(row[6].(string), ",", ".", -1)
				taskManDay, _ := strconv.ParseFloat(manDayStr, 32)
				taskBots := strings.Split(row[7].(string), "&")
				us := UserStory{
					Priority:    int32(priority),
					Number:      int32(number),
					Description: taskDescription,
					ManDay:      float32(taskManDay),
					Bot:         taskBots,
				}
				if len(us.Bot) < 1 {
					panic("Issue has not bot, so it will not be created in any gitlab project")
				}
				result = append(result, us)
			}
		}
	} else {
		fmt.Print("No data found.")
	}
	return result
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("sheets.googleapis.com-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
