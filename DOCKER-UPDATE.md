# Cara Update Docker Compose untuk Backend Saja

## 1. Rebuild dan Restart Backend Service

Untuk rebuild image dan restart container backend:

```bash
docker-compose up -d --build app
```

**Penjelasan:**
- `up -d`: Start container dalam mode detached (background)
- `--build`: Force rebuild image dari Dockerfile
- `app`: Hanya service `app` yang di-rebuild dan restart

## 2. Stop, Rebuild, dan Start Backend

Jika ingin lebih eksplisit:

```bash
# Stop container app
docker-compose stop app

# Rebuild image app (tanpa cache untuk fresh build)
docker-compose build --no-cache app

# Start container app
docker-compose up -d app
```

## 3. Force Recreate Container (Recreate Container dari Scratch)

Jika ada perubahan di docker-compose.yml atau ingin fresh container:

```bash
docker-compose up -d --force-recreate --build app
```

## 4. Rebuild Tanpa Restart Dependencies

Service lain (db, redis, rabbitmq) tetap running, hanya app yang di-rebuild:

```bash
docker-compose build app
docker-compose up -d app
```

## 5. View Logs Setelah Rebuild

Setelah rebuild, cek logs untuk memastikan app berjalan dengan baik:

```bash
docker-compose logs -f app
```

## 6. Restart Saja (Tanpa Rebuild)

Jika hanya ingin restart tanpa rebuild:

```bash
docker-compose restart app
```

## Troubleshooting

### Hapus Container dan Image Lama

Jika ada masalah dengan build cache:

```bash
# Stop dan remove container
docker-compose stop app
docker-compose rm -f app

# Remove image
docker rmi lostmediago_app

# Rebuild dari scratch
docker-compose up -d --build app
```

### Clear Build Cache

Untuk clear semua build cache:

```bash
docker builder prune -f
docker-compose build --no-cache app
docker-compose up -d app
```

### Cek Status Services

```bash
# Cek semua services
docker-compose ps

# Cek hanya app service
docker-compose ps app
```

## Tips

1. **Untuk Development**: Gunakan `docker-compose up -d --build app` untuk rebuild cepat
2. **Untuk Production**: Gunakan `docker-compose build --no-cache app` untuk fresh build
3. **Jika Ada Error**: Cek logs dengan `docker-compose logs app` atau `docker logs lostmediago_app`
4. **Quick Restart**: `docker-compose restart app` untuk restart cepat tanpa rebuild

