package com.voltaserve.watermark.dtos;

import lombok.Data;

import jakarta.validation.constraints.NotBlank;

@Data
public class WatermarkRequest {

  @NotBlank
  private String dateTime;

  @NotBlank
  private String username;

  @NotBlank
  private String workspace;

  @NotBlank
  private String inputFile;

  @NotBlank
  private String outputFile;

  @NotBlank
  private String fileCategory;
}
