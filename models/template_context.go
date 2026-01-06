package models

import (
	"bytes"
	"encoding/base64"
	"net/mail"
	"net/url"
	"path"
	"text/template"

	qrcode "github.com/skip2/go-qrcode"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// TemplateContext is an interface that allows both campaigns and email
// requests to have a PhishingTemplateContext generated for them.
type TemplateContext interface {
	getFromAddress() string
	getBaseURL() string
}

// PhishingTemplateContext is the context that is sent to any template, such
// as the email or landing page content.
type PhishingTemplateContext struct {
	From        string
	URL         string
	Tracker     string
	TrackingURL string
	RId         string
	BaseURL     string
	// QRURL is an inline base64-encoded PNG <img> tag representing a QR code of {{.URL}}.
	// It is safe to embed directly in HTML email templates.
	QRURL       string
	// QRURLCID is a convenience <img> tag that references the QR as an embedded CID attachment.
	// When used, the mail generator should embed a PNG attachment named 'qr.png'.
	QRURLCID    string
	BaseRecipient
}

// NewPhishingTemplateContext returns a populated PhishingTemplateContext,
// parsing the correct fields from the provided TemplateContext and recipient.
func NewPhishingTemplateContext(ctx TemplateContext, r BaseRecipient, rid string) (PhishingTemplateContext, error) {
	f, err := mail.ParseAddress(ctx.getFromAddress())
	if err != nil {
		return PhishingTemplateContext{}, err
	}
	fn := f.Name
	if fn == "" {
		fn = f.Address
	}
	templateURL, err := ExecuteTemplate(ctx.getBaseURL(), r)
	if err != nil {
		return PhishingTemplateContext{}, err
	}

	// For the base URL, we'll reset the the path and the query
	// This will create a URL in the form of http://example.com
	baseURL, err := url.Parse(templateURL)
	if err != nil {
		return PhishingTemplateContext{}, err
	}
	baseURL.Path = ""
	baseURL.RawQuery = ""

	phishURL, _ := url.Parse(templateURL)
	q := phishURL.Query()
	q.Set(RecipientParameter, rid)
	phishURL.RawQuery = q.Encode()

	trackingURL, _ := url.Parse(templateURL)
	trackingURL.Path = path.Join(trackingURL.Path, "/track")
	trackingURL.RawQuery = q.Encode()

	// Generate a QR code for the phishing URL. If QR generation fails for any reason,
	// we keep QRURL empty to avoid interrupting normal email generation.
	var qrTag string
	// Tăng độ phân giải, thêm border, bo góc, đổ bóng để QR đẹp hơn
	if qr, err := qrcode.New(phishURL.String(), qrcode.High); err == nil {
		qr.DisableBorder = false
		img := qr.Image(384)

		rgba := image.NewRGBA(img.Bounds())
		draw.Draw(rgba, img.Bounds(), img, image.Point{}, draw.Src)
		label := "Scan QR"
		col := color.RGBA{128, 0, 192, 255} // tím đậm
		face := basicfont.Face7x13
		labelWidth := len(label) * 7 // 7px/char
		x := (rgba.Bounds().Dx() - labelWidth) / 2
		y := (rgba.Bounds().Dy()+13)/2 // 13px chiều cao font
		d := &font.Drawer{
			Dst:  rgba,
			Src:  image.NewUniform(col),
			Face: face,
			Dot:  fixed.P(x, y),
		}
		d.DrawString(label)

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, rgba); err == nil {
			b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
			qrTag = `<img alt='' style='display:block;margin:16px auto;border-radius:16px;box-shadow:0 2px 12px #0002;border:4px solid #fff;width:256px;height:256px;background:#fff' width='256' height='256' src='data:image/png;base64,` + b64 + `'/>`
		}
	}

	// Convenience CID-based tag (requires embedding an attachment named 'qr.png')
	cidTag := `<img alt='' style='display:block;margin:16px auto;border-radius:16px;box-shadow:0 2px 12px #0002;border:4px solid #fff;width:256px;height:256px;background:#fff' width='256' height='256' src='cid:qr.png'/>`

	return PhishingTemplateContext{
		BaseRecipient: r,
		BaseURL:       baseURL.String(),
		URL:           phishURL.String(),
		TrackingURL:   trackingURL.String(),
		Tracker:       "<img alt='' style='display: none' src='" + trackingURL.String() + "'/>",
		From:          fn,
		RId:           rid,
		QRURL:         qrTag,
		QRURLCID:      cidTag,
	}, nil
}

// GenerateQRBase64 returns a base64-encoded PNG QR image for the given URL.
// The size parameter specifies the width/height in pixels.
func GenerateQRBase64(url string, size int) (string, error) {
	png, err := qrcode.Encode(url, qrcode.Medium, size)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(png), nil
}

// ExecuteTemplate creates a templated string based on the provided
// template body and data.
func ExecuteTemplate(text string, data interface{}) (string, error) {
	buff := bytes.Buffer{}
	tmpl, err := template.New("template").Parse(text)
	if err != nil {
		return buff.String(), err
	}
	err = tmpl.Execute(&buff, data)
	return buff.String(), err
}

// ValidationContext is used for validating templates and pages
type ValidationContext struct {
	FromAddress string
	BaseURL     string
}

func (vc ValidationContext) getFromAddress() string {
	return vc.FromAddress
}

func (vc ValidationContext) getBaseURL() string {
	return vc.BaseURL
}

// ValidateTemplate ensures that the provided text in the page or template
// uses the supported template variables correctly.
func ValidateTemplate(text string) error {
	vc := ValidationContext{
		FromAddress: "foo@bar.com",
		BaseURL:     "http://example.com",
	}
	td := Result{
		BaseRecipient: BaseRecipient{
			Email:     "foo@bar.com",
			FirstName: "Foo",
			LastName:  "Bar",
			Position:  "Test",
		},
		RId: "123456",
	}
	ptx, err := NewPhishingTemplateContext(vc, td.BaseRecipient, td.RId)
	if err != nil {
		return err
	}
	_, err = ExecuteTemplate(text, ptx)
	if err != nil {
		return err
	}
	return nil
}
