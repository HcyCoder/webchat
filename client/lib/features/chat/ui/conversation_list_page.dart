import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/conversation_bloc.dart';
import '../data/chat_repository.dart';
import 'chat_page.dart';

class ConversationListPage extends StatelessWidget {
  const ConversationListPage({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => ConversationBloc(context.read<ChatRepository>())..add(LoadConversations()),
      child: BlocBuilder<ConversationBloc, ConvState>(
        builder: (context, state) {
          if (state.loading) return const Center(child: CircularProgressIndicator());
          final convs = state.conversations;
          if (convs.isEmpty) return const Center(child: Text('no conversations'));
          return ListView.builder(
            itemCount: convs.length,
            itemBuilder: (_, i) {
              final c = convs[i];
              return ListTile(
                leading: CircleAvatar(
                  child: Text((c['target_name'] ?? '?')[0]),
                ),
                title: Text(c['target_name'] ?? ''),
                subtitle: Text(c['last_content'] ?? '', maxLines: 1, overflow: TextOverflow.ellipsis),
                trailing: ((c["unread_count"] as int?) ?? 0) > 0
                    ? Badge(label: Text('${c['unread_count']}'))
                    : null,
                onTap: () => Navigator.push(context, MaterialPageRoute(
                  builder: (_) => ChatPage(convId: '${c['target_id']}', title: c['target_name'] ?? ''),
                )),
              );
            },
          );
        },
      ),
    );
  }
}
