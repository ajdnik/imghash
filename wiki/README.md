# Wiki Source Files

These files mirror the project GitHub Wiki pages.

## Publish to GitHub Wiki

From the repo root:

```sh
# Clone wiki repo next to this project (first time only)
git clone git@github.com:ajdnik/imghash.wiki.git ../imghash.wiki

# Copy pages
cp wiki/*.md ../imghash.wiki/

# Commit and push
cd ../imghash.wiki
git add *.md
git commit -m "docs: sync wiki pages from repo"
git push
```

Notes:

- `Home.md` controls the wiki landing page.
- File names map to page titles (for example, `Similarity-Metrics.md` -> `Similarity Metrics`).
