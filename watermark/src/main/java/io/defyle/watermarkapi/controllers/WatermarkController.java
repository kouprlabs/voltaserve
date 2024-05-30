package io.defyle.watermarkapi.controllers;

import io.defyle.watermarkapi.dtos.WatermarkRequest;
import io.defyle.watermarkapi.services.ImageWatermarkService;
import io.defyle.watermarkapi.services.PdfWatermarkService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.io.IOException;

@RestController
@RequestMapping("/watermark")
public class WatermarkController {

  private final PdfWatermarkService pdfWatermarkService;
  private final ImageWatermarkService imageWatermarkService;

  @Autowired
  public WatermarkController(PdfWatermarkService pdfWatermarkService, ImageWatermarkService imageWatermarkService) {
    this.pdfWatermarkService = pdfWatermarkService;
    this.imageWatermarkService = imageWatermarkService;
  }

  @PostMapping
  public ResponseEntity<?> generate(@RequestBody WatermarkRequest watermarkRequest) throws IOException {
    if (watermarkRequest.getFileCategory().equals("document")) {
      pdfWatermarkService.generate(watermarkRequest);
    } else if (watermarkRequest.getFileCategory().equals("image")) {
      imageWatermarkService.generate(watermarkRequest);
    }

    return ResponseEntity.ok().build();
  }
}
