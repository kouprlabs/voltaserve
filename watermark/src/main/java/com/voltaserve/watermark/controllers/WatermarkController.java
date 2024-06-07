package com.voltaserve.watermark.controllers;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
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
import java.util.ArrayList;
import java.util.UUID;
import java.util.Base64;
import org.springframework.web.bind.annotation.PostMapping;

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
    public ResponseEntity<?> create(
            @RequestParam("file") MultipartFile file,
            @RequestParam("category") String category,
            @RequestParam("s3_key") String s3Key,
            @RequestParam("s3_bucket") String s3Bucket,
            @RequestParam("values") String values)
            throws IOException {

        Path path = Paths.get(
                System.getProperty("java.io.tmpdir"),
                UUID.randomUUID()
                        + "."
                        + FilenameUtils.getExtension(file.getOriginalFilename()));
        file.transferTo(path.toFile());

        WatermarkRequest request = WatermarkRequest.builder()
                .path(path.toString())
                .category(category)
                .s3Key(s3Key)
                .s3Bucket(s3Bucket)
                .values(new ObjectMapper().readValue(
                        new String(Base64.getDecoder().decode(values)),
                        new TypeReference<ArrayList<String>>() {
                        }))
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
