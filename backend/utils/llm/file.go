package llm

import (
	"encoding/base64"
	"strings"
)

func uploadFile2Aliyuns(uploadBaseUrl, uploadKey string, filePath string) (string, error) {
	// path := strings.Replace(filePath, "@local:", "", 1)
	// todo 读取文件，
	// demo:
	// def get_upload_policy(api_key, model_name):
	//    """获取文件上传凭证"""
	//    url = "https://dashscope.aliyuncs.com/api/v1/uploads"
	//    headers = {
	//        "Authorization": f"Bearer {api_key}",
	//        "Content-Type": "application/json"
	//    }
	//    params = {
	//        "action": "getPolicy",
	//        "model": model_name
	//    }
	//
	//    response = requests.get(url, headers=headers, params=params)
	//    if response.status_code != 200:
	//        raise Exception(f"Failed to get upload policy: {response.text}")
	//
	//    return response.json()['data']
	//
	//def upload_file_to_oss(policy_data, file_path):
	//    """将文件上传到临时存储OSS"""
	//    file_name = Path(file_path).name
	//    key = f"{policy_data['upload_dir']}/{file_name}"
	//
	//    with open(file_path, 'rb') as file:
	//        files = {
	//            'OSSAccessKeyId': (None, policy_data['oss_access_key_id']),
	//            'Signature': (None, policy_data['signature']),
	//            'policy': (None, policy_data['policy']),
	//            'x-oss-object-acl': (None, policy_data['x_oss_object_acl']),
	//            'x-oss-forbid-overwrite': (None, policy_data['x_oss_forbid_overwrite']),
	//            'key': (None, key),
	//            'success_action_status': (None, '200'),
	//            'file': (file_name, file)
	//        }
	//
	//        response = requests.post(policy_data['upload_host'], files=files)
	//        if response.status_code != 200:
	//            raise Exception(f"Failed to upload file: {response.text}")
	//
	//    return f"oss://{key}"
	//
	//def upload_file_and_get_url(api_key, model_name, file_path):
	//    """上传文件并获取URL"""
	//    # 1. 获取上传凭证，上传凭证接口有限流，超出限流将导致请求失败
	//    policy_data = get_upload_policy(api_key, model_name)
	//    # 2. 上传文件到OSS
	//    oss_url = upload_file_to_oss(policy_data, file_path)
	//
	//    return oss_url
	return "", nil
}

func imgbase64(filePath string) (string, error) {
	path := strings.Replace(filePath, "@local:", "", 1)
	// todo 读取图片，封装为格式：data:[<mediatype>][;base64],<data> 的字符串返回
	return base64.StdEncoding.EncodeToString([]byte(path)), nil
}
