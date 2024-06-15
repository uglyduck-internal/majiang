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

timestamp=1718365706908
api_version_interceptor=1

sign="a6dec476d681a40038cb3a07cc293533"

map = {
    "timestamp": timestamp,
    "api_version_interceptor": api_version_interceptor,
}

# 将字典转换为列表
array = list(map.items())

# 生成所有可能的排列
for permutation in itertools.permutations(array):
    # 将排列转换回字典
    permuted_dict = dict(permutation)
    if calculate_md5(json.dumps(permuted_dict)) == sign:
        print(permuted_dict)
        break