package com.voltaserve.watermark.services;

import com.voltaserve.watermark.dtos.WatermarkRequest;
import io.minio.MinioClient;
import io.minio.UploadObjectArgs;
import org.apache.commons.io.FilenameUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.io.ResourceLoader;
import org.springframework.core.io.support.ResourcePatternUtils;
import org.springframework.stereotype.Service;

import javax.imageio.ImageIO;
import javax.imageio.ImageReader;
import javax.imageio.stream.ImageInputStream;
import java.awt.*;
import java.awt.geom.AffineTransform;
import java.awt.geom.Rectangle2D;
import java.awt.image.BufferedImage;
import java.io.BufferedOutputStream;
import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Iterator;
import java.util.Objects;
import java.util.UUID;

@Service
public class ImageWatermarkService {

  private final Logger logger = LoggerFactory.getLogger(PdfWatermarkService.class);

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

  public ImageWatermarkService(ResourceLoader resourceLoader) throws IOException, FontFormatException {
    GraphicsEnvironment graphicsEnvironment = GraphicsEnvironment.getLocalGraphicsEnvironment();
    graphicsEnvironment.registerFont(
        Font.createFont(Font.TRUETYPE_FONT, ResourcePatternUtils.getResourcePatternResolver(resourceLoader)
            .getResource("classpath:Unbounded-Medium.ttf").getInputStream()));
  }

  public void generate(WatermarkRequest request) throws IOException {
    File inputFile = new File(request.getPath());

    BufferedImage bufferedImage = ImageIO.read(inputFile);
    int width = bufferedImage.getWidth();
    int height = bufferedImage.getHeight();
    Graphics2D g2d = (Graphics2D) bufferedImage.getGraphics();

    g2d.setComposite(AlphaComposite.getInstance(AlphaComposite.SRC_OVER, 0.2f));
    g2d.setColor(Color.RED);
    g2d.setRenderingHint(RenderingHints.KEY_TEXT_ANTIALIASING, RenderingHints.VALUE_TEXT_ANTIALIAS_ON);
    g2d.setFont(new Font("Unbounded", Font.BOLD, 120));

    Rectangle2D dateTimeBounds = g2d.getFontMetrics().getStringBounds(request.getDateTime(), g2d);
    Rectangle2D workspaceBounds = g2d.getFontMetrics().getStringBounds(request.getWorkspace(), g2d);
    Rectangle2D usernameBounds = g2d.getFontMetrics().getStringBounds(request.getUsername(), g2d);

    // Move our origin to the center
    g2d.translate(width / 2.0f, height / 2.0f);

    /* Create a rotation transform to rotate the text based on
       the diagonal angle of the picture aspect ratio */
    AffineTransform affineTransform = new AffineTransform();
    affineTransform.rotate(Math.atan2(height, width));
    g2d.transform(affineTransform);

    /* Reposition our coordinates based on size (same as we would normally
       do to center on straight line but based on starting at center */
    float x1 = (int) workspaceBounds.getWidth() / 2.0f * -1;
    float y1 = (int) dateTimeBounds.getHeight() / 2.0f;
    y1 -= 240;
    g2d.translate(x1, y1);
    g2d.drawString(request.getWorkspace(), 0.0f, 0.0f);

    float x2 = (int) usernameBounds.getWidth() / 2.0f * -1;
    g2d.translate(x2 - x1, 120);
    g2d.drawString(request.getUsername(), 0.0f, 0.0f);

    float x3 = (int) dateTimeBounds.getWidth() / 2.0f * -1;
    g2d.translate(x3 - x2, 120);
    g2d.drawString(request.getDateTime(), 0.0f, 0.0f);

    g2d.dispose();

    Path output = Paths.get(
            System.getProperty("java.io.tmpdir"),
            UUID.randomUUID() + "." + FilenameUtils.getExtension(request.getPath()));
    ImageIO.write(bufferedImage, Objects.requireNonNull(getFormatName(inputFile)), new BufferedOutputStream(
            new FileOutputStream(output.toString())));

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
                      .object(request.getS3Key())
                      .filename(output.toString())
                      .build());
    } catch (Exception e) {
      this.logger.error(e.getMessage(), e);
    } finally {
      Files.deleteIfExists(output);
    }
  }

  private String getFormatName(File file) throws IOException {
    ImageInputStream stream = ImageIO.createImageInputStream(file);
    Iterator<ImageReader> readers = ImageIO.getImageReaders(stream);
    if (readers.hasNext()) {
      ImageReader imageReader = readers.next();
      return imageReader.getFormatName();
    }

    return null;
  }
}
