import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';
import '../auth/token_manager.dart';
import 'api_paths.dart';

typedef WsMessageHandler = void Function(Map<String, dynamic> message);

class WsClient {
  WebSocketChannel? _channel;
  final TokenManager _tokenManager;
  WsMessageHandler? onMessage;

  WsClient(this._tokenManager);

  Future<void> connect() async {
    final token = await _tokenManager.token;
    if (token == null) return;
    final uri = Uri.parse('${ApiPaths.wsBase}?token=$token');
    _channel = WebSocketChannel.connect(uri);
    _channel!.stream.listen(
      (data) => onMessage?.call(json.decode(data) as Map<String, dynamic>),
      onError: (_) => _reconnect(),
      onDone: () => _reconnect(),
    );
  }

  void send(Map<String, dynamic> message) {
    _channel?.sink.add(json.encode(message));
  }

  void _reconnect() => Future.delayed(const Duration(seconds: 3), connect);

  void dispose() => _channel?.sink.close();
}
