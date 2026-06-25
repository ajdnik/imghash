# Practical Examples

> Real-world examples and use cases for perceptual image hashing

## Overview

This guide provides practical, production-ready examples for common image hashing scenarios.

## Duplicate Detection

### Basic Duplicate Detection

Detect if two images are duplicates or near-duplicates:

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
)

func main() {
    // Create hasher (PDQ is recommended for duplicate detection)
    pdq, err := imghash.NewPDQ()
    if err != nil {
        panic(err)
    }

    // Hash both images
    h1, err := imghash.HashFile(pdq, "image1.jpg")
    if err != nil {
        panic(err)
    }

    h2, err := imghash.HashFile(pdq, "image2.jpg")
    if err != nil {
        panic(err)
    }

    // Compare hashes
    dist, err := pdq.Compare(h1, h2)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Distance: %v\n", dist)
    
    // Determine if duplicate
    if dist < 10 {
        fmt.Println("Images are duplicates or near-duplicates")
    } else {
        fmt.Println("Images are different")
    }
}
```

### Batch Duplicate Detection

Find all duplicates in a collection of images:

```go
package main

import (
    "fmt"
    "path/filepath"
    "github.com/ajdnik/imghash/v2"
)

type ImageHash struct {
    Path string
    Hash imghash.Hash
}

func findDuplicates(imagePaths []string, threshold float64) map[string][]string {
    pdq, _ := imghash.NewPDQ()
    
    // Compute all hashes
    imageHashes := make([]ImageHash, 0, len(imagePaths))
    for _, path := range imagePaths {
        hash, err := imghash.HashFile(pdq, path)
        if err != nil {
            fmt.Printf("Error hashing %s: %v\n", path, err)
            continue
        }
        imageHashes = append(imageHashes, ImageHash{Path: path, Hash: hash})
    }
    
    // Find duplicate groups
    duplicates := make(map[string][]string)
    visited := make(map[int]bool)
    
    for i := 0; i < len(imageHashes); i++ {
        if visited[i] {
            continue
        }
        
        group := []string{imageHashes[i].Path}
        visited[i] = true
        
        for j := i + 1; j < len(imageHashes); j++ {
            if visited[j] {
                continue
            }
            
            dist, _ := pdq.Compare(imageHashes[i].Hash, imageHashes[j].Hash)
            if float64(dist) <= threshold {
                group = append(group, imageHashes[j].Path)
                visited[j] = true
            }
        }
        
        if len(group) > 1 {
            duplicates[imageHashes[i].Path] = group
        }
    }
    
    return duplicates
}

func main() {
    images := []string{
        "photo1.jpg",
        "photo1_edited.jpg",
        "photo2.jpg",
        "photo2_compressed.jpg",
        "unique.jpg",
    }
    
    duplicates := findDuplicates(images, 10.0)
    
    fmt.Println("Duplicate groups found:")
    for original, group := range duplicates {
        fmt.Printf("\n%s:\n", filepath.Base(original))
        for _, dup := range group[1:] {
            fmt.Printf("  - %s\n", filepath.Base(dup))
        }
    }
}
```

### Deduplication Pipeline

Remove duplicates from a directory:

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/ajdnik/imghash/v2"
)

func deduplicateDirectory(dir string, threshold float64, dryRun bool) error {
    pdq, _ := imghash.NewPDQ()
    
    // Find all image files
    var imagePaths []string
    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            ext := filepath.Ext(path)
            if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
                imagePaths = append(imagePaths, path)
            }
        }
        return nil
    })
    if err != nil {
        return err
    }
    
    fmt.Printf("Found %d images\n", len(imagePaths))
    
    // Compute hashes
    type HashInfo struct {
        path string
        hash imghash.Hash
        size int64
    }
    
    hashes := make([]HashInfo, 0, len(imagePaths))
    for _, path := range imagePaths {
        hash, err := imghash.HashFile(pdq, path)
        if err != nil {
            fmt.Printf("Skipping %s: %v\n", path, err)
            continue
        }
        
        info, _ := os.Stat(path)
        hashes = append(hashes, HashInfo{
            path: path,
            hash: hash,
            size: info.Size(),
        })
    }
    
    // Find and remove duplicates (keep largest file)
    visited := make(map[int]bool)
    removed := 0
    
    for i := 0; i < len(hashes); i++ {
        if visited[i] {
            continue
        }
        
        visited[i] = true
        keeper := hashes[i]
        
        for j := i + 1; j < len(hashes); j++ {
            if visited[j] {
                continue
            }
            
            dist, _ := pdq.Compare(hashes[i].hash, hashes[j].hash)
            if float64(dist) <= threshold {
                visited[j] = true
                
                // Keep the larger file
                var toRemove HashInfo
                if hashes[j].size > keeper.size {
                    toRemove = keeper
                    keeper = hashes[j]
                } else {
                    toRemove = hashes[j]
                }
                
                fmt.Printf("Duplicate: %s (keeping %s)\n",
                    filepath.Base(toRemove.path),
                    filepath.Base(keeper.path))
                
                if !dryRun {
                    if err := os.Remove(toRemove.path); err != nil {
                        fmt.Printf("Error removing %s: %v\n", toRemove.path, err)
                    } else {
                        removed++
                    }
                }
            }
        }
    }
    
    if dryRun {
        fmt.Printf("\nDry run: would remove %d duplicates\n", removed)
    } else {
        fmt.Printf("\nRemoved %d duplicates\n", removed)
    }
    
    return nil
}

func main() {
    // Dry run first to preview
    deduplicateDirectory("./photos", 10.0, true)
    
    // Uncomment to actually remove duplicates
    // deduplicateDirectory("./photos", 10.0, false)
}
```

