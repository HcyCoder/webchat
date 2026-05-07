import 'package:dio/dio.dart';
import '../../../core/network/api_paths.dart';

class ChatApi {
  final Dio _dio;
  ChatApi(this._dio);

  Future<List<dynamic>> getConversations() async {
    final resp = await _dio.get(ApiPaths.conversations);
    return resp.data['conversations'] as List<dynamic>;
  }

  Future<List<dynamic>> getMessages(String convId, {int page = 1, int pageSize = 20}) async {
    final resp = await _dio.get('${ApiPaths.messages}/$convId', queryParameters: {'page': page, 'page_size': pageSize});
    return resp.data['messages'] as List<dynamic>;
  }
}
