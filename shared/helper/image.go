package helper

func AspectRatio(width, height, originalWidth, originalHeight int) (int, int) {
	if width == 0 && height == 0 {
		return originalWidth, originalHeight
	}
	if width == 0 {
		width = (height * originalWidth) / originalHeight
	}
	if height == 0 {
		height = (width * originalHeight) / originalWidth
	}
	return width, height
}
