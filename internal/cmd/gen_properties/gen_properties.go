// This program generates a property file in Go file from Unicode Character
// Database auxiliary data files. The command line arguments are as follows:
//
//  1. The name of the Unicode data file (just the filename, without extension).
//     Can be "-" (to skip) if the emoji flag is included.
//  2. The name of the locally generated Go file.
//  3. The name of the slice mapping code points to properties.
//  4. The name of the generator, for logging purposes.
//  5. (Optional) Flags, comma-separated. The following flags are available:
//     - "emojis=<property>": include the specified emoji properties (e.g.
//     "Extended_Pictographic").
//     - "gencat": include general category properties.
package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// We want to test against a specific version rather than the latest. When the
// package is upgraded to a new version, change these to generate new tests.
const (
	unicodeVersion    = "15.1.0"
	propertyURLFormat = `https://www.unicode.org/Public/%s/ucd/%s.txt`
	emojiURLFormat    = `https://unicode.org/Public/%s/ucd/emoji/emoji-data.txt`
)

// The regular expression for a line containing a code point range property.
var propertyPattern = regexp.MustCompile(`^([0-9A-F]{4,6})(\.\.([0-9A-F]{4,6}))?\s*;\s*([A-Za-z0-9_]+)\s*#\s(.+)$`)

type Options struct {
	propertyName   string
	emojis         string
	gencat         bool
	dictionaryName string
	prefix         string
	typeName       string
}

