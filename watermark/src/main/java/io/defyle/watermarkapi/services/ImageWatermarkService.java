package io.defyle.watermarkapi.services;

import io.defyle.watermarkapi.dtos.WatermarkRequest;
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
import java.util.Iterator;
import java.util.Objects;

@Service
public class ImageWatermarkService {

  public ImageWatermarkService(ResourceLoader resourceLoader) throws IOException, FontFormatException {
    GraphicsEnvironment graphicsEnvironment = GraphicsEnvironment.getLocalGraphicsEnvironment();
    graphicsEnvironment.registerFont(Font.createFont(Font.TRUETYPE_FONT, ResourcePatternUtils.getResourcePatternResolver(resourceLoader)
        .getResource("classpath:Montserrat-Bold.ttf").getInputStream()));
  }

  public void generate(WatermarkRequest watermarkRequest) throws IOException {
    File inputFile = new File(watermarkRequest.getInputFile());

    BufferedImage bufferedImage = ImageIO.read(inputFile);
    int width = bufferedImage.getWidth();
    int height = bufferedImage.getHeight();
    Graphics2D g2d = (Graphics2D) bufferedImage.getGraphics();

    g2d.setComposite(AlphaComposite.getInstance(AlphaComposite.SRC_OVER, 0.2f));
    g2d.setColor(Color.RED);
    g2d.setRenderingHint(RenderingHints.KEY_TEXT_ANTIALIASING, RenderingHints.VALUE_TEXT_ANTIALIAS_ON);
    g2d.setFont(new Font("Montserrat", Font.BOLD, 120));

    Rectangle2D dateTimeStringBounds = g2d.getFontMetrics().getStringBounds(watermarkRequest.getDateTime(), g2d);
    Rectangle2D workspaceStringBounds = g2d.getFontMetrics().getStringBounds(watermarkRequest.getWorkspace(), g2d);
    Rectangle2D usernameStringBounds = g2d.getFontMetrics().getStringBounds(watermarkRequest.getUsername(), g2d);

    // Move our origin to the center
    g2d.translate(width / 2.0f, height / 2.0f);

    // Create a rotation transform to rotate the text based on the diagonal angle of the picture aspect ratio
    AffineTransform affineTransform = new AffineTransform();
    affineTransform.rotate(Math.atan2(height, width));
    g2d.transform(affineTransform);

    // Reposition our coordinates based on size (same as we would normally
    // do to center on straight line but based on starting at center
    float x1 = (int) workspaceStringBounds.getWidth() / 2.0f * -1;
    float y1 = (int) dateTimeStringBounds.getHeight() / 2.0f;
    y1 -= 240;
    g2d.translate(x1, y1);
    g2d.drawString(watermarkRequest.getWorkspace(), 0.0f, 0.0f);

    float x2 = (int) usernameStringBounds.getWidth() / 2.0f * -1;
    g2d.translate(x2 - x1, 120);
    g2d.drawString(watermarkRequest.getUsername(), 0.0f, 0.0f);

    float x3 = (int) dateTimeStringBounds.getWidth() / 2.0f * -1;
    g2d.translate(x3 - x2, 120);
    g2d.drawString(watermarkRequest.getDateTime(), 0.0f, 0.0f);

    g2d.dispose();

    if (new File(watermarkRequest.getOutputFile()).isFile() && new File(watermarkRequest.getOutputFile()).exists()) {
      if (!new File(watermarkRequest.getOutputFile()).delete()) {
        throw new RuntimeException(String.format("Could not delete '%s'.", watermarkRequest.getOutputFile()));
      }
    }
    ImageIO.write(bufferedImage, Objects.requireNonNull(getFormatName(inputFile)), new BufferedOutputStream(
        new FileOutputStream(new File(watermarkRequest.getOutputFile()))));
  }

  private String getFormatName(File file) throws IOException {
    ImageInputStream imageInputStream = ImageIO.createImageInputStream(file);
    Iterator<ImageReader> imageReaders = ImageIO.getImageReaders(imageInputStream);
    if (imageReaders.hasNext()) {
      ImageReader imageReader = imageReaders.next();
      return imageReader.getFormatName();
    }

    return null;
  }
}
