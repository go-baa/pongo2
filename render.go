// Package pongo2 providers the pongo2 template engine for baa.
package pongo2

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-baa/baa"
	"github.com/safeie/pongo2"
)

// Render the pongo2 template engine
type Render struct {
	Options
	fileChanges chan notifyItem // notify file changes
}

// Options render options
type Options struct {
	Baa        *baa.Baa                         // baa
	Root       string                           // template root dir
	Extensions []string                         // template file extensions
	Filters    map[string]pongo2.FilterFunction // template filters
	Functions  map[string]interface{}           // template functions
	Context    map[string]interface{}           // template global context
}

// tplIndexes template name path indexes
var tplIndexes map[string]string

// New create a template engine
func New(o Options) *Render {
	// init indexes map
	tplIndexes = map[string]string{}

	r := new(Render)
	r.Baa = o.Baa
	r.Root = o.Root
	r.Extensions = o.Extensions
	r.Context = o.Context

	// check template dir
	if r.Root == "" {
		panic("pongo2.New: template dir is empty!")
	}
	r.Root, _ = filepath.Abs(r.Root)
	slash := "/"
	if runtime.GOOS == "windows" {
		slash = "\\"
	}
	if r.Root[len(r.Root)-1] != slash[0] {
		r.Root += slash
	}
	if f, err := os.Stat(r.Root); err != nil {
		panic("pongo2.New: template dir[" + r.Root + "] open error: " + err.Error())
	} else {
		if !f.IsDir() {
			panic("pongo2.New: template dir[" + r.Root + "] is not s directory!")
		}
	}

	// check extension
	if r.Extensions == nil {
		r.Extensions = []string{".html"}
	}

	// register filter
	for name, filter := range o.Filters {
		pongo2.RegisterOrReplaceFilter(name, filter)
	}

	// merge function into context
	for k, v := range o.Functions {
		if _, ok := r.Context[k]; ok {
			panic("pongo2.New: context key[" + k + "] already exists in functions")
		}
		r.Context[k] = v
	}

	if baa.Env != baa.PROD {
		// enable debug mode
		pongo2.DefaultSet.Debug = true

		r.fileChanges = make(chan notifyItem, 8)
		go r.notify()
		go func() {
			for item := range r.fileChanges {
				if r.Baa != nil && r.Baa.Debug() {
					r.Error("filechanges Receive -> " + item.path)
				}
				if item.event == Create || item.event == Write {
					r.parseFile(item.path)
				}
			}
		}()
	}

	// load templates
	r.loadTpls()

	return r
}

// Render template
func (r *Render) Render(w io.Writer, tpl string, data interface{}) error {
	path, ok := tplIndexes[tpl]
	if !ok {
		return fmt.Errorf("pongo2.Render: tpl [%s] not found", tpl)
	}

	t, err := pongo2.FromCache(path)
	if err != nil {
		return err
	}

	ctx, err := r.buildContext(data)
	if err != nil {
		return err
	}

	return t.ExecuteWriter(ctx, w)
}

// buildContext build pongo2 render context
func (r *Render) buildContext(in interface{}) (pongo2.Context, error) {
	// check data type
	data, ok := in.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("pongo2.buildContext: unsupported render data type [%v]", in)
	}

	// copy from global context
	ctx := map[string]interface{}{}
	for k, v := range r.Context {
		ctx[k] = v
	}

	// fill with render data
	for k, v := range data {
		if _, ok := ctx[k]; ok {
			return nil, fmt.Errorf("pongo2.buildContext: render data key [%s] already exists", k)
		}
		ctx[k] = v
	}

	return pongo2.Context(ctx), nil
}

// loadTpls load all template files
func (r *Render) loadTpls() {
	paths, err := r.readDir(r.Root)
	if err != nil {
		r.Error(err)
		return
	}
	for _, path := range paths {
		err = r.parseFile(path)
		if err != nil {
			r.Error(err)
		}
	}
}

// readDir scan dir load all template files
func (r *Render) readDir(path string) ([]string, error) {
	var paths []string
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fs, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var p string
	for _, f := range fs {
		p = filepath.Clean(path + "/" + f.Name())
		if f.IsDir() {
			fs, err := r.readDir(p)
			if err != nil {
				continue
			}
			for _, f := range fs {
				paths = append(paths, f)
			}
		} else {
			if r.checkExt(p) {
				paths = append(paths, p)
			}
		}
	}
	return paths, nil
}

// tplName get template alias from a template file path
func (r *Render) tplName(path string) string {
	if len(path) > len(r.Root) && path[:len(r.Root)] == r.Root {
		path = path[len(r.Root):]
	}
	ext := filepath.Ext(path)
	path = path[:len(path)-len(ext)]
	if runtime.GOOS == "windows" {
		path = strings.Replace(path, "\\", "/", -1)
	}
	return path
}

// checkExt check path extension allow use
func (r *Render) checkExt(path string) bool {
	ext := filepath.Ext(path)
	if ext == "" {
		return false
	}
	for i := range r.Extensions {
		if r.Extensions[i] == ext {
			return true
		}
	}
	return false
}

// parseFile load file and parse to template
func (r *Render) parseFile(path string) error {
	// parse template
	_, err := pongo2.FromCache(path)
	if err != nil {
		return err
	}

	// update indexes
	tpl := r.tplName(path)
	tplIndexes[tpl] = path

	return nil
}

// Error log error
func (r *Render) Error(v interface{}) {
	if r.Baa != nil {
		r.Baa.Logger().Println(v)
	}
}
