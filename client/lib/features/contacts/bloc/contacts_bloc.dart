import 'package:flutter_bloc/flutter_bloc.dart';
import '../data/contacts_repository.dart';

abstract class ContactsEvent {}
class LoadContacts extends ContactsEvent {}

class ContactsState {
  final bool loading;
  final List<Map<String, dynamic>> contacts;
  final String? error;
  const ContactsState({this.loading = false, this.contacts = const [], this.error});
}

class ContactsBloc extends Bloc<ContactsEvent, ContactsState> {
  final ContactsRepository _repo;
  ContactsBloc(this._repo) : super(const ContactsState()) {
    on<LoadContacts>(_onLoad);
  }

  Future<void> _onLoad(LoadContacts event, Emitter<ContactsState> emit) async {
    emit(const ContactsState(loading: true));
    try {
      final c = await _repo.getContacts();
      emit(ContactsState(contacts: c));
    } catch (e) {
      emit(ContactsState(error: e.toString()));
    }
  }
}
