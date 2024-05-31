package com.voltaserve.watermark.controllers;

import com.voltaserve.watermark.dtos.WatermarkRequest;
import com.voltaserve.watermark.services.ImageWatermarkService;
import com.voltaserve.watermark.services.PdfWatermarkService;
import org.apache.commons.io.FilenameUtils;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.UUID;

@RestController
@RequestMapping("/v2/watermarks")
public class WatermarkController {

    private final PdfWatermarkService pdfWatermarkService;
    private final ImageWatermarkService imageWatermarkService;

    public WatermarkController(
            PdfWatermarkService pdfWatermarkService, ImageWatermarkService imageWatermarkService) {
        this.pdfWatermarkService = pdfWatermarkService;
        this.imageWatermarkService = imageWatermarkService;
    }

    @PostMapping
    public ResponseEntity<?> generate(
            @RequestParam("file") MultipartFile file,
            @RequestParam("category") String category,
            @RequestParam("s3_key") String s3Key,
            @RequestParam("s3_bucket") String s3Bucket,
            @RequestParam("date_time") String dateTime,
            @RequestParam("username") String username,
            @RequestParam("workspace") String workspace)
            throws IOException {

        Path path =
                Paths.get(
                        System.getProperty("java.io.tmpdir"),
                        UUID.randomUUID()
                                + "."
                                + FilenameUtils.getExtension(file.getOriginalFilename()));
        file.transferTo(path.toFile());

        var request =
                WatermarkRequest.builder()
                        .path(path.toString())
                        .category(category)
                        .s3Key(s3Key)
                        .s3Bucket(s3Bucket)
                        .dateTime(dateTime)
                        .username(username)
                        .workspace(workspace)
                        .build();
        try {
            if (request.getCategory().equals("document")) {
                pdfWatermarkService.generate(request);
            } else if (request.getCategory().equals("image")) {
                imageWatermarkService.generate(request);
            }
        } finally {
            Files.delete(path);
        }
        return ResponseEntity.ok().build();
    }
}