## Image Similarity Search

### Find Similar Images

Find the most similar images to a query image:

```go
package main

import (
    "fmt"
    "sort"
    "github.com/ajdnik/imghash/v2"
)

type SearchResult struct {
    Path     string
    Distance float64
}

func searchSimilar(queryPath string, databasePaths []string, topK int) []SearchResult {
    // Use GIST for scene-level similarity
    gist, _ := imghash.NewGIST()
    
    // Hash query image
    queryHash, err := imghash.HashFile(gist, queryPath)
    if err != nil {
        panic(err)
    }
    
    // Compute distances to all database images
    results := make([]SearchResult, 0, len(databasePaths))
    for _, path := range databasePaths {
        hash, err := imghash.HashFile(gist, path)
        if err != nil {
            fmt.Printf("Skipping %s: %v\n", path, err)
            continue
        }
        
        dist, _ := gist.Compare(queryHash, hash)
        results = append(results, SearchResult{
            Path:     path,
            Distance: float64(dist),
        })
    }
    
    // Sort by distance (ascending = most similar first)
    sort.Slice(results, func(i, j int) bool {
        return results[i].Distance < results[j].Distance
    })
    
    // Return top K results
    if len(results) > topK {
        results = results[:topK]
    }
    
    return results
}

func main() {
    database := []string{
        "images/beach1.jpg",
        "images/beach2.jpg",
        "images/mountain1.jpg",
        "images/city1.jpg",
        "images/sunset.jpg",
    }
    
    results := searchSimilar("query.jpg", database, 3)
    
    fmt.Println("Top 3 similar images:")
    for i, r := range results {
        fmt.Printf("%d. %s (distance: %.4f)\n", i+1, r.Path, r.Distance)
    }
}
```

### Reverse Image Search

Build an in-memory index for fast similarity search:

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
)

type ImageIndex struct {
    hasher imghash.HasherComparer
    index  map[string]imghash.Hash
}

func NewImageIndex() *ImageIndex {
    gist, _ := imghash.NewGIST()
    return &ImageIndex{
        hasher: gist,
        index:  make(map[string]imghash.Hash),
    }
}

func (idx *ImageIndex) Add(path string) error {
    hash, err := imghash.HashFile(idx.hasher, path)
    if err != nil {
        return err
    }
    idx.index[path] = hash
    return nil
}

func (idx *ImageIndex) AddBatch(paths []string) error {
    for _, path := range paths {
        if err := idx.Add(path); err != nil {
            fmt.Printf("Warning: failed to index %s: %v\n", path, err)
        }
    }
    return nil
}

