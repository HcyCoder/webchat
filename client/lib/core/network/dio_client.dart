import 'package:dio/dio.dart';
import '../auth/token_manager.dart';
import 'api_paths.dart';

class DioClient {
  late final Dio dio;
  final TokenManager _tokenManager;

  DioClient(this._tokenManager) {
    dio = Dio(BaseOptions(
      baseUrl: ApiPaths.baseUrl,
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 10),
      headers: {'Content-Type': 'application/json'},
    ));
    dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        final token = await _tokenManager.token;
        if (token != null) {
          options.headers['Authorization'] = 'Bearer $token';
        }
        handler.next(options);
      },
      onError: (error, handler) async {
        if (error.response?.statusCode == 401) {
          final token = await _tokenManager.token;
          if (token != null) {
            final newToken = await _refreshToken();
            if (newToken != null) {
              error.requestOptions.headers['Authorization'] = 'Bearer $newToken';
              final retry = await dio.fetch(error.requestOptions);
              return handler.resolve(retry);
            }
          }
          await _tokenManager.clear();
        }
        handler.next(error);
      },
    ));
  }

  Future<String?> _refreshToken() async {
    try {
      final resp = await dio.post(ApiPaths.refreshToken);
      final token = resp.data['token'] as String;
      final userId = resp.data['user_id'] as String;
      await _tokenManager.save(token, userId);
      return token;
    } catch (_) {
      return null;
    }
  }
}
