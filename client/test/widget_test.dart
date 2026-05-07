import 'package:flutter_test/flutter_test.dart';
import 'package:webchat/app.dart';

void main() {
  testWidgets('app renders', (WidgetTester tester) async {
    await tester.pumpWidget(const App());
    expect(find.text('wechat'), findsOneWidget);
  });
}
