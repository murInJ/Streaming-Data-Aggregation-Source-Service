package utils

import (
	"encoding/gob"
	"image"
)

func InitEncoder() {
	gob.Register(image.RGBA{})
}
