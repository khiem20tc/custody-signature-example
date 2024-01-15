package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Replace these with your actual API keys and secrets.
const (
	// 	apiKey     = "ZqyrjFE8DgCT"
	// 	privateKey = `-----BEGIN RSA PRIVATE KEY-----
	// MIIEowIBAAKCAQEAoZPMGY90ek4OzQ7FvS7cPnJNECPvRNpaKJDPSPiHmWLFqo52
	// ZvMYQu9c3r2AjT4VFVJaJoN/qGlruq56ckUGhLp9e/rfv+Q4TtxX8PHjpG5hs7Ft
	// wV+UVg7W5EGshT8qrHtWdN1DmSg03ZfKIGFSI95179K6tnflpAbtZtvFknM0/ECP
	// eXAkPCbF7b2/Xv1MkyV/o/QBeumES4rGF7RnBdlTBPeUBZnUPELgsnI8TxTMjeGP
	// BNrj6R1rCbLh7l6p09lx1/KotZMgtNsGFtuJlhaw5sYhoFvyYeiUpFLpU6vtcuCJ
	// N5SGtNYl4MOoNasM7Mb3XVVSdO0xiwUhZDvfSwIDAQABAoIBAGdAAeGniSAKt2yT
	// 7wooYrdI5TPWMrTF720StFMF9eivdG674K+C0lMbkDYJ1JbtQB3C5TbNOwtMann9
	// uuNAdpzkawGJ2+dMmCrUpSGkAPr3SlnAnMlAIZMoomt0CCGRrtxPaHz/U44QYk/k
	// ClbMueeP5b9d4tBtJ4K8poHfGI6vKg1yzCwpkU8XhkKmFlJZOf4APqYnfjQn9mqI
	// xYMJVjvJEi0hiMd0kbxV6vg1UxdeasLiyzouBHn940+l4xDbZ1VXz6M6GR5Xhrjp
	// v8QesIxk4hokA/g5D9RwwER3TdoN90sgzdFeUdKTUeFmvCkz7lD9zpIxHHLf1grw
	// FATMVjkCgYEAzios5G4ibWaY34H/4qi2cvCCXGpRSVx/VNzS00Rrkxd5oBbcJmml
	// qpELEqQ3ZevGdkpJIRCd3TPrln/4+Nbb1fXITvCTkQ7TJi0P8/JWUMczmg4VqLwf
	// nDgcfz3nzWBwpRjQPHD4uFhv7msfPp6PQ4VN103YVd05CvcDxAdbu7cCgYEAyKJ6
	// HA15qBYCWpN/W8PRyXi2wYCwxT03ZwBZk/3kV9uftPI0W15XGRqESXkl2zgMYng0
	// kC4YUhXUmF8/L0G2Frs5CogZJBUL1SQ+HEwBsZqIs53NtZWmEfiruVVH1Hn7/PFx
	// EQ7vXIjCZUbskKn6dNCn+58vVl2w58BtOeclYQ0CgYA8PqjVq7VVwMhlb+ilhGWk
	// WtHNTagpRuVSmCDnabQBzLdW57c3ZmHp4O6aaPBjUS2yfWy3Q9LNxBFQ7l6D4M1m
	// zabWIokMt4dOPZbO038TpdJXb0w2/ZpDHUZ+jEmDg24HYKPhNaYIwJcc1aLQuqbk
	// tTyU8QOJu9aidKJeE0RkKwKBgBohz3XH64iRFU1m2LfDEZgEOQmLEXsfNhAcY457
	// CzrGSE7xHRCpgP6sDX7kYKHk8vgAYBhHaLOIVGBkR36IOIdNa2iLwXqJozjnt49H
	// 9xCC6Ds82oZEL5U3pmZFTU3HdaLEb82g/Fw5E9jNHBLbkNuWMcr8ONYu7dPBpHhe
	// OughAoGBAMwMlD8EhDVzXZrgOsJM+8mODQaO7p1mGhnIHMVEZKEFHoBY9DrWChPu
	// crU11mL2TBLCMT49pyd+XArK9mNBPwxHDj7rJZkQAcn7AKw2M/dxdZroev4JGsKf
	// yOUJggU3NbOjNMuJQL8JtU5lXV/mNpZpfbQdj1d2CZNJaR/wDM1g
	// -----END RSA PRIVATE KEY-----`
	// )

	apiKey     = "5WWCMPHwIEeB"
	privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAu265+cwNqwGIIRR5Ub/XCQ05fxYGnCGxIRdMM7lYlOvFiaQu
vBbjdcPtIsQuKIsuLfgHPVwibqzoNSMO91V6bN7ArkhrnXOb2MXNtu+ZodQ17MA9
+5MW6+pQHHfycvx/Gk1ldLJgtJiUSFaVujE7O4q10gtcMKT4J+HnrCFlEloQtxlT
ElEafUyYk07hbbz0U4saO935/kAoT+PdMkAxTUD1xMxSXEhk7CociuS/gFSPxg8H
KSzHp6iXgozOtOkvl99LD3yZ87EjxhzMKTyFQGx/WjooNzK5yyruPyqTlrgrfaIy
1CZOlRXAP4Ttk2IoEahI9f0sLdXssxXZuwWWkwIDAQABAoIBAD7CPJNft9Pil2o8
KMMusRney7m57kypG14xJtrK3NZAe8wypVNldpQgHm7dsXbx42yQ+Bublgvo6Xeh
XYmDnZKGo423whDefPiAgvkWESMWo1e6pwZtoecsddaScyP9V7G+6JHCiI7v5/aw
x0Go6mRtdaP3Gc9P7aetBJ2mMOmLnAvZkZz+ePMTys3Y0PpGBs/fnyWFZUyuI+w4
Eouz3oX/RNrOfTAymo9GvZjL0pZ4N9UTu3XEGQYd97eUA1J/9HQ8m9JaGl454bxL
1iZmXlfuPnCEJUF58pqDbZc+eejKVy5G496OQsOUxq3e5KEEuJqBh/+rqtBKecr4
EvmiaIECgYEA2sGbqhtT/N2plviHGoL8s7ULnaQIoH4dE+t1l81IjhRUNhLOD8Kc
EB8RBSpKcT7iiE+Kaz0ibhjIQQusHvd0q4n9LE9odMV2s+R2kk/4YMtpjmA8wBDs
1lHl4PpVQx7fq976PidKNTH7yaUIB19bREBWhGltnlRBpivrHf2S18sCgYEA21fj
C1mK59kqQ5GXijeqTrx5hm0JNRQFxHnBoLgMVJ42QiTvU+GhwMKhS0X80jIDdGip
7ydPJaRCo01ldYU7sOV1q/fpz/fEqC4FNLCFsm31VZldX976tihImNtq/+8pybch
ccPjFVo9lqGEDm3ID5kEdmtLBaMzSGyjzetNk1kCgYEAxEnxmePHqyBjKjp7UEi0
47PSZnNn4ksHYHZpH/tt3T9UiOi6yd2AF98ocJAQGCmrL1DgDXXfzRajqeoFWgwF
Pl8lM3tVaWI+LxETbBoh7wjXAJBOMrF9MppuQT+e/glX/mqn9NlgdvcQzVEuMR9Z
T5bDizDm0akc9zR1VoXQG50CgYAovqSwYQvKka6mKo9p33lFcwFoFS0WrQd9PdjY
EBhKR7FwjAfhHxK7CeyIXRHfweaeYyreAAFVzrOKPkBQmlVCQP2g2kaWmUHws8vH
w9qyEHb4Vargujz8RXNm4at4q2apz9jolyjBuKekKZCsVXxKWXRYwwmGnJBULcon
4EPi0QKBgGSyxEosiG5rKjwYW5RP3Ag5invC6yZ/SRHyTYZ8BObjGyo80zS9rMmx
t0hAvBYb1D8CUdq7zW+9RbLxWPuGnXgoIah/6FqNKQ4liAYbICprA+762+zDxgyF
dZ9Mh43LQ0onvSii1GsYSaIh5E9R2X5H1UXT1b+M4Ix+YXLeSUCZ
-----END RSA PRIVATE KEY-----`
)

type MyStruct struct {
	Key string `json:"key"`
}

func privateKeyFromPEM(privateKeyPEM string) (*rsa.PrivateKey, error) {
	// Parse the PEM-encoded private key
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	// Parse the RSA private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func createSignature(req *http.Request, httpMethod string, requestBody string, path string) (string, int64, string, string) {
	// Create a timestamp for the request.
	timestamp := time.Now().Unix()
	// Generate a nonce (random string) for each request.
	nonce := generateNonce()

	// Define the HTTP method (GET, POST, etc.).

	// Create a string to sign based on your API requirements.
	message := []byte(fmt.Sprintf("%s%d%s%s%s%s", apiKey, timestamp, nonce, path, httpMethod, requestBody))
	signature, _ := signData(message)
	// Set the request headers with API key, timestamp, nonce, and signature.
	return signature, timestamp, nonce, path
}

// func main() {
// 	// Replace with your RSA public key
// 	publicKeyPEM := `-----BEGIN RSA PUBLIC KEY-----
// MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0sFIxMIQLCQMdGEL6C7u
// KEjtWb2zbBf8UOZLaFDc90zay+nOJC8uDVsPtVYr7DBuqGpQzuqn1ZCl0J/SHesb
// ye9OPJzldrYotQvYSXovhrnbdPip6zvbrCg37SRVzgumXVpzyfArk9ETrxaah3Np
// izj3TRl+J6O7mylsgZHRKdPxYFsGNKN32xJBGNlGc5LABY9MqyNLQygriJ4qYiub
// qqzVlxo+t6oDdWBQa4d3Oae9H0Jy94du2louVPoke3pkbHuDpgsBY4AE6gr3ytzz
// mQTLqlKKVncDjDJkwZ9ue/r0lVYR/Jh+p2mzmKlF+sN8kqLDX9yBJbEOCjcYaljm
// uwIDAQAB
// -----END RSA PUBLIC KEY-----
// `

// 	// Remove newline characters from the input data
// 	// publicKeyPEM = strings.ReplaceAll(publicKeyPEM, "\n", "")
// 	fmt.Println("pubPEM", []byte(publicKeyPEM))

// 	// Calculate the SHA-256 hash of the modified input data
// 	sha256Hash := sha256.Sum256([]byte(publicKeyPEM))

// 	// Convert the SHA-256 hash to a hexadecimal string
// 	hashString := hex.EncodeToString(sha256Hash[:])

// 	fmt.Println("SHA-256 Hash: ", hashString)
// }

func main() {
	// Define the API endpoint you want to call.
	// apiURL := "http://127.0.0.1:7011/test"
	apiURL := "http://127.0.0.1:3000/api/v1/address/deposit-address?depositLabel=poolETH1a"

	// Create an HTTP client.
	client := &http.Client{}

	// httpMethod := "POST" // Change this to the desired method.
	httpMethod := "GET" // Change this to the desired method.

	myData := MyStruct{
		Key: "value",
	}
	jsonData, _ := json.Marshal(myData)
	requestBody := string(jsonData)
	// Create an HTTP request.
	getReq, err := http.NewRequest(httpMethod, apiURL, strings.NewReader(requestBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	signature, timestamp, nonce, _ := createSignature(getReq, httpMethod, requestBody, "/api/v1/address/deposit-address")

	getReq.Header.Set("X-Aegis-Api-Key", apiKey)
	getReq.Header.Set("X-Aegis-Api-Timestamp", fmt.Sprintf("%d", timestamp))
	getReq.Header.Set("X-Aegis-Api-Nonce", nonce)
	getReq.Header.Set("X-Aegis-Api-Signature", signature)
	getReq.Header.Set("Content-Type", "application/json") // Set the appropriate content type.

	// Send the HTTP request.
	resp, err := client.Do(getReq)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read and print the API response.
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Printf("API Response:\n%s\n", responseBody)
}

func generateNonce() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func signData(message []byte) (string, error) {
	privateKey, _ := privateKeyFromPEM(privateKey)

	opts := &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.SHA256,
	}
	hashed := sha256.Sum256(message)
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, hashed[:], opts)
	if err != nil {
		return "", err
	}

	signatureBase64 := base64.StdEncoding.EncodeToString(signature)
	return signatureBase64, nil
}
