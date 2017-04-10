# Go Git Wrapper

Simple wrapper to GIT commands.

## Depencencies

This nifty little wrapper uses [GO-GIT](https://github.com/src-d/go-git/) by src-d.

## Example

```go
// Start with 0.0.1
versionStruct = versioning.CreateVersion("0.0.1")
fmt.Printf("%s\n", versionStruct.GetVersionString())

// Make it 0.0.2
versionStruct.Patch()
fmt.Printf("%s\n", versionStruct.GetVersionString())

// Get TAG (for GIT Tagging)
fmt.Printf("%s\n", versionStruct.GetVersionTag())
```
