package services

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/joaquinrovira/infiltra2-returns/app/model"
	"github.com/samber/do/v2"
	"golang.org/x/net/html"
)

type RandomWordService struct {
	ctx        context.Context
	randomWord chan model.Word
}

func NewRandomWordService(svc do.Injector) (*RandomWordService, error) {
	service := &RandomWordService{
		ctx:        do.MustInvoke[context.Context](svc),
		randomWord: make(chan model.Word, 5),
	}
	go service.background()
	return service, nil
}

func (random *RandomWordService) Next() model.Word {
	return <-random.randomWord
}

const NoDefinitionAttempts = 10

func (random *RandomWordService) background() {
	noDefinitionsCountdown := NoDefinitionAttempts
	for random.ctx.Err() == nil {
		if word, err := fetchRandomWord(); err != nil {
			log.Printf("[WARN] unable to fetch random word - %v", err)
		} else {
			log.Printf("fetched random word '%v' with %d definitions", word.Word, len(word.Description))
			if len(word.Description) == 0 && noDefinitionsCountdown > 0 {
				log.Printf("[INFO] ignoring '%v' due to a lack of descriptions", word.Word)

				noDefinitionsCountdown -= 1
				if noDefinitionsCountdown == 0 {
					log.Print("[WARN] word retry limit reached, using next word unconditionally", word.Word)
				}

				continue
			}
			random.randomWord <- word
			noDefinitionsCountdown = NoDefinitionAttempts
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func fetchRandomWord() (model.Word, error) {
	word, err := fetchRandomWordFromFile()
	if err != nil {
		return model.Word{}, err
	}
	meaning := tryFetchMeaning(word)
	return model.Word{Word: word, Description: meaning}, nil
}

func fetchRandomWordFromFile() (string, error) {
	// https://stackoverflow.com/questions/60093768/reading-a-random-line-from-a-file-in-constant-time-in-go
	file, err := os.Open("static/es.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	randsource := rand.NewSource(time.Now().UnixNano())
	randgenerator := rand.New(randsource)

	lineNum := 1
	var pick string
	for scanner.Scan() {
		line := scanner.Text()
		roll := randgenerator.Intn(lineNum)
		if roll == 0 {
			pick = line
		}
		lineNum += 1
	}
	return pick, nil
}

func tryFetchMeaning(word string) []model.WordDescription {
	// Try to get meaning from wordreference API
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		log.Printf("[WARN] %v", err)
	}
	client := &http.Client{Jar: jar}
	url, err := url.Parse("https://www.wordreference.com/definicion/" + word)
	if err != nil {
		log.Printf("[WARN] %v", err)
	}
	client.Jar.SetCookies(url, []*http.Cookie{{Name: "llang", Domain: url.Host, Path: "/", Secure: true, Value: "esesi"}})

	request := http.Request{
		Method: "GET",
		URL:    url,
		Header: map[string][]string{
			"User-Agent": {"Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0"},
			"Accept":     {"text/html"},
			"Cookie":     {"llang=esesi"},
		},
	}
	res, err := client.Do(&request)
	if err != nil {
		log.Printf("[WARN] %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	definitions := doc.Find("#article #otherDicts .small1+div ol li")
	numDescription := definitions.Length()
	descriptions := make([]model.WordDescription, 0, numDescription)

	definitions.Each(func(i int, s *goquery.Selection) {
		description, err := Description(s)
		if err != nil {
			return
		}
		if description == "" {
			log.Printf("[WARN] lost description '%d' for word '%s'", i, word)
		}
		descriptions = append(descriptions, description)
	})

	return descriptions
}

func Description(s *goquery.Selection) (string, error) {
	node := s.Nodes[0].FirstChild
	for node != nil {
		n := node
		data := strings.TrimSpace(n.Data)
		node = node.NextSibling

		if n.Type != html.TextNode {
			continue
		}

		if data == "tr." || data == "f." || len(data) < 15 {
			continue
		}

		return strings.TrimSuffix(data, ":"), nil
	}
	return "", fmt.Errorf("unable to find description")
}
