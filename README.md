# envsync

> Securely sync and diff `.env` files across environments using encrypted remote backends.

---

## Installation

```bash
go install github.com/yourusername/envsync@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envsync.git && cd envsync && go build -o envsync .
```

---

## Usage

Initialize a new envsync configuration:

```bash
envsync init --backend s3 --bucket my-env-bucket
```

Push your local `.env` file to the remote backend:

```bash
envsync push --env production --file .env
```

Pull and decrypt the latest `.env` from remote:

```bash
envsync pull --env production --output .env
```

Diff your local `.env` against the remote version:

```bash
envsync diff --env production --file .env
```

All values are encrypted before leaving your machine using AES-256-GCM. Encryption keys are managed locally via a `.envsync.key` file or an environment variable `ENVSYNC_KEY`.

---

## Supported Backends

- AWS S3
- Google Cloud Storage
- Azure Blob Storage
- Local filesystem (for testing)

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

---

## License

[MIT](LICENSE) © 2024 yourusername