package com.voltaserve.watermark.controllers;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.voltaserve.watermark.dtos.ValuesRequest;

import java.util.Base64;

@RestController
@RequestMapping("/v2/values")
public class ValuesController {

  @PostMapping
  public ResponseEntity<?> create(@RequestBody ValuesRequest request) throws JsonProcessingException {
    return ResponseEntity.ok(
        Base64.getEncoder().encodeToString(
            new ObjectMapper().writeValueAsString(request.getValues()).getBytes()));
  }
}
