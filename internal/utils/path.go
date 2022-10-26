package utils

import "path/filepath"

func getStoragePath(elem ...string) string {
	return filepath.Join("storage", filepath.Join(elem...))
}

func GetFlagPath(name string) string {
	return getStoragePath("flag", name)
}

func GetQrsavePath(name string) string {
	return getStoragePath("qrsave", name)
}

func GetAssetsPath(name string) string {
	return filepath.Join("assets", name)
}
