package files

import (
	"context"
	"fmt"

	"github.com/viant/afs"
	"github.com/viant/afs/url"
	"github.com/viant/mly/service/config"
)

// ModifiedSnapshot checks and updates modified times based on the object in URL
func ModifiedSnapshot(ctx context.Context, fs afs.Service, URL string, resource *config.Modified) (*config.Modified, error) {
	objects, err := fs.List(ctx, URL)
	if err != nil {
		return resource, fmt.Errorf("failed to list URL:%s; error:%w", URL, err)
	}

	if extURL := url.SchemeExtensionURL(URL); extURL != "" {
		object, err := fs.Object(ctx, ex