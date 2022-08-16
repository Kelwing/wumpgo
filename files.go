package objects

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"strings"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

type DiscordFile struct {
	*bytes.Buffer
	Filename    string
	Description string
	ContentType string
	Spoiler     bool
}

func NewDiscordFile(r io.Reader, filename, description string) (*DiscordFile, error) {
	f := &DiscordFile{
		Buffer:      &bytes.Buffer{},
		Filename:    filename,
		Description: description,
	}
	_, err := io.Copy(f.Buffer, r)
	if err != nil {
		return nil, err
	}

	return f, nil
}

const formFieldNameFmt = "files[%d]"

func (f *DiscordFile) GenerateAttachment(index Snowflake, m *multipart.Writer) (*Attachment, error) {
	if f.Spoiler && !strings.HasPrefix(f.Filename, "SPOILER_") {
		f.Filename = "SPOILER_" + f.Filename
	}
	a := &Attachment{
		DiscordBaseObject: DiscordBaseObject{ID: index},
		Filename:          f.Filename,
		Description:       f.Description,
	}
	contentType := "application/octet-stream"
	if f.ContentType != "" {
		contentType = f.ContentType
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fmt.Sprintf(formFieldNameFmt, index)), escapeQuotes(f.Filename)))
	h.Set("Content-Type", contentType)

	w, err := m.CreatePart(h)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(w, f)
	if err != nil {
		return nil, err
	}

	return a, nil
}
