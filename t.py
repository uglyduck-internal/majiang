import itertools
import hashlib
import time
import json

def calculate_md5(input_string):
    # 创建一个 md5 hash 对象
    m = hashlib.md5()
    # 提供要哈希的数据
    m.update(input_string.encode('utf-8'))
    # 获取哈希值
    md5_hash = m.hexdigest()
    return md5_hash

oem_id = "300ab330835844d58a8bccfc1c8b0800"
lat = "114.37409153788359"
lng = "30.45228049075855"
store_name = ""
city_code = "420100"
page = "1"
limit = "5"
is_included_qipai = "1"
is_included_billiards = "0"
is_experience = "0"
timestamp = "1718361676760"
api_version_interceptor = "1"

map = {
    "lat": lat,
    "oem_id": oem_id,
    "lng": lng,
    "store_name": store_name,
    "city_code": city_code,
    "page": page,
    "limit": limit,
    "is_included_qipai": is_included_qipai,
    "is_included_billiards": is_included_billiards,
    "is_experience": is_experience,
    "timestamp": timestamp,
    "api_version_interceptor": api_version_interceptor
}



# 将字典转换为列表
array = list(map.items())

# 生成所有可能的排列
for permutation in itertools.permutations(array):
    # 将排列转换回字典
    permuted_dict = dict(permutation)
    if calculate_md5(json.dumps(permuted_dict)) == "e5921744e98db5205483c2b74613d8d8":
        print(permuted_dict)
        break