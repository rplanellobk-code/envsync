// Package storage provides the Backend interface and concrete implementations
// for persisting encrypted .env snapshots.
//
// The Backend interface abstracts over different storage mechanisms (local
// filesystem, S3, etc.), allowing envsync to read and write opaque byte
// payloads without knowing the encryption details.
//
// Current implementations:
//
//   - FileBackend: stores each environment as an encrypted file on the local
//     filesystem, suitable for single-machine use or version-controlled
//     directories.
package storage
