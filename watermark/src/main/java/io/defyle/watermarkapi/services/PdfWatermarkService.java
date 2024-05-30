package io.defyle.watermarkapi.services;

import io.defyle.watermarkapi.dtos.WatermarkRequest;
import io.defyle.watermarkapi.pojos.StrokeProperties;
import org.apache.pdfbox.cos.COSName;
import org.apache.pdfbox.pdmodel.PDDocument;
import org.apache.pdfbox.pdmodel.PDPage;
import org.apache.pdfbox.pdmodel.PDPageContentStream;
import org.apache.pdfbox.pdmodel.font.PDFont;
import org.apache.pdfbox.pdmodel.font.PDTrueTypeFont;
import org.apache.pdfbox.pdmodel.font.encoding.Encoding;
import org.apache.pdfbox.pdmodel.graphics.blend.BlendMode;
import org.apache.pdfbox.pdmodel.graphics.state.PDExtendedGraphicsState;
import org.apache.pdfbox.util.Matrix;
import org.springframework.core.io.ResourceLoader;
import org.springframework.core.io.support.ResourcePatternUtils;
import org.springframework.stereotype.Service;

import java.io.File;
import java.io.IOException;

@Service
public class PdfWatermarkService {

  private final ResourceLoader resourceLoader;

  public PdfWatermarkService(ResourceLoader resourceLoader) {
    this.resourceLoader = resourceLoader;
  }

  public void generate(WatermarkRequest watermarkRequest) throws IOException {
    PDDocument pdDocument = PDDocument.load(new File(watermarkRequest.getInputFile()));
    pdDocument.setAllSecurityToBeRemoved(true);

    PDFont font = PDTrueTypeFont.load(pdDocument, ResourcePatternUtils.getResourcePatternResolver(resourceLoader)
        .getResource("classpath:Montserrat-Bold.ttf").getInputStream(), Encoding.getInstance(COSName.WIN_ANSI_ENCODING));

    for (PDPage page : pdDocument.getPages()) {
      try (PDPageContentStream pdPageContentStream = new PDPageContentStream(pdDocument, page, PDPageContentStream.AppendMode.APPEND, true, true)) {
        strokeDateTime(watermarkRequest.getDateTime(), page, pdPageContentStream, font);
        strokeUsername(watermarkRequest.getUsername(), page, pdPageContentStream, font);
        strokeWorkspace(watermarkRequest.getWorkspace(), page, pdPageContentStream, font);
      }
    }

    if (new File(watermarkRequest.getOutputFile()).isFile() && new File(watermarkRequest.getOutputFile()).exists()) {
      if (!new File(watermarkRequest.getOutputFile()).delete()) {
        throw new RuntimeException(String.format("Could not delete '%s'.", watermarkRequest.getOutputFile()));
      }
    }
    pdDocument.save(watermarkRequest.getOutputFile());
  }

  private void strokeDateTime(String text, PDPage page, PDPageContentStream pdPageContentStream, PDFont font) throws IOException {
    StrokeProperties strokeProperties = getStrokeProperties(page, font, text);

    pdPageContentStream.transform(Matrix.getRotateInstance(Math.toRadians(270), 0, strokeProperties.getWidth()));
    pdPageContentStream.transform(Matrix.getRotateInstance((float) Math.atan2(strokeProperties.getHeight(), strokeProperties.getWidth()), 0, 0));

    pdPageContentStream.setFont(font, strokeProperties.getFontHeight());
    pdPageContentStream.setGraphicsStateParameters(getPDExtendedGraphicsState(pdPageContentStream));
    setStrokingColor(pdPageContentStream);

    pdPageContentStream.beginText();
    strokeProperties.setY(strokeProperties.getY() - 100);
    pdPageContentStream.newLineAtOffset(strokeProperties.getX(), strokeProperties.getY());
    pdPageContentStream.showText(text);
    pdPageContentStream.endText();
  }

  private void strokeUsername(String text, PDPage page, PDPageContentStream pdPageContentStream, PDFont font) throws IOException {
    StrokeProperties strokeProperties = getStrokeProperties(page, font, text);

    pdPageContentStream.setFont(font, strokeProperties.getFontHeight());
    pdPageContentStream.setGraphicsStateParameters(getPDExtendedGraphicsState(pdPageContentStream));
    setStrokingColor(pdPageContentStream);

    pdPageContentStream.beginText();
    pdPageContentStream.newLineAtOffset(strokeProperties.getX(), strokeProperties.getY());
    pdPageContentStream.showText(text);
    pdPageContentStream.endText();
  }

  private void strokeWorkspace(String text, PDPage page, PDPageContentStream pdPageContentStream, PDFont font) throws IOException {
    StrokeProperties strokeProperties = getStrokeProperties(page, font, text);

    pdPageContentStream.setFont(font, strokeProperties.getFontHeight());
    pdPageContentStream.setGraphicsStateParameters(getPDExtendedGraphicsState(pdPageContentStream));
    setStrokingColor(pdPageContentStream);

    pdPageContentStream.beginText();
    strokeProperties.setY(strokeProperties.getY() + 100);
    pdPageContentStream.newLineAtOffset(strokeProperties.getX(), strokeProperties.getY());
    pdPageContentStream.showText(text);
    pdPageContentStream.endText();
  }

  private PDExtendedGraphicsState getPDExtendedGraphicsState(PDPageContentStream pdPageContentStream) throws IOException {
    PDExtendedGraphicsState pdExtendedGraphicsState = new PDExtendedGraphicsState();
    pdExtendedGraphicsState.setNonStrokingAlphaConstant(0.2f);
    pdExtendedGraphicsState.setStrokingAlphaConstant(0.2f);
    pdExtendedGraphicsState.getCOSObject().setItem(COSName.BM, COSName.MULTIPLY);
    pdExtendedGraphicsState.setBlendMode(BlendMode.MULTIPLY); // will work in 2.0.14
    pdPageContentStream.setGraphicsStateParameters(pdExtendedGraphicsState);

    return pdExtendedGraphicsState;
  }

  private void setStrokingColor(PDPageContentStream pdPageContentStream) throws IOException {
    // Some API weirdness here. When int, range is 0..255
    // When float, this would be 0..1f
    pdPageContentStream.setNonStrokingColor(255, 0, 0);
    pdPageContentStream.setStrokingColor(255, 0, 0);
  }

  private StrokeProperties getStrokeProperties(PDPage page, PDFont font, String text) throws IOException {
    StrokeProperties strokeProperties = new StrokeProperties();
    strokeProperties.setWidth(page.getMediaBox().getHeight());
    strokeProperties.setHeight(page.getMediaBox().getWidth());
    strokeProperties.setFontHeight(72);
    strokeProperties.setStringWidth(font.getStringWidth(text) / 1000 * strokeProperties.getFontHeight());
    strokeProperties.setDiagonalLength((float) Math.sqrt(strokeProperties.getWidth() * strokeProperties.getWidth() +
        strokeProperties.getHeight() * strokeProperties.getHeight()));
    strokeProperties.setX((strokeProperties.getDiagonalLength() - strokeProperties.getStringWidth()) / 2);
    strokeProperties.setY(-strokeProperties.getFontHeight() / 4);

    return strokeProperties;
  }
}
