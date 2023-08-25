package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/raphaelmb/go-lenslocked/controllers"
	"github.com/raphaelmb/go-lenslocked/migrations"
	"github.com/raphaelmb/go-lenslocked/models"
	"github.com/raphaelmb/go-lenslocked/templates"
	"github.com/raphaelmb/go-lenslocked/views"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// host := os.Getenv("SMTP_HOST")
	// portStr := os.Getenv("SMTP_PORT")
	// port, err := strconv.Atoi(portStr)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// username := os.Getenv("SMTP_USERNAME")
	// password := os.Getenv("SMTP_PASSWORD")

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	csrfKey := "secret"
	csrfMiddleware := csrf.Protect([]byte(csrfKey), csrf.Secure(false))

	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.tmpl.html", "tailwind.tmpl.html"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.tmpl.html", "tailwind.tmpl.html"))

	r := chi.NewRouter()

	r.Use(csrfMiddleware)
	r.Use(umw.SetUser)

	tmpl := views.Must(views.ParseFS(templates.FS, "home.tmpl.html", "tailwind.tmpl.html"))
	r.Get("/", controllers.StaticHandler(tmpl))

	tmpl = views.Must(views.ParseFS(templates.FS, "contact.tmpl.html", "tailwind.tmpl.html"))
	r.Get("/contact", controllers.StaticHandler(tmpl))

	tmpl = views.Must(views.ParseFS(templates.FS, "faq.tmpl.html", "tailwind.tmpl.html"))
	r.Get("/faq", controllers.FAQ(tmpl))

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page Not Found", http.StatusNotFound)
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)
}
