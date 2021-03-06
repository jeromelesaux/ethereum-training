package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jeromelesaux/ethereum-training/config"
	"github.com/jeromelesaux/ethereum-training/model"
	"github.com/jeromelesaux/ethereum-training/token"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var conf *oauth2.Config
var cred Credentials

type OauthState struct {
	State     string `json:"state"`
	OriginURL string `json:"originurl"`
}

func (ctr *Controller) AuthHandler(c *gin.Context) {
	fmt.Fprintf(os.Stdout, "Enter in Auth handler.\n")
	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	queryState := c.Request.URL.Query().Get("state")
	reader := bytes.NewReader([]byte(queryState))
	ostate := &OauthState{}
	if err := json.NewDecoder(reader).Decode(ostate); err != nil {
		c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": err.Error()})
		return
	}
	if retrievedState != ostate.State {
		log.Printf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
		c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Invalid session state."})
		return
	}
	code := c.Request.URL.Query().Get("code")
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Login failed. Please try again."})
		return
	}

	client := conf.Client(oauth2.NoContext, tok)
	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)
	u := model.User{}
	if err = json.Unmarshal(data, &u); err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error marshalling response. Please try agian."})
		return
	}
	session.Set("user-id", u.Email)
	err = session.Save()
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving session. Please try again."})
		return
	}

	//seen := false
	model.AddUser(u)
	c.Redirect(http.StatusTemporaryRedirect, ostate.OriginURL)
	return
}

func (ctr *Controller) LoginHandler(c *gin.Context) {
	fmt.Fprintf(os.Stdout, "Enter in Login handler.\n")

	from := c.Query("from")
	state, err := token.RandToken(32)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Error while generating random data."})
		return
	}
	session := sessions.Default(c)
	session.Set("state", state)
	err = session.Save()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Error while saving session."})
		return
	}

	ostate := OauthState{State: state, OriginURL: from}
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(&ostate); err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": err.Error()})
		return
	}

	link := getLoginURL(buffer.String())
	fmt.Printf("redirect link is : (%s) and state : (%s)\n", link, state)
	c.Redirect(http.StatusTemporaryRedirect, link)
	//c.HTML(http.StatusOK, "auth.tmpl", gin.H{"link": link})
}

func (ctr *Controller) LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	err := session.Save() // destroy session for user
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Error while saving session."})
		return
	}
	c.HTML(http.StatusOK, "index.tmpl", nil)
}

// Credentials which stores google ids.
type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

func LoadCredentialsEnv() {
	cid := os.Getenv("GOOGLE_ID")
	csecret := os.Getenv("GOOGLE_SECRET")
	if cid == "" || csecret == "" {
		log.Printf("No google credentials found. Aborting...")
		os.Exit(1)
	}
	callbackUrl := config.MyConfig.CallbackUrl

	conf = &oauth2.Config{
		ClientID:     cid,
		ClientSecret: csecret,
		RedirectURL:  callbackUrl + "/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
	return
}

func LoadCredentials() {
	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		LoadCredentialsEnv()
	}
	if err := json.Unmarshal(file, &cred); err != nil {
		log.Println("unable to marshal data")
		return
	}
	callbackUrl := config.MyConfig.CallbackUrl

	conf = &oauth2.Config{
		ClientID:     cred.Cid,
		ClientSecret: cred.Csecret,
		RedirectURL:  callbackUrl + "/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
}

func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

// AuthorizeRequest is used to authorize a request for a certain end-point group.
func (ctr *Controller) AuthorizeRequest() gin.HandlerFunc {
	fmt.Fprintf(os.Stdout, "Enter in Authorize handler.\n")
	return func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("user-id")
		if v == nil {
			fullpath := c.FullPath()

			fmt.Fprintf(os.Stdout, "fullpath:[%s]\n", fullpath)
			//	referer := c.Request.URL.Scheme + c.Request.Host + c.Request.RequestURI
			//	c.Request.Header.Set("Referer", referer)
			c.Redirect(http.StatusTemporaryRedirect, "/login?from="+url.QueryEscape(fullpath))
			c.Abort()
		}

		c.Next()
	}
}

func (ctr *Controller) Certification(c *gin.Context) {
	c.HTML(http.StatusOK, "certification.tmpl", nil)
	return
}

func (ctr *Controller) Verification(c *gin.Context) {
	c.HTML(http.StatusOK, "checker.tmpl", nil)
	return
}

func (ctr *Controller) Safebox(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.HTML(http.StatusOK, "safebox.tmpl", nil)
	return
}
