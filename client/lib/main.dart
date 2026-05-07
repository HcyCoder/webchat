import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'core/network/dio_client.dart';
import 'core/network/ws_client.dart';
import 'core/auth/token_manager.dart';
import 'features/auth/data/auth_api.dart';
import 'features/auth/data/auth_repository.dart';
import 'features/auth/bloc/auth_bloc.dart';
import 'features/chat/data/chat_api.dart';
import 'features/chat/data/chat_repository.dart';
import 'features/contacts/data/contacts_api.dart';
import 'features/contacts/data/contacts_repository.dart';
import 'app.dart';

void main() {
  final tokenManager = TokenManager();
  final dioClient = DioClient(tokenManager);
  final wsClient = WsClient(tokenManager);
  final authApi = AuthApi(dioClient.dio);
  final authRepo = AuthRepository(authApi, tokenManager);
  final chatApi = ChatApi(dioClient.dio);
  final chatRepo = ChatRepository(chatApi, wsClient);
  final contactsApi = ContactsApi(dioClient.dio);
  final contactsRepo = ContactsRepository(contactsApi);

  runApp(
    MultiRepositoryProvider(
      providers: [
        RepositoryProvider.value(value: chatRepo),
        RepositoryProvider.value(value: contactsRepo),
      ],
      child: BlocProvider(
        create: (_) => AuthBloc(authRepo),
        child: const App(),
      ),
    ),
  );
}
