import 'package:dio/dio.dart';
import '../../../core/network/api_paths.dart';

class ContactsApi {
  final Dio _dio;
  ContactsApi(this._dio);

  Future<List<dynamic>> list() async {
    final resp = await _dio.get(ApiPaths.contacts);
    return resp.data['contacts'] as List<dynamic>;
  }

  Future<void> addRequest(String toUserId, String message) async {
    await _dio.post(ApiPaths.contactRequest, data: {'to_user': toUserId, 'message': message});
  }

  Future<void> handleRequest(String requestId, String action) async {
    await _dio.put('${ApiPaths.contactRequest}/$requestId', data: {'action': action});
  }
}
