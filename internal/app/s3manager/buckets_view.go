package s3manager

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"net/url"

	"github.com/minio/minio-go/v7"
)

// HandleBucketsView shows a list of all buckets.
func HandleBucketsView(s3 S3, templates fs.FS, allowDelete bool) http.HandlerFunc {
	type pageData struct {
		Buckets     []minio.BucketInfo
		AllowDelete bool
	}

	return func(w http.ResponseWriter, r *http.Request) {
		buckets, err := s3.ListBuckets(r.Context())
		if err != nil {
			handleHTTPError(w, fmt.Errorf("error listing buckets: %w", err))
			return
		}
		data := pageData{
			Buckets:     buckets,
			AllowDelete: allowDelete,
		}

		t := template.New("layout.html.tmpl").Funcs(template.FuncMap{
			"urlPathEscape": url.PathEscape,
		})
		t, err = t.ParseFS(templates, "layout.html.tmpl", "buckets.html.tmpl")
		if err != nil {
			handleHTTPError(w, fmt.Errorf("error parsing template files: %w", err))
			return
		}
		err = t.ExecuteTemplate(w, "layout", data)
		if err != nil {
			handleHTTPError(w, fmt.Errorf("error executing template: %w", err))
		}
	}
}
