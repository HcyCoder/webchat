import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class TokenManager {
  final _storage = const FlutterSecureStorage();
  static const _tokenKey = 'auth_token';
  static const _userIdKey = 'user_id';

  Future<void> save(String token, String userId) async {
    await _storage.write(key: _tokenKey, value: token);
    await _storage.write(key: _userIdKey, value: userId);
  }

  Future<String?> get token => _storage.read(key: _tokenKey);
  Future<String?> get userId => _storage.read(key: _userIdKey);
  Future<void> clear() => _storage.deleteAll();
}
