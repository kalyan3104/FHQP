import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:device_info_plus/device_info_plus.dart';
import 'dart:io' show Platform;
import '../../data/auth_repository.dart';

class AuthState {
  final bool isLoading;
  final String? error;
  final String? phoneNumber;
  final bool isOtpSent;
  final bool isAuthenticated;

  AuthState({
    this.isLoading = false,
    this.error,
    this.phoneNumber,
    this.isOtpSent = false,
    this.isAuthenticated = false,
  });

  AuthState copyWith({
    bool? isLoading,
    String? error,
    String? phoneNumber,
    bool? isOtpSent,
    bool? isAuthenticated,
  }) {
    return AuthState(
      isLoading: isLoading ?? this.isLoading,
      error: error,
      phoneNumber: phoneNumber ?? this.phoneNumber,
      isOtpSent: isOtpSent ?? this.isOtpSent,
      isAuthenticated: isAuthenticated ?? this.isAuthenticated,
    );
  }
}

class AuthNotifier extends StateNotifier<AuthState> {
  final AuthRepository _repository;
  final _storage = const FlutterSecureStorage();

  AuthNotifier(this._repository) : super(AuthState());

  Future<void> sendOtp(String phoneNumber) async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      await _repository.sendVerification(phoneNumber);
      state = state.copyWith(isLoading: false, isOtpSent: true, phoneNumber: phoneNumber);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: _handleError(e));
    }
  }

  Future<bool> verifyOtp(String code) async {
    if (state.phoneNumber == null) return false;
    
    state = state.copyWith(isLoading: true, error: null);
    try {
      final deviceId = await _getDeviceId();
      final response = await _repository.checkVerification(
        phoneNumber: state.phoneNumber!,
        code: code,
        deviceId: deviceId,
      );

      final token = response.data['data']['token'] ?? response.data['data']; // Adjust based on exact structure
      if (token is String) {
        await _storage.write(key: 'jwt_token', value: token);
      }

      state = state.copyWith(isLoading: false, isAuthenticated: true);
      return true;
    } catch (e) {
      state = state.copyWith(isLoading: false, error: _handleError(e));
      return false;
    }
  }

  Future<String> _getDeviceId() async {
    if (kIsWeb) return 'web-browser-device';
    final deviceInfo = DeviceInfoPlugin();
    if (Platform.isAndroid) {
      final androidInfo = await deviceInfo.androidInfo;
      return androidInfo.id;
    } else if (Platform.isIOS) {
      final iosInfo = await deviceInfo.iosInfo;
      return iosInfo.identifierForVendor ?? 'ios-device';
    }
    return 'other-device';
  }

  String _handleError(dynamic e) {
    if (e is Exception) {
      // Basic extraction of message if available
      return e.toString();
    }
    return 'An unexpected error occurred';
  }
  
  void clearError() => state = state.copyWith(error: null);
}

final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  final repository = ref.watch(authRepositoryProvider);
  return AuthNotifier(repository);
});
