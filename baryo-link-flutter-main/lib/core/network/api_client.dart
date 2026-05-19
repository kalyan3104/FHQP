import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter/foundation.dart';

final dioProvider = Provider<Dio>((ref) {
  String baseUrl;

  if (kIsWeb) {
    // Flutter Web MUST use machine IP (NOT 127.0.0.1)
    baseUrl = 'http://localhost:8000'; 
    // or: http://192.168.x.x:8000
  } else {
    // Mobile platforms
    baseUrl = 'http://10.0.2.2:8000'; // Android emulator
  }

  final dio = Dio(
    BaseOptions(
      baseUrl: baseUrl,
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 10),
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
    ),
  );

  dio.interceptors.add(
    LogInterceptor(
      requestBody: true,
      responseBody: true,
    ),
  );

  return dio;
});
