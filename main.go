package main

import (
	
	"log"
	"fmt"
	"net/http"
	"encoding/base64"
	"encoding/json"
	"github.com/qiniu/api/rs"
	. "github.com/qiniu/api/conf"
)

const (
	BUCKET = "APPLY YOUR BUCKET NAME HERE"  // change to own space name
	DOMAIN = "APPLY YOUR DOMAIN HERE" // For example: myspace.qiniudn.com
)

// --------------------------------------------------------------------------------

func init() {

	ACCESS_KEY = "" // Apply Access key here
	SECRET_KEY = "" // Apply Secret key here
	if ACCESS_KEY == "" || SECRET_KEY == "" {
		panic("require ACCESS_KEY & SECRET_KEY")
	}
}

// --------------------------------------------------------------------------------
// HTML code that will be shown on the webpage

//Simple upload without assigning key for the image you want to upload
var uploadFormFmt = `
<html>
 <body>
  <form method="post" action="http://up.qiniu.com/" enctype="multipart/form-data">
   <input name="token" type="hidden" value="%s">
   Image to upload: <input name="file" type="file"/>
   <input type="submit" value="Upload">
  </form>
 </body>
</html>
`
//Assign a key for the image you want to upload
var uploadWithKeyFormFmt = `
<html>
 <body>
  <form method="post" action="http://up.qiniu.com/" enctype="multipart/form-data">
   <input name="token" type="hidden" value="%s">
   Image key in qiniu cloud storage: <input name="key" value="foo bar.jpg"><br>
   Image to upload: <input name="file" type="file"/>
   <input type="submit" value="Upload">
  </form>
 </body>
</html>
`
//Assign both key and custom field 
var uploadWithkeyAndCustomFieldFmt = `
<html>
 <body>
  <form method="post" action="http://up.qiniu.com/" enctype="multipart/form-data">
   <input name="token" type="hidden" value="%s">
   <input name="x:custom_field_name" value="x:custom_field_name">
   Image key in qiniu cloud storage: <input name="key" value="foo bar.jpg"><br>
   Image to upload: <input name="file" type="file"/>
   <input type="submit" value="Upload">
  </form>
 </body>
</html>
`

var returnPageFmt = `
<html>
 <body>
  <p>%s
  <p>ImageDownloadUrl: %s
  <p><a href="/upload">Back to upload</a>
  <p><a href="/upload2">Back to uploadWithKey</a>
  <p><a href="/upload3">Back to uploadWithkeyAndCustomField</a>
  <p><img src="%s">
 </body>
</html>
`

type UploadRet struct {
	Key string `json:"key"`
}

func handleReturn(w http.ResponseWriter, req *http.Request) {

	ret := req.FormValue("upload_ret")
	b, err := base64.URLEncoding.DecodeString(ret)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	var uploadRet UploadRet
	err = json.Unmarshal(b, &uploadRet)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	policy := rs.GetPolicy{Scope: "*/" + uploadRet.Key}
	img := policy.MakeRequest(rs.MakeBaseUrl(DOMAIN, uploadRet.Key))
	returnPage := fmt.Sprintf(returnPageFmt, string(b), img, img)
	w.Write([]byte(returnPage))
}

func handleUpload(w http.ResponseWriter, req *http.Request) {

	policy := rs.PutPolicy{Scope: BUCKET, ReturnUrl: "http://localhost:8765/uploaded"}
	token := policy.Token()
	log.Println("token:", token)
	uploadForm := fmt.Sprintf(uploadFormFmt, token)
	w.Write([]byte(uploadForm))
}

func handleUploadWithKey(w http.ResponseWriter, req *http.Request) {

	policy := rs.PutPolicy{Scope: BUCKET, ReturnUrl: "http://localhost:8765/uploaded"}
	token := policy.Token()
	log.Println("token:", token)
	uploadForm := fmt.Sprintf(uploadWithKeyFormFmt, token)
	w.Write([]byte(uploadForm))
}

func handleUploadWithKeyAndCustomField(w http.ResponseWriter, req *http.Request) {

	policy := rs.PutPolicy{Scope: BUCKET, ReturnUrl: "http://localhost:8765/uploaded"}
	token := policy.Token()
	log.Println("token:", token)
	uploadForm := fmt.Sprintf(uploadWithkeyAndCustomFieldFmt, token)
	w.Write([]byte(uploadForm))
}

func handleDefault(w http.ResponseWriter, req *http.Request) {

	http.Redirect(w, req, "/upload", 302)
}

func main() {

	http.HandleFunc("/", handleDefault)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/upload2", handleUploadWithKey)
	http.HandleFunc("/upload3", handleUploadWithKeyAndCustomField)
	http.HandleFunc("/uploaded", handleReturn)
	log.Fatal(http.ListenAndServe(":8765", nil))
}

// --------------------------------------------------------------------------------