func (idx *ImageIndex) Search(queryPath string, maxResults int, maxDistance float64) []SearchResult {
    queryHash, err := imghash.HashFile(idx.hasher, queryPath)
    if err != nil {
        return nil
    }
    
    results := make([]SearchResult, 0)
    for path, hash := range idx.index {
        dist, _ := idx.hasher.Compare(queryHash, hash)
        distFloat := float64(dist)
        
        if distFloat <= maxDistance {
            results = append(results, SearchResult{
                Path:     path,
                Distance: distFloat,
            })
        }
    }
    
    sort.Slice(results, func(i, j int) bool {
        return results[i].Distance < results[j].Distance
    })
    
    if len(results) > maxResults {
        results = results[:maxResults]
    }
    
    return results
}

func main() {
    // Build index
    index := NewImageIndex()
    
    images := []string{
        "db/beach1.jpg", "db/beach2.jpg", "db/mountain1.jpg",
        "db/city1.jpg", "db/sunset1.jpg", "db/forest1.jpg",
    }
    
    fmt.Println("Indexing images...")
    index.AddBatch(images)
    fmt.Printf("Indexed %d images\n\n", len(index.index))
    
    // Search
    results := index.Search("query_beach.jpg", 5, 0.5)
    
    fmt.Println("Search results:")
    for i, r := range results {
        fmt.Printf("%d. %s (distance: %.4f)\n", i+1, r.Path, r.Distance)
    }
}
```

## Content Moderation

### Blocklist Matching

Check if an uploaded image matches known inappropriate content:

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
)

type ContentModerator struct {
    hasher    imghash.HasherComparer
    blocklist map[string]imghash.Hash
    threshold float64
}

func NewContentModerator(threshold float64) *ContentModerator {
    pdq, _ := imghash.NewPDQ()
    return &ContentModerator{
        hasher:    pdq,
        blocklist: make(map[string]imghash.Hash),
        threshold: threshold,
    }
}

func (cm *ContentModerator) AddToBlocklist(id, path string) error {
    hash, err := imghash.HashFile(cm.hasher, path)
    if err != nil {
        return err
    }
    cm.blocklist[id] = hash
    return nil
}

func (cm *ContentModerator) Check(imagePath string) (bool, string, error) {
    hash, err := imghash.HashFile(cm.hasher, imagePath)
    if err != nil {
        return false, "", err
    }
    
    for id, blockHash := range cm.blocklist {
        dist, err := cm.hasher.Compare(hash, blockHash)
        if err != nil {
            continue
        }
        
        if float64(dist) <= cm.threshold {
            return true, id, nil
        }
    }
    
    return false, "", nil
}

func main() {
    moderator := NewContentModerator(10.0)
    
    // Add known bad content to blocklist
    moderator.AddToBlocklist("harmful-001", "blocklist/image1.jpg")
    moderator.AddToBlocklist("harmful-002", "blocklist/image2.jpg")
    
    // Check new upload
    blocked, matchID, err := moderator.Check("uploads/new_image.jpg")
    if err != nil {
        panic(err)
    }
    
    if blocked {
        fmt.Printf("⛔ Content blocked - matches %s\n", matchID)
    } else {
        fmt.Println("✅ Content approved")
    }
}
```

### Adversarial-Robust Verification with DINOHash

PDQ pre-filters at sub-millisecond throughput; DINOHash confirms survivors with deep semantic features that resist crops, recolors, and adversarial perturbations:

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
    "github.com/ajdnik/imghash/dinoweights"
)

type TwoStageModerator struct {
    pdq        imghash.HasherComparer
    dino       imghash.HasherComparer
    pdqBlock   map[string]imghash.Binary
    dinoBlock  map[string]imghash.Binary
    pdqLoose   int
    dinoStrict int
}

func NewTwoStageModerator() *TwoStageModerator {
    pdq, _ := imghash.NewPDQ()
    dn, _ := imghash.NewDINOHash(imghash.WithSafetensorsBlob(dinoweights.Blob))
    return &TwoStageModerator{
        pdq:        pdq,
        dino:       dn,
        pdqBlock:   make(map[string]imghash.Binary),
        dinoBlock:  make(map[string]imghash.Binary),
        pdqLoose:   100, // generous pre-filter
        dinoStrict: 25,  // semantic confirmation
    }
}

func (m *TwoStageModerator) AddToBlocklist(id, path string) error {
    ph, err := imghash.HashFile(m.pdq, path)
    if err != nil {
        return err
    }
    dh, err := imghash.HashFile(m.dino, path)
    if err != nil {
        return err
    }
    m.pdqBlock[id] = ph.(imghash.Binary)
    m.dinoBlock[id] = dh.(imghash.Binary)
    return nil
}

