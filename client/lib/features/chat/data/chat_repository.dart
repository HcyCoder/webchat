import 'chat_api.dart';
import '../../../core/network/ws_client.dart';

class ChatRepository {
  final ChatApi _api;
  final WsClient _wsClient;
  ChatRepository(this._api, this._wsClient);

  Future<List<Map<String, dynamic>>> getConversations() async {
    final list = await _api.getConversations();
    return list.cast<Map<String, dynamic>>();
  }

  Future<List<Map<String, dynamic>>> getMessages(String convId) async {
    final list = await _api.getMessages(convId);
    return list.cast<Map<String, dynamic>>();
  }

  void sendMessage(Map<String, dynamic> msg) => _wsClient.send(msg);
  void connectWebSocket() => _wsClient.connect();
}
