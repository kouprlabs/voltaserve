package com.voltaserve.watermark.controllers;

import com.voltaserve.watermark.dtos.WatermarkRequest;
import com.voltaserve.watermark.services.ImageWatermarkService;
import com.voltaserve.watermark.services.PdfWatermarkService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.io.IOException;

@RestController
@RequestMapping("/v2/watermarks")
public class WatermarkController {

  private final PdfWatermarkService pdfWatermarkService;
  private final ImageWatermarkService imageWatermarkService;

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
