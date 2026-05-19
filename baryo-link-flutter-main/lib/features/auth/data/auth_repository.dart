import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/network/api_client.dart';

class AuthRepository {
  final Dio _dio;

  AuthRepository(this._dio);

  Future<Response> sendVerification(String phoneNumber) async {
    return await _dio.post(
      '/auth/send-verification',
      data: {'phone_number': phoneNumber},
    );
  }

  Future<Response> checkVerification({
    required String phoneNumber,
    required String code,
    required String deviceId,
  }) async {
    return await _dio.post(
      '/auth/check-verification',
      data: {
        'phone_number': phoneNumber,
        'code': code,
        'device_id': deviceId,
      },
    );
  }
}

final authRepositoryProvider = Provider<AuthRepository>((ref) {
  final dio = ref.watch(dioProvider);
  return AuthRepository(dio);
});
