import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/contacts_bloc.dart';
import '../data/contacts_repository.dart';

class ContactsListPage extends StatelessWidget {
  const ContactsListPage({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => ContactsBloc(context.read<ContactsRepository>())..add(LoadContacts()),
      child: BlocBuilder<ContactsBloc, ContactsState>(
        builder: (_, state) {
          if (state.loading) return const Center(child: CircularProgressIndicator());
          return ListView.builder(
            itemCount: state.contacts.length,
            itemBuilder: (_, i) {
              final c = state.contacts[i];
              return ListTile(
                leading: CircleAvatar(child: Text((c['nickname'] ?? '?')[0])),
                title: Text(c['nickname'] ?? ''),
              );
            },
          );
        },
      ),
    );
  }
}
