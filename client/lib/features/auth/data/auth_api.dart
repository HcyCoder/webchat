import 'package:dio/dio.dart';
import '../../../core/network/api_paths.dart';

class AuthApi {
  final Dio _dio;
  AuthApi(this._dio);

  Future<Map<String, dynamic>> login(String phone, String password) async {
    final resp = await _dio.post(ApiPaths.login, data: {'phone': phone, 'password': password});
    return resp.data;
  }

  Future<Map<String, dynamic>> register(String phone, String password, String nickname) async {
    final resp = await _dio.post(ApiPaths.register, data: {'phone': phone, 'password': password, 'nickname': nickname});
    return resp.data;
  }
}
