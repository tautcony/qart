package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"gopkg.in/ini.v1"
	"log"
	"path"
)

type langType struct {
	Lang, Name string
}

var (
	langTags    []language.Tag
	i18nStrings map[string]map[string]string // lang -> key -> value
)

func init() {
	initLocales()
}

func initLocales() {
	// Available languages
	availableLangs := []string{"en", "zh-CN", "zh-TW", "ja"}
	
	langTags = make([]language.Tag, 0, len(availableLangs))
	i18nStrings = make(map[string]map[string]string)
	
	for _, name := range availableLangs {
		l := language.Make(name)
		langTags = append(langTags, l)
		
		// Load locale file
		localeFile := path.Join("conf", "locale", fmt.Sprintf("locale_%s.ini", name))
		cfg, err := ini.Load(localeFile)
		if err != nil {
			log.Printf("Failed to load locale file %s: %v", localeFile, err)
			i18nStrings[name] = make(map[string]string)
			continue
		}
		
		i18nStrings[name] = make(map[string]string)
		for _, section := range cfg.Sections() {
			sectionName := section.Name()
			if sectionName == "DEFAULT" {
				sectionName = ""
			}
			for _, key := range section.Keys() {
				fullKey := key.Name()
				if sectionName != "" {
					fullKey = sectionName + "." + key.Name()
				}
				i18nStrings[name][fullKey] = key.Value()
			}
		}
		log.Printf("Loaded language: %v[%v]", display.Self.Name(l), l.String())
	}
}

func I18n() gin.HandlerFunc {
	return func(c *gin.Context) {
		matcher := language.NewMatcher(langTags)
		
		urlLang := c.Query("lang")
		cookieLang, _ := c.Cookie("lang")
		accept := c.GetHeader("Accept-Language")
		
		curLang, _ := language.MatchStrings(matcher, urlLang, cookieLang, accept)
		
		// Save language in cookie if needed
		if cookieLang == "" || cookieLang != curLang.String() {
			c.SetCookie("lang", curLang.String(), 1<<31-1, "/", "", false, false)
		}
		
		// Set language in context
		c.Set("lang", curLang.String())
		
		// Prepare language data for templates
		restLangs := make([]*langType, 0, len(langTags)-1)
		for _, v := range langTags {
			if curLang != v {
				restLangs = append(restLangs, &langType{
					Lang: v.String(),
					Name: display.Self.Name(v),
				})
			}
		}
		
		c.Set("Lang", curLang.String())
		c.Set("CurLang", &langType{
			Lang: curLang.String(),
			Name: display.Self.Name(curLang),
		})
		c.Set("RestLangs", restLangs)
		
		// Add i18n function to context
		c.Set("i18n", func(key string, args ...interface{}) string {
			return Tr(curLang.String(), key, args...)
		})
		
		c.Next()
	}
}

// Tr translates a key to the target language
func Tr(lang, key string, args ...interface{}) string {
	if strings, ok := i18nStrings[lang]; ok {
		if val, exists := strings[key]; exists {
			if len(args) > 0 {
				return fmt.Sprintf(val, args...)
			}
			return val
		}
	}
	// Fallback to English
	if strings, ok := i18nStrings["en"]; ok {
		if val, exists := strings[key]; exists {
			if len(args) > 0 {
				return fmt.Sprintf(val, args...)
			}
			return val
		}
	}
	return key
}

// Helper to get available languages config
func GetAvailableLangs() string {
	langs := []string{"en", "zh-CN", "zh-TW", "ja"}
	data, _ := json.Marshal(langs)
	return string(data)
}
