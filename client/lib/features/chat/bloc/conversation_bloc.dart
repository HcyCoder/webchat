import 'package:flutter_bloc/flutter_bloc.dart';
import '../data/chat_repository.dart';

abstract class ConvEvent {}
class LoadConversations extends ConvEvent {}

class ConvState {
  final bool loading;
  final List<Map<String, dynamic>> conversations;
  final String? error;
  const ConvState({this.loading = false, this.conversations = const [], this.error});
}

class ConversationBloc extends Bloc<ConvEvent, ConvState> {
  final ChatRepository _repo;
  ConversationBloc(this._repo) : super(const ConvState()) {
    on<LoadConversations>(_onLoad);
  }

  Future<void> _onLoad(LoadConversations event, Emitter<ConvState> emit) async {
    emit(const ConvState(loading: true));
    try {
      final convs = await _repo.getConversations();
      emit(ConvState(conversations: convs));
    } catch (e) {
      emit(ConvState(error: e.toString()));
    }
  }
}
