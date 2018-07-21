package main

import (
	"flag"
	"fmt"
	"github.com/ahmdrz/goinsta"
	"image"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"sort"
	"time"
)

var (
	subreddit = flag.String("sub", "ik_ihe", "The Subreddit to pull from")
	username  = flag.String("username", "ik_ihe_bot", "Instagram Username")
	password  = flag.String("password", "gk7ax84w", "Instagram Password")
	caption   = flag.String("caption", "ik_ihe", "The post caption")
)

func init() {
	rand.Seed(time.Now().Unix())
	flag.Parse()
}

func main() {
	if err := DoPost(); err != nil {
		log.Fatal(err)
	}
}

func DoPost() error {
	ss, err := GetSubmissions(*subreddit)
	if err != nil {
		return err
	}
	sort.Sort(ByScore(ss))
	for _, s := range ss {
		used, err := IsUsed(s.Id)
		if err != nil {
			return err
		}
		if used {
			continue
		}
		im, err := MakeImage(s.Title)
		if err != nil {
			return err
		}
		if err := MarkUsed(s.Id, s.Title); err != nil {
			return err
		}
		fmt.Printf("%d: %s\n", s.Score, s.Title)
		if err := SaveJPEG(im, "out.jpeg"); err != nil {
			return err
		}
		return PostImage("out.jpeg")
	}
	return fmt.Errorf("all %d submissions are used", len(ss))
}

func SaveJPEG(m image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return jpeg.Encode(f, m, &jpeg.Options{jpeg.DefaultQuality})
}

func PostImage(imgpath string) error {
	insta := goinsta.New(*username, *password)
	if err := insta.Login(); err != nil {
		return fmt.Errorf("failed to login: %s", err)
	}
	defer insta.Logout()
	f, err := os.Open(imgpath)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := insta.UploadPhoto(f, *caption, 87, 0); err != nil {
		return fmt.Errorf("failed to upload:", err)
	}
	return nil
}
