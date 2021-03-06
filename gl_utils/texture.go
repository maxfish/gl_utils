package gl_utils

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	// Used only to initialize the JPEG subsystem
	_ "image/jpeg"
	// Used only to initialize the PNG subsystem
	_ "image/png"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// Texture a representation of an image file in memory
type Texture struct {
	id     uint32
	width  int32
	height int32
}

// NewTextureFromFile loads the image from a file into a texture
func NewTextureFromFile(filePath string) *Texture {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error loading texture. %s\n", err)
		return nil
	}
	defer file.Close()

	decodedImage, format, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Error decoding <%s> image: '%s'\n", format, filePath)
		return nil
	}
	return NewTextureFromImage(decodedImage)
}

// NewTextureFromImage uses the data from an Image struct to create a texture
func NewTextureFromImage(imageData image.Image) *Texture {
	texture := &Texture{
		width:  int32(imageData.Bounds().Dx()),
		height: int32(imageData.Bounds().Dy()),
	}
	gl.GenTextures(1, &texture.id)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture.id)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	switch imageData.(type) {
	case *image.Gray16:
		// 16-bit monochrome image --> Gray
		grayImage := image.NewGray(imageData.Bounds())
		if grayImage.Stride != grayImage.Rect.Size().X*1 {
			fmt.Println("Error creating texture: unsupported stride")
			return nil
		}
		draw.Draw(grayImage, grayImage.Bounds(), imageData, image.Point{0, 0}, draw.Src)
		gl.TexImage2D(
			gl.TEXTURE_2D, 0, gl.RED, texture.width, texture.height,
			0, gl.RED, gl.UNSIGNED_BYTE, gl.Ptr(grayImage.Pix),
		)
	case *image.NRGBA:
		// non-alpha-premultiplied 32-bit color image --> RGBA
		pixelData := imageData.(*image.NRGBA).Pix
		gl.TexImage2D(
			gl.TEXTURE_2D, 0, gl.RGBA, texture.width, texture.height,
			0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixelData),
		)
	default:
		// All the other formats -->  RGBA
		rgba := image.NewRGBA(imageData.Bounds())
		if rgba.Stride != rgba.Rect.Size().X*4 {
			fmt.Println("Error creating texture: unsupported stride")
			return nil
		}
		draw.Draw(rgba, rgba.Bounds(), imageData, image.Point{0, 0}, draw.Src)
		gl.TexImage2D(
			gl.TEXTURE_2D, 0, gl.RGBA, texture.width, texture.height,
			0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix),
		)
	}

	gl.BindTexture(gl.TEXTURE_2D, 0)

	return texture
}

// NewEmptyTexture creates an empty texture with a specified size
func NewEmptyTexture(width int, height int, pixelFormat int32) (*Texture, error) {
	bounds := image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: width, Y: height},
	}
	imageData := image.NewRGBA(bounds)

	texture := &Texture{
		width:  int32(imageData.Bounds().Dx()),
		height: int32(imageData.Bounds().Dy()),
	}
	gl.GenTextures(1, &texture.id)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture.id)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, pixelFormat, texture.width, texture.height,
		0, uint32(pixelFormat), gl.UNSIGNED_BYTE, gl.Ptr(imageData.Pix),
	)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return texture, nil
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.id)
}

func (t *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// ID returns the unique OpenGL ID of this texture
func (t *Texture) ID() uint32 {
	return t.id
}

// Width returns the texture width in pixels
func (t *Texture) Width() int32 {
	return t.width
}

// Height returns the texture width in pixels
func (t *Texture) Height() int32 {
	return t.height
}
