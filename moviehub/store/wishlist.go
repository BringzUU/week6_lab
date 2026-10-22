package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type WishlistItem struct {
	Movie
	Note    *string   `json:"note"`
	AddedAt time.Time `json:"addedAt"`
}

// ListWishlist คืนเฉพาะรายการของสมาชิกคนนั้น เลขสมาชิกมาจากตั๋วเท่านั้น
func (s *Store) ListWishlist(ctx context.Context, memberID int64) ([]WishlistItem, error) {
	rows, err := s.Pool.Query(ctx,
		`SELECT m.id, m.title, m.title_th,
                (SELECT coalesce(g.name_th, g.name_en) FROM movie_genres mg
                   JOIN genres g ON g.id = mg.genre_id
                  WHERE mg.movie_id = m.id AND mg.position = 0),
                extract(year FROM m.release_date)::int,
                round(m.vote_average, 1)::float8,
                `+posterSQL+`,
                w.note, w.added_at
         FROM wishlist_items w
         JOIN movies m ON m.id = w.movie_id
         WHERE w.member_id = $1
         ORDER BY w.added_at DESC`, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []WishlistItem{}
	for rows.Next() {
		var it WishlistItem
		if err := rows.Scan(&it.ID, &it.Title, &it.TitleTh, &it.Genre, &it.Year,
			&it.Rating, &it.Poster, &it.Note, &it.AddedAt); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

// PutWishlist กดซ้ำได้ผลเหมือนเดิม ด้วย INSERT ... ON CONFLICT
func (s *Store) PutWishlist(ctx context.Context, memberID, movieID int64, note *string) error {
	_, err := s.Pool.Exec(ctx,
		`INSERT INTO wishlist_items (member_id, movie_id, note)
         VALUES ($1, $2, $3)
         ON CONFLICT (member_id, movie_id) DO UPDATE SET note = EXCLUDED.note`,
		memberID, movieID, note)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return ErrNoMovie
	}
	return err
}

// DeleteWishlist ลบเฉพาะแถวของเจ้าของตั๋ว ลบซ้ำก็ไม่เป็นไร
func (s *Store) DeleteWishlist(ctx context.Context, memberID, movieID int64) error {
	_, err := s.Pool.Exec(ctx,
		`DELETE FROM wishlist_items WHERE member_id = $1 AND movie_id = $2`,
		memberID, movieID)
	return err
}