func (m *TwoStageModerator) Check(path string) (bool, string, error) {
    candPDQ, err := imghash.HashFile(m.pdq, path)
    if err != nil {
        return false, "", err
    }

    // Stage 1: cheap PDQ pre-filter to narrow candidates.
    var dinoCandidates []string
    for id, ph := range m.pdqBlock {
        dist, _ := m.pdq.Compare(candPDQ, ph)
        if int(dist) <= m.pdqLoose {
            dinoCandidates = append(dinoCandidates, id)
        }
    }
    if len(dinoCandidates) == 0 {
        return false, "", nil
    }

    // Stage 2: DINOHash on survivors only.
    candDINO, err := imghash.HashFile(m.dino, path)
    if err != nil {
        return false, "", err
    }
    for _, id := range dinoCandidates {
        dist, _ := m.dino.Compare(candDINO, m.dinoBlock[id])
        if int(dist) <= m.dinoStrict {
            return true, id, nil
        }
    }
    return false, "", nil
}

func main() {
    mod := NewTwoStageModerator()
    mod.AddToBlocklist("harmful-001", "blocklist/image1.jpg")

    blocked, matchID, err := mod.Check("uploads/new_image.jpg")
    if err != nil {
        panic(err)
    }
    if blocked {
        fmt.Printf("⛔ Content blocked - matches %s\n", matchID)
    } else {
        fmt.Println("✅ Content approved")
    }
}
```

> [!TIP]
> The two-stage layout keeps throughput close to PDQ on the bulk of traffic while catching adversarial near-duplicates the classical hash misses. Loosen the PDQ threshold and tighten the DINOHash threshold to trade recall for precision.

## Working with Different Image Sources

### Reading from File

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
)

func main() {
    avg, err := imghash.NewAverage()
    if err != nil {
        panic(err)
    }
    
    hash, err := imghash.HashFile(avg, "image.jpg")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(hash)
}
```

### Reading from io.Reader

```go
package main

import (
    "fmt"
    "os"
    "github.com/ajdnik/imghash/v2"
)

func main() {
    f, err := os.Open("image.jpg")
    if err != nil {
        panic(err)
    }
    defer f.Close()
    
    avg, _ := imghash.NewAverage()
    hash, err := imghash.HashReader(avg, f)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(hash)
}
```

### Hashing HTTP Response

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/ajdnik/imghash/v2"
)

func main() {
    resp, err := http.Get("https://example.com/image.jpg")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    avg, _ := imghash.NewAverage()
    hash, err := imghash.HashReader(avg, resp.Body)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(hash)
}
```

### Hashing image.Image

```go
package main

import (
    "fmt"
    "image"
    _ "image/jpeg"
    "os"
    "github.com/ajdnik/imghash/v2"
)

func main() {
    // Open and decode image
    f, _ := os.Open("image.jpg")
    defer f.Close()
    
    img, _, err := image.Decode(f)
    if err != nil {
        panic(err)
    }
    
    // Hash the image.Image directly
    avg, _ := imghash.NewAverage()
    hash, err := avg.Calculate(img)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(hash)
}
```

## Advanced Techniques

### Multi-Algorithm Consensus

Use multiple algorithms and require agreement:

```go
package main

import (
    "fmt"
    "github.com/ajdnik/imghash/v2"
)

func robustDuplicateCheck(img1, img2 string) bool {
    type AlgorithmCheck struct {
        name      string
        threshold float64
        match     bool
    }
    
    checks := []AlgorithmCheck{
        {name: "Average", threshold: 5.0},
        {name: "PDQ", threshold: 10.0},
        {name: "ColorMoment", threshold: 15.0},
    }
    
    // Average
    avg, _ := imghash.NewAverage()
    h1a, _ := imghash.HashFile(avg, img1)
    h2a, _ := imghash.HashFile(avg, img2)
    distAvg, _ := avg.Compare(h1a, h2a)
    checks[0].match = float64(distAvg) <= checks[0].threshold
    
    // PDQ
    pdq, _ := imghash.NewPDQ()
    h1p, _ := imghash.HashFile(pdq, img1)
    h2p, _ := imghash.HashFile(pdq, img2)
    distPDQ, _ := pdq.Compare(h1p, h2p)
    checks[1].match = float64(distPDQ) <= checks[1].threshold
    
    // ColorMoment
    cm, _ := imghash.NewColorMoment()
    h1c, _ := imghash.HashFile(cm, img1)
    h2c, _ := imghash.HashFile(cm, img2)
    distCM, _ := cm.Compare(h1c, h2c)
    checks[2].match = float64(distCM) <= checks[2].threshold
    
    // Require at least 2 out of 3 to agree
    votes := 0
    for _, check := range checks {
        fmt.Printf("%s: %v\n", check.name, check.match)
        if check.match {
            votes++
        }
    }
    
    return votes >= 2
}

