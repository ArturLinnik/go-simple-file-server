package schemas

import "time"

type File struct {
	Name    string    `json:"name"`
	IsDir   bool      `json:"is_dir"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
}

type Response struct {
	Files []File `json:"files"`
	Path  string `json:"path"`
}
