# Migration Guide: v1 to v2

This guide covers the breaking changes introduced in v2 and what you need to
change in your code when migrating from v1.

## 1. Update the dependency

```bash
go get github.com/ajdnik/imghash/v2@latest
go mod tidy
```

## 2. Replace `New*WithParams` constructors with functional options

v1 constructors with explicit parameter lists were removed.

| v1 | v2 |
|---|---|
| `NewAverageWithParams(w, h, resize)` | `NewAverage(WithSize(w, h), WithInterpolation(resize))` |
| `NewDifferenceWithParams(w, h, resize)` | `NewDifference(WithSize(w, h), WithInterpolation(resize))` |
| `NewMedianWithParams(w, h, resize)` | `NewMedian(WithSize(w, h), WithInterpolation(resize))` |
| `NewPHashWithParams(w, h, resize)` | `NewPHash(WithSize(w, h), WithInterpolation(resize))` |
| `NewColorMomentWithParams(w, h, resize, kernel, sigma)` | `NewColorMoment(WithSize(w, h), WithInterpolation(resize), WithKernelSize(kernel), WithSigma(sigma))` |
| `NewMarrHildrethWithParams(scale, alpha, w, h, resize, kernel, sigma)` | `NewMarrHildreth(WithScale(scale), WithAlpha(alpha), WithSize(w, h), WithInterpolation(resize), WithKernelSize(kernel), WithSigma(sigma))` |
| `NewBlockMeanWithParams(w, h, resize, bw, bh, method)` | `NewBlockMean(WithSize(w, h), WithInterpolation(resize), WithBlockSize(bw, bh), WithBlockMeanMethod(method))` |
| `NewRadialVarianceWithParams(sigma, angles)` | `NewRadialVariance(WithSigma(sigma), WithAngles(angles))` |

## 3. Handle `New*` constructor errors

All `New*` constructors now return `(Type, error)` and validate their options.
Invalid configurations (e.g. zero dimensions) are rejected at construction time
instead of panicking during `Calculate`.

```go
// v1 / early v2
avg := imghash.NewAverage(imghash.WithSize(16, 16))

// current
avg, err := imghash.NewAverage(imghash.WithSize(16, 16))
if err != nil {
	return err
}
```

Sentinel errors for programmatic checking:

| Error | Meaning |
|-------|---------|
| `ErrInvalidSize` | width or height is zero |
| `ErrInvalidBlockSize` | block width or height is zero |
| `ErrInvalidAngles` | angles is not positive |
| `ErrInvalidKernelSize` | kernel size is not positive |
| `ErrInvalidScale` | scale is not positive |
| `ErrInvalidAlpha` | alpha is not positive |
| `ErrInvalidSigma` | sigma is negative |

## 4. Update interpolation constants and imports

`imgproc` is no longer part of the public API surface.

- v1: `github.com/ajdnik/imghash/imgproc`
- v2: use `imghash` interpolation constants directly (`imghash.Bilinear`,
  `imghash.Bicubic`, `imghash.BilinearExact`, ...).

Example:

```go
// v1
// import "github.com/ajdnik/imghash/imgproc"
// hash := imghash.NewAverageWithParams(16, 16, imgproc.Bilinear)

// v2
hash, err := imghash.NewAverage(
	imghash.WithSize(16, 16),
	imghash.WithInterpolation(imghash.Bilinear),
)
```

## 5. Handle `Calculate` errors and generic hash return type

All hashers now implement:

```go
Calculate(image.Image) (hashtype.Hash, error)
```

In v1, `Calculate` returned a concrete hash type and no error.

```go
// v1
h1 := ph.Calculate(img1)
h2 := ph.Calculate(img2)
d := similarity.Hamming(h1, h2)

// v2
h1, err := ph.Calculate(img1)
if err != nil {
	return err
}
h2, err := ph.Calculate(img2)
if err != nil {
	return err
}
d, err := imghash.Compare(h1, h2)
if err != nil {
	return err
}
_ = d
```

Use type assertions only when you need hash-type-specific behavior:

```go
bin := h1.(imghash.Binary)
```

## 6. Update similarity API usage

`similarity` functions are now generic over `hashtype.Hash`.

| v1 | v2 |
|---|---|
| `similarity.Hamming(b1, b2)` | `similarity.Hamming(h1, h2)` (returns `(Distance, error)`) |
| `similarity.L2Float64(f1, f2)` | `similarity.L2(h1, h2)` |
| `similarity.L2UInt8(u1, u2)` | `similarity.L2(h1, h2)` |
| `similarity.PCCFloat64(f1, f2)` | `similarity.PCC(h1, h2)` |
| `similarity.PCCUInt8(u1, u2)` | `similarity.PCC(h1, h2)` |

For most callers, the new top-level helper is simpler:

```go
dist, err := imghash.Compare(h1, h2)
```

## 7. Optional: adopt new convenience helpers

v2 adds:

- `OpenImage(path)`
- `DecodeImage(r)`
- `HashFile(hasher, path)`
- `HashReader(hasher, r)`
- `Compare(h1, h2)`

These remove most decoding and comparison boilerplate from v1-style code.
