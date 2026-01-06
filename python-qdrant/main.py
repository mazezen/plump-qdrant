from  qdrant_client import QdrantClient
from pydantic import BaseModel
from typing import Any, Type, TypeVar
from pydantic.version import VERSION as PYDANTIC_VERSION

from qdrant_client.models import (VectorParams, Distance, PointStruct, models)

def main():
    
    client = QdrantClient(
        url="",
        grpc_port=6333,
        api_key="",
    )

    collection_name: str = "python-qdrant-test"
    flag = False

    # ------------------- 集合列表 ----------------------
    collections = client.get_collections().collections
    for collection in collections:
        # print(to_dict(collection)['name'])
        if (collection_name == to_dict(collection)['name']):
            flag = True

    # ------------------- 创建集合 ----------------------
    if not flag:
        client.create_collection(collection_name=collection_name,
                                     vectors_config=VectorParams(size=100, distance=Distance.DOT),
                                     timeout=60)

    collections_response = client.get_collections().collections      
    for collection in collections_response:
        print(to_dict(collection)['name']) 


    # ------------------- 获取集合 ----------------------
    collection_info = client.get_collection(collection_name=collection_name)
    print("集合状态: ", collection_info.status)
    print("集合point数量: ", collection_info.points_count)

    # ------------------- 插入数据 ----------------------
    updateResult = client.upsert(collection_name=collection_name, wait=True, points=[PointStruct(id=2, payload={"hello": "world"}, vector=[0.1, 0.2, 0.3, 0.4, 0.5]*20)])
    print(updateResult)

    # ------------------- 检索数据 ----------------------
    listRecord = client.retrieve(
        collection_name=collection_name,
        ids=[1, 2, 3],
        with_payload=True,
        with_vectors=True,
    )
    for record in listRecord:
        print("record: ", record)


     # ------------------- 检索数据 ----------------------
    query_response = client.query_points(
        collection_name=collection_name,
        query=2
    )
    print("query_response: ", query_response)


    sampled = client.query_points(
        collection_name=collection_name,
        query=models.SampleQuery(sample=models.Sample.RANDOM)
    )
    print("sampled: ",sampled)

    return

PYDANTIC_V2 = PYDANTIC_VERSION.startswith("2.")
Model = TypeVar("Model", bound="BaseModel")


def to_dict(model: BaseModel, *args: Any, **kwargs: Any) -> dict[Any, Any]:
    if PYDANTIC_V2:
        return model.model_dump(*args, **kwargs)
    else:
        return model.dict(*args, **kwargs)


if __name__ == '__main__':
    main()