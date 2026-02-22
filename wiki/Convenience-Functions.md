# Convenience Functions

The library provides helpers so you do not need image-decoding boilerplate.

| Function | Description |
|----------|-------------|
| `OpenImage(path)` | Reads and decodes an image file (JPEG, PNG, GIF) |
| `DecodeImage(r)` | Decodes an image from any `io.Reader` |
| `HashFile(hasher, path)` | Opens a file and computes its hash in one call |
| `HashReader(hasher, r)` | Decodes from a reader and computes the hash |
| `Compare(h1, h2)` | Computes distance using the natural metric for the hash type |

Use the algorithm's `Compare` method for its recommended metric, or call top-level `imghash.Compare(h1, h2)` for generic type-based comparison.