func main() {
    if robustDuplicateCheck("img1.jpg", "img2.jpg") {
        fmt.Println("\n✅ Consensus: Images are duplicates")
    } else {
        fmt.Println("\n❌ Consensus: Images are different")
    }
}
```

### Hash Persistence

Save and load hashes for faster subsequent comparisons:

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "github.com/ajdnik/imghash/v2"
)

type HashDatabase struct {
    Hashes map[string][]byte `json:"hashes"`
}

func (db *HashDatabase) Save(path string) error {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()
    
    return json.NewEncoder(f).Encode(db)
}

func LoadHashDatabase(path string) (*HashDatabase, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    
    var db HashDatabase
    err = json.NewDecoder(f).Decode(&db)
    return &db, err
}

func main() {
    pdq, _ := imghash.NewPDQ()
    
    // Create and populate database
    db := &HashDatabase{
        Hashes: make(map[string][]byte),
    }
    
    images := []string{"img1.jpg", "img2.jpg", "img3.jpg"}
    for _, img := range images {
        hash, _ := imghash.HashFile(pdq, img)
        // Binary hashes can be cast to []byte for storage
        if binHash, ok := hash.(imghash.Binary); ok {
            db.Hashes[img] = []byte(binHash)
        }
    }
    
    // Save to disk
    db.Save("hashes.json")
    fmt.Println("Saved hash database")
    
    // Load from disk
    loaded, _ := LoadHashDatabase("hashes.json")
    fmt.Printf("Loaded %d hashes\n", len(loaded.Hashes))
    
    // Use loaded hashes
    for path, hashBytes := range loaded.Hashes {
        fmt.Printf("%s: %v\n", path, imghash.Binary(hashBytes))
    }
}
```

## Error Handling

### Comprehensive Error Handling

```go
package main

import (
    "errors"
    "fmt"
    "os"
    "github.com/ajdnik/imghash/v2"
)

func hashImageSafely(path string) (imghash.Hash, error) {
    // Check file exists
    if _, err := os.Stat(path); err != nil {
        if os.IsNotExist(err) {
            return nil, fmt.Errorf("image file not found: %s", path)
        }
        return nil, fmt.Errorf("error accessing file: %w", err)
    }
    
    // Create hasher
    pdq, err := imghash.NewPDQ()
    if err != nil {
        return nil, fmt.Errorf("failed to create hasher: %w", err)
    }
    
    // Compute hash
    hash, err := imghash.HashFile(pdq, path)
    if err != nil {
        return nil, fmt.Errorf("failed to hash image: %w", err)
    }
    
    return hash, nil
}

func compareImagesSafely(path1, path2 string) error {
    pdq, _ := imghash.NewPDQ()
    
    h1, err := hashImageSafely(path1)
    if err != nil {
        return err
    }
    
    h2, err := hashImageSafely(path2)
    if err != nil {
        return err
    }
    
    dist, err := pdq.Compare(h1, h2)
    if err != nil {
        if errors.Is(err, imghash.ErrIncompatibleHash) {
            return fmt.Errorf("hashes are incompatible types")
        }
        return fmt.Errorf("comparison failed: %w", err)
    }
    
    fmt.Printf("Distance: %v\n", dist)
    return nil
}

func main() {
    if err := compareImagesSafely("img1.jpg", "img2.jpg"); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}
```

## Next Steps

- **[API Reference](hasher)** — Explore the complete API documentation
- **[Algorithm Guide](choosing-algorithm)** — Learn how to choose the right algorithm
