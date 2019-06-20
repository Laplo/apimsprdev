package actions

import (
	"fmt"
	"github.com/apimsprdev/dao"
	"github.com/apimsprdev/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
	"net/http"
	"time"
)

type UserResource struct{}

type Credentials struct {
	Password string `json:"password", db:"password"`
	Email    string `json:"email", db:"email"`
}

var cache = models.NewCache()

func (ur UserResource) SignUp(c buffalo.Context) error {
	// Parse and decode the request body into a new `Credentials` instance
	creds := &Credentials{
		c.Request().FormValue("password"),
		c.Request().FormValue("email"),
	}
	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	id := dao.CreateUser(creds.Password, creds.Email)
	if id == "" {
		return c.Render(500, r.String("unable to create user"))
	}
	return c.Render(201, r.String("User created : "+id))
	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
}

func (ur UserResource) SignIn(c buffalo.Context) error {
	// Parse and decode the request body into a new `Credentials` instance
	creds := &Credentials{
		c.Request().FormValue("password"),
		c.Request().FormValue("email"),
	}

	user := dao.GetUserByEmailAndPassword(creds.Password, creds.Email)

	if user.Email == "" {
		return c.Render(500, r.String("wrong credentials"))
	}

	sessT := user.ID
	sessionToken := sessT.String()

	if !cache.CheckIsIn(sessionToken) {
		cache.Check(sessionToken)
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds, the same as the cache
	http.SetCookie(c.Response(), &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})
	// If we reach this point, that means the users password was correct, and that they are authorized
	// The default 200 status is sent
	return c.Render(200, r.JSON(user))
}

func (ur UserResource) Refresh(c buffalo.Context) error {
	// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route
	cont, err := c.Request().Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return c.Render(401, r.JSON(err))
		}
		return c.Render(400, r.JSON(err))
	}
	sessionToken := cont.Value
	if !cache.CheckIsIn(sessionToken) {
		return c.Render(401, r.String("Cache empty"))
	}
	newSessionToken, _ := uuid.NewV4()
	cache.Check(newSessionToken.String())

	// Set the new token as the users `session_token` cookie
	http.SetCookie(c.Response(), &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken.String(),
		Expires: time.Now().Add(120 * time.Second),
	})

	return c.Render(200, r.String("Session refreshed"))
}

func (ur UserResource) Welcome(c buffalo.Context) error {
	// We can obtain the session token from the requests cookies, which come with every request
	cont, err := c.Request().Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			return c.Render(401, r.JSON(err))
		}
		// For any other type of error, return a bad request status
		return c.Render(400, r.JSON(err))
	}
	sessionToken := cont.Value
	if cache.CheckIsIn(sessionToken) {
		return c.Render(200, r.String(fmt.Sprintf("Welcome %s!", sessionToken)))
	}
	return c.Render(401, r.String("Disconnected"))
}

/*
func (ur UserResource) List(c buffalo.Context) error {
	return c.Render(200, r.JSON(db))
}

func (ur UserResource) Create(c buffalo.Context) error {
	// on crée un nouvel utilisateur
	id, _ := uuid.NewV4()
	user := &models.User{
		// on génère un nouvel id
		ID: id,
	}
	// on l'ajoute à notre base de données
	db[user.ID] = *user

	return c.Render(201, r.JSON(user))
}

func (ur UserResource) Show(c buffalo.Context) error {
	// on récupère l'id de la requête qu'on formate en uuid
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		// si l'id n'est pas un uuid on génère une erreur
		return c.Render(500, r.String("id is not uuid v4"))
	}

	// on récupère notre utilisateur dans notre base de données
	user, ok := db[id]
	if ok {
		// si il existe, on le retourne
		return c.Render(200, r.JSON(user))
	}

	// si il n'existe pas, on retourne une erreur 404
	return c.Render(404, r.String("user not found"))
}
*/
