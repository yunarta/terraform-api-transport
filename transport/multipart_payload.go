package transport

import (
	"bytes"
	"fmt"
	"mime/multipart"
)

// MultipartFile represents a file for multipart upload with a specific key, name, and content.
type MultipartFile struct {
	Key     string
	Name    string
	Content string
}

// MultipartPayload Define a struct to hold all data for a multipart payload, including other form fields and the file
type MultipartPayload struct {
	Form     map[string]string
	File     *MultipartFile
	Boundary string
}

// Accept defines a method that returns the "Accept" header value for the MultipartPayload.
// In this case, the method always returns "application/json".
func (m *MultipartPayload) Accept() string {
	return "application/json"
}

// ContentType defines a method to generate the "Content-Type" header value for a multipart payload.
// The value is computed by combining the content type "multipart/form-data" with the boundary value of the payload.
func (m *MultipartPayload) ContentType() string {
	return fmt.Sprintf("multipart/form-data; boundary=%s", m.Boundary)
}

// Content Define method to create the actual content payload of the multipart/form-data.
func (m *MultipartPayload) Content() ([]byte, error) {
	// Initialize a bytes buffer and a multipart writer for constructing the payload
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Set the boundary value for multipart.
	m.Boundary = writer.Boundary()

	// Loop through the form map and add all fields to the writer
	for key, value := range m.Form {
		_ = writer.WriteField(key, value)
	}

	// If a file is present in the payload, add it to the writer
	if m.File != nil {
		// Create a form file field in the writer
		fileWriter, err := writer.CreateFormFile(m.File.Key, m.File.Name)
		if err != nil {
			return nil, err
		}
		// Write the content to the newly created form file field
		_, err = fileWriter.Write([]byte(m.File.Content))
		if err != nil {
			return nil, err
		}
	}

	// Done with adding all form fields and files, close the multipart writer
	err := writer.Close()
	if err != nil {
		return nil, err
	}

	// Convert the byte buffer to byte slice and return as final payload
	return body.Bytes(), nil
}

func (m *MultipartPayload) ContentMust() []byte {
	payload, err := m.Content()
	if err != nil {
		panic(err)
	}

	return payload
}

// _ is an assertion to ensure that our MultipartPayload struct complies with the PayloadData interface
var _ PayloadData = &MultipartPayload{}
