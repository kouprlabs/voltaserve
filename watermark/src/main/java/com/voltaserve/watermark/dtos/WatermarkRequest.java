package com.voltaserve.watermark.dtos;

import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class WatermarkRequest {

  private String category;
  private String path;
  private String s3Key;
  private String s3Bucket;
  private String dateTime;
  private String username;
  private String workspace;
}
