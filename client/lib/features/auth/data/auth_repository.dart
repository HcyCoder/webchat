import 'auth_api.dart';
import '../../../core/auth/token_manager.dart';

class AuthRepository {
  final AuthApi _api;
  final TokenManager _tokenManager;
  AuthRepository(this._api, this._tokenManager);

  Future<void> login(String phone, String password) async {
    final data = await _api.login(phone, password);
    await _tokenManager.save(data['token'], data['user_id']);
  }

  Future<void> register(String phone, String password, String nickname) async {
    final data = await _api.register(phone, password, nickname);
    await _tokenManager.save(data['token'], data['user_id']);
  }

  Future<bool> isLoggedIn() async => await _tokenManager.token != null;
  Future<void> logout() => _tokenManager.clear();
}
