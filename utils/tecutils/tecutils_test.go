package tecutils

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestUUIDNoDuplicates(t *testing.T) {
	iterations := 100
	log.Printf("Processing %v iterations", iterations)
	values := make(map[string]bool)
	dups := 0

	for i := 0; i < iterations; i++ {
		v := UUID()
		if _, ok := values[v]; ok {
			dups++
		}
		values[v] = true
	}
	if dups != 0 {
		t.Error()
	}
}

func TestEncrypt(t *testing.T) {
	pwd := UUID()
	encripted := Encrypt(pwd)
	log.Printf("Original string is %v long. Encripted string is %v", len(pwd), len(encripted))
	if encripted == pwd {
		t.Error()
	}
}

func TestUrlParser(t *testing.T) {
	url := "https://localhost:8080/mau/testing?something=ok"
	parsed, err := ParseBaseUrl(url)
	log.Println(parsed)
	if err != nil {
		t.Error(err)
	}
	if parsed != "https://localhost:8080" {
		t.Error()
	}
}

func TestSubdirectoriesInfo(t *testing.T) {
	t.Skip("TODO: make sure it can run in CI environment")
	downloads := filepath.Join(os.Getenv("HOME"), "Downloads")
	dirs, err := SubdirectoriesInfo(downloads)
	if err != nil {
		t.Error(err)
		return
	}
	if len(dirs) == 0 {
		t.Error()
		return
	}
	counter := 0
	for _, d := range dirs {
		counter++
		log.Printf("%#v\n", d)
		if counter > 5 {
			log.Println("There are more directories...")
			break
		}
	}
}

func TestNumberFormats(t *testing.T) {
	amount := 3512834.746
	res := FormatCurrency(amount)
	if res != "$3,512,834.75" {
		t.Error()
	}
}
