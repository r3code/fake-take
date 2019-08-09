package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
 	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const apiVersion = "1"
const apiRootDefault = "/api/v" + apiVersion

func main() {
	var addr, apiRoot, contentType, ext string
	var port int
	var help bool
	flag.BoolVar(&help, "help", false, "show this usage info")
	flag.StringVar(&addr, "addr", "localhost", "server IP addres to listen at")
	flag.IntVar(&port, "port", 3000, "server port to listen at")
	flag.StringVar(&apiRoot, "apiroot", apiRootDefault, "API base URL, default "+apiRootDefault)
	flag.StringVar(&contentType, "contentType", "application/json", "content type for responses")
	flag.StringVar(&ext, "ext", "resp", "extension for data files")
	flag.Parse()

	fmt.Println("HTTP API test server")
	if help {
		flag.Usage()
		os.Exit(2)
	}

	_, err := url.Parse(apiRoot)
	if err != nil {
		fmt.Println("-apiroot wrong URL format, should be realtive URL (e.g. /api/v2)")
		os.Exit(1)
	}

	router := setupRouter(apiRoot, contentType, ext)
	address := net.JoinHostPort(addr, strconv.Itoa(port))

	fmt.Printf("Listening and serving HTTP on %s\n", address)
	http.ListenAndServe(address, router)
}

func fetchJSONFileMiddleware(apiRoot, contentType, ext string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiRelativePath := strings.Replace(c.Request.URL.Path, apiRoot, "", -1)
		if apiRelativePath == "" || !strings.Contains(c.Request.URL.Path, apiRoot) {
			c.Next()
			return
		}

		filename := strings.Replace(apiRelativePath, "/", "_", -1) + "." + ext
		filename = strings.TrimPrefix(filename, "_")
		log.Println("Read filename: " + filename)
		dat, err := ioutil.ReadFile(filename)
		if err != nil {
			c.Error(err)
			return
		}
		text := string(dat)
		c.Header("Content-Type", contentType)
		c.String(http.StatusOK, text)
		c.Next()
	}
}

func doAPIWelcome(apiRoot, contentType, ext string) gin.HandlerFunc {
	return func(c *gin.Context) {

		paths, err2 := listAPIPaths(".", ext)
		if err2 != nil {
			c.AbortWithError(500, errors.Errorf("Error loading API data: %s\n", err2))
		}

		apiPaths := []string{}
		for _, path := range paths {
			apiPath := apiRoot + "/" + path
			apiPath = strings.Replace(apiPath, "//", "/", -1)
			apiPaths = append(apiPaths, apiPath)
		}

		// c.String(http.StatusOK, fmt.Sprintf("HTTP API TEST Server: API root at '%s', content type '%s'", apiRoot, contentType))

		c.HTML(200, "tpl", gin.H{
			"apiRoot":  apiRoot,
			"apiPaths": apiPaths,
		})

		// c.JSON(200, apiPaths)
	}
}

func showInfo(apiRoot string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, apiRoot)
	}
}

func convertFileNameToPath(name, ext string) (path string, err error) {
	if !strings.Contains(name, "."+ext) {
		return "", errors.Errorf("File name '%s' has invalid extension, must be '.%s'", name, ext)
	}
	onlyName := strings.Replace(name, "."+ext, "", -1)
	path = strings.Replace(onlyName, "_", "/", -1)
	return path, nil
}

func listAPIPaths(dir, ext string) (paths []string, err error) {
	paths = []string{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return paths, errors.WithMessage(err, "list API paths failed")
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		filename := f.Name()
		if strings.Contains(filename, "."+ext) {
			convertedName, err2 := convertFileNameToPath(filename, ext)
			if err2 != nil {
				fmt.Println("ERROR: " + err2.Error())
				continue
			}
			paths = append(paths, convertedName)
		}
	}
	return paths, nil
}

func setupRouter(apiRoot, contentType, ext string) *gin.Engine {
	routes := gin.Default()
	routes.Use(cors.Default()) // allows all origins
	routes.Use(gin.ErrorLogger())
	routes.Use(fetchJSONFileMiddleware(apiRoot, contentType, ext))

	var html = template.Must(template.New("tpl").Parse(`
<h1>Test API Server <sup><small>{{ .apiRoot }}</small></sup></h1>
<h2>Available paths to GET:</h2>
<ol>
{{ range .apiPaths}} 
	<li><a href="{{.}}">{{ . }}</a></li>
{{ end }}
</ol>
`))

	routes.SetHTMLTemplate(html)

	routes.GET("/", showInfo(apiRoot))
	routes.GET(apiRoot, doAPIWelcome(apiRoot, contentType, ext))
	return routes
}
