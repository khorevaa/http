package client

import (
	"bytes"
	"log"
	"math"
	"net"
	"net/url"

	jsoniter "github.com/json-iterator/go"
)

// Client represents a single
type Client struct {
	request  request
	response Response
}

// strings.Replace(path, " ", "%20", -1)

// Get builds a GET request.
func Get(path string) *Client {
	parsedURL, _ := url.Parse(path)

	http := &Client{
		request: request{
			method: "GET",
			url:    parsedURL,
			headers: Headers{
				"Accept-Encoding": "gzip",
			},
		},
	}

	return http
}

// Post builds a POST request.
func Post(path string) *Client {
	parsedURL, _ := url.Parse(path)

	http := &Client{
		request: request{
			method: "POST",
			url:    parsedURL,
			headers: Headers{
				"Accept-Encoding": "gzip",
			},
		},
	}

	return http
}

// Header sets one HTTP header for the request.
func (http *Client) Header(key string, value string) *Client {
	http.request.headers[key] = value
	return http
}

// Headers sets the HTTP headers for the request.
func (http *Client) Headers(headers Headers) *Client {
	for key, value := range headers {
		http.request.headers[key] = value
	}

	return http
}

// Body sets the request body.
func (http *Client) Body(raw string) *Client {
	http.request.body = raw
	return http
}

// BodyJSON sets the request body by converting the object to JSON.
func (http *Client) BodyJSON(obj interface{}) *Client {
	data, err := jsoniter.MarshalToString(obj)

	if err != nil {
		log.Printf("Error converting request body to JSON: %v", err)
		return http
	}

	http.request.body = data
	return http
}

// BodyBytes sets the request body as a byte slice.
func (http *Client) BodyBytes(raw []byte) *Client {
	http.request.body = string(raw)
	return http
}

// Response returns the response object.
func (http *Client) Response() *Response {
	return &http.response
}

// Do executes the request and returns the response.
func (http *Client) Do() error {
	remoteAddress := net.TCPAddr{
		IP:   net.ParseIP("165.22.146.88"),
		Port: 80,
	}

	connection, err := net.DialTCP("tcp", nil, &remoteAddress)

	if err != nil {
		return err
	}

	connection.SetNoDelay(true)
	connection.Write([]byte("GET / HTTP/1.1\r\nHost: notify.moe\r\nAccept: */*\r\nUser-Agent: Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)\r\n\r\n"))

	var header bytes.Buffer
	var body bytes.Buffer
	current := &header
	tmp := make([]byte, 16384)
	contentLength := 0

	for {
		n, err := connection.Read(tmp)
		headerEnd := bytes.Index(tmp, []byte{'\r', '\n', '\r', '\n'})

		if headerEnd != -1 {
			header.Write(tmp[:headerEnd])
			body.Write(tmp[headerEnd+4 : n])
			current = &body
			println(header.String())

			// Find content length
			http.response.header = header.Bytes()
			lengthSlice := http.response.Header([]byte("Content-Length"))

			// Convert it to an integer
			for i := 0; i < len(lengthSlice); i++ {
				contentLength += (int(lengthSlice[i]) - 48) * int(math.Pow10(len(lengthSlice)-i-1))
			}

			body.Grow(contentLength)
		} else {
			current.Write(tmp[:n])
		}

		if err != nil {
			http.response.body = body.Bytes()
			return err
		}

		if body.Len() >= contentLength {
			http.response.body = body.Bytes()
			return nil
		}
	}
}

// End executes the request and returns the response.
func (http *Client) End() (*Response, error) {
	err := http.Do()
	return &http.response, err
}

// // EndStruct executes the request, unmarshals the response body into a struct and returns the response.
// func (http *Client) EndStruct(obj interface{}) (*Response, error) {
// 	err := http.Do()

// 	if err != nil {
// 		return &http.response, err
// 	}

// 	err = http.response.Unmarshal(obj)
// 	return &http.response, err
// }
