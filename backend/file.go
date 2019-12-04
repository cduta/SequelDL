package backend

import (
  "os"
  "path/filepath"
)

func ToAbsolutePath(relativePath string) (string, error) {
  var (
    err              error 
    absoluteRootPath string 
  )

  absoluteRootPath, err = os.Getwd()
  if err != nil {
    return "", err 
  }   

  return filepath.Join(absoluteRootPath, relativePath), err
}