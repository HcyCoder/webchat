import 'package:flutter/material.dart';

class MessageBubble extends StatelessWidget {
  final bool isMe;
  final String content;
  const MessageBubble({super.key, required this.isMe, required this.content});

  @override
  Widget build(BuildContext context) {
    return Align(
      alignment: isMe ? Alignment.centerRight : Alignment.centerLeft,
      child: Container(
        margin: const EdgeInsets.symmetric(vertical: 4, horizontal: 12),
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        decoration: BoxDecoration(
          color: isMe ? const Color(0xFF95EC69) : Colors.white,
          borderRadius: BorderRadius.circular(4),
        ),
        child: Text(content, style: const TextStyle(fontSize: 16)),
      ),
    );
  }
}
