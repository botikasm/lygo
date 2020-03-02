package lygo_images

import (
	"fmt"
	"github.com/botikasm/lygo/_tests"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/libs/lygo_images"
	"path/filepath"
	"sync"
	"testing"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Test PDF -> JPG conversion
func TestConvertPdfToJpeg(t *testing.T) {

	_tests.InitContext()

	var group sync.WaitGroup

	group.Add(3)

	/**/
	go func() {
		defer group.Done()
		fmt.Println("CONVERT")
		pdfName := lygo_paths.WorkspacePath("./resources/sample.pdf")
		imageName := lygo_paths.WorkspacePath("./resources/output/sample.jpg")
		lygo_paths.Mkdir(imageName)
		_, err := convertPdfToJpg(pdfName, imageName)
		if err != nil {
			t.Error(err)
		}
	}()

	/**/
	go func() {
		defer group.Done()
		fmt.Println("ROTATE")
		imageName := lygo_paths.WorkspacePath("./resources/image001.jpg")
		lygo_paths.Mkdir(imageName)
		err := rotate(imageName)
		if err != nil {
			t.Error(err)
		}
	}()

	/**/
	go func() {
		defer group.Done()
		fmt.Println("CROP")
		imageName := lygo_paths.WorkspacePath("./resources/image001.jpg")
		lygo_paths.Mkdir(imageName)
		if err := crop(imageName); err != nil {
			t.Error(err)
		}
	}()

	group.Wait()
}

func TestRotate(t *testing.T) {

	_tests.InitContext()

	imageName := lygo_paths.WorkspacePath("./resources/image001.jpg")

	fmt.Println("ROTATE")
	if err := rotate(imageName); err != nil {
		t.Error(err)
	}

}

func TestAutoOrient(t *testing.T) {

	_tests.InitContext()

	imageName := "../_workspace/output/rotate_cv_converted.jpeg"

	fmt.Println("AUTO ORIENT")
	if err := autoOrient(imageName); err != nil {
		t.Error(err)
	}

}

func TestAutoLevel(t *testing.T) {

	_tests.InitContext()

	imageName := "../_workspace/output/cv_converted.jpeg"

	fmt.Println("AUTO LEVEL")
	if err := autoLevel(imageName); err != nil {
		t.Error(err)
	}

}

func TestAutoGamma(t *testing.T) {

	_tests.InitContext()

	imageName := "../_workspace/output/cv_converted.jpeg"

	fmt.Println("AUTO GAMMA")
	if err := autoGamma(imageName); err != nil {
		t.Error(err)
	}

}

func TestContrast(t *testing.T) {

	_tests.InitContext()

	imageName := "../_workspace/output/cv_converted.jpeg"

	fmt.Println("CONTRAST")
	if err := contrast(imageName); err != nil {
		t.Error(err)
	}

}

func TestSharpen(t *testing.T) {

	_tests.InitContext()

	imageName := "../_workspace/output/cv_converted.jpeg"

	fmt.Println("SHARPEN")
	if err := sharpen(imageName); err != nil {
		t.Error(err)
	}

}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

// ConvertPdfToJpg will take a filename of a pdf file and convert the file into an
// image which will be saved back to the same location. It will save the image as a
// high resolution jpg file with minimal compression.
func convertPdfToJpg(pdfName string, imageName string) ([]string, error) {

	params := lygo_images.NewImageConvertParams()
	params.Source = pdfName
	params.Target = imageName
	params.Format = "jpg"
	params.XRes = 200
	params.YRes = 200

	pages, err := lygo_images.Convert(params)

	return pages, err
}

func crop(imageName string) error {
	params := lygo_images.NewImageCropParams()
	params.Source = imageName
	params.Target = filepath.Join(filepath.Dir(imageName), "crop_"+filepath.Base(imageName))
	params.Width = 60000
	params.Height = 300

	return lygo_images.Crop(params)
}

func rotate(imageName string) error {
	params := lygo_images.NewImageRotateParams()
	params.Source = imageName
	params.Target = filepath.Join(filepath.Dir(imageName), "rotate_"+filepath.Base(imageName))
	params.Degree = 90

	return lygo_images.Rotate(params)
}

func autoOrient(imageName string) error {
	params := lygo_images.NewImageParams()
	params.Source = imageName
	params.Target = filepath.Join(filepath.Dir(imageName), "autoorient_"+filepath.Base(imageName))

	return lygo_images.AutoOrient(params)
}

func autoLevel(imageName string) error {
	params := lygo_images.NewImageParams()
	params.Source = imageName
	params.Target = filepath.Join(filepath.Dir(imageName), "autolevel_"+filepath.Base(imageName))

	return lygo_images.AutoLevel(params)
}

func autoGamma(imageName string) error {
	params := lygo_images.NewImageParams()
	params.Source = imageName
	params.Target = filepath.Join(filepath.Dir(imageName), "autogamma_"+filepath.Base(imageName))

	return lygo_images.AutoGamma(params)
}

func contrast(imageName string) error {
	params := lygo_images.NewImageParams()
	params.Source = imageName
	params.Target = filepath.Join(filepath.Dir(imageName), "contrast_"+filepath.Base(imageName))

	return lygo_images.ContrastSharpen(params)
}

func sharpen(imageName string) error {
	params := lygo_images.NewImageSharpenParams()
	params.Source = imageName
	params.Target = filepath.Join(filepath.Dir(imageName), "sharpen_"+filepath.Base(imageName))

	return lygo_images.Sharpen(params)
}
