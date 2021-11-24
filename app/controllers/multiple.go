package controllers

import (
	"io/ioutil"

	"github.com/revel/revel"
)

type Multiple struct {
	App
}

func (c *Multiple) Upload() revel.Result {
	return c.Render()
}

func (c *Multiple) HandleUpload() revel.Result {

	var files [][]byte
	c.Params.Bind(&files, "file")

	// Make sure at least 2 but no more than 3 files are submitted.
	c.Validation.MinSize(files, 2).Message("You cannot submit less than 2 files")
	c.Validation.MaxSize(files, 3).Message("You cannot submit more than 3 files")

	// Handle errors.
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect((*Multiple).Upload)
	}

	// Prepare result.
	filesInfo := make([]FileInfo, len(files))
	for i, _ := range files {
		f, _ := c.Params.Files["file[]"][i].Open()
		b, _ := ioutil.ReadAll(f)
		ioutil.WriteFile(c.Params.Files["file[]"][i].Filename, b, 0666)
		filesInfo[i] = FileInfo{
			ContentType: c.Params.Files["file[]"][i].Header.Get("Content-Type"),
			Filename:    c.Params.Files["file[]"][i].Filename,
			RealFormat:  "",
			Resolution:  "",
			Size:        len(files[i]),
			Status:      "",
		}
	}

	s := make(map[string]interface{})
	s["Count"] = len(files)
	s["Files"] = filesInfo
	s["Status"] = "Successfully uploaded"
	return c.RenderJSON(s)
}
