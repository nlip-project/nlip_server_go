package models

type Subformat string

const (
	English   Subformat = "english"
	Spanish   Subformat = "spanish"
	German    Subformat = "german"
	JSON      Subformat = "json"
	URI       Subformat = "uri"
	XML       Subformat = "xml"
	HTML      Subformat = "html"
	ImageBMP  Subformat = "image/bmp"
	ImageGIF  Subformat = "image/gif"
	ImageJPEG Subformat = "image/jpeg"
	ImageJPG  Subformat = "image/jpg"
	ImagePNG  Subformat = "image/png"
	ImageTIFF Subformat = "image/tiff"
	AudioMP3  Subformat = "audio/mp3"
	TextSF    Subformat = "text" // need SF for subformat. better way?
	GPS       Subformat = "gps"
	// TODO: Think about the generic ones.
)
