import 'package:flutter/material.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/chat_view.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:lemon_tea/utils/system.dart';

import '../../controls/window_title_bar.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<StatefulWidget> createState() => _HomePage();
}

class _HomePage extends State<HomePage> {
  List<Message> historyMessages = [
    Message(
      role: MessageRole.assistant,
      content: """# 欢迎使用 Markdown

这是一个简单的 Markdown 示例文档，展示常用语法：

## 标题层级
二级标题 (`##`) 到六级标题 (`######`)

## 文字样式
- **加粗文本** (`**加粗**`)
- *斜体文本* (`*斜体*`)
- ~~删除线~~ (`~~删除线~~`)
- `行内代码` (`` `行内代码` ``)

## 列表
### 无序列表
- 项目一
- 项目二
  - 子项目 (缩进两个空格)

### 有序列表
1. 第一项
2. 第二项
   1. 子项 (缩进三个空格)

## 链接与图片
[百度链接](https://www.baidu.com)  
![示例图片](https://via.placeholder.com/150 "悬浮提示")

## 代码块
```python
def hello():
    print("代码高亮示例")""",
    ),
  ];
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      // appBar: ,
      body: Container(
        child: Column(
          children: [
            WindowTitleBar(title: "title"),
            Expanded(
              child: Row(
                children: [
                  Center(
                    child: ChatView(
                      historyMessages: historyMessages,
                      onSend: (msg) {
                        setState(() {
                          historyMessages.add(
                            Message(role: MessageRole.user, content: msg),
                          );
                        });
                      },
                    ),
                  ),
                  Text("data"),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
