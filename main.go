package main

import (
    "archive/tar"
    "compress/gzip"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"

    "github.com/cheggaaa/pb/v3"
    "gopkg.in/yaml.v3"
)

type Config struct {
    OutputDir  string   `yaml:"output_dir"`
    SourceDirs []string `yaml:"source_dirs"`
}

func main() {
    // Parse command-line arguments
    flag.Parse()

    // Load configuration from file
    var config Config
    configFile, err := os.Open("config.yaml")
    if err != nil {
        log.Fatal(err)
    }
    defer configFile.Close()
    decoder := yaml.NewDecoder(configFile)
    if err = decoder.Decode(&config); err != nil {
        log.Fatal(err)
    }

    // Check that output_dir is specified
    if config.OutputDir == "" {
        log.Fatal("Error: output_dir not specified in configuration file")
    }

    // Define destination directory
    destDir := config.OutputDir

    // Create destination directory if it doesn't exist
    os.MkdirAll(destDir, 0755)

    for _, sourceDir := range config.SourceDirs {
        // Print source directory being backed up
        fmt.Printf("Backing up %s...\n", sourceDir)

        // Define archive name
        archiveName := fmt.Sprintf("%s_%s.tar.gz", filepath.Base(sourceDir), time.Now().Format("2006-01-02_15-04-05"))

        // Create destination file for writing
        destFile, err := os.Create(filepath.Join(destDir, archiveName))
        if err != nil {
            log.Fatal(err)
        }
        defer destFile.Close()

        // Create gzip writer
        gzipWriter := gzip.NewWriter(destFile)
        defer gzipWriter.Close()

        // Create tar writer
        tarWriter := tar.NewWriter(gzipWriter)
        defer tarWriter.Close()

        // Count number of files in source directory
        fileCount := 0
        filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
            if !info.IsDir() && filepath.Base(path)[0] != '.' {
                fileCount++
            }
            return nil
        })

        // Create progress bar
        bar := pb.Full.Start(fileCount)
        bar.SetTemplateString(`{{ green "Backup Progress:" }} {{ bar . "[" "=" ">" "-" "]"}} {{ percent . }} {{speed .}} {{etime .}} / {{rtime .}}`)
        bar.SetWidth(80)

        // Walk through source directory recursively
        filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
            if err != nil || info.IsDir() || filepath.Base(path)[0] == '.' {
                return err
            }

            // Create tar header
            header, err := tar.FileInfoHeader(info, "")
            if err != nil {
                return err
            }
            header.Name = path[len(sourceDir)+1:]

            // Write header to tar archive
            if err = tarWriter.WriteHeader(header); err != nil {
                return err
            }

            // Open source file for reading
            sourceFile, err := os.Open(path)
            if err != nil {
                return err
            }
            defer sourceFile.Close()

            // Copy source file contents to tar archive
            if _, err = io.Copy(tarWriter, sourceFile); err != nil {
                return err
            }

            // Increment progress bar
            bar.Increment()

            return nil
        })

        // Finish progress bar
        bar.Finish()

        // Print success message
        fmt.Printf("Backup created successfully! Archive saved to %s\n", filepath.Join(destDir, archiveName))
    }
}
