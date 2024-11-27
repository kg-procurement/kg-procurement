package mailer

import (
	"bytes"
	"io"
	"mime/multipart"
	"testing"

	"github.com/onsi/gomega"
)

func TestBulkFromMultipartForm(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	// Arrange: Create a mock multipart form with files
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Create first file
	part1, err := writer.CreateFormFile("files", "file1.txt")
	g.Expect(err).To(gomega.BeNil())
	_, err = io.WriteString(part1, "This is the content of file1.")
	g.Expect(err).To(gomega.BeNil())

	// Create second file
	part2, err := writer.CreateFormFile("files", "file2.txt")
	g.Expect(err).To(gomega.BeNil())
	_, err = io.WriteString(part2, "This is the content of file2.")
	g.Expect(err).To(gomega.BeNil())

	// Close the writer to finalize the form
	writer.Close()

	req := multipart.NewReader(&buffer, writer.Boundary())
	form, err := req.ReadForm(16 << 20)
	g.Expect(err).To(gomega.BeNil())

	files := form.File["files"]
	attachments, err := BulkFromMultipartForm(files)

	// Assert: Check for no errors
	g.Expect(err).To(gomega.BeNil())
	g.Expect(attachments).To(gomega.HaveLen(2))

	// Verify the first attachment
	g.Expect(attachments[0].Filename).To(gomega.Equal("file1.txt"))
	g.Expect(string(attachments[0].Data)).To(gomega.Equal("This is the content of file1."))
	g.Expect(attachments[0].MIMEType).To(gomega.Equal("application/octet-stream"))

	// Verify the second attachment
	g.Expect(attachments[1].Filename).To(gomega.Equal("file2.txt"))
	g.Expect(string(attachments[1].Data)).To(gomega.Equal("This is the content of file2."))
	g.Expect(attachments[1].MIMEType).To(gomega.Equal("application/octet-stream"))
}
