package com.voltaserve.watermark.services;

import com.voltaserve.watermark.dtos.WatermarkRequest;
import com.voltaserve.watermark.pojos.StrokeProperties;
import io.minio.MinioClient;
import io.minio.UploadObjectArgs;
import io.minio.errors.MinioException;
import org.apache.commons.io.FilenameUtils;
import org.apache.pdfbox.Loader;
import org.apache.pdfbox.cos.COSName;
import org.apache.pdfbox.io.RandomAccessReadBufferedFile;
import org.apache.pdfbox.pdmodel.PDDocument;
import org.apache.pdfbox.pdmodel.PDPage;
import org.apache.pdfbox.pdmodel.PDPageContentStream;
import org.apache.pdfbox.pdmodel.font.PDFont;
import org.apache.pdfbox.pdmodel.font.PDTrueTypeFont;
import org.apache.pdfbox.pdmodel.font.encoding.Encoding;
import org.apache.pdfbox.pdmodel.graphics.blend.BlendMode;
import org.apache.pdfbox.pdmodel.graphics.state.PDExtendedGraphicsState;
import org.apache.pdfbox.util.Matrix;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.io.ResourceLoader;
import org.springframework.core.io.support.ResourcePatternUtils;
import org.springframework.stereotype.Service;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.UUID;

@Service
public class PdfWatermarkService {

  private final Logger logger = LoggerFactory.getLogger(PdfWatermarkService.class);

  @Autowired
  private ResourceLoader resourceLoader;

  @Value("${app.s3.url}")
  private String s3Url;

  @Value("${app.s3.access-key}")
  private String s3AccessKey;

  @Value("${app.s3.secret-key}")
  private String s3SecretKey;

  @Value("${app.s3.region}")
  private String s3Region;

  @Value("${app.s3.secure}")
  private Boolean s3Secure;

  public void generate(WatermarkRequest request) throws IOException {
    PDDocument document = Loader.loadPDF(new RandomAccessReadBufferedFile(new File(request.getPath())));
    document.setAllSecurityToBeRemoved(true);

    PDFont font = PDTrueTypeFont.load(document, ResourcePatternUtils.getResourcePatternResolver(resourceLoader)
        .getResource("classpath:Unbounded-Medium.ttf").getInputStream(),
        Encoding.getInstance(COSName.WIN_ANSI_ENCODING));

    for (PDPage page : document.getPages()) {
      try (PDPageContentStream pdPageContentStream = new PDPageContentStream(document, page,
          PDPageContentStream.AppendMode.APPEND, true, true)) {
        strokeDateTime(request.getDateTime(), page, pdPageContentStream, font);
        strokeUsername(request.getUsername(), page, pdPageContentStream, font);
        strokeWorkspace(request.getWorkspace(), page, pdPageContentStream, font);
      }
    }

    Path output = Paths.get(
            System.getProperty("java.io.tmpdir"),
            UUID.randomUUID() + "." + FilenameUtils.getExtension(request.getPath()));
    document.save(output.toString());

    String endpoint = this.s3Url.split(":")[0];
    int port = Integer.parseInt(this.s3Url.split(":")[1]);
    try(MinioClient minioClient = MinioClient.builder()
            .endpoint(endpoint, port, this.s3Secure)
            .credentials(this.s3AccessKey, this.s3SecretKey)
            .region(this.s3Region)
            .build()) {
      minioClient.uploadObject(
              UploadObjectArgs.builder()
                      .bucket(request.getS3Bucket())
                      .object(request.getS3Key() + "/watermark." + FilenameUtils.getExtension(request.getPath()))
                      .filename(output.toString())
                      .build());
    } catch (Exception e) {
      this.logger.error(e.getMessage(), e);
    } finally {
      Files.deleteIfExists(output);
    }
  }

  private void strokeDateTime(String text, PDPage page, PDPageContentStream stream, PDFont font)
      throws IOException {
    StrokeProperties props = getStrokeProperties(page, font, text);

    stream.transform(Matrix.getRotateInstance(Math.toRadians(270), 0, props.getWidth()));
    stream.transform(
        Matrix.getRotateInstance((float) Math.atan2(props.getHeight(), props.getWidth()), 0, 0));

    stream.setFont(font, props.getFontHeight());
    stream.setGraphicsStateParameters(getPDExtendedGraphicsState(stream));
    setStrokingColor(stream);

    stream.beginText();
    props.setY(props.getY() - 100);
    stream.newLineAtOffset(props.getX(), props.getY());
    stream.showText(text);
    stream.endText();
  }

  private void strokeUsername(String text, PDPage page, PDPageContentStream stream, PDFont font)
      throws IOException {
    StrokeProperties props = getStrokeProperties(page, font, text);

    stream.setFont(font, props.getFontHeight());
    stream.setGraphicsStateParameters(getPDExtendedGraphicsState(stream));
    setStrokingColor(stream);

    stream.beginText();
    stream.newLineAtOffset(props.getX(), props.getY());
    stream.showText(text);
    stream.endText();
  }

  private void strokeWorkspace(String text, PDPage page, PDPageContentStream stream, PDFont font)
      throws IOException {
    StrokeProperties props = getStrokeProperties(page, font, text);

    stream.setFont(font, props.getFontHeight());
    stream.setGraphicsStateParameters(getPDExtendedGraphicsState(stream));
    setStrokingColor(stream);

    stream.beginText();
    props.setY(props.getY() + 100);
    stream.newLineAtOffset(props.getX(), props.getY());
    stream.showText(text);
    stream.endText();
  }

  private PDExtendedGraphicsState getPDExtendedGraphicsState(PDPageContentStream stream)
      throws IOException {
    PDExtendedGraphicsState state = new PDExtendedGraphicsState();
    state.setNonStrokingAlphaConstant(0.2f);
    state.setStrokingAlphaConstant(0.2f);
    state.getCOSObject().setItem(COSName.BM, COSName.MULTIPLY);
    state.setBlendMode(BlendMode.MULTIPLY); // will work in 2.0.14
    stream.setGraphicsStateParameters(state);

    return state;
  }

  private void setStrokingColor(PDPageContentStream stream) throws IOException {
    stream.setNonStrokingColor(1, 0, 0);
    stream.setStrokingColor(1, 0, 0);
  }

  private StrokeProperties getStrokeProperties(PDPage page, PDFont font, String text) throws IOException {
    StrokeProperties props = new StrokeProperties();
    props.setWidth(page.getMediaBox().getHeight());
    props.setHeight(page.getMediaBox().getWidth());
    props.setFontHeight(72);
    props.setStringWidth(font.getStringWidth(text) / 1000 * props.getFontHeight());
    props.setDiagonalLength((float) Math.sqrt(props.getWidth() * props.getWidth() +
        props.getHeight() * props.getHeight()));
    props.setX((props.getDiagonalLength() - props.getStringWidth()) / 2);
    props.setY(-props.getFontHeight() / 4);

    return props;
  }
}