func main() {
	var propertyName string
	var emojis string
	var gencat bool
	var logPrefix string
	var prefix string
	var typeName string
	flag.StringVar(&propertyName, "property", "", "name of the property")
	flag.StringVar(&emojis, "emojis", "", "emoji properties to include")
	flag.BoolVar(&gencat, "gencat", false, "include general category properties")
	flag.StringVar(&logPrefix, "logprefix", "", "prefix for log messages")
	flag.StringVar(&prefix, "prefix", "pr", "prefix for property names")
	flag.StringVar(&typeName, "type", "property", "name of the property type")
	flag.Parse()

	outputFilename := flag.Arg(0)
	dictionaryName := flag.Arg(1)

	log.SetPrefix("gen_properties (" + logPrefix + "): ")
	log.SetFlags(0)

	src, err := parse(&Options{
		propertyName:   propertyName,
		emojis:         emojis,
		gencat:         gencat,
		dictionaryName: dictionaryName,
		prefix:         prefix,
		typeName:       typeName,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Format the Go code.
	formatted, err := format.Source([]byte(src))
	if err != nil {
		log.Fatal("gofmt:", err)
	}

	// Save it to the (local) target file.
	log.Print("Writing to ", outputFilename)
	if err := os.WriteFile(outputFilename, formatted, 0644); err != nil {
		log.Fatal(err)
	}
}

// parse parses the Unicode Properties text files located at the given URLs and
// returns their equivalent Go source code to be used in the uniseg package. If
// "emojiProperty" is not an empty string, emoji code points for that emoji
// property (e.g. "Extended_Pictographic") will be included. In those cases, you
// may pass an empty "propertyURL" to skip parsing the main properties file. If
// "includeGeneralCategory" is true, the Unicode General Category property will
// be extracted from the comments and included in the output.
func parse(opts *Options) (string, error) {
	if opts.propertyName == "" && opts.emojis == "" {
		return "", errors.New("no properties to parse")
	}
	var propertyURL string
	var emojiURL string

	// Temporary buffer to hold properties.
	var properties [][4]string

	// Open the first URL.
	if opts.propertyName != "" {
		propertyURL = fmt.Sprintf(propertyURLFormat, unicodeVersion, opts.propertyName)
		log.Printf("Parsing %s", propertyURL)
		res, err := http.Get(propertyURL)
		if err != nil {
			return "", err
		}
		in1 := res.Body
		defer in1.Close()

		// Parse it.
		scanner := bufio.NewScanner(in1)
		num := 0
		for scanner.Scan() {
			num++
			line := strings.TrimSpace(scanner.Text())

			// Skip comments and empty lines.
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}

			// Everything else must be a code point range, a property and a comment.
			from, to, property, comment, err := parseProperty(line)
			if err != nil {
				return "", fmt.Errorf("%s line %d: %v", os.Args[4], num, err)
			}
			properties = append(properties, [4]string{from, to, property, comment})
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
	}

	// Open the second URL.
	if opts.emojis != "" {
		emojiURL = fmt.Sprintf(emojiURLFormat, unicodeVersion)
		log.Printf("Parsing %s", emojiURL)
		res, err := http.Get(emojiURL)
		if err != nil {
			return "", err
		}
		in2 := res.Body
		defer in2.Close()

		// Parse it.
		scanner := bufio.NewScanner(in2)
		num := 0
		for scanner.Scan() {
			num++
			line := scanner.Text()

			// Skip comments, empty lines
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}

			// Everything else must be a code point range, a property and a comment.
			from, to, property, comment, err := parseProperty(line)
			if err != nil {
				return "", fmt.Errorf("emojis line %d: %v", num, err)
			}
			if property != opts.emojis {
				continue
			}
			properties = append(properties, [4]string{from, to, property, comment})
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
	}

	// Sort properties.
	sort.Slice(properties, func(i, j int) bool {
		left, _ := strconv.ParseUint(properties[i][0], 16, 64)
		right, _ := strconv.ParseUint(properties[j][0], 16, 64)
		return left < right
	})
	properties = eytzinger(properties)

	// Header.
	var (
		buf          bytes.Buffer
		emojiComment string
	)
	if emojiURL != "" {
		emojiComment = `
// and
// ` + emojiURL + `
// ("Extended_Pictographic" only)`
	}

	buf.WriteString(`// Code generated by ./internal/cmd/gen_properties/gen_properties.go; DO NOT EDIT.

package uniseg

// ` + opts.dictionaryName + ` are taken from
// ` + propertyURL + emojiComment + `
// See https://www.unicode.org/license.html for the Unicode license agreement.
var ` + opts.dictionaryName + ` = dictionary[` + opts.typeName + `]{
`)

	// Properties.
	for _, prop := range properties {
		if opts.gencat {
			generalCategory := "0"
			if len(prop[3]) >= 2 {
				generalCategory = "gc" + prop[3][:2]
				if generalCategory == "gcL&" {
					generalCategory = "gcLC"
				}
				prop[3] = prop[3][3:]
			}
			fmt.Fprintf(
				&buf,
				"{runeRange{%s,%s}, propertyGeneralCategory{%s, %s}}, // %s\n",
				formatRune(prop[0]), formatRune(prop[1]), translateProperty(opts.prefix, prop[2]), generalCategory, prop[3],
			)
		} else {
			fmt.Fprintf(
				&buf,
				"{runeRange{%s,%s}, %s}, // %s\n",
				formatRune(prop[0]), formatRune(prop[1]), translateProperty(opts.prefix, prop[2]), prop[3],
			)
		}
	}

	// Tail.
	buf.WriteString("}")

	return buf.String(), nil
}

func formatRune(s string) string {
	if s == "" {
		return "0"
	}
	return "0x" + s
}

// parseProperty parses a line of the Unicode properties text file containing a
// property for a code point range and returns it along with its comment.
func parseProperty(line string) (from, to, property, comment string, err error) {
	fields := propertyPattern.FindStringSubmatch(line)
	if fields == nil {
		err = errors.New("no property found")
		return
	}
	from = fields[1]
	to = fields[3]
	if to == "" {
		to = from
	}
	property = fields[4]
	comment = fields[5]
	return
}

// translateProperty translates a property name as used in the Unicode data file
// to a variable used in the Go code.
func translateProperty(prefix, property string) string {
	if property == "" {
		return "0"
	}
	return prefix + strings.ReplaceAll(property, "_", "")
}

func eytzinger[S ~[]E, E any](a S) S {
	b := make(S, len(a))
	ctx := &eytzingerContext[S, E]{a: a, b: b}
	ctx.eytzinger(0, 0)
	return b
}

type eytzingerContext[S ~[]E, E any] struct {
	a S // input slice
	b S // output slice
}

func (ctx *eytzingerContext[S, E]) eytzinger(i, k int) int {
	if k < len(ctx.a) {
		i = ctx.eytzinger(i, 2*k+1)
		ctx.b[k] = ctx.a[i]
		i++
		i = ctx.eytzinger(i, 2*k+2)
	}
	return i
}
