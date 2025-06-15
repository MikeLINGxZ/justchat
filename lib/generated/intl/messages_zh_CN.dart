// DO NOT EDIT. This is code generated via package:intl/generate_localized.dart
// This is a library that provides messages for a zh_CN locale. All the
// messages from the main program should be duplicated here with the same
// function name.

// Ignore issues from commonly used lints in this file.
// ignore_for_file:unnecessary_brace_in_string_interps, unnecessary_new
// ignore_for_file:prefer_single_quotes,comment_references, directives_ordering
// ignore_for_file:annotate_overrides,prefer_generic_function_type_aliases
// ignore_for_file:unused_import, file_names, avoid_escaping_inner_quotes
// ignore_for_file:unnecessary_string_interpolations, unnecessary_string_escapes

import 'package:intl/intl.dart';
import 'package:intl/message_lookup_by_library.dart';

final messages = new MessageLookup();

typedef String MessageIfAbsent(String messageStr, List<dynamic> args);

class MessageLookup extends MessageLookupByLibrary {
  String get localeName => 'zh_CN';

  final messages = _notInlinedMessages(_notInlinedMessages);
  static Map<String, Function> _notInlinedMessages(_) => <String, Function>{
    "addAccount": MessageLookupByLibrary.simpleMessage("添加账户"),
    "advancedSettings": MessageLookupByLibrary.simpleMessage("高级设置"),
    "advancedSettingsMissing": MessageLookupByLibrary.simpleMessage(
      "高级设置信息错误，请检查设置",
    ),
    "appName": MessageLookupByLibrary.simpleMessage("Teamail"),
    "authenticationFailedMsg": MessageLookupByLibrary.simpleMessage(
      "身份验证失败，请检查邮箱和密码",
    ),
    "back": MessageLookupByLibrary.simpleMessage("返回"),
    "close": MessageLookupByLibrary.simpleMessage("关闭"),
    "continuE": MessageLookupByLibrary.simpleMessage("继续"),
    "editSettings": MessageLookupByLibrary.simpleMessage("修改设置"),
    "email": MessageLookupByLibrary.simpleMessage("邮箱"),
    "emailFormatError": MessageLookupByLibrary.simpleMessage("邮件格式错误"),
    "emailInvalid": MessageLookupByLibrary.simpleMessage("邮箱格式错误"),
    "emailSettingsNotMatch": MessageLookupByLibrary.simpleMessage(
      "登陆的邮箱与所选邮箱的设置不一致，继续登陆或修改设置？",
    ),
    "enableProxy": MessageLookupByLibrary.simpleMessage("启用代理"),
    "login": MessageLookupByLibrary.simpleMessage("登陆"),
    "loginFailed": MessageLookupByLibrary.simpleMessage("登陆失败"),
    "loginFailedMsg": MessageLookupByLibrary.simpleMessage(
      "登陆失败，请检查邮箱和密码并确认邮件服务器设置正确",
    ),
    "loginMessageMissing": MessageLookupByLibrary.simpleMessage("登陆信息不完整"),
    "noSupportedEmailServer": MessageLookupByLibrary.simpleMessage("不支持的邮件服务"),
    "otherEmail": MessageLookupByLibrary.simpleMessage("其他邮箱"),
    "password": MessageLookupByLibrary.simpleMessage("密码"),
    "passwordMustNotEmpty": MessageLookupByLibrary.simpleMessage("密码不能为空"),
    "pleaseEnterEmail": MessageLookupByLibrary.simpleMessage("请输入邮箱"),
    "pleaseEnterPassword": MessageLookupByLibrary.simpleMessage("请输入密码或授权码"),
    "port": MessageLookupByLibrary.simpleMessage("端口"),
    "portFormatErr": MessageLookupByLibrary.simpleMessage("端口格式错误"),
    "proxyAddress": MessageLookupByLibrary.simpleMessage("代理服务地址"),
    "secureConnection": MessageLookupByLibrary.simpleMessage("安全连接"),
    "server": MessageLookupByLibrary.simpleMessage("服务器"),
    "serverAddressFormatErr": MessageLookupByLibrary.simpleMessage(
      "服务器或IP格式错误",
    ),
    "serverType": MessageLookupByLibrary.simpleMessage("服务器类型"),
    "smtpErrCode211": MessageLookupByLibrary.simpleMessage("系统状态或系统帮助回复"),
    "smtpErrCode214": MessageLookupByLibrary.simpleMessage("帮助信息（HELP命令的响应）"),
    "smtpErrCode220": MessageLookupByLibrary.simpleMessage("服务准备就绪"),
    "smtpErrCode221": MessageLookupByLibrary.simpleMessage("服务关闭传输通道"),
    "smtpErrCode221_2_0_0": MessageLookupByLibrary.simpleMessage("再见"),
    "smtpErrCode235_2_7_0": MessageLookupByLibrary.simpleMessage("认证成功"),
    "smtpErrCode240": MessageLookupByLibrary.simpleMessage("退出"),
    "smtpErrCode250": MessageLookupByLibrary.simpleMessage("请求的邮件操作完成"),
    "smtpErrCode251": MessageLookupByLibrary.simpleMessage("用户不在本地；将转发邮件"),
    "smtpErrCode252": MessageLookupByLibrary.simpleMessage("无法验证用户，但仍会尝试投递邮件"),
    "smtpErrCode334": MessageLookupByLibrary.simpleMessage(
      "服务器质询 - 文本部分包含Base64编码的质询",
    ),
    "smtpErrCode354": MessageLookupByLibrary.simpleMessage("开始邮件输入"),
    "smtpErrCode421": MessageLookupByLibrary.simpleMessage("服务不可用，关闭传输通道"),
    "smtpErrCode432_4_7_12": MessageLookupByLibrary.simpleMessage("需要密码转换"),
    "smtpErrCode450": MessageLookupByLibrary.simpleMessage("请求的邮件操作未执行：邮箱不可用"),
    "smtpErrCode451": MessageLookupByLibrary.simpleMessage(
      "请求的操作中止：处理过程中出现本地错误",
    ),
    "smtpErrCode451_4_4_1": MessageLookupByLibrary.simpleMessage("IMAP服务器不可用"),
    "smtpErrCode452": MessageLookupByLibrary.simpleMessage("请求的操作未执行：系统存储不足"),
    "smtpErrCode454_4_7_0": MessageLookupByLibrary.simpleMessage("临时认证失败"),
    "smtpErrCode455": MessageLookupByLibrary.simpleMessage("服务器无法接受参数"),
    "smtpErrCode500": MessageLookupByLibrary.simpleMessage("语法错误，命令无法识别"),
    "smtpErrCode500_5_5_6": MessageLookupByLibrary.simpleMessage("认证交换行过长"),
    "smtpErrCode501": MessageLookupByLibrary.simpleMessage("参数或参数语法错误"),
    "smtpErrCode501_5_5_2": MessageLookupByLibrary.simpleMessage(
      "无法Base64解码客户端响应",
    ),
    "smtpErrCode501_5_7_0": MessageLookupByLibrary.simpleMessage("客户端发起的认证交换"),
    "smtpErrCode502": MessageLookupByLibrary.simpleMessage("命令未实现"),
    "smtpErrCode503": MessageLookupByLibrary.simpleMessage("命令顺序错误"),
    "smtpErrCode504": MessageLookupByLibrary.simpleMessage("命令参数未实现"),
    "smtpErrCode504_5_5_4": MessageLookupByLibrary.simpleMessage("无法识别的认证类型"),
    "smtpErrCode521": MessageLookupByLibrary.simpleMessage("服务器不接受邮件"),
    "smtpErrCode523": MessageLookupByLibrary.simpleMessage("需要加密"),
    "smtpErrCode530_5_7_0": MessageLookupByLibrary.simpleMessage("需要认证"),
    "smtpErrCode534_5_7_9": MessageLookupByLibrary.simpleMessage("认证机制强度不足"),
    "smtpErrCode535_5_7_8": MessageLookupByLibrary.simpleMessage("认证凭据无效"),
    "smtpErrCode550": MessageLookupByLibrary.simpleMessage("请求的操作未执行：邮箱不可用"),
    "smtpErrCode551": MessageLookupByLibrary.simpleMessage("用户不在本地；请尝试<转发路径>"),
    "smtpErrCode552": MessageLookupByLibrary.simpleMessage("请求的邮件操作中止：超出存储分配"),
    "smtpErrCode553": MessageLookupByLibrary.simpleMessage("请求的操作未执行：邮箱名称不允许"),
    "smtpErrCode554": MessageLookupByLibrary.simpleMessage(
      "交易失败（或\'此处无SMTP服务\'）",
    ),
    "smtpErrCode554_5_3_4": MessageLookupByLibrary.simpleMessage("邮件过大超出系统限制"),
    "smtpErrCode556": MessageLookupByLibrary.simpleMessage("域名不接受邮件"),
    "unknownError": MessageLookupByLibrary.simpleMessage("未知错误"),
  };
}
