package docx

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/xakepp35/pkg/xerrors"
)

// ImageManager управляет изображениями в документе
type ImageManager struct {
	images []*bytes.Buffer
	count  int
}

// NewImageManager создает новый менеджер изображений
func NewImageManager() *ImageManager {
	return &ImageManager{
		images: make([]*bytes.Buffer, 0),
		count:  0,
	}
}

// AddImage добавляет изображение и возвращает его ID и filename
func (im *ImageManager) AddImage(imageData []byte) (string, string, error) {
	if len(imageData) == 0 {
		return "", "", xerrors.Err(nil).Msg("image data is empty")
	}

	im.count++
	imageBuffer := bytes.NewBuffer(imageData)
	imageFormat := DetectImageFormat(imageData)
	filename := fmt.Sprintf("image_embed%d.%s", im.count, imageFormat)
	imageID := fmt.Sprintf("rImageId%d", im.count)

	im.images = append(im.images, imageBuffer)

	return imageID, filename, nil
}

// GetImages возвращает все добавленные изображения
func (im *ImageManager) GetImages() []*bytes.Buffer {
	return im.images
}

// Reset сбрасывает состояние менеджера
func (im *ImageManager) Reset() {
	im.images = make([]*bytes.Buffer, 0)
	im.count = 0
}

// Count возвращает количество добавленных изображений
func (im *ImageManager) Count() int {
	return im.count
}

// DetectImageFormat определяет формат изображения по содержимому
func DetectImageFormat(data []byte) string {
	switch http.DetectContentType(data) {
	case "image/png":
		return "png"
	case "image/jpeg":
		return "jpg"
	case "image/gif":
		return "gif"
	case "image/webp":
		return "webp"
	case "image/bmp":
		return "bmp"
	default:
		return "png"
	}
}

// CreateImageXML создает XML для вставки изображения
func CreateImageXML(relID string, width, height int) string {
	const dpi = 96
	const emuPerInch = 914400
	cx := width * emuPerInch / dpi
	cy := height * emuPerInch / dpi

	return fmt.Sprintf(`
<w:p>
  <w:r>
    <w:drawing>
      <wp:anchor behindDoc="0" distT="0" distB="0" distL="0" distR="0" simplePos="0" locked="0"
                 layoutInCell="0" allowOverlap="1" relativeHeight="0">
        <wp:simplePos x="0" y="0"/>
        <wp:positionH relativeFrom="column">
          <wp:align>center</wp:align>
        </wp:positionH>
        <wp:positionV relativeFrom="paragraph">
          <wp:align>bottom</wp:align>
        </wp:positionV>
        <wp:extent cx="%[2]d" cy="%[3]d"/>
        <wp:effectExtent l="0" t="0" r="0" b="0"/>
        <wp:wrapTopAndBottom/>
        <wp:docPr id="1" name="Image1"/>
        <wp:cNvGraphicFramePr>
          <a:graphicFrameLocks xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
                               noChangeAspect="1"/>
        </wp:cNvGraphicFramePr>
        <a:graphic xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
          <a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture">
            <pic:pic xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">
              <pic:nvPicPr>
                <pic:cNvPr id="1" name="Image1"/>
                <pic:cNvPicPr/>
              </pic:nvPicPr>
              <pic:blipFill>
                <a:blip r:embed="%[1]s"/>
                <a:stretch><a:fillRect/></a:stretch>
              </pic:blipFill>
              <pic:spPr>
                <a:xfrm>
                  <a:off x="0" y="0"/>
                  <a:ext cx="%[2]d" cy="%[3]d"/>
                </a:xfrm>
                <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
              </pic:spPr>
            </pic:pic>
          </a:graphicData>
        </a:graphic>
      </wp:anchor>
    </w:drawing>
  </w:r>
</w:p>
`, relID, cx, cy)
}
