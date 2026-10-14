# MovieHub API (Lab สัปดาห์ที่ 12)

ติดตั้งไลบรารีก่อนรันครั้งแรก:

    go mod tidy

รัน:

    export DATABASE_URL="postgres://postgres:secret@localhost:5432/moviehub"
    openssl rand -hex 32 > .jwt_secret   # ครั้งเดียว
    export JWT_SECRET=$(cat .jwt_secret)
    export ALLOWED_ORIGINS=http://localhost:3000
    go run .

ทดสอบตั๋วปลอม: `go test -count=1 ./auth -v`
