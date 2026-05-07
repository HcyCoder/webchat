import 'package:flutter/material.dart';
import '../chat/ui/conversation_list_page.dart';
import '../contacts/ui/contacts_list_page.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});
  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  int _idx = 0;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: IndexedStack(
        index: _idx,
        children: const [
          ConversationListPage(),
          ContactsListPage(),
          Center(child: Text('discover')),
          Center(child: Text('me')),
        ],
      ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: _idx,
        onTap: (i) => setState(() => _idx = i),
        type: BottomNavigationBarType.fixed,
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.chat_bubble_outline), label: 'chat'),
          BottomNavigationBarItem(icon: Icon(Icons.contacts_outlined), label: 'contacts'),
          BottomNavigationBarItem(icon: Icon(Icons.explore_outlined), label: 'discover'),
          BottomNavigationBarItem(icon: Icon(Icons.person_outline), label: 'me'),
        ],
      ),
    );
  }
}
