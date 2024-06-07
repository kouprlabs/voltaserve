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

import java.awt.*;
import java.awt.geom.AffineTransform;
import java.io.BufferedOutputStream;
import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.Objects;
import java.util.UUID;

import javax.imageio.ImageIO;

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
    GraphicsEnvironment.getLocalGraphicsEnvironment().registerFont(
        Font.createFont(Font.TRUETYPE_FONT, ResourcePatternUtils.getResourcePatternResolver(resourceLoader)
            .getResource("classpath:Unbounded-Medium.ttf").getInputStream()));
  }

  public void generate(WatermarkRequest request) throws IOException {
    var inputFile = new File(request.getPath());

    var bufferedImage = ImageIO.read(inputFile);
    int width = bufferedImage.getWidth();
    int height = bufferedImage.getHeight();
    var g2d = (Graphics2D) bufferedImage.getGraphics();

    g2d.setComposite(AlphaComposite.getInstance(AlphaComposite.SRC_OVER, 0.2f));
    g2d.setColor(Color.RED);
    g2d.setRenderingHint(RenderingHints.KEY_TEXT_ANTIALIASING, RenderingHints.VALUE_TEXT_ANTIALIAS_ON);
    g2d.setFont(new Font("Unbounded", Font.BOLD, 120));

    // Move our origin to the center
    g2d.translate(width / 2.0f, height / 2.0f);

    /*
     * Create a rotation transform to rotate the text based on
     * the diagonal angle of the picture aspect ratio
     */
    var affineTransform = new AffineTransform();
    affineTransform.rotate(Math.atan2(height, width));
    g2d.transform(affineTransform);

    /* Calculate total height of words */
    int totalHeightOfWords = 0;
    for (var value : request.getValues()) {
      totalHeightOfWords += g2d.getFontMetrics().getStringBounds(value, g2d).getHeight();
    }

    // Move our origin slightly to the top based on total height of words
    g2d.translate(0, -totalHeightOfWords / 2);

    /*
     * Reposition our coordinates based on size (same as we would normally
     * do to center on straight line but based on starting at center
     */
    for (var value : request.getValues()) {
      var bounds = g2d.getFontMetrics().getStringBounds(value, g2d);
      float x = (int) bounds.getWidth() / 2.0f;

      g2d.translate(-x, bounds.getHeight());
      g2d.drawString(value, 0.0f, 0.0f);

      // Reset x coordinate
      g2d.translate(x, 0);
    }
    g2d.dispose();

    var output = Paths.get(
        System.getProperty("java.io.tmpdir"),
        UUID.randomUUID() + "." + FilenameUtils.getExtension(request.getPath()));
    ImageIO.write(bufferedImage, Objects.requireNonNull(getFormatName(inputFile)), new BufferedOutputStream(
        new FileOutputStream(output.toString())));

    var endpoint = this.s3Url.split(":")[0];
    int port = Integer.parseInt(this.s3Url.split(":")[1]);
    try (var minioClient = MinioClient.builder()
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
    var stream = ImageIO.createImageInputStream(file);
    var readers = ImageIO.getImageReaders(stream);
    if (readers.hasNext()) {
      var imageReader = readers.next();
      return imageReader.getFormatName();
    }

    return null;
  }
}
