package lygo_images

var ALPHA_IMAGES = [...]string{"png", "gif"}


type ImageSource struct {
	Source string
}
type ImageTarget struct {
	Target string
}
type ImageFormat struct {
	Format string
	Quality uint
}
type ImageResolution struct {
	XRes float64
	YRes float64
}
type ImageRotation struct {
	Degree float64
}
type ImageCrop struct {
	X int
	Y int
	Width uint
	Height uint
}
type ImageSharpen struct {
	Radius float64
	Sigma float64
}

type ImageParams struct {
	ImageSource
	ImageTarget
}
type ImageSharpenParams struct {
	ImageSource
	ImageTarget
	ImageSharpen
}
type ImageConvertParams struct {
	ImageSource
	ImageTarget
	ImageFormat
	ImageResolution
}

type ImageCropParams struct {
	ImageSource
	ImageTarget
	ImageCrop
}

type ImageRotateParams struct {
	ImageSource
	ImageTarget
	ImageRotation
}

func NewImageParams() *ImageParams {
	result := new (ImageParams)

	return result
}

func NewImageSharpenParams() *ImageSharpenParams {
	result := new (ImageSharpenParams)
	result.Sigma = 1
	result.Radius = 0 // (auto) result.Sigma*2
	return result
}

func NewImageConvertParams() *ImageConvertParams {
	result := new (ImageConvertParams)

	// ImageFormat
	result.Format = "jpg"
	result.Quality = 95

	// ImageResolution
	result.XRes = 300.0
	result.YRes = 300.0

	return result
}

func NewImageCropParams() *ImageCropParams {
	result := new (ImageCropParams)

	result.X = 0
	result.Y = 0
	result.Width = 100
	result.Height = 100

	return result
}

func NewImageRotateParams() *ImageRotateParams {
	result := new (ImageRotateParams)

	result.Degree = 90.0

	return result
}