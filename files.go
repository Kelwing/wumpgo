package objects

import (
	"fmt"
	"io"
	"mime/multipart"
	"strings"
)

type DiscordFile struct {
	io.Reader
	Filename    string
	Description string
	Spoiler     bool
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
