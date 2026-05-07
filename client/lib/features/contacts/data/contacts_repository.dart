import 'contacts_api.dart';

class ContactsRepository {
  final ContactsApi _api;
  ContactsRepository(this._api);

  Future<List<Map<String, dynamic>>> getContacts() async {
    final list = await _api.list();
    return list.cast<Map<String, dynamic>>();
  }

  Future<void> addContact(String userId, String msg) => _api.addRequest(userId, msg);
}
