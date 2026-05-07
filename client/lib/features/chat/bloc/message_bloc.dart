import 'package:flutter_bloc/flutter_bloc.dart';
import '../data/chat_repository.dart';

abstract class MsgEvent {}
class LoadMessages extends MsgEvent {
  final String convId;
  LoadMessages(this.convId);
}
class SendTextMessage extends MsgEvent {
  final String convId, toId, content;
  SendTextMessage(this.convId, this.toId, this.content);
}
class ReceiveMessage extends MsgEvent {
  final Map<String, dynamic> message;
  ReceiveMessage(this.message);
}

class MsgState {
  final bool loading;
  final List<Map<String, dynamic>> messages;
  final String? error;
  const MsgState({this.loading = false, this.messages = const [], this.error});
}

class MessageBloc extends Bloc<MsgEvent, MsgState> {
  final ChatRepository _repo;
  MessageBloc(this._repo) : super(const MsgState()) {
    on<LoadMessages>(_onLoad);
    on<SendTextMessage>(_onSend);
    on<ReceiveMessage>(_onReceive);
  }

  Future<void> _onLoad(LoadMessages event, Emitter<MsgState> emit) async {
    emit(const MsgState(loading: true));
    try {
      final msgs = await _repo.getMessages(event.convId);
      emit(MsgState(messages: msgs));
    } catch (e) {
      emit(MsgState(error: e.toString()));
    }
  }

  Future<void> _onSend(SendTextMessage event, Emitter<MsgState> emit) async {
    _repo.sendMessage({
      'type': 'chat.message',
      'data': {'chat_type': 'single', 'to_id': event.toId, 'msg_type': 'text', 'content': event.content},
    });
    emit(MsgState(messages: [...state.messages, {
      'from_user': 'me', 'msg_type': 'text', 'content': event.content,
      'created_at': DateTime.now().millisecondsSinceEpoch,
    }]));
  }

  void _onReceive(ReceiveMessage event, Emitter<MsgState> emit) {
    emit(MsgState(messages: [...state.messages, event.message]));
  }
}
