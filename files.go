package objects

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
)

type DiscordFile struct {
	*bytes.Buffer
	Filename    string
	Description string
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

	w, err := m.CreateFormFile(fmt.Sprintf(formFieldNameFmt, index), f.Filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(w, f)
	if err != nil {
		return nil, err
	}

	return a, nil
}
