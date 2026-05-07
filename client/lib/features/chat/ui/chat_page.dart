import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/message_bloc.dart';
import 'widgets/message_bubble.dart';
import 'widgets/chat_input_bar.dart';

class ChatPage extends StatefulWidget {
  final String convId, title;
  const ChatPage({super.key, required this.convId, required this.title});

  @override
  State<ChatPage> createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final _ctrl = TextEditingController();

  @override
  void initState() {
    super.initState();
    context.read<MessageBloc>().add(LoadMessages(widget.convId));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(widget.title)),
      body: Column(children: [
        Expanded(
          child: BlocBuilder<MessageBloc, MsgState>(
            builder: (_, state) {
              if (state.loading) return const Center(child: CircularProgressIndicator());
              return ListView.builder(
                itemCount: state.messages.length,
                itemBuilder: (_, i) {
                  final m = state.messages[i];
                  return MessageBubble(
                    isMe: m['from_user'] == 'me',
                    content: m['content'] ?? '',
                  );
                },
              );
            },
          ),
        ),
        ChatInputBar(
          controller: _ctrl,
          onSend: () {
            if (_ctrl.text.trim().isEmpty) return;
            context.read<MessageBloc>().add(SendTextMessage(widget.convId, widget.convId, _ctrl.text.trim()));
            _ctrl.clear();
          },
        ),
      ]),
    );
  }
}
