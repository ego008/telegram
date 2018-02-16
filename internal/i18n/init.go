package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/HentaiDB/HentaiDBot/internal/errors"
	log "github.com/kirillDanshin/dlog"
	"github.com/nicksnyder/go-i18n/i18n"
)

var (
	Names = make(map[string]string)
	Tags  []string
)

func Initialize(pathToStrings string) {
	err := filepath.Walk(pathToStrings, func(path string, file os.FileInfo, err error) error {
		log.Ln("Walk to", path)
		if !strings.HasSuffix(file.Name(), ".all.yaml") {
			return nil
		}

		i18n.MustLoadTranslationFile(path)
		return nil
	})
	errors.Check(err)

	Tags = i18n.LanguageTags()
	for _, tag := range Tags {
		T, err := SwitchTo(tag)
		errors.Check(err)

		Names[tag] = fmt.Sprint(
			T("language_flag"), " ", strings.Title(T("language_name")),
		)
		log.Ln("Tag", tag, ":", Names[tag])
	}

	sort.Strings(Tags)
}
