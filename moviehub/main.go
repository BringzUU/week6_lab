package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"moviehub/auth"
	"moviehub/handler"
	"moviehub/mw"
	"moviehub/store"
)

func main() {
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		log.Fatal("JWT_SECRET ต้องยาวอย่างน้อย 32 ตัวอักษร")
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("ต่อฐานข้อมูลไม่ติด ", err)
	}

	tickets := &auth.Issuer{Secret: []byte(secret), TTL: time.Hour}
	h := &handler.Handler{Store: &store.Store{Pool: pool}, Tickets: tickets}

	login := func(f http.HandlerFunc) http.HandlerFunc { return mw.RequireLogin(tickets, f) }

	mux := http.NewServeMux()
	// ประตูสาธารณะ
	mux.HandleFunc("POST /api/auth/register", h.Register)
	mux.HandleFunc("POST /api/auth/login", h.Login)
	mux.HandleFunc("GET /api/movies", h.SearchMovies)
	mux.HandleFunc("GET /api/movies/{id}", h.GetMovie)
	mux.HandleFunc("GET /api/movies/{id}/reviews", h.ListReviews)
	// ต้องเป็นสมาชิก
	mux.HandleFunc("GET /api/me", login(h.Me))
	mux.HandleFunc("POST /api/movies/{id}/reviews", login(h.AddReview))
	mux.HandleFunc("DELETE /api/reviews/{reviewId}", login(h.DeleteReview))
	mux.HandleFunc("GET /api/me/wishlist", login(h.ListWishlist))
	mux.HandleFunc("PUT /api/me/wishlist/{movieId}", login(h.PutWishlist))
	mux.HandleFunc("DELETE /api/me/wishlist/{movieId}", login(h.DeleteWishlist))
	// ต้องเป็นแอดมิน (สวมสองประตูซ้อนกัน)
	mux.HandleFunc("GET /api/admin/members", login(mw.RequireRole("admin", h.AdminMembers)))

	origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("MovieHub API ฟังอยู่ที่ :"+port, "ต้อนรับ", origins)
	log.Fatal(http.ListenAndServe(":"+port, mw.Logging(mw.CORS(origins, mux))))
}
