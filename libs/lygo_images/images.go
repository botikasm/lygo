package lygo_images

import (
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/base/lygo_strings"
	"gopkg.in/gographics/imagick.v3/imagick"
	"strings"
	"sync"
)

var mutex sync.Mutex

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Check if passed extension has an ALPHA channel
// @param extension
func HasAlphaChannel(extensionOrFilename string) bool {
	tokens := strings.Split(extensionOrFilename, ".")
	var extension string
	if extension = tokens[0]; len(tokens) == 2 {
		extension = tokens[1]
	}
	for _, v := range ALPHA_IMAGES {
		if v == extension {
			return true
		}
	}
	return false
}

// ConvertPdfToJpg will take a filename of a pdf file and convert the file into an
// image which will be saved back to the same location. It will save the image as a
// high resolution jpg file with minimal compression.
func Convert(params *ImageConvertParams) (pages []string, err error) {

	defer mutex.Unlock()
	mutex.Lock()

	pages = make([]string, 0)

	// Setup
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// Must be *before* ReadImageFile
	// Make sure our image is high quality
	if err := mw.SetResolution(params.XRes, params.YRes); err != nil {
		return pages, err
	}

	// Load the image file into imagick
	if err := mw.ReadImage(params.Source); err != nil {
		return pages, err
	}

	// Must be *after* ReadImageFile
	// Flatten image and remove alpha channel, to prevent alpha turning black in jpg
	if HasAlphaChannel(params.Format) {
		if err := mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_REMOVE); err != nil {
			return pages, err
		}
	}

	// Set any compression (100 = max quality)
	if err := mw.SetCompressionQuality(params.Quality); err != nil {
		return pages, err
	}

	// Convert into JPG
	if err := mw.SetFormat(params.Format); err != nil {
		return pages, err
	}

	// ensure target folder
	lygo_paths.Mkdir(params.Target)

	count := int(mw.GetIteratorIndex()) + 1
	for i := 0; i < count; i++ {
		// Select only first page of pdf
		mw.SetIteratorIndex(i)

		target := lygo_paths.ChangeFileNameWithSuffix(params.Target, lygo_strings.Format("-%s", i))
		err = mw.WriteImage(target)
		if nil==err {
			pages = append(pages, target)
		}else {
			pages = append(pages, "")
		}
	}

	return pages, err
}

func Crop(params *ImageCropParams) (err error) {

	defer mutex.Unlock()
	mutex.Lock()

	// Setup
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// Load the image file into imagick
	if err = mw.ReadImage(params.Source); err != nil {
		return err
	}

	err = mw.CropImage(params.Width, params.Height, params.X, params.Y)

	return write(params.Source, params.Target, mw)
}

func Rotate(params *ImageRotateParams) (err error) {

	defer mutex.Unlock()
	mutex.Lock()

	if params.Degree != 0 {
		// Setup
		imagick.Initialize()
		defer imagick.Terminate()

		mw := imagick.NewMagickWand()
		defer mw.Destroy()

		// Load the image file into imagick
		if err = mw.ReadImage(params.Source); err != nil {
			return err
		}

		px := imagick.NewPixelWand()
		err = mw.RotateImage(px, params.Degree)

		return write(params.Source, params.Target, mw)
	} else {
		// no rotation
		if params.Target != params.Source {
			// just copy file
			_, err = lygo_io.CopyFile(params.Source, params.Target)
		}
	}
	return err
}

func AutoOrient(params *ImageParams) (err error) {

	defer mutex.Unlock()
	mutex.Lock()

	// Setup
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// Load the image file into imagick
	if err = mw.ReadImage(params.Source); err != nil {
		return err
	}

	err = mw.AutoOrientImage()

	return write(params.Source, params.Target, mw)
}

func AutoLevel(params *ImageParams) (err error) {

	defer mutex.Unlock()
	mutex.Lock()

	// Setup
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// Load the image file into imagick
	if err = mw.ReadImage(params.Source); err != nil {
		return err
	}

	err = mw.AutoLevelImage()

	return write(params.Source, params.Target, mw)
}

func AutoGamma(params *ImageParams) (err error) {

	defer mutex.Unlock()
	mutex.Lock()

	// Setup
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// Load the image file into imagick
	if err = mw.ReadImage(params.Source); err != nil {
		return err
	}

	err = mw.AutoGammaImage()

	return write(params.Source, params.Target, mw)
}

func Contrast(params *ImageParams) (err error) {

	defer mutex.Unlock()
	mutex.Lock()

	// Setup
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// Load the image file into imagick
	if err = mw.ReadImage(params.Source); err != nil {
		return err
	}

	err = mw.ContrastImage(false)

	return write(params.Source, params.Target, mw)
}

func ContrastSharpen(params *ImageParams) (err error) {

	defer mutex.Unlock()
	mutex.Lock()

	// Setup
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// Load the image file into imagick
	if err = mw.ReadImage(params.Source); err != nil {
		return err
	}

	err = mw.ContrastImage(true)

	return write(params.Source, params.Target, mw)
}

func Sharpen(params *ImageSharpenParams) (err error) {

	defer mutex.Unlock()
	mutex.Lock()

	// Setup
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// Load the image file into imagick
	if err = mw.ReadImage(params.Source); err != nil {
		return err
	}

	// adjust Radius
	if params.Radius == 0 {
		params.Radius = params.Sigma * 2
	}

	err = mw.SharpenImage(params.Radius, params.Sigma)

	return write(params.Source, params.Target, mw)
}

func AdaptiveSharpen(params *ImageSharpenParams) (err error) {

	defer mutex.Unlock()
	mutex.Lock()

	// Setup
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// Load the image file into imagick
	if err = mw.ReadImage(params.Source); err != nil {
		return err
	}

	// adjust Radius
	if params.Radius == 0 {
		params.Radius = params.Sigma * 2
	}

	err = mw.AdaptiveSharpenImage(params.Radius, params.Sigma)

	return write(params.Source, params.Target, mw)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func write(source, target string, mw *imagick.MagickWand) (err error) {
	if len(target) > 0 {
		lygo_paths.Mkdir(target)
		err = mw.WriteImage(target)
	} else {
		err = mw.WriteImage(source)
	}
	return err
}
