import 'package:flutter_test/flutter_test.dart';
import 'package:lemon_tea/utils/conversation_manager.dart';

void main() {
  test('generateTitleFromMessage should work correctly', () {
    final manager = ConversationManager();
    
    final multiLineMessage = '第一行\n第二行\n第三行';
    final result = manager.generateTitleFromMessage(multiLineMessage);
    
    print('Input: "$multiLineMessage"');
    print('Output: "$result"');
    print('Expected: "第一行"');
    
    expect(result, equals('第一行'));
  });
} 