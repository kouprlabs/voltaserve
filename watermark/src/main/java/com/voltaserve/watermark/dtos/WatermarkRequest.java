package com.voltaserve.watermark.dtos;

import java.util.ArrayList;
import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class WatermarkRequest {

    private String category;
    private String path;
    private String s3Key;
    private String s3Bucket;
    private ArrayList<String> values;
}
