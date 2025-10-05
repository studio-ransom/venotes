# S3 Storage Configuration

Venotes supports both local file storage and S3-compatible cloud storage for uploaded files.

## Environment Variables

### Storage Type
Set `STORAGE_TYPE` to choose your storage backend:
- `local` - Local file system storage (default)
- `s3` - S3-compatible cloud storage

### Local Storage (STORAGE_TYPE=local)
```bash
STORAGE_TYPE=local
STORAGE_LOCAL_PATH=data/uploads  # Optional, defaults to data/uploads
```

### S3 Storage (STORAGE_TYPE=s3)

#### Required Variables
```bash
STORAGE_TYPE=s3
S3_ACCESS_KEY=your_access_key_here
S3_SECRET_KEY=your_secret_key_here
S3_BUCKET=your_bucket_name_here
```

#### Optional Variables
```bash
S3_ENDPOINT=https://s3.amazonaws.com  # Default AWS S3 endpoint
S3_REGION=us-east-1                   # Default region
S3_BASE_PATH=uploads                  # Optional base path within bucket
```

## S3-Compatible Services

### AWS S3
```bash
STORAGE_TYPE=s3
S3_ACCESS_KEY=AKIA...
S3_SECRET_KEY=...
S3_BUCKET=my-bucket
S3_REGION=us-east-1
# S3_ENDPOINT defaults to https://s3.amazonaws.com
```

### MinIO
```bash
STORAGE_TYPE=s3
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=my-bucket
S3_ENDPOINT=http://localhost:9000
S3_REGION=us-east-1
```

### DigitalOcean Spaces
```bash
STORAGE_TYPE=s3
S3_ACCESS_KEY=DO...
S3_SECRET_KEY=...
S3_BUCKET=my-space
S3_ENDPOINT=https://nyc3.digitaloceanspaces.com
S3_REGION=nyc3
```

### Cloudflare R2
```bash
STORAGE_TYPE=s3
S3_ACCESS_KEY=...
S3_SECRET_KEY=...
S3_BUCKET=my-bucket
S3_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com
S3_REGION=auto
```

### Backblaze B2
```bash
STORAGE_TYPE=s3
S3_ACCESS_KEY=...
S3_SECRET_KEY=...
S3_BUCKET=my-bucket
S3_ENDPOINT=https://s3.us-west-000.backblazeb2.com
S3_REGION=us-west-000
```

## File Deduplication

The application uses SHA256 hashing for file deduplication:
- Files with identical content are stored only once
- Multiple database entries can reference the same file
- Files appear in logs every time they're uploaded
- Storage space is optimized by avoiding duplicates

## Migration

To migrate from local storage to S3:
1. Set up your S3-compatible storage
2. Configure the environment variables
3. Restart the application
4. New uploads will use S3 storage
5. Existing local files remain accessible until migrated

## Security

- Store credentials securely using environment variables
- Use IAM policies to limit S3 access
- Consider using temporary credentials for production
- Enable S3 bucket encryption
